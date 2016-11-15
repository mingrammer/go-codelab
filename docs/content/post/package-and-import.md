+++
title = "패키지와 임포트"
date = "2016-11-14T10:20:44+09:00"
weight = 4
+++

이전에도 살펴봤듯이 이 튜토리얼은 다음과 같이 구성되어 있습니다.

```
go-codelab/
        faker/
                range.go
        models/
                sensor.go
        sensor_client.go
        sensor_server.go
``` 

`go-codelab`은 튜토리얼 프로젝트의 루트 디렉토리이며 모든 소스코드는 이 안에 위치합니다. `go-codelab` 하위에는 `faker` 디렉토리와 `models` 디렉토리가 존재하며 `sensor_client.go`와 `sensor_server.go`가 단일 파일로 존재합니다. 

Go 프로그램은 **패키지**라고 하는 일련의 파일 그룹에 의해 관리됩니다. 패키지는 재사용이 가능하며 작은 단위의 `.go` 확장자를 갖는 프로그램들로 구성됩니다. 이 튜토리얼엔 2개의 패키지가 존재합니다. 물론 정확히 말하자면 메인 패키지를 포함하여 총 3개의 패키지가 존재하지만 메인 패키지는 조금 특별하게 다뤄지므로 패키지라고 하면 메인 패키지를 제외한 일반적인 패키지를 말합니다. 이 튜토리얼엔 2개의 패키지가 존재하는데 바로 `faker`와 `models` 패키지입니다.

```
faker/
        range.go
models/
        sensor.go
```

우리가 다룰 패키지는 비록 단일 파일만 갖고 있지만 일반적인 경우엔 한 패키지에 작은 단위로 나뉜 여러 파일과 디렉토리가 존재합니다. 쉬운 예로 나중에 살펴보게될 `net/http` 패키지의 경우 다음과 같이 구성되어 있습니다.

```
net/http/
        cgi/
        cookiejar/
                testdata/
        fcgi/
        httptest/
        httputil/
        pprof/
        testdata/
``` 

그럼 이제 우리 튜토리얼의 패키지들을 살펴봅시다. `faker`의 `range.go` 파일과 `models`의 `sensor.go`를 보면 상단에 다음과 같은 코드를 볼 수 있습니다.

```go
package faker

import (
	"math/rand"
	"time"
)
```

```go
package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mingrammer/go-codelab/faker"
)
```

최상단을 보면 `package <name>`형태의 코드를 볼 수 있습니다. 한 패키지에 속한 모든 Go 파일들은 이렇게 최상단에 자신이 속한 패키지명을 선언해야 합니다. 이 패키지명은 보통 패키지가 속한 디렉토리명을 따릅니다. 따라서, `range.go`나 `sensor.go`를 보면 각 파일이 `faker`와 `models` 패키지에 속한다는걸 바로 알 수 있습니다.
