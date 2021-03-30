package lock

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"time"
)

type DistributeLock struct {
	context context.Context
	name string
}

func NewDistributeLock(ctx context.Context, name string, id uint64) DistributeLock {
	return DistributeLock{ctx, fmt.Sprintf("%s_lock:%d", name, id)}
}

func (l DistributeLock) Lock(expire time.Duration) bool {
	ok, err := utils.GetRedis().SetNX(l.context, l.name, 0, expire).Result()
	if err != nil {
		log.WithField("lock", l).Fatal("Lock failed: ", err)
	}
	return ok
}

func (l DistributeLock) UnLock() {
	utils.GetRedis().Del(l.context, l.name)
}
