package mq

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chenbaoding2818/chainly/config"
	iface "github.com/chenbaoding2818/chainly/interface"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProducer struct {
	cfg  *config.RabbitMQConfig
	conn *amqp.Connection
	ch   *amqp.Channel
	lock sync.Mutex
	ctx  context.Context
}

func NewRabbitMQProducer(ctx context.Context, rabbitMQCfg *config.RabbitMQConfig) iface.IProducer {
	return &RabbitMQProducer{
		cfg: rabbitMQCfg,
		ctx: ctx,
	}
}

func (r *RabbitMQProducer) Connect() error {
	// conn, err := amqp.Dial(r.cfg.Urls[0])
	return nil
}

func (r *RabbitMQProducer) Reconnect() error {
	return nil
}

func (r *RabbitMQProducer) Send(ctx context.Context, msg []byte) error {
	return nil
}

func (r *RabbitMQProducer) SendAsync(ctx context.Context, msg []byte) error {
	publishing := r.NewMQMsg(msg)
	fmt.Printf("publishing: %v\n", publishing)
	return nil
}

func (r *RabbitMQProducer) CombineMessages(msgs [][]byte) []byte {
	var buf bytes.Buffer
	for _, msg := range msgs {
		buf.Write(msg)
	}
	return buf.Bytes()
}

func (r *RabbitMQProducer) NewMQMsg(body []byte) amqp.Publishing {
	// 空数据检测
	if len(body) == 0 {
		return amqp.Publishing{}
	}
	return amqp.Publishing{
		ContentType:  "text/plain",
		Body:         body,
		DeliveryMode: amqp.Persistent,
		MessageId:    fmt.Sprintf("%d", time.Now().UnixMilli()),
	}
}

func (r *RabbitMQProducer) Close() {
}
