+++
title = "프로젝트 구조"
date = "2016-11-14T10:20:44+09:00"
weight = 3
+++


### 센서들의 종류
본 프로젝트는 크게 3부분(Server, Sensor, Log Handler)으로 나뉠 수 있습니다.

![Alt project_structure](/img/project_structure.png)

`Sensor`는 IoT 기기의 `Sensor`로, 이 구조에선 Client를 의미하며, `sensor_client.go`에서 정의되어있습니다. 이 프로젝트에서는 다음의 3개의 `Sensor`들이 존재합니다.

> #### **자이로 센서 (Gyroscope Sensor)**
> #### **가속도 센서 (Accelometer Sensor)**
> #### **온도 및 습도계 (Temperature Sensor)**

 이들 센서의 종류는 잠시 후 설명할 `model` 패키지에 정의되어 있으며, 센서들은 일정 시간 간격마다 `Server`에게 `Sensor`에서 발생한 측정값을 보내줍니다. 이 측정값으로는 각속도(Angule Velocity), 선속도(Gravitational Velocity), 온도 및 습도(Temperature and Humidity)가 있습니다.

 그러나, 우리는 가상의 센서를 구현하는 것이기 때문에, 실제 측정값을 넣지 못하는 관계로, `faker`라는 것을 이용해, 가상의 난수값을 측정값으로 설정하여 서버에 보내줄 것입니다.

 `Server`는 각 `Sensor`에 대한 측정값을 HTTP Protocol을 통해서 받으며, `sensor_server.go`에 정의되어 있습니다. 이 때, 각 측정값 종류에 대한 포트 번호는 달리해서 보내게 됩니다. `Server`는 `Sensor`로부터 받은 데이터의 내용으로 측정값의 종류와 그 값을 확인합니다.

 이렇게 확인된 데이터는 `Log Handler`에게 넘겨지며, 이를 받은 `Log Handler`는 실시간으로 각 `Sensor`가 받은 측정값을 `log` 폴더의 센서별 파일에 로그 형식으로 남깁니다.

 이러한 프로젝트 구조를 트리 형식으로 보면 다음과 같습니다.

```
go-codelab/
        faker/
                range.go
        models/
                sensor.go
        sensor_client.go
        sensor_server.go
```