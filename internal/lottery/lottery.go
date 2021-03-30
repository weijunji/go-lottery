package lottery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fastrand"
	"github.com/weijunji/go-lottery/pkgs/lock"
	"github.com/weijunji/go-lottery/pkgs/middleware"
	"github.com/weijunji/go-lottery/pkgs/utils"
	pb "github.com/weijunji/go-lottery/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"strconv"
	"time"
)

const RATE_DENOMINATOR = 1000000

func SetupRouter(_ *gin.RouterGroup, authGroup *gin.RouterGroup) {
	authGroup.GET("/once", lotteryOnce)
}

func lotteryOnce(c *gin.Context) {
	info, _ := c.Get("userinfo")
	user := info.(middleware.Userinfo).ID
	idStr := c.Query("id")
	if idStr == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	li, err := strconv.Atoi(idStr)
	lottery := uint64(li)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// TODO: lottery is started ?
	duration := getLotteryDuration(c, lottery)
	if duration == nil {
		// no such lottery
		c.Status(http.StatusNotFound)
		return
	}
	now := time.Now()
	if !(now.After(duration.GetStart().AsTime()) && now.Before(duration.GetEnd().AsTime())) {
		// lottery not start
		c.Status(http.StatusNotFound)
		return
	}

	userLock := lock.NewDistributeLock(c, "user", user)
	if !userLock.Lock(time.Second * 30) {
		// user have been locked
		c.Status(http.StatusConflict)
		return
	}
	defer userLock.UnLock()

	log.WithFields(log.Fields{"user": user, "lottery": lottery}).Info("Lottery once")

	// get rate from redis
	rates := getRate(c, lottery)
	if rates == nil {
		c.Status(http.StatusNotFound)
		return
	}

	// process lottery
	if !decreaseTimes(c, lottery, user) {
		c.JSON(200, gin.H{"win": false, "message": "no lottery times"})
		return
	}
	award := processLottery(rates)
	log.WithFields(log.Fields{"user": user, "lottery": lottery, "result": award}).Info("lottery result")
	if award == nil {
		c.JSON(200, gin.H{"win": false, "message": "no award"})
	} else {
		var decrOk bool
		if award.GetValue() == pb.LotteryRates_HIGH_VAL {
			decrOk = decreaseAwardHighVal(c, award.GetId())
		} else {
			decrOk = decreaseAwardLowVal(c, award.GetId())
		}
		if !decrOk {
			c.JSON(200, gin.H{"win": false, "message": "no award"})
			return
		}
		log.WithFields(log.Fields{"user": user, "lottery": lottery, "award": award}).Infof("Win an award")
		sendMessage(user, lottery, award.GetId())
		//db.ExecContext(c, "INSERT INTO winning_infos(user, award, lottery) VALUES (?, ?, ?);", user, award, lottery)
		c.JSON(200, gin.H{"win": true, "award": award})
	}
}

func decreaseAwardHighVal(ctx context.Context, award uint64) bool {
	logger := log.WithField("award", award)
	db, err := utils.GetMysql().DB()
	if err != nil {
		logger.Fatal("Get db failed", err)
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		logger.Fatal("Begin tx failed", err)
	}
	var remain uint64
	err = tx.QueryRow("SELECT remain FROM awards WHERE award = ? FOR UPDATE", award).Scan(&remain)
	if err != nil {
		logger.Fatal("Select failed", err)
	}
	log.Info("remain: ", remain)
	if remain > 0 {
		_, err = tx.Exec("UPDATE awards SET remain = ? WHERE award = ?", remain - 1, award)
		if err != nil {
			_ = tx.Rollback()
			logger.Fatal("Update failed")
		}
		if err = tx.Commit(); err != nil {
			logger.Fatal("Commit failed")
		}
		return true
	} else {
		_ = tx.Rollback()
		return false
	}
}

func decreaseAwardLowVal(ctx context.Context, award uint64) bool {
	key := fmt.Sprintf("awards:%d", award)
	rds := utils.GetRedis()
	remain, err := rds.Decr(ctx, key).Result()
	if err != nil {
		log.WithField("award", award).Fatal("Decrease award failed", err)
	}
	return remain >= 0
}

