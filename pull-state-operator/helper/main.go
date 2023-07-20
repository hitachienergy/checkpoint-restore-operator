package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"hitachienergy.com/pull-state-operator/helper/generated"
	"hitachienergy.com/pull-state-operator/helper/server"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	generated.RegisterHelperServer(s, &server.Server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
