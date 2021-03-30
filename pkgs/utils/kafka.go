package utils

import (
	"github.com/Shopify/sarama"
	"sync"
)

var producer sarama.SyncProducer
var producer_once sync.Once

func GetKafkaProducer() sarama.SyncProducer {
	producer_once.Do(func() {
		config := sarama.NewConfig()
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Partitioner = sarama.NewRandomPartitioner
		config.Producer.Return.Successes = true

		p, err := sarama.NewSyncProducer(getKafkaAddr(), config)
		if err != nil {
			panic("Connect kafka failed")
		}

		producer = p

		// msg := &sarama.ProducerMessage{}
		// msg.Topic = "TestTopic"
		// msg.Value = sarama.StringEncoder("this is a test")
		// pid, offset, err := client.SendMessage(msg)
	})
	return producer
}
