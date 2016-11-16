+++
title = "서버에서의 채널 동작 방식 및 로그 핸들러"
date = "2016-11-14T10:20:44+09:00"
weight = 9
+++

### Data Pipeline으로서 Channel
이렇게 `Server`는 HTTP Protocol을 통해 받은 데이터를 로그 파일에 남길 형식으로 가공하여 `Log Handler`에게 전달해줍니다. 여기서 `Log Handler`는 sensor_server.go에서 함수로 구현되어 현재 Goroutine 상에서 실행되고 있습니다. 그렇다면 각 HTTP 핸들러는 `Log Handler`인 `fileLogger()` 메서드에게 데이터를 어떻게 보내야할까요? 여기서 사용되는 것이 Channel입니다.

sensor_server.go 내에서 로그컨텐트에 대한 채널을 만들어줍니다. 그리고 이 채널을 `ServeHTTP()`를 구현한 각 핸들러 구조체를 생성할 때 인자로 같이 넘겨줍니다. 이렇게 되면 각 구조체는 해당 채널에 대한 사용권을 받게 됩니다. 마찬가지로, 로그를 남기는 역할을 하는 `fileLogger()` 메서드도 채널 사용권을 인자로 넘겨받습니다. 이렇게 되면 데이터를 보내야할 `ServeHTTP()` 메서드와 데이터를 받아야할 `Log Handler`가 채널을 통해 데이터를 교환할 수 있는 환경이 만들어졌습니다!

### 양방향 채널? 단방향 채널!
기본적으로 채널을 생성하게 되면 채널은 양방향의 성격을 가지게 됩니다. 다시 말해, 채널은 특별한 언급이 없으면 어디서든 데이터를 저장하고, 불러올 수 있도록 전제하고 있습니다. 이해가 잘 안 간다고요? 그럼 예시로 설명해보죠.

```go
func ShowMeTheChannel() {
	data := make(chan int)

	for i := 0; i < 10; i++ {
		data <- i
	}

	for i := 0; i < 10; i ++ {
		fmt.Println(<-data)
	}
}
```
채널을 이용한 간단한 메서드를 만들어 봤습니다. `data`라는 채널은 int형 파이프라인으로 선언되었습니다. 그리고 `ShowMeTheChannel()`이라는 메서드는 반복문 10번씩 두 번 돌면서 각각 데이터를 저장하고, 불러옵니다. 이 메서드는 채널의 양 끝에 대한 사용권을 모두 가지고 있죠. 이것을 우리는 양방향 채널(Bidirectional Channel)이라고 합니다.

여려분이 만들 프로그램은 각 데이터 종류에 대한 핸들러와 `fileLogger()`메서드가 채널을 사용합니다. 만약 각 요소에서 데이터를 활용할 방향성에 대해 정의하지 않으면, 자신의 권한과 상관없이 데이터를 마음껏 저장하고 불러올 수 있겠죠. 이건 우리가 이상적으로 생각하는 모습이 아닙니다. 그래서 우리는 각 핸들러와 메서드가 사용할 채널의 방향성을 지정해줌으로서 파이프라인에서의 데이터의 흐름을 한 방향으로 보장하도록 하려고 합니다. 고맙게도 고언어에서는 이에 대해서도 미리 [준비](https://golang.org/ref/spec#Channel_types)를 해놨습니다.
```go
chan T          // 이런 형식으로 사용권을 받으면 채널에서 T 타입 데이터를 불러오고 저장할 수 있습니다.
chan<- float64  // 이런 형식으로 받게 되면 채널에서 T 타입 데이터를 불러오기만 할 수 있습니다.
<-chan int      // 이런 형식으로 받으면 T 타입 데이터를 채널에 저장만 할 수 있습니다.
```
그렇다면 우리는 `ServeHTTP()` 메서드에선 채널에서 데이터를 보내기만 하도록, `fileLogger()` 메서드에선 채널에서 데이터를 받기만 하도록 지정해줘야겠죠. 따라서 우리는 각 메서드의 정의에서 채널의 방향성을 아래와 같이 지정해줘야합니다. 이렇게 바꿔주면 에러가 날까봐 걱정인가요? 지금까지 잘 따라왔다면 에러없이 잘 작동할 것입니다. :)

```go
type GyroHandler struct {
	buf chan<- logContent
}

type AccelHandler struct {
	buf chan<- logContent
}

type TempHandler struct {
	buf chan<- logContent
}

func fileLogger(logStream <-chan logContent) {
}
```

