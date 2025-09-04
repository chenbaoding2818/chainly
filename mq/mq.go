package mq

import (
	"context"
	"sync"

	"github.com/chenbaoding2818/chainly/config"
	iface "github.com/chenbaoding2818/chainly/interface"
)

var (
	Producer iface.IProducer
)

func NewProducer(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config) {
	switch cfg.MQCfg.MQType {
	case iface.MQTypeRabbitMQ:
		Producer = NewRabbitMQProducer(ctx, wg, cfg.MQCfg.RabbitMQ)
	case iface.MQTypeKafka: // TODO: 待实现
		// Producer = NewKafkaProducer(ctx, wg, cfg.MQCfg.Kafka)
	}
}
