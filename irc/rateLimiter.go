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
	return &rl
}

type RateLimiter struct {
	baseWriter    io.Writer
	limitedWriter io.Writer
	capacity      float64
	maxCapacity   float64
	refreshRate   float64
	refreshAmount float64
}

func (r *RateLimiter) Init(profile string) {
	switch profile {
	case "unlimited":
		r.limitedWriter = r.baseWriter
	case "restrictive":
		r.refreshRate = 1
		r.refreshAmount = 0.5
		r.maxCapacity = 5
	}
	r.capacity = r.maxCapacity
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
			r.capacity = math.Max(0, math.Min(r.maxCapacity, r.capacity+r.refreshAmount))
		}
	}
}

func (r *RateLimiter) Write(p []byte) (n int, err error) {
	needed := math.Min(math.Ceil(float64(len(p))/128), r.maxCapacity)
	for {
		if needed > r.capacity {
			time.Sleep(time.Duration(250) * time.Millisecond)
			break
		} else {
			break
		}
	}
	r.capacity -= needed
	return r.baseWriter.Write(p)
}
