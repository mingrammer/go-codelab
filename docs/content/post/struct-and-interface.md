+++
title = "구조체와 인터페이스"
date = "2016-11-14T10:20:44+09:00"
weight = 5
+++

앞에서 잠깐 살펴본 모델 패키지 `models`에는 클라이언트와 서버에서 사용할 센서들의 공용 메서드를 정의하는 `Sensor interface`와 센서들의 공통 필드를 갖는 `SensorInfo struct`를 정의해놓은 `sensor.go` 파일이 있습니다.

구조체와 인터페이스가 무엇인지 살펴보기 전에 실제 `sensor.go`의 코드의 일부를 봅시다.

```go
// Sensor is common interface for any sensors
type Sensor interface {
    ...
}

// SensorInfo has common fields for any sensors
type SensorInfo struct {
    ...
}
```

위 코드에서 보이는 `SensorInfo`가 센서들의 공통 필드를 정의한 `struct`이며 `Sensor`가 센서들의 공용 메서드를 정의하는 `interface`입니다. 그럼 이 `struct`와 `interface`가 무엇인지 우리 코드를 보며 자세히 살펴봅시다.

<br>

### 구조체 (Struct)

구조체 `struct`는 여러 필드를 가질 수 있는 일종의 확장 타입입니다. C 언어에서의 `struct`와 거의 유사합니다.`struct`는 다음과 같이 구성됩니다.

```go
type StructName struct {
    <embededStruct>
    <varname> <vartype> `<tag>`
    ...
}
```

`StructName`은 구조체의 이름을 뜻하며 `varname`은 필드명, `vartype`은 필드의 타입, `tag`는 잠시후 다시 설명할 해당 필드가 인코딩 되었을 때의 키값을 정의하는 태그입니다. 참고로 Go에서는 변수를 선언할 때 `var name type`처럼 타입을 변수명 뒤에 선언합니다. 필드의 타입은 그 어떤 타입도 가능합니다. `embededStruct`는 조금 이따 살펴보겠습니다.
 
우선 실제 구조체를 살펴보기 위해 `SensorInfo` 코드를 봅시다.

```go
type SensorInfo struct {
	Name    string    `json:"name"`
	Type    string    `json:"type"`
	GenTime time.Time `json:"gen_time"`
}
```

`StructInfo` 구조체는 `Name`, `Type`, `Gentime`이라는 필드를 가지며 각각 `string`, `string`, `time.Time`의 타입을 갖습니다. 그렇다면 뒤에 `json:"<tag>"`는 무엇일까요? 위에서 이미 말했듯이 이는 인코딩 되었을 때의 키값을 정의합니다. 즉, 이 구조체를 `JSON` 타입으로 인코딩 했을 때 해당 필드의 키값이 `tag`로 설정된다는 뜻입니다. 위 구조체를 `JSON`으로 인코딩하기 되면 다음과 같이 인코딩 됩니다.

```
{
    'name': 'name value',
    'type': 'type value',
    'gen_time': 'gentime value'
}
```

이 태그는 필수는 아니며 필요에 따라 선택적으로 설정할 수 있습니다. 

이제 아까 언급한 `embededStruct`를 다시 살펴봅시다. `sensor.go`에 있는 실제 코드를 보겠습니다.

```go
// GyroSensor produces x-y-z axes angle velocity values
type GyroSensor struct {
	SensorInfo
	AngleVelocityX float64 `json:"x_axis_angle_velocity"`
	AngleVelocityY float64 `json:"y_axis_angle_velocity"`
	AngleVelocityZ float64 `json:"z_axis_angle_velocity"`
}
```

`sensor.go`에 정의된 센서 구조체중 하나입니다. 여기에 선언되어 있는 `SensorInfo`가 `Embedded struct`이며 이는 `Embedding`을 뜻합니다. 이는 다른 구조체의 정보를 그대로 가져와 사용하겠다는 뜻이며, 따라서 해당 구조체는 임베딩한 구조체의 모든 필드를 그대로 가질 수 있게됩니다. 즉, `GyroSensor`는 `SensorInfo`에 선언되어있는 모든 필드를 갖게됩니다. 

