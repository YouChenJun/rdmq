// @Author Chen_dark
// @Date 2024/10/31 11:26:00
// @Desc
package examples

import (
	"context"
	"github.com/YouChenJun/rdmq"
	"github.com/YouChenJun/rdmq/redis"
	"testing"
	"time"
)

// 自定义实现的死信队列
type DemoDeadLetterMailbox struct {
	do func(msg *redis.MsgEntity)
}

func NewDemoDeadLetterMailbox(do func(msg *redis.MsgEntity)) *DemoDeadLetterMailbox {
	return &DemoDeadLetterMailbox{
		do: do,
	}
}

// 死信队列接收消息的处理方法
func (d *DemoDeadLetterMailbox) Deliver(ctx context.Context, msg *redis.MsgEntity) error {
	d.do(msg)
	return nil
}

func Test_Consumer(t *testing.T) {
	config := redis.Config{
		Network:  "tcp",
		Address:  ":6379",
		Password: "",
	}
	ConsumerConifg := &redis.ConsumerConfig{
		RedisConfig: config,
		GroupID:     "",
		ConsumerID:  "22",
		Topic:       "",
	}

	client := redis.NewClient(config)

	// 接收到消息后的处理函数
	callbackFunc := func(ctx context.Context, msg *redis.MsgEntity) (error, bool) {
		// 此处可以设置业务逻辑 是否回复ACK
		if msg.Key == "111" {
			t.Logf("要消费receive msg, msg id: %s, msg key: %s, msg val: %s", msg.MsgID, msg.Key, msg.Val)
			return nil, true
		}
		t.Logf("消息%v不消费", msg.MsgID)
		return nil, false
	}

	// 自定义实现的死信队列
	demoDeadLetterMailbox := NewDemoDeadLetterMailbox(func(msg *redis.MsgEntity) {
		t.Logf("receive dead letter, msg id: %s, msg key: %s, msg val: %s", msg.MsgID, msg.Key, msg.Val)
	})

	// 构造并启动消费者
	consumer, err := redmq.NewConsumer(client, ConsumerConifg, callbackFunc,
		// 每条消息最多重试 2 次
		redmq.WithMaxRetryLimit(2),
		// 每轮接收消息的超时时间为 2 s
		redmq.WithReceiveTimeout(2*time.Second),
		// 注入自定义实现的死信队列
		redmq.WithDeadLetterMailbox(demoDeadLetterMailbox),
		// 轮询等待时间 默认为1S
		redmq.WaitTime(3*time.Second))
	if err != nil {
		t.Error(err)
		return
	}
	defer consumer.Stop()

	// 十秒后退出单测程序
	<-time.After(1000 * time.Second)
}

func Test_Consumer2(t *testing.T) {
	config := redis.Config{
		Network:  "tcp",
		Address:  ":6379",
		Password: "",
	}
	ConsumerConifg := &redis.ConsumerConfig{
		RedisConfig: config,
		GroupID:     "",
		ConsumerID:  "22",
		Topic:       "",
	}

	client := redis.NewClient(config)

	// 接收到消息后的处理函数
	callbackFunc := func(ctx context.Context, msg *redis.MsgEntity) (error, bool) {
		// 此处可以设置业务逻辑 是否回复ACK
		if msg.Key == "112" {
			t.Logf("要消费receive msg, msg id: %s, msg key: %s, msg val: %s", msg.MsgID, msg.Key, msg.Val)
			return nil, true
		}
		t.Logf("消息%v不消费", msg.MsgID)
		return nil, false
	}

	// 自定义实现的死信队列
	demoDeadLetterMailbox := NewDemoDeadLetterMailbox(func(msg *redis.MsgEntity) {
		t.Logf("receive dead letter, msg id: %s, msg key: %s, msg val: %s", msg.MsgID, msg.Key, msg.Val)
	})

	// 构造并启动消费者
	consumer, err := redmq.NewConsumer(client, ConsumerConifg, callbackFunc,
		// 每条消息最多重试 2 次
		redmq.WithMaxRetryLimit(2),
		// 每轮接收消息的超时时间为 2 s
		redmq.WithReceiveTimeout(2*time.Second),
		// 注入自定义实现的死信队列
		redmq.WithDeadLetterMailbox(demoDeadLetterMailbox))
	if err != nil {
		t.Error(err)
		return
	}
	defer consumer.Stop()

	// 十秒后退出单测程序
	<-time.After(1000 * time.Second)
}
