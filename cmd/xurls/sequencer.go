// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The code below is borrowed from Go's cmd/gofmt as of 1.18beta1.
// We tweaked it to add the "broken URLs" result and removed what we don't use.

package main

import (
	"context"
	"io"

	"golang.org/x/sync/semaphore"
)

// A sequencer performs concurrent tasks that may write output, but emits that
// output in a deterministic order.
type sequencer struct {
	maxWeight int64
	sem       *semaphore.Weighted   // weighted by input bytes (an approximate proxy for memory overhead)
	prev      <-chan *reporterState // 1-buffered
}

// newSequencer returns a sequencer that allows concurrent tasks up to maxWeight
// and writes tasks' output to out.
func newSequencer(maxWeight int64, out io.Writer) *sequencer {
	sem := semaphore.NewWeighted(maxWeight)
	prev := make(chan *reporterState, 1)
	prev <- &reporterState{out: out}
	return &sequencer{
		maxWeight: maxWeight,
		sem:       sem,
		prev:      prev,
	}
}

// Add blocks until the sequencer has enough weight to spare, then adds f as a
// task to be executed concurrently.
//
// If the weight is either negative or larger than the sequencer's maximum
// weight, Add blocks until all other tasks have completed, then the task
// executes exclusively (blocking all other calls to Add until it completes).
//
// f may run concurrently in a goroutine, but its output to the passed-in
// reporter will be sequential relative to the other tasks in the sequencer.
//
// If f invokes a method on the reporter, execution of that method may block
// until the previous task has finished. (To maximize concurrency, f should
// avoid invoking the reporter until it has finished any parallelizable work.)
func (s *sequencer) Add(weight int64, f func(*reporter)) {
	if weight < 0 || weight > s.maxWeight {
		weight = s.maxWeight
	}
	s.sem.Acquire(context.Background(), weight) // cannot fail without a deadline

	r := &reporter{prev: s.prev}
	next := make(chan *reporterState, 1)
	s.prev = next

	// Start f in parallel: it can run until it invokes a method on r, at which
	// point it will block until the previous task releases the output state.
	go func() {
		f(r)
		next <- r.getState() // Release the next task.
		s.sem.Release(weight)
	}()
}

func (s *sequencer) finalState() reporterState {
	c := make(chan reporterState, 1)
	s.Add(0, func(r *reporter) {
		c <- *r.getState()
	})
	return <-c
}

// A reporter reports output and broken URLs.
type reporter struct {
	prev  <-chan *reporterState
	state *reporterState
}

// reporterState carries the state of a reporter instance.
//
// Only one reporter at a time may have access to a reporterState.
type reporterState struct {
	out        io.Writer
	brokenURLs []brokenURL
}

type brokenURL struct {
	url    string
	reason string
}

// getState blocks until any prior reporters are finished with the reporter
// state, then returns the state for manipulation.
func (r *reporter) getState() *reporterState {
	if r.state == nil {
		r.state = <-r.prev
	}
	return r.state
}

// Write emits a slice to the reporter's output stream.
func (r *reporter) Write(p []byte) (int, error) {
	return r.getState().out.Write(p)
}

func (r *reporter) appendBroken(url, reason string) {
	state := r.getState()
	state.brokenURLs = append(state.brokenURLs, brokenURL{url, reason})
}