### 로그, 로그를 보자!
이제 로그를 남겨봅시다! 로그는 각 센서별로 정보를 모아서 저장을 할 것입니다. 이미 눈치 채신 분들이 있겠지만, 이미 프로젝트 폴더에 로그라는 폴더가 존재합니다.
```
log/
		Accel.log
		Gyro.log
		Temp.log
```
약어가 의미하듯이 엑셀.로그는 가속도 센서의 측정값을, 자이로.로그는 자이로 센서의 측정값을, 그리고 템프.로그는 온도 및 습도계의 측정값을 저장합니다. 로그를 본격적으로 저장하기 전에, 우리가 로그를 저장할 폴더가 존재하는지 확인해야합니다. 이를 외해 [`os.Open('log')`](https://golang.org/pkg/os/#Open)이라는 메서드를 사용할 것입니다.만약 이 메서드가 로그 폴더를 정상적으로 찾아내면, 그에 대한 파일 타입을 반환할 것입니다. 아닐 경우, 로그 폴더를 [`os.Mkdir('log', os.ModePerm)`](https://golang.org/pkg/os/#Mkdir)이라는 메서드를 이용해 생성할 것입니다.

로그 폴더가 정상적으로 존재하는 경우, `Log Handler`인 `fileLogger()` 메서드는 채널에 값이 들어오기를 기다립니다. 아래의 `for := range` 문구가 그 역할을 합니다.
```
for logData := range logStream
```
이 `for, range` 문은 main 메서드에서 전달받은 채널에 데이터가 존재할지 기다리고, 존재하는 경우 바로 내부 코드를 실행합니다. 이런 과정은 프로그램이 강제로 종료하거나, 외부에서 채널을 종료(Close)하기까지 계속 반복됩니다.

채널에 데이터가 들어오게 되면 받은 데이터 구조체에서 해당 데이터가 저장될 위치를 읽어와 이를 디렉토리 주소로 조합하고 해당 로그파일을 열게됩니다. 
```go
joinee := []string{logDir, logData.location}
filePath := strings.Join(joinee, "/")

fileHandle, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
```
해당 파일이 정상적으로 열릴 경우, 다음 단계로 넘어가면 됩니다. (만약 에러가 발생하면 이에 대한 예외처리가 필요합니다.) 이제 파일에 로그를 남겨줄 `logger`를 만들겁니다. 여기서 우리는 [log 패키지](https://golang.org/pkg/log)를 이용해 직접 해당 파일로 로그를 남기는 로거를 직접 만들어보겠습니다. 아래의 [New()](https://golang.org/pkg/log/#New) 메서드는 특정 `Writer` 인터페이스로 로그를 남길 수 있도록 합니다. 여러분의 경우, 파일에 로그를 남기기로 했으니, 위에서 열어준 `fileHandle`을 이용하면 됩니다.

```go
logger := log.New(fileHandle, "", log.LstdFlags)
logger.Printf("[%s Data Received]\n%s\n", logData.sensorName, logData.content)
```
이렇게 생성한 `logger`를 이용해 로그 내용을 출력하게 되면 다음과 같은 결과를 로그파일에서 확인할 수 있습니다. 결과값은 `main()` 메서드에서 실제로 핸들러와 로거를 실행해서 결과를 확인해봅시다!

<br>
### 일해라 서버야!
이제 `main()` 메서드에서 우리가 정의해줬던 것들을 실제로 실행할 수 있게 해보겠습니다.

먼저, HTTP 핸들러 메서드에서 로그 핸들러로 기록할 로그의 내용을 전달해주는 채널, 각 센서에 대한 핸들러를 정의해줍니다.
```go
logBuf := make(chan logContent)
gyroHander := &GyroHandler{buf: logBuf}
accelHandler := &AccelHandler{buf: logBuf}
tempHandler := &TempHandler{buf: logBuf}
```
그 다음엔 실제 HTTP 서버를 실행하는 `ListenAndServe()` 메서드를 실행합니다. 여기서 주의할 것은, 그냥 실행하면 맨 위의 한 개 서버만 생성이 됩니다. `ListenAndServe()` 메서드가 실행되면 그 즉시 해당 핸들러에 대해 HTTP 요청을 대기 시작하기 때문입니다. 우리는 3개의 센서에 대한 서버가 동시에 돌아가는 것을 윈하는 것이죠? 이 때 고루틴을 적용해줍니다! 단순히 실행할 메서드 앞에 `go`를 붙여주면 됩니다.

```go
go http.ListenAndServe(":8001", gyroHander)
go http.ListenAndServe(":8002", accelHandler)
go http.ListenAndServe(":8003", tempHandler)
```

그리고 동시에 파일로거도 실행이 되어야겠죠? `fileLogger()` 메서드도 고루틴을 이용해 실행해줍니다.

```go
go fileLogger(logBuf)
```

이제 바로 실행을 해볼까요? 실행이 잘 되나요? 아마 안 될 겁니다. 왜냐하면 고루틴들이 일을 하고 있지만, 프로그램 메인 프로세스가 먼저 종료되었기 때문입니다. 기본적으로 고루틴들이 종료되기를 기다리지 않는다는 것이죠. 그래서 우리는 메인 스레드가 고루틴을 기다릴 수 있도록 한 가지를 더 추가할 것입니다. 바로 [`WaitGroup`](https://golang.org/pkg/sync/#WaitGroup)입니다. 이 `WaitGroup`에서 우리는 `Wait()`이라는 메서드를 사용해서 고루틴들이 끝나기기 까지 기다려야하지만, 그 전에 얼마나 많은 고루틴을 기다려야하는지를 정의해야합니다. 이 때 우리는 `Add()` 메서드를 이용하며, 얼마나 더 많은 고루틴을 기다릴지 `WaitGroup`에 알려줍니다. 우리의 경우 4개의 고루틴을 기다려야겠죠? (주의할 것은 `WaitGroup`은 `sync` 패키지에 정의되어 있는 타입이기 때문에 이 또한 `import`해줘야 합니다!!)

```go
import "sync"

...

wg.Add(4)
```
그리고 고루틴들을 실행하고 기다려주도록 `Wait()` 메서드를 실행해줍니다.

```go
wg.Wait()
```
이제 `WaitGroup`까지 만들어줬으면 프로그램을 실행해봅니다! 만약 정상적으로 프로그램이 실행된다면 로그 파일에 다음과 같이 표시 될 것 입니다!

```
2016/11/12 17:03:06 [GyroSensor Data Received]
Measured on 2016-11-12 17:03:06.032046895 +0900 KST
Angle Velocity of X-axis : 12.025082
Angle Velocity of Y-axis : 182.022138
Angle Velocity of Z-axis : 122.692864
```
