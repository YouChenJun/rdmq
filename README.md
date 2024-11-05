# rdmq-SDK

> 使用Redis实现mq的SDK

# 接入SOP

​	首先需要创建topic和consumer group 

- 创建topic

```
xadd xscan_topic * first_key first_val
```

- 创建 consumer group

```
XGROUP CREATE test_topic test_group 0-0
```

- 构造redis客户端实例

```go
	config := redis.Config{
		Network:  "tcp",
		Address:  "",
		Password: "",
	}
	
	client := redis.NewClient(config)
```

- 启动生产者producer

```go
	producer := redmq.NewProducer(client, redmq.WithMsgQueueLen(10))//最多保留10条消息

	ctx := context.Background()
	msgID, err := producer.SendMsg(ctx, producers.Topic, "test_kk", "test_vv")
	if err != nil {
		t.Error(err)
		return
	}
```

# 完整示例

## consumer

```go
// @Author Chen_dark
// @Date 2024/10/31 11:26:00
// @Desc
package examples

import (
	"context"
	"github.com/YouChenJun/redmq"
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
		Address:  "",
		Password: "",
	}
	ConsumerConifg := &redis.ConsumerConfig{
		RedisConfig: config,
		GroupID:     "",
		ConsumerID:  "",
		Topic:       "",
	}

	client := redis.NewClient(config)

	// 接收到消息后的处理函数
	callbackFunc := func(ctx context.Context, msg *redis.MsgEntity) error {
		t.Logf("receive msg, msg id: %s, msg key: %s, msg val: %s", msg.MsgID, msg.Key, msg.Val)
		return nil
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
	<-time.After(10 * time.Second)
}

```

## producer

```go
// @Author Chen_dark
// @Date 2024/10/31 11:33:00
// @Desc
package examples

import (
	"context"
	"github.com/YouChenJun/redmq"
	"github.com/YouChenJun/rdmq/redis"
	"testing"
)

func Test_Producer(t *testing.T) {
	config := redis.Config{
		Network:  "tcp",
		Address:  "192.168.8.189:6379",
		Password: "",
	}
	producers := &redis.ProducerConfig{
		RedisConfig: config,
		Topic:       "",
	}

	client := redis.NewClient(producers.RedisConfig)

	// 最多保留十条消息
	producer := redmq.NewProducer(client, redmq.WithMsgQueueLen(10))

	ctx := context.Background()
	msgID, err := producer.SendMsg(ctx, producers.Topic, "test_kk", "test_vv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(msgID)

}

```

