package irc

import (
	"context"
	"golang.org/x/time/rate"
	"io"
)

func (irc *Connection) NewRateLimiter(w io.Writer, floodProfile string) *RateLimiter {
	rl := RateLimiter{
		baseWriter: w,
	}
	rl.Init(floodProfile)
	irc.AddInboundHandler("001", rl.handle001)
	return &rl
}

type RateLimiter struct {
	baseWriter  io.Writer
	limiter     *rate.Limiter
	received001 bool
	byteToToken int
}

func (r *RateLimiter) Init(profile string) {
	switch profile {
	case "unlimited":
		r.limiter = rate.NewLimiter(rate.Inf, 0)
		r.byteToToken = 1
	case "restrictive":
		r.limiter = rate.NewLimiter(rate.Limit(0.4), 4)
		r.byteToToken = 128
	}
}

func (r *RateLimiter) Write(p []byte) (n int, err error) {
	needed := len(p) / r.byteToToken
	if r.received001 {
		_ = r.limiter.WaitN(context.Background(), needed)
	}
	return r.baseWriter.Write(p)
}

func (r *RateLimiter) handle001(*EventManager, *Connection, *Message) {
	r.received001 = true
}
