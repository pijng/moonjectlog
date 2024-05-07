# moonjectlog
Inject logging to each function call

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

