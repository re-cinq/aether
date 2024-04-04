# Aethers Plugin Systems

We have based our plugin system off of [hashicorps go-plugins](https://github.com/hashicorp/go-plugin/tree/main). This was mainly so we could detach the need for matching dependencies. We use gRPC for communication

## Example Plugin

There is an example plugin in the [example directory](./example/example.go) with comments on the moving parts. 

## Install a Plugin

// TODO //

## Create a Docker Image with your Plugin

// TODO //

## Requirements for local setup

as we use gRPC and protobuffers, if you want to make changes to the `.proto`
files you will need to install buf
```bash
go install github.com/bufbuild/buf/cmd/buf@v1.30.1
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0

buf generate
```
