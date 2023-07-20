package api

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hitachienergy.com/cr-operator/agent/config"
	"hitachienergy.com/cr-operator/agent/generated"
	"hitachienergy.com/cr-operator/agent/internal"
	"hitachienergy.com/cr-operator/agent/util"
	"io"
	"os"
)

type AgentServer struct {
	generated.UnimplementedAgentServer
}

func (s *AgentServer) AcceptCheckpoint(stream generated.Agent_AcceptCheckpointServer) error {
	fmt.Println("Starting checkpoint Image transfer")
	var file *os.File
	var fileName string
	var receivedBytes int
	for {
		checkpointArchive, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Checkpoint Image Transfer completed. Received %s (%d)\n", util.ToHumanReadable(receivedBytes), receivedBytes)
				err := stream.SendAndClose(&generated.AcceptCheckpointResponse{})
				if err != nil {
					return err
				}
				err = file.Close()
				if err != nil {
					return err
				}
				return internal.ImportImage(fileName)
			}

			fmt.Println("recv error ", err)
			return err
		}

		if file == nil && checkpointArchive.CheckpointImageName != "" {
			fileName = fmt.Sprintf("%s/%s", config.TempDir, checkpointArchive.CheckpointImageName)
			file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			fmt.Printf("Writing to %s\n", fileName)
		}

		written, err := file.Write(checkpointArchive.CheckpointOCIArchive)
		if err != nil {
			return err
		}
		receivedBytes += written
	}
}

func (s *AgentServer) CreateCheckpointImage(ctx context.Context, request *generated.CreateCheckpointImageRequest) (*generated.CreateCheckpointImageResponse, error) {
	fmt.Println("Got CreateCheckpointImage request")
	//_, _, err := internal.CreateCheckpointImage(ctx, request.CheckpointArchiveLocation, request.ContainerName, request.CheckpointName)
	err := internal.CreateOCIImage(request.CheckpointArchiveLocation, request.ContainerName, request.CheckpointName)
	if err != nil {
		return nil, err
	}
	return &generated.CreateCheckpointImageResponse{}, nil
}

func (s *AgentServer) TransferCheckpoint(ctx context.Context, request *generated.TransferCheckpointRequest) (*generated.TransferCheckpointResponse, error) {
	conn, err := grpc.Dial(request.TransferTo, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err, "did not connect")
		return nil, err
	}
	c := generated.NewAgentClient(conn)
	client, err := c.AcceptCheckpoint(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println("opening checkpoint archive")
	checkpointArchive, err := os.Open(fmt.Sprintf("%s/%s", config.TempDir, request.CheckpointName))
	if err != nil {
		return nil, err
	}
	defer checkpointArchive.Close()

	bytesSent := 0
	buffer := make([]byte, 4000)
	first := true
	for {
		bytesRead, err := checkpointArchive.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("ignoring EOF")
				break
			}
			return nil, err
		}

		bytesSent += bytesRead

		if bytesRead < len(buffer) {
			// only send slice here to avoid sending superfluous data
			err = client.Send(&generated.AcceptCheckpointRequest{CheckpointOCIArchive: buffer[:bytesRead]})
			if err != nil {
				fmt.Println("error while sending last slice ", err)
				return nil, err
			}
			break
		}

		acceptRequest := &generated.AcceptCheckpointRequest{CheckpointOCIArchive: buffer}
		if first {
			acceptRequest.CheckpointImageName = request.CheckpointName
			first = false
		}
		err = client.Send(acceptRequest)
		if err != nil {
			fmt.Println("error while sending ", err)
			return nil, err
		}
	}
	fmt.Println("sent bytes", bytesSent)

	fmt.Println("executing close send")
	response, err := client.CloseAndRecv()
	if err != nil {
		fmt.Println("CloseSend error ", err)
		return nil, err
	}
	fmt.Println(response)

	return &generated.TransferCheckpointResponse{}, nil
}

func (s *AgentServer) ExtractStats(ctx context.Context, request *generated.ExtractStatsRequest) (*generated.StatsEntry, error) {
	fmt.Println("Got ExtractStats request")
	response := internal.ReadStats(request.CheckpointArchiveLocation)

	// TODO inefficient: unmarshalling to marshal again...
	var statsEntry generated.StatsEntry
	err := proto.Unmarshal(response, &statsEntry)
	if err != nil {
		return nil, err
	}
	return &statsEntry, nil
}
