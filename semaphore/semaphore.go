// Package semaphore provide a semaphore implementation
package semaphore

import (
	"context"

	sema "golang.org/x/sync/semaphore"
)

// Semaphore represent a semaphore object
type Semaphore struct {
	s   *sema.Weighted
	ctx context.Context
}

// New create a semaphore object
func New(n int) *Semaphore {
	return &Semaphore{
		s:   sema.NewWeighted(int64(n)),
		ctx: context.TODO(),
	}
}

// Acquire reference increased
func (s *Semaphore) Acquire() {
	s.s.Acquire(s.ctx, 1)
}

// Release reference decreased
func (s *Semaphore) Release() {
	s.s.Release(1)
}

// AcquireN reference increased
func (s *Semaphore) AcquireN(n int64) {
	s.s.Acquire(s.ctx, n)
}

// ReleaseN reference decreased
func (s *Semaphore) ReleaseN(n int64) {
	s.s.Release(n)
}
