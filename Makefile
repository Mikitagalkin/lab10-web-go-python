.PHONY: proto build run-http run-grpc run-python test

proto:
protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       go-service/proto/calculator.proto
python -m grpc_tools.protoc -I./go-service/proto `
       --python_out=./python-service `
       --grpc_python_out=./python-service `
       ./go-service/proto/calculator.proto

build:
cd go-service && go build ./...
cd go-service/grpc_server && go build .

run-http:
cd go-service && go run main.go

run-grpc:
cd go-service/grpc_server && go run main.go

run-python:
cd python-client && python http_client.py
cd python-client && python tcp_client.py

test:
cd python-client && pytest -v
