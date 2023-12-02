#!/bin/bash
echo 'generate proto model'
protoc --proto_path=$GOPATH/src:. --go_out=. micro/proto/model/*.proto 
echo 'generate proto rpcapi'
protoc --proto_path=$GOPATH/src:. --micro_out=. micro/proto/rpcapi/*.proto

echo 'generate proto python model'
python3 -m grpc_tools.protoc -I micro/proto/model/ --python_out=micro/proto/model/ pocsuite.proto
echo 'generate proto python rpcapi'

# 生成Python-GRPC
python3 -m grpc_tools.protoc -I . --grpc_python_out=.  --python_out=. example.proto
protoc -I . --go_out=plugins=grpc:. example.proto