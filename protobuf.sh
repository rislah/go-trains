#!/usr/bin/env bash
set -e

# train
protoc -I. \
       -I $GOPATH/src \
       -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
       --go_out=plugins=grpc:../train/protobuf --grpc-gateway_out=logtostderr=true:../train/protobuf --swagger_out=logtostderr=true:../train/protobuf \
       train.proto

# user
protoc -I. \
       -I $GOPATH/src \
       -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
       --go_out=plugins=grpc:../user/protobuf --grpc-gateway_out=logtostderr=true:../user/protobuf --swagger_out=logtostderr=true:../user/protobuf \
       user.proto


# email
protoc -I. \
       -I $GOPATH/src \
       -I $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
       --go_out=plugins=grpc:../email/protobuf \
       email.proto
