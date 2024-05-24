package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter(5, time.Second)
	identifier := "user1"

	for i := 0; i < 5; i++ {
		if !rl.Allow(identifier) {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}

	if rl.Allow(identifier) {
		t.Error("Expected request to be denied")
	}

	// Test requests after time duration
	time.Sleep(time.Second)

	for i := 0; i < 5; i++ {
		if !rl.Allow(identifier) {
			t.Errorf("Expected request %d to be allowed after duration", i+1)
		}
	}

	if rl.Allow(identifier) {
		t.Error("Expected request to be denied after exceeding limit again")
	}
}

func TestMultipleIdentifiers(t *testing.T) {
	rl := NewRateLimiter(3, time.Second)
	id1 := "user1"
	id2 := "user2"

	for i := 0; i < 3; i++ {
		if !rl.Allow(id1) {
			t.Errorf("Expected request %d for %s to be allowed", i+1, id1)
		}
		if !rl.Allow(id2) {
			t.Errorf("Expected request %d for %s to be allowed", i+1, id2)
		}
	}

	if rl.Allow(id1) {
		t.Error("Expected request to be denied for first identifier")
	}

	if rl.Allow(id2) {
		t.Error("Expected request to be allowed for second identifier")
	}
}

func TestConcarentAccess(t *testing.T) {
	rate := 5
	rl := NewRateLimiter(rate, time.Second)
	identifier := "user1"
	var wg sync.WaitGroup
	wg.Add(rate)

	makeRequests := func() {
		defer wg.Done()
		rl.Allow(identifier)

	}

	for i := 0; i < rate; i++ {
		go makeRequests()
	}

	wg.Wait()

	if rl.Allow(identifier) {
		t.Errorf("Expected request to be restricted after concurent access")
	}

	// Ensure that after a second, we can make new requests
	time.Sleep(time.Second)

	if !rl.Allow(identifier) {
		t.Error("Expected request to be allowed after exceeding limit again")
	}
}