func sendMessage(user, lottery, award uint64) {
	info := pb.WinningInfo{User: user, Lottery: lottery, Award: award}
	bs, _ := proto.Marshal(&info)

	producer := utils.GetKafkaProducer()
	log.Info(producer)
	msg := &sarama.ProducerMessage{}
	msg.Topic = "WinningTopic"
	msg.Value = sarama.ByteEncoder(bs)
	pid, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatal("Send message failed", err)
	} else {
		log.WithFields(log.Fields{"pid": pid, "offset": offset}).Info("Send message success")
	}
}

func getRate(ctx context.Context,id uint64) *pb.LotteryRates {
	logger := log.WithField("lottery", id)
	key := fmt.Sprintf("rate:%d", id)
	rds := utils.GetRedis()
	bytes, err := rds.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存击穿
			logger.Info("lottery rate cache miss")
			l := lock.NewDistributeLock(ctx, "rate", id)
			if l.Lock(time.Millisecond * 100) {
				defer l.UnLock()
				logger.Info("lottery rate lock success")
				db, _ := utils.GetMysql().DB()
				rows, err := db.QueryContext(ctx, "SELECT id, rate, value FROM award_infos WHERE lottery = ?", id)
				if err != nil {
					l.UnLock()
					logger.Fatal("Get rate from mysql failed: ", err)
				}
				defer rows.Close()

				rates := new(pb.LotteryRates)
				for rows.Next() {
					var id uint64
					var rate, value uint32
					err := rows.Scan(&id, &rate, &value)
					if err != nil {
						logger.Fatal("Get rate failed", err)
					}
					rates.Total += rate

					var r pb.LotteryRates_AwardRate
					r.Id = id
					r.Rate = rate
					if value == 1 {
						r.Value = pb.LotteryRates_HIGH_VAL
					} else {
						r.Value = pb.LotteryRates_LOW_VAL
					}
					rates.Rates = append(rates.Rates, &r)
				}

				if rates.GetTotal() == 0 {
					return nil
				}
				// update redis
				bs, _ := proto.Marshal(rates)
				rds.Set(ctx, key, bs, time.Minute * 11)
				return rates
			} else {
				logger.Info("lottery times lock failed")
				time.Sleep(time.Millisecond * 50)
				return getRate(ctx, id)
			}
		} else {
			log.WithField("id", id).Fatal("Get rate from redis failed: ", err)
		}
	}

	rate := new(pb.LotteryRates)
	if err := proto.Unmarshal(bytes, rate); err != nil {
		log.WithField("id", id).Fatal("Unmarshal rate message failed: ", err)
	}
	return rate
}

// lottery, return award id, 0 for no award
func processLottery(rates *pb.LotteryRates) *pb.LotteryRates_AwardRate {
	randNum := fastrand.Uint32n(RATE_DENOMINATOR)
	if randNum < rates.GetTotal() {
		var count uint32 = 0
		for _, award := range rates.GetRates() {
			count += award.GetRate()
			if count >= randNum {
				return award
			}
		}
		log.Error("Total rate is not consist with awards' rate")
		return nil
	}
	return nil
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func decreaseTimes(ctx context.Context, lottery uint64, user uint64) bool {
	logger := log.WithFields(log.Fields{"lottery": lottery, "user": user})
	rds := utils.GetRedis()
	key := fmt.Sprintf("remain:%d:%d", lottery, user)
	bytes, err := rds.Get(ctx, key).Bytes()
	times := new(pb.UserTimes)
	var lt *pb.LotteryTimes
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logger.Info("no times in redis")
			lt = getLotteryTimes(ctx, lottery)
			times.Permanent = lt.Permanent
			times.Temporary = lt.Temporary
			times.Update = timestamppb.Now()
		} else {
			logger.Fatal("Get times from redis failed: ", err)
		}
	} else {
		if err := proto.Unmarshal(bytes, times); err != nil {
			logger.Fatal("Unmarshal times failed: ", err)
		}
	}

	update := times.GetUpdate().AsTime()
	if !DateEqual(time.Now(), update) {
		// update temporary times
		if lt == nil {
			lt = getLotteryTimes(ctx, lottery)
		}
		times.Update = timestamppb.Now()
		times.Temporary = lt.Temporary
	}
	// decrease times
	if times.GetTemporary() > 0 {
		times.Temporary--
	} else if times.GetPermanent() > 0 {
		times.Permanent--
	} else {
		return false
	}

	// update redis
	nb, err := proto.Marshal(times)
	if err != nil {
		log.Fatal("Marshal times failed")
	}
	err = rds.Set(ctx, key, nb, 0).Err()
	if err != nil {
		log.Fatal("Update times failed")
	}
	return true
}

