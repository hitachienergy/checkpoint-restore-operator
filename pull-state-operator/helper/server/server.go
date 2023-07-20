package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hitachienergy.com/pull-state-operator/helper/generated"
	"hitachienergy.com/pull-state-operator/helper/storage"
	"io"
	"net/http"
	"strings"
)

type Server struct {
	generated.UnimplementedHelperServer
}

func (server *Server) NewState(_ context.Context, state *generated.State) (*generated.Empty, error) {
	storage.StoreState(state)
	return &generated.Empty{}, nil
}

func (server *Server) DeletePod(_ context.Context, id *generated.PodId) (*generated.Empty, error) {
	fmt.Println("Deleting Pod", id.Id)
	storage.Delete(id)
	return &generated.Empty{}, nil
}

func (server *Server) Restore(_ context.Context, restoreSpec *generated.RestoreSpec) (*generated.Empty, error) {
	address := fmt.Sprintf("http://%s:%d/%s", restoreSpec.Ip, restoreSpec.Port, strings.TrimPrefix(restoreSpec.Path, "/"))
	fmt.Println("Restoring state: ", restoreSpec.FromId, " => ", address)

	state := storage.GetState(restoreSpec.FromId)

	if state == nil {
		fmt.Println("no state found!")
		return nil, errors.New("no state found")
	}

	r := bytes.NewReader(state.State)
	resp, err := http.Post(address, state.ContentType, r)
	if err != nil {
		fmt.Println("http post error: ", err)
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error while reading body: ", err)
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		fmt.Printf("http post returned code: %d status: %s\nresponse: %s\n", resp.StatusCode, resp.Status, body)
		return nil, err
	}

	return &generated.Empty{}, nil
}
