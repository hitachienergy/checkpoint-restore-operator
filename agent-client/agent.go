package agent_client

import (
	"fmt"
	"google.golang.org/grpc"
	"hitachienergy.com/cr-operator/generated"
)

type Agent struct {
	Ip     string
	HostIp string
}

// PORT TODO: make this configurable
const PORT = 50051

func (h *Agent) CreateCheckpointImage(createCheckpointImageRequest *generated.CreateCheckpointImageRequest, opts ...grpc.CallOption) (*generated.CreateCheckpointImageResponse, error) {
	client, err := GetClient(fmt.Sprintf("%s:%d", h.Ip, PORT))
	if err != nil {
		return nil, err
	}
	ctx, cancel := DefaultContext()
	defer cancel()
	return client.CreateCheckpointImage(ctx, createCheckpointImageRequest, opts...)
}

func (h *Agent) TransferCheckpointImage(id *generated.TransferCheckpointRequest, opts ...grpc.CallOption) (*generated.TransferCheckpointResponse, error) {
	client, err := GetClient(fmt.Sprintf("%s:%d", h.Ip, PORT))
	if err != nil {
		return nil, err
	}
	ctx, cancel := DefaultContext()
	defer cancel()
	return client.TransferCheckpoint(ctx, id, opts...)
}

func (h *Agent) ExtractStats(request *generated.ExtractStatsRequest, opts ...grpc.CallOption) (*generated.StatsEntry, error) {
	client, err := GetClient(fmt.Sprintf("%s:%d", h.Ip, PORT))
	if err != nil {
		return nil, err
	}
	ctx, cancel := DefaultContext()
	defer cancel()
	return client.ExtractStats(ctx, request, opts...)
}

var agentMap = make(map[string]Agent)

func GetAgent(hostIp string) Agent {
	return agentMap[hostIp]
}

func SetAgents(agents *[]Agent) {
	var newAgents = make(map[string]Agent)
	for _, agent := range *agents {
		newAgents[agent.HostIp] = agent
	}
	agentMap = newAgents
}

func AgentExists(agent Agent) bool {
	_, ok := agentMap[agent.HostIp]
	return ok
}
