#!bin/bash

protoc --go_out=server/. \
    --go-grpc_out=require_unimplemented_servers=false:server/. \
    protoc/blog.proto
