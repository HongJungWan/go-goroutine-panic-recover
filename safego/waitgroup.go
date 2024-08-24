package safego

import (
	"errors"
	"sync"
	"sync/atomic"
)

const (
	ErrPanicOccurred          = "panic occurred"
	ErrMultipleErrorsOccurred = "one or more errors occurred"
)

type WaitGroupInterface interface {
	Add(int)
	Done()
	Wait() error
}

type WaitGroup struct {
	wg   sync.WaitGroup
	mu   sync.Mutex
	once sync.Once

	ch     chan error
	errCnt int32
}

func (wg *WaitGroup) Add(n int) {
	wg.once.Do(func() {
		wg.ch = make(chan error)
	})
	wg.wg.Add(n)
}

func (wg *WaitGroup) Done() {
	defer wg.wg.Done()
	if r := recover(); r != nil {
		wg.mu.Lock()
		defer wg.mu.Unlock()

		wg.ch <- errors.New(ErrPanicOccurred)
		atomic.AddInt32(&wg.errCnt, 1)
	}
}

func (wg *WaitGroup) Wait() error {
	go func() {
		wg.wg.Wait()
		close(wg.ch)
	}()

	var errs []error
	for err := range wg.ch {
		errs = append(errs, err)
	}

	if atomic.LoadInt32(&wg.errCnt) > 0 {
		return errors.New(ErrMultipleErrorsOccurred)
	}

	return nil
}
