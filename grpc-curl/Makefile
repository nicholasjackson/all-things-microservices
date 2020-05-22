.PHONY: protos

protos:
	 protoc -I protos/ protos/service.proto --go_out=plugins=grpc:protos/service

run_service:
	go run main.go