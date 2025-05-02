package service

import (
	"cloud/Balancer/internal/models"
	"context"
	"net/http"
	"strings"
	"sync"
	"time"
)

type RateLimiter struct {
	buckets       map[models.ClientIdentifier]*TokenBucket
	defaultCap    int
	defaultRate   float64
	Storage       models.Storage
	mu            sync.RWMutex
	cleanupTicker *time.Ticker
}

func NewRateLimiterWithStorage(defaultCapacity int, defaultRate float64, storage models.Storage) (*RateLimiter, error) {
	rl := &RateLimiter{
		buckets:     make(map[models.ClientIdentifier]*TokenBucket),
		defaultCap:  defaultCapacity,
		defaultRate: defaultRate,
		Storage:     storage,
	}

	if err := rl.loadLimitsFromStorage(); err != nil {
		return nil, err
	}

	rl.startCleanupRoutine()
	return rl, nil
}

func (rl *RateLimiter) GetClientID(r *http.Request) models.ClientIdentifier {
	if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
		return models.ClientIdentifier(models.APIKeyPrefix + apiKey)
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]
	return models.ClientIdentifier(models.IPPrefix + ip)
}

func (rl *RateLimiter) loadLimitsFromStorage() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	limits, err := rl.Storage.GetAllClientLimits(ctx)
	if err != nil {
		return err
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	for _, limit := range limits {
		rl.buckets[models.ClientIdentifier(limit.ClientID)] = NewTokenBucket(
			limit.Capacity,
			limit.Rate,
		)
	}

	return nil
}

func (rl *RateLimiter) Allow(clientID models.ClientIdentifier) bool {
	rl.mu.RLock()
	bucket, exists := rl.buckets[clientID]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		bucket = NewTokenBucket(rl.defaultCap, rl.defaultRate)
		rl.buckets[clientID] = bucket
		rl.mu.Unlock()
	}

	return bucket.Allow()
}

func (rl *RateLimiter) SetLimit(clientID models.ClientIdentifier, capacity int, rate float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rl.Storage.SetClientLimit(ctx, string(clientID), capacity, rate); err != nil {
		return err
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.buckets[clientID] = NewTokenBucket(capacity, rate)
	return nil
}

func (rl *RateLimiter) startCleanupRoutine() {
	rl.cleanupTicker = time.NewTicker(1 * time.Hour)
	go func() {
		for range rl.cleanupTicker.C {
			rl.cleanupOldBuckets()
		}
	}()
}

func (rl *RateLimiter) cleanupOldBuckets() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for id, bucket := range rl.buckets {
		if time.Since(bucket.lastRefill) > 24*time.Hour {
			delete(rl.buckets, id)
		}
	}
}

func (rl *RateLimiter) Stop() {
	if rl.cleanupTicker != nil {
		rl.cleanupTicker.Stop()
	}
}
