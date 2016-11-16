+++
title = "클라이언트 구조"
date = "2016-11-14T10:20:44+09:00"
weight = 6
+++

이전에 패키지를 다룰 때 우리는 메인 패키지라고 하는 조금 특별한 패키지에 대해 언급한 적이 있습니다. Go 프로그램이 실행가능한 바이너리 파일로 컴파일되기 위해선 메인 패키지가 반드시 필요합니다. 즉 실행 가능한 Go 프로그램들은 항상 메인 패키지를 가지고 있어야 하며 메인 패키지를 가진 프로그램은 반드시 `main` 함수를 선언해야합니다. `main` 함수는 조금 이따 다루겠습니다.

<br>

### 클라이언트의 기본 동작 방식

클라이언트 `sensor_client.go`가 하는 일은 간단합니다. `models`에 정의된 센서들로부터 일정 간격으로 센서 데이터를 생성하고 이를 서버에 보내는 역할을 합니다. 기본적인 플로우는 다음과 같습니다.

* 센서값을 서버로 보내기위한 워커(worker)를 구동하기 위한 정보를 가진 워커 구조체를 정의합니다.
* 고루틴(Goroutine)을 사용해 워커 정보의 갯수만큼 워커를 구동합니다.
* 각 워커는 일정 간격으로 데이터를 생성해 서로 독립적으로 서버에 데이터를 전송합니다.
* 서버 데이터 전송이 일정횟수가 넘어가면 모든 워커를 종료시켜 클라이언트 프로그램을 종료합니다.

위 플로우를 기반으로 실제 코드(`sensor_client.go`)의 각 부분에 대해 자세히 살펴보도록 하겠습니다.

<br>

### 워커(worker)와 카운터(counter) 구조체 정의하기

우리는 각 센서 데이터를 독립적으로 서버에 보낼 수 있도록 고루틴 위에서 센서 워커를 돌릴 것 입니다. 따라서 워커를 돌리기 위해 필요한 정보들을 갖는 `worker` 구조체를 정의할 것입니다. 고루틴에 대해선 다음 장에서 더 자세히 살펴볼 것이며 지금은 일단 경량 스레드쯤으로 이해하시면 됩니다.

```go
type worker struct {
	ticker      *time.Ticker
	sensor      models.Sensor
	sensorError float64
	serverPort  int
}
```

각 워커는 일정 간격마다 데이터를 생성 및 전송하므로 일정 간격마다 신호를 생성할 수 있는 `ticker`를 필드로 넣습니다. 이 `ticker`의 타입은 `*time.Ticker`인데 이는 `time`이라는 표준 라이브러리 패키지에 정의되어 있으며 원하는 주기마다 `Ticker`로부터 신호를 받을 수 있습니다.

`sensor`는 각 센서 구조체들을 저장히기 위한 필드로 `Sensor`의 타입으로 선언되어 있습니다. `Sensor` 구조체는 `models` 패키지에 존재하므로 `models.Sensor`로 접근합니다. [구조체와 인터페이스](/struct-and-interface)에서 보았듯이 `Sensor` 인터페이스를 만족하면 그 어떤 구조체든 `Sensor` 타입으로 선언이 가능합니다. 따라서, 종류와 상관없이 각 센서들을 저장할 수 있습니다. `serverError`는 임의의 센서값을 생성할 때 사용되는 센서 오차값을 뜻하며 `serverPort`는 각 센서별로 서버의 어느 포트로 데이터를 보낼지 설정하기 위한 필드입니다.

또한 우리는 데이터 전송 횟수를 제한하기 위해 `counter` 구조체를 정의할 것입니다.

```go
type counter struct {
	mutex sync.Mutex
	n     int
}

func (c *counter) count() {
	c.mutex.Lock()
	c.n++
	c.mutex.Unlock()
}

func (c *counter) value() int {
	return c.n
}
```

바로 아래에서 살펴볼거지만 `counter`의 경우 여러개의 고루틴에서 동시에 사용하기 때문에 필드값을 변경할 때 동기화 처리를 해줘야합니다. 다행히도 Go는 동기화와 관련된 기능들을 `sync`라는 표준 라이브러리 패키지로 제공하고 있으며 `sycn.Mutex`를 통해 뮤텍스 처리를 할 수 있도록 해줍니다. 따라서, `sync.Mutex` 타입을 갖는 `mutex`를 필드로 가지면 `count()` 함수에서 볼 수 있듯이 특정값을 변경할 때 여러개의 고루틴이 동시에 값을 변경 할 수 없도록 `Lock`을 걸 수 있습니다. 

그럼 이제 `worker`와 `counter`가 실제로 어떻게 사용되고 있는지 살펴봅시다.

<br>

### 센서 워커 (sensorWorker) 함수 살펴보기

`worker`에 정의된 워커 정보를 가지고 실제 센서 데이터를 서버로 보내는 워커 함수인 `sensorWorker`를 살펴봅시다.

```go
func sensorWorker(done <-chan struct{}, w worker, c *counter) {
	for {
		select {
                // done 채널이 닫히기까지 대기
		case <-done:
			return
		// ticker 신호 대기
		case <-w.ticker.C:
                        // 센서 데이터 생성
			sensorData := w.sensor.GenerateSensorData(w.sensorError)
			 // serverPort 값으로 서버 URL 생성
			url := getRequestServerURL(w.serverPort)

			fmt.Println(sensorData.SendingOutputString())

                        // 서버로 데이터 전송
			sendJSONSensorData(url, sensorData)

                        // 전송할 때 마다 카운팅
			c.count()
		}
	}
}
```

