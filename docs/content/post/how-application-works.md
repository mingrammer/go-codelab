+++
title = "애플리케이션 동작 과정"
date = "2016-11-14T10:20:44+09:00"
weight = 11
+++

### Go 어플리케이션 실행 방법
이제 우리 튜토리얼이 모두 완성되었습니다! 이제 직접 실행해볼까요? Go로 제작한 어플리케이션은 두 가지 방식으로 실행할 수 있습니다. 하나는 컴파일과 동시에 실행, 다른 하나는 실행 파일을 제작하고 실행하는 것입니다. 컴파일과 동시에 실행할 경우, 실행파일을 남기지 않고 실행합니다.
```
lkaybob@lkaybob-530u4c:~$ go run sample.go
```
만약 실행파일을 따로 만들고 싶다면 `build`를 하면 됩니다. 이럴 경우, 원본 소스파일의 이름으로 바이너리 형태의 실행파일이 만들어집니다.
```
lkaybob@lkaybob-530u4c:~$ go build sample.go
lkaybob@lkaybob-530u4c:~$ ls
sample.go 			sample*
```
<br>
### 우와 실행된다...!
먼저 터미널(혹은 명령 프롬프트) 창을 2개 열어놓습니다. 먼저 한 쪽에서 클라이언트쪽 코드를 실행해봅니다. 그러나 실행해보면 에러를 반환하면서 프로그램이 강제로 종료됩니다.
```
lkaybob@lkaybob-530u4c:~/GoDev/src/github.com/mingrammer/go-codelab$ go run sensor_client.go 
[VelocitySensor] SENT : AccelerometerSensor : 238.999997, 480.007518, 529.467091
[VelocitySensor] SENT : GyroSensor : 221.278028, 196.978569, 237.693392
2016/11/17 22:52:18 Error occurs when request the post data
exit status 
```
이는 클라이언트 코드를 작성할 당시, HTTP 프로토콜로 요청을 보냈지만, 현재 데이터가 존재하지 않기 때문에 에러를 반환하고 종료하도록 정의했기 때문입니다.

그럼 이번엔 다른 한 쪽에서 서버를 먼저 실행해보겠습니다. 실행할 경우, 아무것도 안 뜨면서 커서만 깜박거릴 것입니다. 에러가 난 것이나고요? 아닙니다. 지극히 정상적인 것입니다. 클라이언트 측에서 데이터를 받으면 서버는 이에 반응을 할 것입니다.
```
lkaybob@lkaybob-530u4c:~/GoDev/src/github.com/mingrammer/go-codelab$ go run sensor_server.go 

```
이제 클라이언트 코드를 실행해볼까요? 실행을 하면 잠시 후 콘솔창에 로그들이 찍히는 것을 확인할 수 있습니다. 아래는 클라이언트 측의 실행결과입니다. 각 센서의 측정값들이 정상적으로 전송되었음을 확인할 수 있습니다.

```
[VelocitySensor] SENT : GyroSensor : 191.977586, 245.389473, 236.617520
[VelocitySensor] SENT : AccelerometerSensor : 554.004670, 155.890607, 639.557394
[VelocitySensor] SENT : GyroSensor : 239.264445, 175.081075, 25.573417
[VelocitySensor] SENT : AccelerometerSensor : 670.593904, 114.496005, 272.511751
[VelocitySensor] SENT : AccelerometerSensor : 1033.903791, 441.048931, 454.012025
[VelocitySensor] SENT : GyroSensor : 73.746588, 200.838994, 363.733532
[AtmosphericSensor] SENT : TemperatureSensor : 88.013658, 32.216430
...
```

다음은 서버측 콘솔 화면을 보겠습니다. 실행을 했을 당시만 해도 꿈쩍을 안 하던 콘솔 화면이 변화를 보였습니다! 서버 측은  먼저 도달한 데이터 순으로 차례대로 처리하며 해당 데이터를 처리했음을 콘솔 화면에서 보여줍니다.
```
[VelocitySensor] RECEIVED : GyroSensor : 191.977586, 245.389473, 236.617520
[VelocitySensor] RECEIVED : AccelerometerSensor : 554.004670, 155.890607, 639.557394
[VelocitySensor] RECEIVED : GyroSensor : 239.264445, 175.081075, 25.573417
[VelocitySensor] RECEIVED : AccelerometerSensor : 670.593904, 114.496005, 272.511751
[VelocitySensor] RECEIVED : AccelerometerSensor : 1033.903791, 441.048931, 454.012025
[VelocitySensor] RECEIVED : GyroSensor : 73.746588, 200.838994, 363.733532
[VelocitySensor] RECEIVED : GyroSensor : 154.033246, 335.525676, 24.761997
[VelocitySensor] RECEIVED : AccelerometerSensor : 480.081296, 148.277650, 491.804754
[AtmosphericSensor] RECEIVED : TemperatureSensor : 88.013658, 32.216430
...
```

다만 표시되는 데이터 순서가 다를 수 있지만, 이것은 서버 측 채널에 먼저 도달한 순서대로 표시하기 때문에 일부 차이가 있을 수 있습니다. 이제 로그에 어떻게 기록되었는지 확인해볼까요? `Accel.log`를 먼저 열어봤습니다.

```
2016/11/17 23:02:06 [AccelerometerSensor Data Received]
Measured on 2016-11-17 23:02:06.996888976 +0900 KST
Gravitational Velocity of X-axis : 554.004670
Gravitational Velocity of Y-axis : 155.890607
Gravitational Velocity of Z-axis : 639.557394

2016/11/17 23:02:07 [AccelerometerSensor Data Received]
Measured on 2016-11-17 23:02:07.496897415 +0900 KST
Gravitational Velocity of X-axis : 670.593904
Gravitational Velocity of Y-axis : 114.496005
Gravitational Velocity of Z-axis : 272.511751

2016/11/17 23:02:07 [AccelerometerSensor Data Received]
Measured on 2016-11-17 23:02:07.996894498 +0900 KST
Gravitational Velocity of X-axis : 1033.903791
Gravitational Velocity of Y-axis : 441.048931
Gravitational Velocity of Z-axis : 454.012025

2016/11/17 23:02:08 [AccelerometerSensor Data Received]
Measured on 2016-11-17 23:02:08.49689645 +0900 KST
Gravitational Velocity of X-axis : 480.081296
Gravitational Velocity of Y-axis : 148.277650
Gravitational Velocity of Z-axis : 491.804754

.....
```
아주 잘 정돈된 형식으로 기록되어 있습니다! 다른 센서들에 대한 로그 파일도 열어보면 각 센서에 맞는 형식으로 아주 정돈된 형태로 기록이 남아있을 겁니다.

<br>
### 도전
여러분이 직접 정의했던 센서들에 대해서도 로그가 잘 남는지 확인해보세요. 안 된다면 어디서 문제인지 한 번 찾아보고 해결해보세요!