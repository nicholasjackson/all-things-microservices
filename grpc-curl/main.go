package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/nicholasjackson/all-things-microservices/grpc-curl/protos/service"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var log hclog.Logger

func main() {
	log = hclog.Default()

	// create a new gRPC server, use WithInsecure to allow http connections
	gs := grpc.NewServer()

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(gs)

	cs := &CurrencyServer{}
	service.RegisterCurrencyServer(gs, cs)

	// create a TCP socket for inbound server connections
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 9092))
	if err != nil {
		log.Error("Unable to create listener", "error", err)
		os.Exit(1)
	}

	// listen for requests
	log.Info("Starting service on 0.0.0.0:9092")
	gs.Serve(l)
}

// CurrencyServer implements definition from Proto
type CurrencyServer struct{}

// GetRate is a unary gRPC function which returns a currency rate for the given currencies
func (c *CurrencyServer) GetRate(ctx context.Context, rr *service.RateRequest) (*service.RateResponse, error) {
	log.Info("GetRate called", "base", rr.GetBase().String(), "dest", rr.GetDestination().String())

	return &service.RateResponse{Rate: 23.12}, nil
}

// SubscribeRates is a gRPC bidirectional streaming endpoint which allows updates to currencies
// to be pushed to the client
func (c *CurrencyServer) SubscribeRates(sub service.Currency_SubscribeRatesServer) error {
	log.Info("SubscribeRates called")

	go func() {
		for {
			rr, err := sub.Recv()
			if err != nil {
				log.Error("Error reading from client", "error", err)
			}

			log.Info("New message from client", "base", rr.GetBase().String(), "dest", rr.GetDestination().String())
		}
	}()

	for {
		log.Info("Send message to client")
		time.Sleep(1 * time.Second)

		sub.Send(&service.RateResponse{Rate: 12.12})
	}

	return nil
}
