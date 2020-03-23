package irc

import (
	"io"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	baseWriter    io.Writer
	limitedWriter io.Writer
	capacity      float64
	maxCapacity   float64
	refreshRate   float64
	refreshAmount float64
	refreshUnit   float64
	initialRampUp bool
	initialAmount float64
	received001   bool
}

func (r *RateLimiter) Init(profile string) {
	switch profile {
	case "unlimited":
		r.limitedWriter = r.baseWriter
	case "restrictive":
		r.refreshRate = 1
		r.refreshAmount = 0.5
		r.refreshUnit = 128
		r.initialRampUp = true
		r.initialAmount = 0.3
		r.maxCapacity = 4
	}
	r.capacity = 2
	go r.refilTimer()
}

func (r *RateLimiter) refilTimer() {
	ticker := time.NewTicker(time.Duration(r.refreshRate) * time.Second)
	sigWait := make(chan os.Signal, 1)
	signal.Notify(sigWait, os.Interrupt)
	signal.Notify(sigWait, syscall.SIGTERM)
	for {
		select {
		case <-sigWait:
			return
		case <-ticker.C:
			amount := float64(0)
			if r.initialRampUp {
				amount = math.Min(r.maxCapacity, r.capacity+r.initialAmount)
				if amount == r.maxCapacity {
					r.initialRampUp = false
				}
			} else {
				amount = math.Max(0, math.Min(r.maxCapacity, r.capacity+r.refreshAmount))
			}
			r.capacity = amount
		}
	}
}

func (r *RateLimiter) Write(p []byte) (n int, err error) {
	needed := math.Min(math.Ceil(float64(len(p))/r.refreshUnit), r.maxCapacity)
	if r.received001 {
		for {
			if needed > r.capacity {
				time.Sleep(time.Duration(250) * time.Millisecond)
			} else {
				break
			}
		}
		r.capacity -= needed
	}
	return r.baseWriter.Write(p)
}

func (r *RateLimiter) handle001(*Connection, *Message) {
	r.received001 = true
}
