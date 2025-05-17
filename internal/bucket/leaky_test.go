package bucket

import (
	"testing"
	"time"
)

type action string

const (
	actionAdd   action = "add"
	actionWait  action = "wait"
	actionReset action = "reset"
)

func TestLeakyBucket(t *testing.T) {
	tests := []struct {
		name       string
		capacity   int
		leakRate   float64
		operations []struct {
			action action
			wait   time.Duration // wait time for operation "wait"
			expect bool          // expected result for action="add"
		}
	}{
		{
			name:     "capacity 3, 1 leak per second",
			capacity: 3,
			leakRate: 1.0,
			operations: []struct {
				action action
				wait   time.Duration
				expect bool
			}{
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, false},                        // bucket is full
				{actionWait, 1100 * time.Millisecond, false}, // wait for 1.1 seconds for one token leak
				{actionAdd, 0, true},                         // OK
				{actionWait, 2100 * time.Millisecond, false}, // leak 2 tokens
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, false},                        // bucket is full
			},
		},
		{
			name:     "capacity 2, 0.5 leak per second",
			capacity: 2,
			leakRate: 0.5,
			operations: []struct {
				action action
				wait   time.Duration
				expect bool
			}{
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, false},                        // bucket is full
				{actionWait, 2100 * time.Millisecond, false}, // wait for 2.1 seconds for one token leak (0.5 * 2 = 1)
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, false},                        // bucket is full
			},
		},
		{
			name:     "reset after full",
			capacity: 2,
			leakRate: 1.0,
			operations: []struct {
				action action
				wait   time.Duration
				expect bool
			}{
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, false},                        // bucket is full
				{actionWait, 1100 * time.Millisecond, false}, // wait for 1.1 seconds for one token leak
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, false},                        // bucket is full
				{actionReset, 0, false},                      // reset bucket
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, false},                        // bucket is full
			},
		},
		{
			name:     "reset with leaking",
			capacity: 3,
			leakRate: 1.0,
			operations: []struct {
				action action
				wait   time.Duration
				expect bool
			}{
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, true},                         // OK
				{actionWait, 1100 * time.Millisecond, false}, // wait for leak
				{actionAdd, 0, true},                         // OK (after leak, should have 2 tokens)
				{actionReset, 0, false},                      // reset bucket
				{actionAdd, 0, true},                         // OK after reset (0 tokens)
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, false},                        // bucket is full
				{actionWait, 2100 * time.Millisecond, false}, // wait for leak
				{actionAdd, 0, true},                         // should be OK (2 tokens after leak)
			},
		},
		{
			name:     "reset with zero capacity",
			capacity: 0,
			leakRate: 1.0,
			operations: []struct {
				action action
				wait   time.Duration
				expect bool
			}{
				{actionAdd, 0, false},   // bucket is full (capacity 0)
				{actionReset, 0, false}, // reset bucket
				{actionAdd, 0, false},   // still rejected (capacity still 0)
			},
		},
		{
			name:     "multiple resets",
			capacity: 2,
			leakRate: 1.0,
			operations: []struct {
				action action
				wait   time.Duration
				expect bool
			}{
				{actionAdd, 0, true},    // OK
				{actionAdd, 0, true},    // OK
				{actionAdd, 0, false},   // bucket is full
				{actionReset, 0, false}, // reset bucket
				{actionAdd, 0, true},    // OK after reset
				{actionReset, 0, false}, // reset again
				{actionAdd, 0, true},    // OK
				{actionAdd, 0, true},    // OK
				{actionAdd, 0, false},   // bucket is full
			},
		},
		{
			name:     "reset after leak",
			capacity: 3,
			leakRate: 0.5,
			operations: []struct {
				action action
				wait   time.Duration
				expect bool
			}{
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, true},                         // OK
				{actionAdd, 0, true},                         // OK
				{actionWait, 1100 * time.Millisecond, false}, // wait for 1.1 seconds (leak 0.5 tokens)
				{actionReset, 0, false},                      // reset bucket
				// After reset, all tokens are gone regardless of previous leaks
				{actionAdd, 0, true},  // OK
				{actionAdd, 0, true},  // OK
				{actionAdd, 0, true},  // OK
				{actionAdd, 0, false}, // bucket is full
			},
		},
		{
			name:     "leakRate = 0 (no leaking)",
			capacity: 2,
			leakRate: 0,
			operations: []struct {
				action action
				wait   time.Duration
				expect bool
			}{
				{actionAdd, 0, true},             // OK
				{actionAdd, 0, true},             // OK
				{actionAdd, 0, false},            // bucket is full
				{actionWait, time.Second, false}, // wait for 1 second (no leaking)
				{actionAdd, 0, false},            // query still rejected
			},
		},
		{
			name:     "zero capacity",
			capacity: 0,
			leakRate: 1.0,
			operations: []struct {
				action action
				wait   time.Duration
				expect bool
			}{
				{actionAdd, 0, false},            // bucket is full
				{actionAdd, 0, false},            // bucket is full
				{actionWait, time.Second, false}, // wait for 1 second (capacity is 0)
				{actionAdd, 0, false},            // query still rejected
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			bucket := NewLeakyBucket(test.capacity, test.leakRate)

			for i, op := range test.operations {
				switch op.action {
				case actionWait:
					time.Sleep(op.wait)
				case actionReset:
					bucket.Reset()
				case actionAdd:
					result := bucket.Add()
					if result != op.expect {
						t.Errorf("operation %d: expected Add() to return %v, got %v",
							i, op.expect, result)
					}
				}
			}
		})
	}
}
