#!/usr/bin/env bash
# Generate grpc code

# Required:
# protobuf
# go-proto-validators
#

current_path=$(cd "$(dirname "${0}")" || exit 1; pwd)
cd "${current_path}"/.. || exit 1

GOPATH=$(go env GOPATH)

if [ -z "${GOPATH}" ]; then
    echo "Error: the environment variable GOPATH is not set, please set it before running"
    exit 1
fi

if ! type protoc > /dev/null 2>&1; then
  echo "Error: the command <protoc> not install, please install it before running"
  exit 1
fi

# generate code
protoc -I=. -I="${GOPATH}"/src --go_out=. qerror/errors.proto

exit $?
