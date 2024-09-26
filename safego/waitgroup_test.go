package safego

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWaitGroup_NoError_Case(t *testing.T) {
	// Given
	as := assert.New(t)
	var wg WaitGroup

	// When
	wg.SafeGo(func() {
		// Ex: 정상적으로 실행되는 함수, 패닉 없음
	})

	// Then
	as.NotPanics(wg.Wait)
}

func TestWaitGroup_WithPanic_Case(t *testing.T) {
	// Given
	as := assert.New(t)
	var wg WaitGroup

	// When
	wg.SafeGo(func() {
		panic("test panic") // Ex: 패닉 발생
	})

	// Then
	as.Panics(wg.Wait) // Wait 호출 시 패닉이 발생하는지 테스트
}
