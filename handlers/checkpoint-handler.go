package handlers

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	agentclient "hitachienergy.com/cr-operator/agent-client"
	"hitachienergy.com/cr-operator/generated"
	"hitachienergy.com/cr-operator/util"
	"io"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
	"time"
)

type kubeletCheckpointResponse struct {
	Items []string `json:"items"`
}

var configUpdate = make(map[types.UID]chan bool)

func NotifyCheckpointHandler(deploymentId types.UID) {
	if _, ok := configUpdate[deploymentId]; !ok {
		// we only need to start a new goroutine if the config is new
		configUpdate[deploymentId] = make(chan bool)
		log.Log.Info("starting checkpoint handler")
		go handleCheckpointRequests(util.ConfigMap[deploymentId])
	}
}

func handleCheckpointRequests(config *util.ConfigMapEntry) {
	startTime := time.Now().Unix()
	for {
		currentTime := time.Now().Unix()
		if int(currentTime-startTime) < config.Interval {
			duration := time.Duration(((int64)(config.Interval) - (currentTime - startTime)) * 1000_000_000)
			time.Sleep(duration)
			continue
		}
		log.Log.Info("Using", "interval", config.Interval)
		// reset start time for next loop
		startTime = currentTime

		for _, pod := range config.Pods {
			if pod.Deleted || pod.HostIp == "" || pod.Ip == "" {
				continue
			}

			log.Log.Info("creating checkpoint for", "pod", pod.Name)
			containerName := pod.KRef.Spec.Containers[0].Name
			address := fmt.Sprintf(
				"https://%s:%d/checkpoint/%s/%s/%s",
				util.KubeletAddress[pod.HostIp],
				util.KubeletPorts[pod.HostIp],
				pod.KRef.Namespace,
				pod.Name,
				containerName,
			)
			log.Log.Info("checkpoint kubelet", "address", address)
			client := getKubeleteClient()
			kubeletCheckpointStart := time.Now()
			resp, err := client.Post(address, "application/json", strings.NewReader(""))
			kubeletCheckpointEnd := time.Now()
			log.Log.Info("checkpoint request took", "time(ms)", kubeletCheckpointEnd.UnixMilli()-kubeletCheckpointStart.UnixMilli())
			if err != nil {
				log.Log.Info("http post error: ", "error", err)
				continue
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Log.Info("error while reading body: ", "error", err)
				continue
			}
			err = resp.Body.Close()
			if err != nil {
				continue
			}

			if resp.StatusCode >= 300 || resp.StatusCode < 200 {
				log.Log.Info("http post returned: ", "code", resp.StatusCode, "status", resp.Status, "body", string(body))
				continue
			}

			kubeletResponse := &kubeletCheckpointResponse{}
			err = json.Unmarshal(body, kubeletResponse)
			if err != nil {
				log.Log.Error(err, "Error unmarshalling kubelet response")
				continue
			}
			log.Log.Info("got response", "response", kubeletResponse, "body", string(body))
			agent := agentclient.GetAgent(pod.HostIp)

			saveMetrics := true
			stats, err := agent.ExtractStats(&generated.ExtractStatsRequest{
				CheckpointArchiveLocation: kubeletResponse.Items[0],
			})
			if err != nil {
				log.Log.Error(err, "Error getting stats")
				saveMetrics = false
			}

			checkpointName := pod.Name + "-checkpoint"
			createCheckpointImageStart := time.Now()
			_, err = agent.CreateCheckpointImage(&generated.CreateCheckpointImageRequest{
				CheckpointName:            checkpointName,
				CheckpointArchiveLocation: kubeletResponse.Items[0],
				ContainerName:             containerName,
			})
			createCheckpointImageEnd := time.Now()
			log.Log.Info("checkpoint image creation took", "time(ms)", createCheckpointImageEnd.UnixMilli()-createCheckpointImageStart.UnixMilli())
			if err != nil {
				log.Log.Error(err, "Unable to create checkpoint image")
				continue
			}
			log.Log.Info("transferring image to ", "recovery", pod.RecoveryNode)
			transferCheckpointImageStart := time.Now()
			_, err = agent.TransferCheckpointImage(&generated.TransferCheckpointRequest{
				CheckpointName: checkpointName,
				TransferTo:     agentclient.GetAgent(pod.RecoveryNode).Ip + ":50051",
			})
			transferCheckpointImageEnd := time.Now()
			log.Log.Info("checkpoint image transfer took", "time(ms)", transferCheckpointImageEnd.UnixMilli()-transferCheckpointImageStart.UnixMilli())
			if err != nil {
				log.Log.Error(err, "Unable to transfer checkpoint image")
				continue
			}

			if saveMetrics {
				util.SaveMetrics(stats,
					kubeletCheckpointStart.UnixMilli(),
					int(kubeletCheckpointEnd.UnixMilli()-kubeletCheckpointStart.UnixMilli()),
					createCheckpointImageStart.UnixMilli(),
					int(createCheckpointImageEnd.UnixMilli()-createCheckpointImageStart.UnixMilli()),
					createCheckpointImageStart.UnixMilli(),
					int(transferCheckpointImageEnd.UnixMilli()-createCheckpointImageStart.UnixMilli()),
					pod.Name,
					containerName,
					pod.HostIp,
				)
			}
		}
	}
}

var clientCache *http.Client

func getKubeleteClient() *http.Client {
	if clientCache != nil {
		return clientCache
	}
	clientCertPrefix := "/var/run/secrets/kubelet-certs"
	clientCert, err := tls.LoadX509KeyPair(
		fmt.Sprintf("%s/client.crt", clientCertPrefix),
		fmt.Sprintf("%s/client.key", clientCertPrefix),
	)
	if err != nil {
		log.Log.Error(err, "could not read client cert key pair")
	}
	certs := x509.NewCertPool()

	// We should really load this path dynamically as this depends on deep internals of kubernetes
	pemData, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	if err != nil {
		log.Log.Error(err, "could not read ca file")
	}
	certs.AppendCertsFromPEM(pemData)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            certs,
			Certificates:       []tls.Certificate{clientCert},
		},
	}
	return &http.Client{Transport: tr}
}
