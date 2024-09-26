package safego

import (
	"fmt"
	"go-goroutine-panic-recover/helper"
	"log"
	"runtime"
	"sync"
)

// WaitGroup은 sync.WaitGroup을 wrapping한 구조체로,
// 고루틴 내에서 발생하는 패닉을 안전하게 처리할 수 있습니다.
// 고루틴이 완료된 후 패닉 정보를 수집하고 다시 발생시킵니다.
type WaitGroup struct {
	wg sync.WaitGroup
	ch chan interface{}

	once sync.Once
}

// Add는 WaitGroup 카운터를 지정된 숫자만큼 증가시킵니다.
// 패닉 수집 채널은 이 함수가 처음 호출될 때 한 번만 초기화됩니다.
func (wg *WaitGroup) Add(n int) {
	wg.once.Do(func() {
		wg.ch = make(chan interface{})
	})

	wg.wg.Add(n)
}

// Done은 WaitGroup 카운터를 감소시킵니다.
// 고루틴이 실행을 완료했을 때 호출해야 합니다.
func (wg *WaitGroup) Done() {
	wg.wg.Done()
}

// SafeGo는 주어진 함수를 고루틴에서 실행합니다.
// 함수 내에서 패닉이 발생하면, 패닉 정보가 복구되어 로깅되고
// 패닉 수집 채널에 전송됩니다. 함수 실행은 recover 블록으로 보호됩니다.
func (wg *WaitGroup) SafeGo(fn func()) {
	wg.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				// 스택 트레이스 수집
				buf := make([]byte, 1024)
				stackSize := runtime.Stack(buf, false)

				// 고루틴 ID 추적
				goroutineID := helper.GetGoroutineID()

				// 패닉 메시지, 스택 트레이스, 고루틴 ID 로깅
				log.Printf("Recovered in goroutine %d: %v\nStack trace: %s", goroutineID, r, buf[:stackSize])
				wg.ch <- fmt.Sprintf("goroutine %d panic: %v\n%s", goroutineID, r, buf[:stackSize])
			}
			wg.Done()
		}()
		fn()
	}()
}

// Wait는 모든 고루틴이 완료될 때까지 블록합니다.
// 고루틴 중 하나라도 패닉이 발생하면, 모든 고루틴이 종료된 후 패닉이 다시 발생합니다.
func (wg *WaitGroup) Wait() {
	go func() {
		wg.wg.Wait()
		close(wg.ch)
	}()

	var rs []interface{}

	for r := range wg.ch {
		rs = append(rs, r)
	}

	if len(rs) > 0 {
		panic(rs)
	}
}
