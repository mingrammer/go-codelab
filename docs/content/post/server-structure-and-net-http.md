+++
title = "서버 구조 및 net/http"
date = "2016-11-14T10:20:44+09:00"
weight = 8
+++

### `Server`의 기본 동작 방식

`Server`는 기본적으로 3개의 포트 번호에 대한 HTTP 서버를 생성합니다. 각 포트 번호를 달리한 것은 서로 다른 종류의 센서값을 받기 위함이며, 각 번호는 다음의 의미를 갖습니다.

#### 8001번은 자이로스코프 센서의 데이터를,
#### 8002번은 가속도 센서의 데이터를,
#### 8003번은 온도 및 습도계의 데이터를 받습니다.

`Sensor`의 코드에서 봤듯이, HTTP 프로토콜을 이용하기 위해선 [`net/http`](https://golang.org/pkg/net/http/) 패키지가 필요합니다. `Sensor` 클라이언트에서는 단순히 POST 요청을 보내면 되므로 패키지를 `import`하고 필요한 API를 사용하는 것으로 끝났습니다.

그러나 `net/http` 패키지에서 서버를 생성하는 메서드인 [`ListenAndServe()`](https://golang.org/pkg/net/http/#ListenAndServe) 메서드를 확인해보면, 인자로 [`Handler Interface`](https://golang.org/pkg/net/http/#Handler)를 요구하고 있습니다. 이 `Handler Interface`를 다시 확인해보면, 결국 [`ServeHTTP()`](https://golang.org/pkg/net/http/#HandlerFunc)라는 메서드를 요구합니다.

<br>

### `Handler Interface` 구현하기

우리는 `ServeHTTP()` 메서드가 구현된 `Handler Interface `가 필요하므로  `Handler Interface `를 만족시키기 위한 각 센서에 대한 핸들러를 정의해줘야 합니다.

```go
type GyroHandler struct {

}

type AccelHandler struct {

}

type TempHandler struct {

}
```

위의 핸들러들이 `Handler Interface`를 만족하기 위해선 `ServeHTTP()` 메서드를 구현해야 합니다. 구조체의 메서드를 정의하기위해 `Pointer Receiver`를 사용합니다.
```go
func (m *TempHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 온도계 및 습도계에서 받은 데이터를 처리합니다.
}

func (m *GyroHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 자이로 센서에서 받은 데이터를 처리합니다.
}

func (m *AccelHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 가속도 센서에서 받은 데이터를 처리합니다.
}
```

<br>

### 데이터 까보기

이제 HTTP 프로토콜을 통해서 받은 데이터를 디코딩해봅시다. 먼저, HTTP를 통해 받은 데이터는 `Request`에 들어있겠죠? 그래서 우리는 `*http.Request` 객체를 사용하겠습니다. 클라이언트 측에선 센서에 대한 데이터를 JSON 형식으로 보냈습니다. 따라서 우리는 Go언어에서 제공해주는 [JSON 패키지](https://golang.org/pkg/encoding/json/#Unmarshal)의 `NewDecoder()` 메서드를 이용해 다음과 같이 데이터를 디코딩하겠습니다.

```go
var data models.TempSensor					// 해독한 데이터를 저장할 공간을 만들어줍니다.

decoder := json.NewDecoder(req.Body)	// Request의 Body를 해독하기 위해 JSON 패키지에서 Decoder를 새로 생성해줍니다.
err := decoder.Decode(&data)				// 생성한 Decoder를 이용해서 데이터를 해독하고 저장합니다.
if err != nil {										// 해독하는 중 발생할 수 있는 에러를 처리하기 위한 문구입니다. 이건 다음 페이지에서 설명할게요 :)
	fmt.Println("Error Occurred When Parsing Temperature Data")
}
defer req.Body.Close()							// Request의 Body를 더 이상 읽을 필요가 없으면 닫아줍니다.
```

그럼 디코딩된 데이터를 어떻게 로그 핸들러에게 전달해줘야할까요? 이것은 다음 페이지에서 자세하게 얘기할게요. 대신 `ServeHTTP()` 메서드 마지막에 이 코드만 추가해주세요.

```go
m.buf <- logContent{content: fmt.Sprintf("%s", data), location: tempLog, sensorName: data.Name}
```

이제 다음 페이지에선 디코딩된 데이터를 로그 핸들러에게 전달해보겠습니다.

<br>

### 도전!

남은 기본 센서들과 이전 단계에서 만들었던 여러분들만의 센서에 대해 `Handler Struct`와 `ServeHTTP()` 메서드를 정의해보세요!
