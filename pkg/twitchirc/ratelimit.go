package twitchirc

import (
	"errors"
	"sync"
	"time"
)

// Ratelimiter holds data for ratelimiting purposes
// The basic idea is having buckets for endpoints that are rate limited.
// These buckets can be locked so other requests won't interfere with ratelimit counting.
type Ratelimiter struct {
	Buckets map[string]*bucket
}

type bucket struct {
	Reset         time.Time
	Limit         uint8
	Current       uint8
	ResetInterval time.Duration
	Mutex         *sync.Mutex
}

// RateLimitInput acts as data wrapper for Call Limit and Reset Interval
type RateLimitInput struct {
	Limit         uint8
	ResetInterval time.Duration
}

// NewRatelimiter creates a new rate limiter with buckets corresponding
// to given input map
func NewRatelimiter(input map[string]RateLimitInput) *Ratelimiter {
	r := &Ratelimiter{
		Buckets: make(map[string]*bucket),
	}
	for k, v := range input {
		r.Buckets[k] = &bucket{
			Reset:         time.Unix(0, 0),
			Limit:         v.Limit,
			Current:       0,
			ResetInterval: v.ResetInterval,
			Mutex:         &sync.Mutex{},
		}
	}
	return r
}

// CaptureEndpoint locks the given endpoints' bucket for other requests until released
func (r *Ratelimiter) CaptureEndpoint(endpoint string) (*bucket, error) {
	b, ok := r.Buckets[endpoint]
	b.Mutex.Lock()
	if !ok {
		return nil, errors.New("Unknown endpoint")
	}
	return b, nil
}

// Increment performs the logic of resetting the counter and refreshing the
// reset time as well as upping the counter
func (b *bucket) Increment() error {
	// If bucket has never been used or not used in last interval renew the the ratelimit reset point
	if b.Reset == time.Unix(0, 0) || b.Reset.Before(time.Now()) && b.Current == 0 {
		b.Reset = time.Now().Add(b.ResetInterval)
	}
	if b.Current == b.Limit && b.Reset.After(time.Now()) {
		return errors.New("Ratelimit exceeded")
	}
	if b.Reset.Before(time.Now()) {
		b.Current = 0
		b.Reset.Add(b.ResetInterval)
	}
	b.Current++
	return nil
}

// Release unlocks the bucket -  to be used after incrementing
func (b *bucket) Release() {
	b.Mutex.Unlock()
}
