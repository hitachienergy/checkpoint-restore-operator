package helper_client

import (
	"fmt"
	"google.golang.org/grpc"
	"hitachienergy.com/pull-state-operator/generated"
)

type Helper struct {
	Ip     string
	HostIp string
}

const PORT = 50051

func (h *Helper) NewState(state *generated.State, opts ...grpc.CallOption) (*generated.Empty, error) {
	client, err := GetClient(fmt.Sprintf("%s:%d", h.Ip, PORT))
	if err != nil {
		return nil, err
	}
	ctx, cancel := DefaultContext()
	defer cancel()
	return client.NewState(ctx, state, opts...)
}

func (h *Helper) Restore(restoreSpec *generated.RestoreSpec, opts ...grpc.CallOption) (*generated.Empty, error) {
	client, err := GetClient(fmt.Sprintf("%s:%d", h.Ip, PORT))
	if err != nil {
		return nil, err
	}
	ctx, cancel := DefaultContext()
	defer cancel()
	return client.Restore(ctx, restoreSpec, opts...)
}

func (h *Helper) DeletePod(id *generated.PodId, opts ...grpc.CallOption) (*generated.Empty, error) {
	client, err := GetClient(fmt.Sprintf("%s:%d", h.Ip, PORT))
	if err != nil {
		return nil, err
	}
	ctx, cancel := DefaultContext()
	defer cancel()
	return client.DeletePod(ctx, id, opts...)
}

var helperMap = make(map[Helper]struct{})

func GetHelpers() []Helper {
	ips := make([]Helper, 0, len(helperMap))
	for k := range helperMap {
		ips = append(ips, k)
	}
	return ips
}

func SetHelpers(helpers *[]Helper) {
	var newHelpers = make(map[Helper]struct{})
	for _, helper := range *helpers {
		newHelpers[helper] = struct{}{}
	}
	helperMap = newHelpers
}

func HelperExists(helper Helper) bool {
	_, ok := helperMap[helper]
	return ok
}
