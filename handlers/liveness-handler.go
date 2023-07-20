package handlers

import (
	"context"
	"fmt"
	"hitachienergy.com/cr-operator/util"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
	"time"
)

var livenessHandler = make(map[types.UID]struct{})
var kClient client.Client

func NotifyLivenessHandler(client client.Client, newId types.UID) {
	kClient = client
	if _, ok := livenessHandler[newId]; !ok {
		// we only need to start a new goroutine if the config is new
		livenessHandler[newId] = struct{}{}
		log.Log.Info("starting liveness handler")
		go handleLivenessProbe(util.ConfigMap[newId])
	}
}

func handleLivenessProbe(config *util.ConfigMapEntry) {
	probeFailures := make(map[string]int)
	startTime := time.Now().Unix()
	for {
		currentTime := time.Now().Unix()
		if int(currentTime-startTime) < config.LivenessProbe.Interval {
			duration := time.Duration(((int64)(config.LivenessProbe.Interval) - (currentTime - startTime)) * 1_000_000_000)
			time.Sleep(duration)
			continue
		}
		// reset start time for next loop
		startTime = currentTime

		for _, pod := range config.Pods {
			if pod.Deleted {
				delete(probeFailures, pod.Uid)
				continue
			}

			if pod.Ip == "" {
				continue
			}

			address := fmt.Sprintf("http://%s:%d/%s", pod.Ip, config.LivenessProbe.Port, strings.TrimPrefix(config.LivenessProbe.Path, "/"))
			resp, err := livenessRequest(address)
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 400 {
				probeFailures[pod.Uid] = 0
				continue
			}
			if err != nil {
				log.Log.Info("error while checking liveness", "error", err)
			} else {
				log.Log.Info("error while checking liveness", "code", resp.StatusCode, "status", resp.Status)
			}

			probeFailures[pod.Uid]++
			if probeFailures[pod.Uid] != 3 {
				continue
			}

			log.Log.Info("deleting Pod", "pod", pod.Name, "host", pod.HostIp, "pod.kref", pod.KRef.Name)

			podName := util.CreatePod(kClient, context.Background(), config.LabelSelector, pod)
			util.RestorePods[podName].OldPod = pod.KRef
		}
	}
}

func livenessRequest(address string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", address, nil)
	if err != nil {
		log.Log.Error(err, "error creating request")
		return nil, err
	}
	return http.DefaultClient.Do(req)
}
