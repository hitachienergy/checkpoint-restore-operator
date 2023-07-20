package main

import (
	"flag"
	"fmt"
	"github.com/containers/buildah"
	_ "github.com/containers/storage/drivers/overlay"
	"google.golang.org/grpc"
	"hitachienergy.com/cr-operator/agent/api"
	"hitachienergy.com/cr-operator/agent/config"
	"hitachienergy.com/cr-operator/agent/generated"
	"log"
	"net"
	"os"
	"os/exec"
)

var port = flag.Int("port", config.Port, "The server port")
var tmpDir = flag.String("tmp", config.TempDir, "The temp dir to use for checkpoints")

func main() {
	if buildah.InitReexec() {
		return
	}
	flag.Parse()
	config.TempDir = *tmpDir

	cmd := exec.Command("test")
	cmd.Start()
	err := os.MkdirAll(config.TempDir, os.ModeDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.RemoveAll(config.TempDir)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	generated.RegisterAgentServer(s, &api.AgentServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
