# protoc-gen-go-gin
Protobuf compiler plugin to generate Go code providing HTTP handlers and basic service implementation from .proto files

## Install

```shell
$ go install wolterskluwer.com/cwm/protoc-gen-go-gin
```
## Usage

Specify the `service=true` option to generate a Go file containing sample implementation code
for each .proto file, this also requires a `genpath=<path>` option that matches the `--go_out` flag

```shell
protoc --go_out=./internal/generated --go_opt=paths=source_relative \
    --plugin=protoc-gen-go-gin="$(which protoc-gen-go-gin)" \
    --go-gin_out=./internal/generated \
    --go-gin_opt=paths=source_relative,validate=true,service=true,genpath=./internal/generated \
    -I=./proto -I="$THIRD_PARTY_LIB"  ./proto/**/*.proto
```