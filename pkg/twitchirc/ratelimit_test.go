package twitchirc

import (
	"testing"
	"time"
)

func TestRatelimiter(t *testing.T) {
	input := make(map[string]RateLimitInput)
	input["chat"] = *&RateLimitInput{
		Limit:         2,
		ResetInterval: 2 * time.Second,
	}
	limiter := NewRatelimiter(input)
	b, err := limiter.CaptureEndpoint("chat")
	if err != nil {
		t.Error(err)
	}
	b.Increment()
	b.Increment()
	err = b.Increment()
	if err == nil {
		t.Error("Expected error: \"Ratelimit exceeded\" but received none")
	}
	time.Sleep(2 * time.Second)
	err = b.Increment()
	if err != nil {
		t.Error("Expected no error but got one when making an increment after reset")
	}
}
