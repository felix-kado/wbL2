Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
error

Как и в 3 задании

Когда функция test возвращает nil как *customError, переменная err в main будет содержать интерфейс типа error, значение которого равно nil, но тип все еще остается *customError. В Go интерфейс считается nil, если и только если и его типовая часть, и значение являются nil
```
...

```
