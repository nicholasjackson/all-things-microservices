FROM ubuntu:20.04
RUN DEBIAN_FRONTEND=noninteractive apt update 
RUN DEBIAN_FRONTEND=noninteractive apt install git protobuf-compiler wget -y 
# Install Golang 1.17
RUN wget https://dl.google.com/go/go1.17.7.linux-amd64.tar.gz && \
    tar -xvf go1.17.7.linux-amd64.tar.gz && \
    mv go /usr/local && \
    export GOROOT=/usr/local/go
RUN git clone https://github.com/nicholasjackson/all-things-microservices && export GOPATH=/all-things-microservices/grpc-curl 
RUN export PATH=$GOPATH/bin:$GOROOT/bin:$PATH 
WORKDIR all-things-microservices 
WORKDIR grpc-curl
RUN /usr/local/go/bin/go get google.golang.org/grpc
RUN apt install golang-goprotobuf-dev -y
RUN protoc -I protos/ protos/service.proto --go_out=plugins=grpc:protos/service
RUN exec /bin/bash
CMD /usr/local/go/bin/go run /all-things-microservices/grpc-curl/main.go
