package safego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWaitGroup_NoError(t *testing.T) {
	// Given
	as := assert.New(t)
	var wg WaitGroupInterface = &WaitGroup{}

	// When
	wg.Add(1)
	go func() {
		defer wg.Done()
	}()

	// Then
	as.NoError(wg.Wait())
}

func TestWaitGroup_Errors(t *testing.T) {
	// Given
	as := assert.New(t)
	var wg WaitGroupInterface = &WaitGroup{}

	// When
	wg.Add(3)
	go func() {
		defer wg.Done()
	}()

	go func() {
		defer wg.Done()
		panic(ErrPanicOccurred) // panic 발생
	}()

	go func() {
		defer wg.Done()
		panic(ErrPanicOccurred) // panic 발생
	}()

	// Then
	err := wg.Wait()
	as.EqualError(err, ErrMultipleErrorsOccurred)
}
