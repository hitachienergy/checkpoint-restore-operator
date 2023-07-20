package helper_client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hitachienergy.com/pull-state-operator/generated"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

var clientPool = make(map[string]generated.HelperClient)
var ctx = context.Background()

func GetClient(address string) (generated.HelperClient, error) {
	if true {
		conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Log.Error(err, "did not connect")
			return nil, err
		}
		c := generated.NewHelperClient(conn)
		clientPool[address] = c
	}
	return clientPool[address], nil
}

func DefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 5*time.Second)
}
