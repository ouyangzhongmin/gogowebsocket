package gogowebsocket

import "time"

type Option func(o *options)
type options struct {
	maxMessageSize  int64
	writeWaitTime   time.Duration
	pongWaitTime    time.Duration
	pingPeriod      time.Duration
	writeBufferSize int
	readBufferSize  int
}

func MaxMessageSize(value int64) Option {
	return func(o *options) {
		o.maxMessageSize = value
	}
}

func WriteWaitTime(value time.Duration) Option {
	return func(o *options) {
		o.writeWaitTime = value
	}
}

func PongWaitTime(value time.Duration) Option {
	return func(o *options) {
		o.pongWaitTime = value
	}
}

func PingPeriod(value time.Duration) Option {
	return func(o *options) {
		o.pingPeriod = value
	}
}

func WriteBufferSize(value int) Option {
	return func(o *options) {
		o.writeBufferSize = value
	}
}

func ReadBufferSize(value int) Option {
	return func(o *options) {
		o.readBufferSize = value
	}
}
