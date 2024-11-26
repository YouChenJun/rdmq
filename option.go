package redmq

import "time"

type ProducerOptions struct {
	msgQueueLen int // 消息队列长度
}

type ProducerOption func(opts *ProducerOptions)

// WithMsgQueueLen 设置消息队列长度的选项函数
func WithMsgQueueLen(len int) ProducerOption {
	return func(opts *ProducerOptions) {
		opts.msgQueueLen = len
	}
}

// repairProducer 修复 ProducerOptions 中的 msgQueueLen，确保其大于 0
func repairProducer(opts *ProducerOptions) {
	if opts.msgQueueLen <= 0 {
		opts.msgQueueLen = 500
	}
}

type ConsumerOptions struct {
	receiveTimeout           time.Duration     // 每轮接收消息的超时时长
	maxRetryLimit            int               // 处理消息的最大重试次数，超过此次数时，消息会被投递到死信队列
	deadLetterMailbox        DeadLetterMailbox // 死信队列，可以由使用方自定义实现
	deadLetterDeliverTimeout time.Duration     // 投递死信流程超时阈值
	handleMsgsTimeout        time.Duration     // 处理消息流程超时阈值
	waitTime                 time.Duration     //消息等待时间-避免一直轮询
}

type ConsumerOption func(opts *ConsumerOptions)

// WaitTime 设置消息等待时间的选项函数
func WaitTime(waittime time.Duration) ConsumerOption {
	return func(opts *ConsumerOptions) {
		opts.waitTime = waittime
	}
}

// WithReceiveTimeout 设置接收消息超时时间的选项函数
func WithReceiveTimeout(timeout time.Duration) ConsumerOption {
	return func(opts *ConsumerOptions) {
		opts.receiveTimeout = timeout
	}
}

// WithMaxRetryLimit 设置处理消息最大重试次数的选项函数
func WithMaxRetryLimit(maxRetryLimit int) ConsumerOption {
	return func(opts *ConsumerOptions) {
		opts.maxRetryLimit = maxRetryLimit
	}
}

// WithDeadLetterMailbox 设置死信队列的选项函数
func WithDeadLetterMailbox(mailbox DeadLetterMailbox) ConsumerOption {
	return func(opts *ConsumerOptions) {
		opts.deadLetterMailbox = mailbox
	}
}

// WithDeadLetterDeliverTimeout 设置投递死信流程超时阈值的选项函数
func WithDeadLetterDeliverTimeout(timeout time.Duration) ConsumerOption {
	return func(opts *ConsumerOptions) {
		opts.deadLetterDeliverTimeout = timeout
	}
}

// WithHandleMsgsTimeout 设置处理消息流程超时阈值的选项函数
func WithHandleMsgsTimeout(timeout time.Duration) ConsumerOption {
	return func(opts *ConsumerOptions) {
		opts.handleMsgsTimeout = timeout
	}
}

// repairConsumer 修复 ConsumerOptions 中的各个字段，确保每个字段有数值
func repairConsumer(opts *ConsumerOptions) {
	if opts.receiveTimeout < 0 {
		opts.receiveTimeout = 2 * time.Second
	}

	if opts.maxRetryLimit < 0 {
		opts.maxRetryLimit = 3
	}

	if opts.deadLetterMailbox == nil {
		opts.deadLetterMailbox = NewDeadLetterLogger()
	}

	if opts.deadLetterDeliverTimeout <= 0 {
		opts.deadLetterDeliverTimeout = time.Second
	}

	if opts.handleMsgsTimeout <= 0 {
		opts.handleMsgsTimeout = time.Second
	}

	if opts.waitTime <= 0 {
		opts.waitTime = 1
	}
}
