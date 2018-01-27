// Package semaphore provide a semaphore implementation
package semaphore

import (
	"context"

	sema "golang.org/x/sync/semaphore"
)

// Semaphore represent a semaphore object
type Semaphore struct {
	w   *sema.Weighted
	ctx context.Context
}

// New create a semaphore object
func New(n int) *Semaphore {
	return &Semaphore{
		w:   sema.NewWeighted(int64(n)),
		ctx: context.TODO(),
	}
}

// Acquire reference increased
func (s *Semaphore) Acquire() {
	s.w.Acquire(s.ctx, 1)
}

// Release reference decreased
func (s *Semaphore) Release() {
	s.w.Release(1)
}

// AcquireN reference increased
func (s *Semaphore) AcquireN(n int64) {
	s.w.Acquire(s.ctx, n)
}

// ReleaseN reference decreased
func (s *Semaphore) ReleaseN(n int64) {
	s.w.Release(n)
}
