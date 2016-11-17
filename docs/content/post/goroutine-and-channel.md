+++
title = "고루틴과 채널"
date = "2016-11-14T10:20:44+09:00"
weight = 7
+++

이번에는 이전 장에서 많이 등장했던 고루틴과 채널에 대해 간단히 다뤄보도록 하겠습니다. 고루틴과 채널은 Go에서 굉장히 중요한 개념입니다. 이 장에서는 우리 튜토리얼에서 필요한 만큼의 설명만 하도록 하겠습니다.

<br>

### 고루틴

Go에서 `고루틴(Goroutine`)이란, Go 런타임에 의해 관리되는 경량 스레드입니다. 고루틴을 사용하면 비동기적으로 여러개의 함수를 실행할 수 있으며 우리는 이를 활용해 각 센서 데이터를 동시에 서버에 전송할 수 있습니다.

고루틴은 `go` 키워드를 사용해 생성할 수 있는데 두 가지 방법으로 생성할 수 있습니다. 하나는 일반 함수를 사용하는 것이며 다른 하나는 익명 함수를 사용하는 것입니다.

```go
func work(n int) {
        for i := 0; i < n; i++ {
                fmt.Println(i)
        }
}

func main() {
        // 일반 함수로 고루틴 생성
        go work(3)
        
        // 익명 함수로 고루틴 생성
        go func(n int) {
                for i := 0; i < n; i++ {
                        fmt.Println(i)
                }
        }(3)
}
```

고루틴은 `main` 함수와는 독립적으로 실행되지만 `main` 함수가 종료되면 모든 고루틴들이 종료됩니다. 따라서 고루틴보다 `main`이 먼저 종료되는걸 방지하기위해 `sync` 라이브러리에 있는 `WaitGroup`이라는 세마포어를 활용해 고루틴의 종료를 대기하는 방법이 있는데 이는 서버 섹션에서 다룰 것입니다. 

<br>

### 채널

근데 만약에 고루틴을 사용하다가 고루틴끼리 데이터를 주고 받아야하는 경우가 생기면 어떻게 해야할까요? 바로 `채널(Channel)`을 사용하면 됩니다. 채널이란 동시에 실행되는 고루틴들을 연결해주는 일종의 파이프(pipe)입니다.

채널을 사용하면 고루틴에서 다른 고루틴으로 값을 전달할 수 있으며 다른 고루틴으로부터 값을 전달받을 수도 있습니다. 채널은 `chan` 키워드로 생성할 수 있으며 채널에 들어가는 데어터는 그 어떤 타입이라도 가능합니다. 예를 들면 정수값을 주고받는 채널의 경우 다음과 같이 선언할 수 있습니다.

```go
ourChannel chan int
```

다음과 같이 고루틴끼리 채널을 통해 데이터를 주고 받을 수 있습니다.

```go
func routine1(ourChannel chan string) {
        // 채널에 값을 넣습니다.
        ourChannel<- "data"
}

func routine2(ourChannel chan string) {
        // 채널로부터 값을 받습니다.
        fmt.Println(<-ourChannel)
        // 출력값 : data
}

func main() {
        // string 채널을 위한 메모리를 할당합니다.
        ourChannel := make(chan string)
        
        go routine1(ourChannel)
        go routine2(ourChannel)
} 
```

위에서 볼 수 있듯이 채널은 기본적으로 데이터를 주고 받을 수 있는 양방향 파이프입니다. 그런데 프로그램을 개발하다보면 어떤 고루틴은 받기만, 또 어떤 고루틴은 보내기만 하는 경우가 생길 수 있습니다. 이 경우엔 양방향 채널을 쓰는 대신 단방향 채널을 사용할 수 있습니다.

```go
// 보내기 전용 단방향 채널을 사용합니다.
func routine1(ourChannel<- chan string) {
        // 채널에 값을 넣습니다.
        ourChannel<- "data"
}

// 받기 전용 단방향 채널을 사용합니다.
func routine2(<-ourChannel chan string) {
        // 채널로부터 값을 받습니다.
        fmt.Println(<-ourChannel)
        // 출력값 : data
}

func main() {
        ...
} 
```

마지막으로, 어떤 고루틴이 특정 채널로부터 값을 받기까지 대기해야 하는 상황을 생각해봅시다.  바로 이 때 우리가 이전 장에서 보았던 `switch/case`문을 사용할 수 있습니다. Go에서의 `switch/case`문은 우리가 원래 알고 있던 그것과 동일합니다. 다만, Go에서는 이를 채널과 함께 사용할 수도 있습니다.

```go
func waitFromChannel(<-ourChannel chan int, <-yourChannel chan string>) {
        switch {
        case <-ourChannel:
                fmt.Println("Received from ourChannel")
        case val := <-yourChannel:
                fmt.Printf("Received %s from yourChannel", val)
        }
}
```

위 코드는 `waitFromChannel` 고루틴이 현재 `ourChannel`과 `yourChannel`로부터 값이 들어올 때까지 대기하고 있음을 나타냅니다. `case`에 채널을 사용하면 해당 `case`는 채널로부터 값이 들어오기까지 대기합니다. 대기하다가 해당 채널에 값이 들어오면 해당 `case`문 아래의 코드가 실행됩니다. 만약 다른 고루틴에서 `yourChannel`에 **"go"**라는 값을 넣게되면 `waitFromChannel`의 `case val := yourChannel`에서 값을 받게되고 아래의 출력문이 실행되어 **"Received go from yourChannel"**가 출력됩니다. 
