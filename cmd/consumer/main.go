package main

import (
	"sync"

	"github.com/Shopify/sarama"
	pb "github.com/weijunji/go-lottery/proto"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func main() {
	wg := sync.WaitGroup{}
	topic := "WinningTopic"
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		log.Fatalf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		log.Fatalf("fail to get list of partition:err%v\n", err)
		return
	}
	log.Info(partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		wg.Add(1)
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			log.Info("getting message...\n")
			for msg := range pc.Messages() {
				info := &pb.WinningInfo{}
				if err := proto.Unmarshal(msg.Value, info); err != nil {
					log.Fatal("Unmarshal failed")
				}
				log.Infof("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, info)
				db := utils.GetMysql()
				db.Exec("INSERT INTO winning_infos(user, award, lottery) VALUES (?, ?, ?);", info.User, info.Award, info.Lottery)
			}
			wg.Done()
		}(pc)
	}
	wg.Wait()
}
