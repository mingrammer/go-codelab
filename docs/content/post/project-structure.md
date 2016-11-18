+++
title = "튜토리얼 구조"
date = "2016-11-14T10:20:44+09:00"
weight = 3
+++


이 튜토리얼은 크게 3부분(Client, Server, Log Handler)으로 나뉩니다.

![Alt project_structure](/img/project_structure.png)

기본적인 플로우는 위와 같으며 클라이언트에서 센서들이 데이터를 생성하면 이를 서버로 전송하고 서버는 이 데이터를 로깅하여 저장합니다. 

위에서 보이는 `Sensor`는 우리가 구현할 가상의 센서로 아래의 3가지 센서들이 존재합니다.

* **자이로 센서 (Gyroscope Sensor)**
* **가속도 센서 (Accelometer Sensor)**
* **온도 및 습도계 (Temperature Sensor)**

이들 센서의 종류는 잠시 후 설명할 `models` 패키지에 정의되어 있으며, 센서들은 일정 시간 간격마다 `Server`에게 `Sensor`에서 발생한 측정값을 보내줍니다. 이 측정값으로는 각속도(Angle Velocity), 선속도(Gravitational Velocity), 온도 및 습도(Temperature and Humidity)가 있습니다.

그러나, 우리는 가상의 센서를 구현하는 것이기 때문에, 실제 측정값을 넣지 못하므로 `faker`라는 랜덤 생성기를 사용하여, 랜덤값을 서버에 보내줄 것입니다.

`Server`는 각 `Sensor`에서 발생한 값을 HTTP를 통해 전달 받습니다. 이 때, 각 센서는 서로 다른 서버 포트 번호로 데이터를 전송하며 `Server`는 받은 데이터를 통해 센서의 종류를 판별하고 값들을 로깅합니다.

로깅할 데이터는 `Log Handler`에게 넘겨지며, 이를 받은 `Log Handler`는 실시간으로 각 `Sensor`가 받은 센서값을 `log` 폴더의 센서별 파일에 로그 형식으로 남깁니다.

전체 튜토리얼의 구조는 다음과 같습니다.

```
go-codelab/
        faker/
                range.go
        models/
                sensor.go
        sensor_client.go
        sensor_server.go
```
