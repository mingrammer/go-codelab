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

<br>

### 패키지

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
// faker/range.go
package faker

import (
	"math/rand"
	"time"
)
```

```go
// models/sensor.go
package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/mingrammer/go-codelab/faker"
)
```

최상단을 보면 `package <name>`형태의 코드를 볼 수 있습니다. 한 패키지에 속한 모든 Go 파일들은 이렇게 최상단에 자신이 속한 패키지명을 선언해야 합니다. 이 패키지명은 보통 패키지가 속한 디렉토리명을 따릅니다. 따라서, `range.go`나 `sensor.go`를 보면 각 파일이 `faker`와 `models` 패키지에 속한다는걸 바로 알 수 있습니다.

<br>

### 임포트

위의 `range.go`와 `sensor.go` 파일을 보면 `package` 아래에서 `import`문을 볼 수 있습니다.  `import`문을 통해 표준 또는 외부 패키지를 불러와 사용할 수 있습니다. 위 코드에서 보이는 `math/rand`, `time`, `fmt`, `strings`는 Go의 표준 라이브러리 패키지입니다. 표준 라이브러리 패키지의 경우 해당 패키지명만 그대로 선언해주면 바로 사용할 수 있습니다. 다음은 콘솔에 값을 출력하는 코드입니다.

```go
// 단일 import문은 이렇게 쓸 수 있습니다.
import "fmt"

// 다중 import문은 이렇게 쓸 수 있습니다.
import (
    "fmt"
    ...
)

func main() {
    fmt.Println("Hello, World!")
}
```

그런데 `sensor.go`의 `import`문을 보면 `github.com/mingrammer/go-codelab/faker`와 같이 조금 특이한 형태의 임포트문을 볼 수 있습니다. Go는 Github나 Bitbucket와 같은 곳에 호스팅 되어있는 외부 패키지와 소스코드를 패키지로써 사용할 수 있게 해줍니다. 즉, 위 임포트문은 github.com에 호스팅 되어있는 mingrammer 유저의 `go-codelab/faker` 패키지를 사용하겠다는 의미입니다. Go에서는 표준 라이브러리가 아닌 패키지의 경우 `$GOPATH/src`를 기준으로 패키지를 불러오기 때문에 이런 외부 패키지가 `import`문에 선언되어 있을 때  `go get`이라는 명령어를 사용하면 해당 위치의 패키지 구조를 `$GOPATH/src` 하위로 다운로드 받아 해당 패키지를 사용할 수 있게 해줍니다.  

우리 프로젝트의 경우 전체 Go 프로젝트 구조는 다음과 같습니다.

```
GOPATH/
        bin/
        pkg/
        src/
                github.com/
                        mingrammer/
                                go-codelab/
                                        faker/
                                                range.go
                                        models/
                                                sensor.go
                                        sensor_client.go
                                        sensor_server.go
```

이렇게 되면 어떤 Go 파일에서도 `github.com/mingrammer/go-codelab/<package name>`를 `import`에 선언하게되면 우리 프로젝트의 패키지를 사용할 수 있게 됩니다.

Go의 또 하나 특이한 점은 현재 사용되지 않는 패키지가 `import`문에 선언되어 있을 경우엔 컴파일이 되지 않습니다. 즉, 미사용 패키지를 임포트하는걸 미연에 방지할 수 있습니다. 물론 `_ (blank indentifier)`를 사용해 해당 패키지를 직접적으로 사용하지 않으면서 초기화를 위해 임포트하는 방법이 있긴 하지만 여기서는 다루지 않을 것입니다. 기본적으로 미사용 패키지의 임포트는 컴파일 에러를 발생시킨다는 것만 알아두면 됩니다. 
