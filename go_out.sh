#!/bin/bash

set -euxo pipefail

SERVICEPATH="$GOPATH/src/github.com/dougfort/ipdnode"

protoc --proto_path=$SERVICEPATH/protobuf \
    --plugin=$GOPATH/bin/protoc-gen-go \
    --go_out=plugins=grpc:$SERVICEPATH/protobuf \
    $SERVICEPATH/protobuf/ipdnode.proto