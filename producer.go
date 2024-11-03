package redmq

import (
	"context"
	"github.com/YouChenJun/redmq/redis"
)

type Producer struct {
	client *redis.Client
	opts   *ProducerOptions
}

func NewProducer(client *redis.Client, opts ...ProducerOption) *Producer {
	p := Producer{
		client: client,
		opts:   &ProducerOptions{},
	}

	for _, opt := range opts {
		opt(p.opts)
	}

	repairProducer(p.opts)
	return &p
}

// 生产一条消息
func (p *Producer) SendMsg(ctx context.Context, topic, key, val string) (string, error) {
	if topic == "" {
		topic = p.client.Topic
	}
	return p.client.XADD(ctx, topic, p.opts.msgQueueLen, key, val)
}

// 生产一条消息
func (p *Producer) SendJsonMsg(ctx context.Context, topic, key, val string) (string, error) {
	//获取配置文件中的topic
	if topic == "" {
		topic = p.client.Topic
	}
	topic = p.client.Topic
	return p.client.XADDJson(ctx, topic, p.opts.msgQueueLen, key, val)
}
