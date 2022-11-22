// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The code below is borrowed from Go's cmd/gofmt as of 1.18beta1.
// We tweaked it slightly to add the "broken URLs" result.

package main

import (
	"context"
	"go/scanner"
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
// and writes tasks' output to out and err.
func newSequencer(maxWeight int64, out, err io.Writer) *sequencer {
	sem := semaphore.NewWeighted(maxWeight)
	prev := make(chan *reporterState, 1)
	prev <- &reporterState{out: out, err: err}
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
//
// If f returns a non-nil error, that error will be reported after f's output
// (if any) and will cause a nonzero final exit code.
func (s *sequencer) Add(weight int64, f func(*reporter) error) {
	if weight < 0 || weight > s.maxWeight {
		weight = s.maxWeight
	}
	if err := s.sem.Acquire(context.TODO(), weight); err != nil {
		// Change the task from "execute f" to "report err".
		weight = 0
		f = func(*reporter) error { return err }
	}

	r := &reporter{prev: s.prev}
	next := make(chan *reporterState, 1)
	s.prev = next

	// Start f in parallel: it can run until it invokes a method on r, at which
	// point it will block until the previous task releases the output state.
	go func() {
		if err := f(r); err != nil {
			r.Report(err)
		}
		next <- r.getState() // Release the next task.
		s.sem.Release(weight)
	}()
}

// GetExitCode waits for all previously-added tasks to complete, then returns an
// exit code for the sequence suitable for passing to os.Exit.
func (s *sequencer) GetExitCode() int {
	c := make(chan int, 1)
	s.Add(0, func(r *reporter) error {
		c <- r.ExitCode()
		return nil
	})
	return <-c
}

func (s *sequencer) finalState() reporterState {
	c := make(chan reporterState, 1)
	s.Add(0, func(r *reporter) error {
		c <- *r.getState()
		return nil
	})
	return <-c
}

// A reporter reports output, warnings, and errors.
type reporter struct {
	prev  <-chan *reporterState
	state *reporterState
}

// reporterState carries the state of a reporter instance.
//
// Only one reporter at a time may have access to a reporterState.
type reporterState struct {
	out, err io.Writer
	exitCode int

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
//
// Any error is returned to the caller, and does not otherwise affect the
// reporter's exit code.
func (r *reporter) Write(p []byte) (int, error) {
	return r.getState().out.Write(p)
}

func (r *reporter) appendBroken(url, reason string) {
	state := r.getState()
	state.brokenURLs = append(state.brokenURLs, brokenURL{url, reason})
}

// Report emits a non-nil error to the reporter's error stream,
// changing its exit code to a nonzero value.
func (r *reporter) Report(err error) {
	if err == nil {
		panic("Report with nil error")
	}
	st := r.getState()
	scanner.PrintError(st.err, err)
	st.exitCode = 2
}

func (r *reporter) ExitCode() int {
	return r.getState().exitCode
}
