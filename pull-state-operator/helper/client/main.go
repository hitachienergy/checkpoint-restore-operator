package main

import (
	"context"
	"flag"
	"hitachienergy.com/pull-state-operator/helper/generated"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := generated.NewHelperClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.NewState(ctx, &generated.State{Id: "asdf-gr", State: []byte("Test-Pod")})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	_, err = c.NewState(ctx, &generated.State{Id: "another", State: []byte("Test-Pod2")})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
}
