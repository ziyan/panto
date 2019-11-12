# panto

## Build

```
# build and strip executable
CGO_ENABLED=0 go build github.com/ziyan/panto
objcopy --strip-all panto
```

## Format

```bash
# format source code
gofmt -l -w .
```
