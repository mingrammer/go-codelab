+++
title = "에러 핸들링"
date = "2016-11-14T10:20:44+09:00"
weight = 10
+++

### 리턴값이 여러개!
Go언어에서 에러 핸들링을 논하기 전에, Go 언어의 가장 강력한 특징을 하나 소개하려고 합니다. Go 언어는 `Multiple Return`(다중 리턴)을 지원합니다! 대부분의 프로그래밍 언어들은 단일 리턴만을 지원했지만, Go 언어에선 다중 리턴을 지원해줍니다. 아래의 코드를 한 번 볼까요?

```go
func functionName (a ...<Parameters>) (<ReturnType1>, <ReturnType2>) {
	return <RetrunValue1>, <ReturnValue2>
}
```
위의 함수 정에서 뒤쪽에 리턴값을 정의하는 부분을 볼까요? 괄호를 이용해서 두 개의 리턴타입을 정해줬습니다. 그리고 리턴 문구에서는 반환할 두 개의 값을 지정해주고 있습니다. 그렇다면 이 메서드를 불렀을 때 두 개의 리턴값을 어떻게 받아야할까요?
```go
<returnValue1>, <returnValue2>, ... := functionName(input1, input2, ...);
```
간단합니다! 메서드를 정의할 때 의미했던 리턴값들을 저장할 변수를 컴마(,)로 구분하여 차례대로 할당을 받으면 됩니다. (보통 짧은 선언문을 이용해서 바로 할당받아 사용합니다.)

<br>
### 에러 핸들링의 실체
그럼 이전에 로그파일을 열어줄 때 사용했던 `os.OpenFile()` 메서드를 다시 볼까요?

```go
fileHandle, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
```
`fileHandle`이라는 값 이외에도 err라는 값을 할당받을 수 있군요. 그럼 이것은 어떤 것을 의미할까요? `os.OpenFile()` 메서드의 정의를 한 번 확인해보죠!

```go
func OpenFile(name string, flag int, perm FileMode) (*File, error)
```
리턴값의 종류를 확인해보면 `*File`이라는 타입과 `error`라는 타입을 받을 수 있음을 확인할 수 있습니다. `*File` 타입은 특정파일에 대한 `FileHandle`임을 알 수 있습니다. 그렇다면 `error`라는 타입은 어떤 것일까요? 이는 [Go언어 자체에서 에러들을 표현하기 위해 만든 `Interface`](https://golang.org/pkg/builtin/#error)입니다.

error 인터페이스는 `Error()`라는 메서드를 가지며, 인터페이스가 불리게 될 경우, 다시 말해 에러가 발생할 경우 `Error()`라는 인터페이스는 해당 에러의 정보를 반환합니다. 에러가 발생하지 않는다면 이러한 에러 인터페이스는 발생하지 않고, 이를 리턴 타입으로 가지는 메서드에서는 이에 대해 nil값을 반환하게 됩니다. (에러가 없으니 당연히 존재하지 않겠죠?) 따라서 우리가 특정 파일을 열고, 여기에 문제가 있는지 확인하기 위해선 다음과 같이 구성하면 됩니다.

```go
fileHandle, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

if err != nil {
	log.Fatal("Error Opening File\n", err)
}

logger := log.New(fileHandle, "", log.LstdFlags)
```

<br>
### 도전!
이 프로그램에선 `erorr`를 반환하는 메서드가 많이 있습니다. 이 메서드들을 찾아서 발생할 수 있는 각 에러에 대한 메세지를 직접 정의해보세요. 