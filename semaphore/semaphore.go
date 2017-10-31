// Package semaphore provide a semaphore implementation
package semaphore

// Semaphore represent a semaphore object
type Semaphore struct {
	c chan int
}

// New create a semaphore object
func New(n int) *Semaphore {
	s := &Semaphore{
		c: make(chan int, n),
	}
	return s
}

// Acquire reference increased
func (s *Semaphore) Acquire() {
	s.c <- 0
}

// Release reference decreased
func (s *Semaphore) Release() {
	<-s.c
}