`Embedding`을 사용하게되면 코드 중복을 피할 수 있으며, 코드의 재사용성이 증가하게 됩니다. 얼핏보면 Java, C++, Python에서의 상속과 비슷해 보이지만 다릅니다. 상속은 하위 클래스가 상위 클래스의 모든걸 가져오는 구조인 반면 `Embedding`은 필요한 구조체의 필드들을 가져와 재사용하겠다는 의미를 가집니다. 즉, 일종의 구조체 모듈인셈입니다. Go 프로그래밍을 하게되면 앞으로 이러한 `Embedding`을 많이 볼 수 있을 것입니다.

구조체는 다음과 같이 사용할 수 있으며, 앞으로도 계속 보게될 것입니다.

```go
gyroSensor := GyroSensor{
                // Embedded 구조체 필드의 경우 다음과 같이 초기화합니다.
		SensorInfo: SensorInfo{
			Name:    "GyroSensor",
			Type:    "VelocitySensor",
			GenTime: time.Now(),
		},
		// 자체 필드는 다음과 같이 초기화합니다.
		AngleVelocityX: faker.GenerateAngleVelocity(epsilon),
		AngleVelocityY: faker.GenerateAngleVelocity(epsilon),
		AngleVelocityZ: faker.GenerateAngleVelocity(epsilon),
	}
```

<br>

### 인터페이스 (Interface)

인터페이스 `interface`는 일종의 메서드 시그니쳐의 모음입니다. Java에서의 인터페이스와 유사하며, 어떤 타입이 특정 인터페이스의 메서드들을 모두 구현하고 있으면 그 타입은 해당 인터페이스 타입을 갖게됩니다. `interface`는 다음과 같이 구성됩니다.

```go
type InterfaceName interface {
    <method signature>
    ...
}
``` 

`InterfaceName`은 인터페이스의 이름을 뜻하며, `method signature`는 메서드의 시그니쳐입니다. 

실제 인터페이스를 살펴보기 위해 `Sensor` 코드를 봅시다.

```go
type Sensor interface {
	SendingOutputString() string
	ReceivingOutputString() string
	GenerateSensorData(epsilon float64) Sensor
}
```

`Sensor`라는 인터페이스에 `SendingOutputString() string`, `ReceivingOutputString string`, `GenerateSensorData(epsilon float64) Sensor`라는 메서드 시그니쳐들이 있습니다. 만약 어떤 타입이 이 시그니쳐들을 갖는 메서드들을 모두 구현한다면 그 타입은 `Sensor interface` 타입을 갖게되며, `Sensor`를 인자로 받는 그 어떤 함수의 인자로도 들어갈 수 있습니다. 

다음은 실제 `GyroSensor` 구조체에 구현된 메서드들입니다.

```go
func (s GyroSensor) SendingOutputString() string {
    ...
}

func (s GyroSensor) ReceivingOutputString() string {
    ...
}

func (s GyroSensor) GenerateSensorData(epsilon float64) Sensor {
    ...
}
```

`GyroSensor`는 `Sensor` 인터페이스의 모든 메서드들을 구현하고 있으므로 `Sensor interface` 타입이 됩니다. 즉, `Sensor` 타입으로 사용할 수 있게 됩니다. 참고로 Go에서는 특정 타입의 메서드를 다음과 같이 선언할 수 있습니다.

```go
// StructType 타입의 st 변수를 '리시버 (receiver)'라고 하며, 리시버를 특정 구조체 타입으로 선언하면 그 타입의 메서드가 됩니다.  
func (st StructType) functionName(args) (returnTypes) {
    // procssing with 'st'
}
```

결론적으로, 그 어떤 구조체라도 `Sensor` 인터페이스의 세가지 메서드 시그니쳐들만 구현하면 `Sensor`타입으로 사용할 수 있게됩니다. Go에서 인터페이스는 타입 호환성 및 확장성 측면에서 굉장히 중요한 개념으로 Go 프로그래밍을 하게되면 인터페이스를 자주 접할 것입니다.

<br>

### 도전

`struct`와 `interface`를 사용해 다른 센서들처럼 호환 가능한 여러분만의 센서를 만들어보세요.
