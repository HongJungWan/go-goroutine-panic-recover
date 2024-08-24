# go-goroutine-panic-recover

[![Go](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org)

<br>

### 불편함 🤔

Go에서 고루틴 내에서 panic이 발생하면, 이 panic은 고루틴을 생성한 부모 함수로 전파되지 않는다. 

부모 함수에 recover 로직이 존재해도 고루틴 내에서 발생한 panic을 잡지 못하고 프로그램이 비정상적으로 종료된다.

위와 같은 현상이 발생하는 이유는 고루틴 내에서 panic이 발생할 경우 panicking 개념이 동작하지 않기 때문인데 해결 방법은 다음과 같다.

<br>

각 고루틴마다 panic을 개별적으로 처리하고 적절히 복구하면 된다. 

하지만, 고루틴마다 recover 로직을 구현하는 것은 매우 귀찮고 번거롭다. 아마 당신도 평소에 귀찮다고 느끼고 있었을 것이다.

그런 당신을 위해 `go-goroutine-panic-recover` 패키지를 만들었다.

`go-goroutine-panic-recover` 패키지는 고루틴 내에서 발생한 panic을 부모 함수로 안전하게 전달할 수 있는 방법을 제공한다.

<br><br><br>

### 해결 과정 📌

이 패키지는 Go의 표준 sync.WaitGroup을 확장하여, 고루틴 내에서 발생한 panic을 자동으로 복구하고, 이를 부모 함수로 전달한다. 

<br>

| 주요 스펙           | 설명                                                                                                                              |
|-----------------|---------------------------------------------------------------------------------------------------------------------------------|
| **자동 Panic 복구** | `safego.WaitGroup`의 `Done()` 메서드는 고루틴 내에서 발생한 `panic`을 자동으로 복구하고, 에러 채널에 기록한.                                                   |
| **에러 전파**       | `Wait()` 메서드는 고루틴에서 발생한 모든 `panic`을 집계하여, 부모 함수로 반환한다.                                                                          |
| **동기화**         | `sync.Mutex`와 `sync/atomic`을 사용하여 고루틴 간의 동기화와 Race Condition 문제를 관리한다.                                                          |
| **주의 사항**       | ‼ `Wait()` 메서드를 호출한 후에는 해당 `WaitGroup` 인스턴스를 재사용해서는 안 된다. <br> ‼ `Wait()` 메서드는 에러를 수신하기 위한 채널을 닫기 때문에, 재사용 시 `panic`이 발생할 수 있다. |

<br><br><br>

### 사용 방법 📌

다음은 safego.WaitGroup을 사용하는 예제 코드다.

아주 간단하다.

| 예상 사용 시나리오                                                 |
|------------------------------------------------------------|
| 1. `safego.WaitGroup` 인스턴스를 생성하고, `panic`을 발생시키는 고루틴을 추가한다. |
| 2. `Wait()` 메서드는 발생한 모든 `panic`을 집계하여 반환한다.                |
| 3. `safego.WaitGroup`을 사용함으로써, 모든 고루틴이 안전하게 종료되고, 발생한 `panic`이 부모 함수로 전달한다. |

<br>

```go
import "github.com/HongJungWan/go-goroutine-panic-recover/safego"

func GoroutinePanicWithExample() {
    defer func() {
        if r := recover(); r != nil {
            // 여기서 panic 복구 처리
        }
    }()
    
    var wg safego.WaitGroup
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        panic(safego.ErrPanicOccurred)
    }()
    
    err := wg.Wait()
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

<br><br>