`sensorWorker`는 세 개의 인자를 받는데 각 인자는 다음을 의미합니다.

* `done <-chan struct{}` : 고루틴을 종료하기 위한 신호를 받는 받기 전용 채널(channel)입니다. 채널에 대한 자세한 내용은 다음 장에서 살펴보겠습니다. 지금은 이 `done`이 고루틴을 종료시키기 위한 채널 변수라고만 알아두면 됩니다. 
* `w worker` : 좀 전에 위에서 만든 `worker`를 받는 인자입니다. 워커 정보를 받기 위한 인자입니다.
* `c *counter` : 전송 횟수를 카운팅 하기 위한 `*counter` 변수입니다. 여러개의 고루틴에서 사용하기 위해 포인터 타입으로 선언하였습니다.

이 함수는 무한 루프 안에서 `select/case`문을 실행합니다. `switch/case`문에 대해서는 다음 장에서 더 자세히 살펴볼 것이며 지금은 `case`에 있는 채널에 값이 들어올 경우 해당 `case` 아래의 코드가 실행된다는 것만 알아두시면 됩니다. 따라서 이는 `worker`에 정의된 `ticker`의 `C` 채널로부터 신호를 받을 때마다 센서의 데이터를 생성하고 이를 서버로 전송합니다. 또한 `done`에 값이 들어올 경우엔 `return`을 실행하며 함수가 종료됩니다.

<br>

### main 함수 살펴보기

위에서 언급했듯이 메인 패키지는 반드시 `main` 함수를 필요로하며 메인 패키지를 가진 Go 프로그램이 실행될 때 바로 이 `main` 함수가 실행됩니다.

```go
func main() {
        // 전송 횟수 카운팅을 위한 카운터 생성 
	var sendCounter counter
        // 동기화를 위한 일종의 세마포어. 이에 대해선 서버 섹션에서 자세히 설명하겠습니다.
	var wg sync.WaitGroup
        // 상수는 const로 정의합니다.
	const numWorkers = 3

        // 고루틴 종료 신호를 위한 채널 생성
	done := make(chan struct{})

        // 워커의 갯수만큼의 카운트값을 갖는 세마포어 생성
        wg.Add(numWorkers)    
    
        // 각각의 Ticker와 Sensor의 인스턴스를 생성해 worker에 저장합니다.
	workerList := [numWorkers]worker{
		worker{
			ticker:      time.NewTicker(500 * time.Millisecond), // 0.5초 간격
			sensor:      models.GyroSensor{}, // GyroSensor 인스턴스
			sensorError: 4.0,
			serverPort:  8001,
		},
		worker{
			ticker:      time.NewTicker(500 * time.Millisecond), // 0.5초 간격
			sensor:      models.AccelSensor{}, // AccelSensor 인스턴스
			sensorError: 12.0,
			serverPort:  8002,
		},
		worker{
			ticker:      time.NewTicker(2 * time.Second), // 2초 간격
			sensor:      models.TempSensor{}, // TempSensor 인스턴스
			sensorError: 1.5,
			serverPort:  8003,
		},
	}

        // workerList를 순회하며 각 워커를 가져옵니다.
        // range로 순회하게되면 index와 value가 리턴되는데 현재는 index를 사용하지 않으므로 _로 무시합니다.
	for _, w := range workerList {
		go func(w worker) {
			sensorWorker(done, w, &sendCounter)
			// 세마포어 카운트값 감소
			wg.Done()
		}(w)
	}
	
	// 전송 횟수를 체크하기 위한 고루틴
	go func() {
		for {
			if sendCounter.value() > 100 {
				close(done)
				return
			}
		}
	}()
        
        // 세마포어 대기
        wg.Wait()	
}
```

`main` 함수에선 3개의 `worker`를 만들고 `go`키워드를 사용해 `sensorWorker`를 실행하는 고루틴을 생성합니다. `go` 키워드에 대해선 다음 장에서 고루틴/채널과 함께 자세히 설명할 것입니다.

이렇게 3개의 고루틴을 생성하게 되면 각 고루틴은 `sensorWorker`를 각각 독립적으로 실행하게되며, 각 워커마다 정의된 `Ticker`의 주기에 따라 서버에 센서 데이터를 전송한 후 전송 횟수를 카운팅 할 것입니다.

그런데 아래를 보면 또 하나의 고루틴을 생성하고 있는걸 볼 수 있습니다. 이 고루틴은 무한 루프를 돌며 카운터의 값을 체크하고 있는데 값이 100을 넘어가면 `done` 채널을 닫은 후 `return`으로 함수를 종료시킵니다. 이는 직관적으로 위에서 생성한 `sensorWorker` 고루틴들의 데이터 전송 횟수가 100회를 넘어가면 고루틴들을 종료 시킨다는 것을 알 수 있습니다. 조금 더 자세히 말하자면, `done` 채널을 닫게되면 `sensorWorker`의 `case <-done:`이 작동하게 되고 따라서 이 `case`문의 코드인 `return`이 실행되면서 모든 `sensorWorker` 고루틴들이 종료됩니다.

즉, 클라이언트 프로그램은 센서의 갯수만큼 워커 고루틴들을 생성하여 주기적으로 센서 데이터를 서버로 전송하며,  총 전송 횟수가 100회를 넘어가면 모든 고루틴을 종료하고 프로그램을 종료하는 방식으로 동작함을 알 수 있습니다!

<br>

### 도전!

이전 장에서 만든 센서를 위한 워커를 생성하여 여러분만의 센서 데이터를 전송해보세요!
