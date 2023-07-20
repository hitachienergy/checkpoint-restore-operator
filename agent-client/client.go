package agent_client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hitachienergy.com/cr-operator/generated"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

var clientPool = make(map[string]generated.AgentClient)
var ctx = context.Background()

func GetClient(address string) (generated.AgentClient, error) {
	if true {
		conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Log.Error(err, "did not connect")
			return nil, err
		}
		c := generated.NewAgentClient(conn)
		clientPool[address] = c
	}
	return clientPool[address], nil
}

func DefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 30*time.Second)
}
