// @Author Chen_dark
// @Date 2024/10/31 11:33:00
// @Desc
package examples

import (
	"context"
	"github.com/YouChenJun/redmq"
	"github.com/YouChenJun/redmq/redis"
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
