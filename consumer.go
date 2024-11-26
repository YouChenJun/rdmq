package redmq

import (
	"context"
	"errors"
	"github.com/YouChenJun/rdmq/log"
	"github.com/YouChenJun/rdmq/redis"
	"time"
)

// 接收到消息后执行的回调函数-此处的bool为判断，如果true则回复ACK
type MsgCallback func(ctx context.Context, msg *redis.MsgEntity) (error, bool)

// 消费者
type Consumer struct {
	ctx          context.Context // consumer 生命周期管理
	stop         context.CancelFunc
	callbackFunc MsgCallback             // 接收到 msg 时执行的回调函数，由使用方定义
	client       *redis.Client           // redis 客户端，基于 redis 实现 message queue
	topic        string                  // 消费的 topic
	groupID      string                  // 所属的消费者组
	consumerID   string                  // 当前节点的消费者 id
	failureCnts  map[redis.MsgEntity]int // 各消息累计失败次数
	waitTime     time.Duration           //消息等待轮询时间
	opts         *ConsumerOptions        // 一些用户自定义的配置
}

func NewConsumer(client *redis.Client, consumerConfig *redis.ConsumerConfig, callbackFunc MsgCallback, opts ...ConsumerOption) (*Consumer, error) {
	ctx, stop := context.WithCancel(context.Background())
	c := Consumer{
		client:       client,
		ctx:          ctx,
		stop:         stop,
		callbackFunc: callbackFunc,
		topic:        consumerConfig.Topic,
		groupID:      consumerConfig.GroupID,
		consumerID:   consumerConfig.ConsumerID,

		opts: &ConsumerOptions{},

		failureCnts: make(map[redis.MsgEntity]int),
	}

	if err := c.checkParam(); err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(c.opts)
	}

	repairConsumer(c.opts)

	go c.run()
	return &c, nil
}

func (c *Consumer) checkParam() error {
	if c.callbackFunc == nil {
		return errors.New("callback function can't be empty")
	}

	if c.client == nil {
		return errors.New("redis client can't be empty")
	}

	if c.topic == "" || c.consumerID == "" || c.groupID == "" {
		return errors.New("topic | group_id | consumer_id can't be empty")
	}

	return nil
}

// 停止 consumer
func (c *Consumer) Stop() {
	c.stop()
}

// 运行消费者
func (c *Consumer) run() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		// 新消息接收处理
		msgs, err := c.receive()
		if err != nil {
			log.ErrorContextf(c.ctx, "receive msg failed, err: %v", err)
			continue
		}

		tctx, _ := context.WithTimeout(c.ctx, c.opts.handleMsgsTimeout)
		c.handlerMsgs(tctx, msgs)

		// 死信队列投递
		tctx, _ = context.WithTimeout(c.ctx, c.opts.deadLetterDeliverTimeout)
		c.deliverDeadLetter(tctx)

		// pending 消息接收处理
		pendingMsgs, err := c.receivePending()
		if err != nil {
			log.ErrorContextf(c.ctx, "pending msg received failed, err: %v", err)
			time.Sleep(c.waitTime * time.Second)
			continue
		}

		tctx, _ = context.WithTimeout(c.ctx, c.opts.handleMsgsTimeout)
		time.Sleep(c.waitTime * time.Second)
		c.handlerMsgs(tctx, pendingMsgs)
	}
}

func (c *Consumer) receive() ([]*redis.MsgEntity, error) {
	msgs, err := c.client.XReadGroup(c.ctx, c.groupID, c.consumerID, c.topic, int(c.opts.receiveTimeout.Milliseconds()))
	if err != nil && !errors.Is(err, redis.ErrNoMsg) {
		return nil, err
	}

	return msgs, nil
}

func (c *Consumer) receivePending() ([]*redis.MsgEntity, error) {
	pendingMsgs, err := c.client.XReadGroupPending(c.ctx, c.groupID, c.consumerID, c.topic)
	if err != nil && !errors.Is(err, redis.ErrNoMsg) {
		return nil, err
	}

	return pendingMsgs, nil
}

func (c *Consumer) handlerMsgs(ctx context.Context, msgs []*redis.MsgEntity) {
	for _, msg := range msgs {
		err, ok := c.callbackFunc(ctx, msg)
		if err != nil {
			// 失败计数器累加
			c.failureCnts[*msg]++
			continue
		}
		//当回调函数返回的ok为true时，才发送ACK
		if ok {
			// callback 执行成功，进行 ack
			if err := c.client.XACK(ctx, c.topic, c.groupID, msg.MsgID); err != nil {
				log.ErrorContextf(ctx, "msg ack failed, msg id: %s, err: %v", msg.MsgID, err)
				continue
			}
			// 如果ACK成功，从失败计数器中删除该消息
			delete(c.failureCnts, *msg)
		} else {
			// 如果回调函数返回的ok为false，则不发送ACK
			log.Infof("callback function returned false, msg id: %s, not sending ACK", msg.MsgID)
			//log.WarnContextf(ctx, "callback function returned false, msg id: %s, not sending ACK", msg.MsgID)
		}
	}
}

func (c *Consumer) deliverDeadLetter(ctx context.Context) {
	// 对于失败达到指定次数的消息，投递到死信中，然后执行 ack
	for msg, failureCnt := range c.failureCnts {
		if failureCnt < c.opts.maxRetryLimit {
			continue
		}

		// 投递死信队列
		if err := c.opts.deadLetterMailbox.Deliver(ctx, &msg); err != nil {
			log.ErrorContextf(c.ctx, "dead letter deliver failed, msg id: %s, err: %v", msg.MsgID, err)
		}

		// 执行 ack 响应
		if err := c.client.XACK(ctx, c.topic, c.groupID, msg.MsgID); err != nil {
			log.ErrorContextf(c.ctx, "msg ack failed, msg id: %s, err: %v", msg.MsgID, err)
			continue
		}

		// 对于 ack 成功的消息，将其从 failure map 中删除
		delete(c.failureCnts, msg)
	}
}
