package middleware

import (
	"sync"
	"time"
)

type TokenBlacklist struct {
	mu     sync.RWMutex
	tokens map[string]time.Time
}

func NewTokenBlacklist() *TokenBlacklist {
	b := &TokenBlacklist{
		tokens: make(map[string]time.Time),
	}
	go b.cleanup()
	return b
}

func (b *TokenBlacklist) Add(token string, exp time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tokens[token] = exp
}

func (b *TokenBlacklist) IsBlacklisted(token string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	exp, ok := b.tokens[token]
	if !ok {
		return false
	}

	return time.Now().Before(exp)
}

func (b *TokenBlacklist) cleanup() {
	for {
		time.Sleep(1 * time.Minute)
		now := time.Now()

		b.mu.Lock()
		for t, exp := range b.tokens {
			if now.After(exp) {
				delete(b.tokens, t)
			}
		}
		b.mu.Unlock()
	}
}