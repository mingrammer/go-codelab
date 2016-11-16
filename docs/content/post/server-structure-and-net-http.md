+++
title = "서버 구조 및 net/http"
date = "2016-11-14T10:20:44+09:00"
weight = 9
+++

### 서버의 기본 동작 방식
서버는 기본적으로 3개의 포트 번호에 대한 HTTP Server를 생성합니다. 각 포트 번호를 달리한 것은 서로 다른 측정값 종류를 받기 위함이여, 각 번호는 다음의 의미를 갖습니다.

#### 8001번은 자이로스코프 센서의 데이터를,
#### 8002번은 가속도 센서의 데이터를,
#### 8003번은 온도 및 습도계의 데이터를 받습니다.

센서의 코드에서 봤듯이, HTTP 프로토콜을 이용하기 위해선 [`net/http`](https://golang.org/pkg/net/http/) 패키지가 필요합니다. 센서 클라이언트에서는 단순히 POST 요청을 보내면 됐기에 패키지를 `import`하고 추가작업이 필요하지 않았습니다.

그러나, `net/http` 패키지에서 서버를 생성하는 메서드인 [`ListenAndServe()`](https://golang.org/pkg/net/http/#ListenAndServe) 메서드를 확인해보면, 인자로 [핸들러 인터페이스](https://golang.org/pkg/net/http/#Handler)를 요구하고 있습니다. 이 핸들러 인터페이스를 다시 확인해보면, 결국 [`ServeHTTP()`](https://golang.org/pkg/net/http/#HandlerFunc)라는 메서드를 요구합니다.

따라서, 우리는 ServeHTTP() 메서드가 구현된 핸들러 인터페이스가 필요합니다. 이를 위해, 우리는 각 데이터 종류에 대한 인터페이스를 다음과 같이 만들어줘야 합니다.

~~~go
// GyroHandler : Gyroscopte sensor handler to implement ServeHTTP method
type GyroHandler struct {

}

// AccelHandler : Accelerator sensro handler to implement ServeHTTP method
type AccelHandler struct {

}

// TempHandler 	: Temperature sensor handler to implement ServeHTTP method
type TempHandler struct {

}
~~~
