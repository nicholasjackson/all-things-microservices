# Examples for using gRPCurl to testing gRPC microservices
This example shows how you can use gRPCurl to interact with gRPC servers.

For full details on gRPCurl please see the offical website [https://github.com/fullstorydev/grpcurl](https://github.com/fullstorydev/grpcurl)

## Building protos for the example application
To build the gRPC client and server interfaces, first install protoc:

### Linux
```shell
sudo apt install protobuf-compiler
```

### Mac
```shell
brew install protoc
```

Then install the Go gRPC plugin:

```shell
go get google.golang.org/grpc
```

Then run the build command:

```shell
protoc -I protos/ protos/currency.proto --go_out=plugins=grpc:protos/currency
```

## Running the example application

The example service can be run using the following command

```
➜ go run main.go 
2020-05-22T16:34:40.139+0100 [INFO]  Starting service on 0.0.0.0:9092
```

## Installing gRPCurl
gRPC Curl can be installed from the following URL:

https://github.com/fullstorydev/grpcurl

or by running the following `go install` command

```shell
go install github.com/fullstorydev/grpcurl/cmd/grpcurl
```

### Reflection

By default gRPCurl will attempt to use [server reflection](https://github.com/grpc/grpc/blob/master/src/proto/grpc/reflection/v1alpha/reflection.proto) to determine the methods and types.
To enable this reflection must be explictly implemented in the API. If reflection is not availabe, gRPCurl can also use the protos
for the service. 

```shell
grpcurl -import-path ./protos -proto service.proto list 
```

### Listing Services

To list services defined by the API.

```shell
grpcurl --plaintext localhost:9092 list
Currency
grpc.reflection.v1alpha.ServerReflection
```

Parameter `--plaintext` tells gRPCurl to not use Insecure mode for the connection

```shell
-plaintext
  Use plain-text HTTP/2 when connecting to server (no TLS).
```

### Listing Methods for a Service

```shell
grpcurl --plaintext localhost:9092 list Currency        
Currency.GetRate
Currency.SubscribeRates
```

### Method detail for GetRate Service

```shell
grpcurl --plaintext localhost:9092 describe Currency.GetRate

Currency.GetRate is a method:
rpc GetRate ( .RateRequest ) returns ( .RateResponse );
```

### Method detail for SubscribeRates Message
```
grpcurl --plaintext localhost:9092 describe Currency.SubscribeRates

Currency.GetRate is a method:
rpc GetRate ( .RateRequest ) returns ( .RateResponse );
```

### RateRequest detail

Parameter `--msg-template` displays a 

```shell
-msg-template
  When describing messages, show a template of input data.
```

```shell
grpcurl --plaintext --msg-template localhost:9092 describe .RateRequest    

RateRequest is a message:
message RateRequest {
  .Currencies Base = 1;
  .Currencies Destination = 2;
}

Message template:
{
  "Base": "EUR",
  "Destination": "EUR"
}
RateRequest is a message:
```

### Execute a request with a payload

```
grpcurl --plaintext -d '{"Base": "GBP", "Destination": "USD"}' localhost:9092 Currency/GetRate
{
  "rate": 1.2229967868538965
}
```

It is also possible to read the data from stdin.

```shell
grpcurl --plaintext -d @  localhost:9092 Currency/GetRate <<EOM
{
  "Base": "GBP", 
  "Destination": "USD"
}
EOM
```

### Execute a bi-directional streaming request

The SubscribeRates service is a bi-directional streaming API, the following example will send a single
message then block while receiving messages from the server.

```shell
grpcurl --plaintext -d @  localhost:9092 Currency/GetRate <<EOM
{
  "Base": "GBP", 
  "Destination": "USD"
}
EOM
```

Client logs:

```
{
  "rate": 12.12
}
{
  "rate": 12.12
}
{
  "rate": 12.12
}
```

Server Logs:
```shell
2020-05-22T21:20:07.572+0100 [INFO]  Starting service on 0.0.0.0:9092
2020-05-22T21:20:12.053+0100 [INFO]  SubscribeRates called
2020-05-22T21:20:12.053+0100 [INFO]  Send message to client
2020-05-22T21:20:12.053+0100 [INFO]  New message from client: base=EUR dest=EUR
2020-05-22T21:20:12.053+0100 [ERROR] Client write closed
2020-05-22T21:20:13.053+0100 [INFO]  Send message to client
2020-05-22T21:20:14.054+0100 [INFO]  Send message to client
```

To keep the client stream open you can use the same call as previous but this time ommit the message.

```
grpcurl --plaintext -d @  localhost:9092 Currency/GetRate
```

You can then paste the message payload to stdin and gRPCurl will send it to the server.