func getLotteryTimes(ctx context.Context, lottery uint64) (times *pb.LotteryTimes) {
	logger := log.WithField("lottery", lottery)
	rds := utils.GetRedis()
	key := fmt.Sprintf("lottery_times:%d", lottery)
	res, err := rds.Get(ctx, key).Bytes()
	times = new(pb.LotteryTimes)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存击穿
			log.WithField("lottery", lottery).Info("lottery times cache miss")
			l := lock.NewDistributeLock(ctx, "lt", lottery)
			if l.Lock(time.Millisecond * 100) {
				defer l.UnLock()
				log.WithField("lottery", lottery).Info("lottery times lock success")
				db, _ := utils.GetMysql().DB()
				var perm, temp uint32
				err := db.QueryRowContext(ctx, "SELECT permanent, temporary FROM lotteries WHERE id = ?", lottery).Scan(&perm, &temp)
				if err != nil {
					l.UnLock()
					log.Fatal("Get times from mysql failed: ", err)
				}
				times.Temporary = temp
				times.Permanent = perm
				// update redis
				bs, _ := proto.Marshal(times)
				rds.Set(ctx, key, bs, time.Minute * 11)
				return
			} else {
				logger.Info("lottery times lock failed")
				time.Sleep(time.Millisecond * 50)
				return getLotteryTimes(ctx, lottery)
			}
		} else {
			logger.Fatal("Get lottery times failed")
		}
	}
	if err = proto.Unmarshal(res, times); err != nil {
		logger.Fatal("Unmarshal times failed")
	}
	return
}

func getLotteryDuration(ctx context.Context, lottery uint64) (duration *pb.LotteryDuration) {
	logger := log.WithField("lottery", lottery)
	rds := utils.GetRedis()
	key := fmt.Sprintf("lottery_duration:%d", lottery)
	res, err := rds.Get(ctx, key).Bytes()
	duration = new(pb.LotteryDuration)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			// 缓存击穿
			log.WithField("lottery", lottery).Info("lottery duration cache miss")
			l := lock.NewDistributeLock(ctx, "ld", lottery)
			if l.Lock(time.Millisecond * 100) {
				defer l.UnLock()
				log.WithField("lottery", lottery).Info("lottery duration lock success")
				db, _ := utils.GetMysql().DB()
				var start, end time.Time
				err := db.QueryRowContext(ctx, "SELECT start_time, end_time FROM lotteries WHERE id = ?", lottery).Scan(&start, &end)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						// TODO: 缓存穿透
						return nil
					} else {
						log.Fatal("Get duration from mysql failed: ", err)
					}
				}
				duration.Start = timestamppb.New(start)
				duration.End = timestamppb.New(end)
				// update redis
				bs, _ := proto.Marshal(duration)
				rds.Set(ctx, key, bs, time.Minute * 10)
				return
			} else {
				logger.Info("lottery duration lock failed")
				time.Sleep(time.Millisecond * 50)
				return getLotteryDuration(ctx, lottery)
			}
		} else {
			logger.Fatal("Get duration failed")
		}
	}
	if err = proto.Unmarshal(res, duration); err != nil {
		logger.Fatal("Unmarshal times failed")
	}
	return
}
