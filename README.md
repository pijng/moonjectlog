# moonjectlog
Inject fmt.Println at the beginning of each function call.

### Usage

1. Install `moonjectlog` with:

```bash
go install github.com/pijng/moonjectlog@latest
```

2. Build your project with `go build` while specifying moonjectlog preprocessor:

```bash
go build -o output -a -toolexec="moonjectlog <absolute/path/to/project>" main.go
```

**Important:**
  * `-a` flag is required to recompile all your project, otherwise go compiler might do nothing and use cached build
  * `<absolute/path/to/project>` is and absolute path to the root of your project. If you run `go build` from the root – simply specify `$PWD` as an argument.

3. Run the final binary:

```bash
./output
```

### Demonstration

Suppose we have this code:

```go
package main

import "fmt"

func main() {
	fmt.Println(someWork())
	fmt.Println(anotherWork())
	fmt.Println(finalWork())

	fmt.Println("End of main")
}

func someWork() int {
	return 1
}

func anotherWork() int {
	return 2
}

func finalWork() int {
	return 3
}
```

If we compile and run it, we get the expected output:

```bash
$ go build main.go
$ ./main
1
2
3
End of main
$
```

But if we apply `moonjectlog` as a preprocessor at compile time, we get the following result:

```bash
$ go build -a -toolexec="moonjectlog $PWD" main.go
$ ./main
Calling [main] func
Calling [someWork] func
1
Calling [anotherWork] func
2
Calling [finalWork] func
3
End of main
$
```
Thus, the preprocessor added a call to `fmt.Println("Calling [%s] func")` to the body of each function, while the source code remained unchanged.