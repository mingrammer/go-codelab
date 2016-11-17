+++
title = "설치 및 환경셋팅"
date = "2016-11-14T10:20:44+09:00"
weight = 2
+++

### 이모! 여기 순대국(SDK) 하나요!

어떤 언어를 쓰던 간에 해당 언어에 대한 SDK가 필요하겠죠. 그리고 해당 언어를 잘 쓰기 위해 환경설정이 필요하고요. 이번 장에선 설치와 환경설정 하는 방법을 설명하려고 합니다.

<br>

### SDK 설치

공식 홈페이지 [https://golang.org/dl/](https://golang.org/dl/)에서 인스톨러를 다운로드받아 설치할 수 있습니다.

`OS X/Linux`에선 패키지 매니저로도 설치할 수 있습니다.

```bash
// OS X
brew install go

// Ubuntu
sudo apt install golang
```

설치가 완료되면 터미널(혹은 명령 프롬프트)에서 `go`를 입력하여 결과를 확인합니다. 다음과 같은 결과가 나오면 기본 설치는 잘 된 것입니다.

```bash
lkaybob@lkaybob-530u4c:/$ go
Go is a tool for managing Go source code.

Usage:

	go command [arguments]
...	
```

<br>

### 환경변수(GOPATH) 설정
Go언어에서는 Workspace라는 개념이 있습니다. 쉽게 말하자면, 여러분이 Go언어만을 사용해서 개발할 수 있는 공간입니다. 이 Workspace를 구성하면, Go언어로 개발하면서 사용하게 될 패키지들을 한 곳에서 관리할 수 있는 장점이 존재합니다. 이 Workspace는 예시로 다음과 같은 구조로 구성되어 있습니다. (Go언어 공식 홈페이지의 예시입니다.)
```
$GOPATH/
        bin/
                hello
                outyet                      
        pkg/
                linux_amd64/
                        github.com/golang/example/
                                stringutil.a
        src/
                github.com/golang/example/
                        .git/   
                        hello/
                	           hello.go 
                        outyet/
                	           main.go  
                	           main_test.go
                golang.org/x/image/
                        .git/   
                ... (many more repositories and packages omitted) ...
```
* **GOPATH** : Workspace 루트 디렉토리를 의미합니다. Go의 커맨드라인 명령어들은 이 디렉토리를 기준으로 동작합니다.
* **src** : Go 소스 파일들이 위치합니다.
* **pkg** : 패키지 오브젝트들이 위치합니다.
* **bin** : 실행 가능한 커맨드들이 위치합니다.

이 Workspace를 Go 프로젝트 디렉토리로 사용하기 위해서는 **GOPATH**라는 환경변수에 등록해야하므로 이를 위한 설정이 필요합니다.

`Windows`의 경우 '내 컴퓨터' 아이콘을 우클릭해서 속성으로 들어가면 '시스템 속성' 창이 뜹니다. 이 창의 좌측 목록을 보면 '고급 시스템 설정'이라는 항목을 누르면 또 다른 창이 뜹니다. 이 창의 아래쪽을 보면 '환경변수'로 들어갑니다.

여기서 조심해야할 것은 환경변수를 잘못 건드리면 Windows를 정상적으로 사용 못할 수 있으니 조심하셔야합니다. 여기서 우리는 환경변수를 새로 만들 것이기 때문에 '새로 만들기'를 클릭합니다. (위와 아래 목록 중 아래 목록에 하는 것을 추천합니다.) 변수 이름은 `GOPATH`로 지정하고, 값은 여러분이 지정할 Workspace의 위치를 입력합니다. 만약 바탕화면에서 `GoDev`라는 폴더를 만들었으면 다음과 같은 경로를 입력해주면 됩니다.

**Tip** : Workspace로 지정할 폴더를 탐색기에서 열고, 주소창의 경로를 복사 붙여넣기하면 쉽게 됩니다.

`OS X/Linux`의 경우, 터미널에서 `export` 명령어를 통해 설정할 수 있습니다.

```bash
cd /path/to     # 원하는 위치로 이동합니다.
mkdir work     # 원하는 위치에 Workspace 디렉토리를 만듭니다.
export GOPATH=/path/to/work      # 해당 디렉토리를 GOPATH로 설정합니다.
mkdir work/src      # 소스 파일 디렉토리를 생성합니다.
```

또한 `.profile`, `.bash_profile`, `.zsh_profile`에 경로를 등록할 수도 있습니다.

```bash
// .profile 
// .bash_profile
// .zsh_profile
// 여러분이 사용하는 쉘에 따라 설정하시면 됩니다.
...

export GOPATH=/path/to/work

...
```
