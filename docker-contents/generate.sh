#!/bin/bash

usage()
{
cat << EOF
Call with flags:
    -n name   full module name, eg: github.com/mystuff/myapp
    -i path   input path containing .proto files
    -o path   output path
    -t path   any third party libraries to include
    -x        build executable (requires main.go to be present in input path)
EOF
}

shopt -s globstar

BUILD=0
OUT="."
no_args="true"

while getopts :n:i:o:t:x option
do
    case "${option}"
        in
        n)APP=${OPTARG};;
        i)IN=${OPTARG};;
        o)OUT=${OPTARG};;
        t)TPLIB=${OPTARG};;
        x)BUILD=true;;
        :)
          echo "Option -${OPTARG} requires an argument."
          exit 1
          ;;
        *)
          echo "Invalid option: -${OPTARG}."
          usage
          exit 1
          ;;
    esac
    no_args="false"
done

[[ "$no_args" == "true" ]] && { usage; exit 1; }

THIRD_PARTY_LIB1="."
THIRD_PARTY_LIB2="."
THIRD_PARTY_LIB3="."

if [ "$TPLIB" != "" ]; then
  THIRD_PARTY_LIB1=$TPLIB
fi

if [ -d "thirdparty" ]; then
  THIRD_PARTY_LIB2="$(pwd)/thirdparty"
fi

OUT_PATH="$OUT/generated/"
IN_PATH="${IN:-.}"

pushd $IN_PATH > /dev/null

if [ -d "thirdparty" ]; then
  THIRD_PARTY_LIB3="./thirdparty"
fi

if [ -d $OUT_PATH ]; then
  rm -rf $OUT_PATH
fi

mkdir -p $OUT_PATH

protoc --go_out=$OUT_PATH --go_opt=paths=source_relative \
  --plugin=protoc-gen-go-gin="$(which protoc-gen-go-gin)" \
  --go-gin_out=$OUT_PATH --go-gin_opt=paths=source_relative,validate=true,service=true,genpath=$OUT_PATH \
  -I=./proto -I="$THIRD_PARTY_LIB1" -I="$THIRD_PARTY_LIB2" -I="$THIRD_PARTY_LIB3" ./proto/**/*.proto

go install golang.org/x/tools/cmd/goimports@latest
goimports -w $OUT_PATH/**/*.go

if [ $? -eq 0 ]; then
  mkdir -p $OUT/service
  mv $OUT_PATH/service/* $OUT/service
  rm -rf $OUT_PATH/service/
  if [ "$BUILD" = "true" ]; then
    go mod tidy
    go build
  fi
fi

popd > /dev/null