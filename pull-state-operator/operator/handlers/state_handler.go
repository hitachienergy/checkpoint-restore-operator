package handlers

import (
	"fmt"
	"hitachienergy.com/pull-state-operator/generated"
	helperclient "hitachienergy.com/pull-state-operator/helper-client"
	"hitachienergy.com/pull-state-operator/util"
	"io"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
	"time"
)

var configUpdate = make(map[types.UID]chan bool)

func NotifyStateHandler(newId types.UID) {
	if _, ok := configUpdate[newId]; !ok {
		// we only need to start a new goroutine if the config is new
		configUpdate[newId] = make(chan bool)
		go handleStateRequests(util.ConfigMap[newId], configUpdate[newId])
	}
	configUpdate[newId] <- true
}

func handleStateRequests(config *util.ConfigMapEntry, configUpdateChannel chan bool) {
	myPods := make(map[string]struct{})

	startTime := time.Now().Unix()
	for {
		select {
		case _ = <-configUpdateChannel:
			newPods := make(map[string]struct{})
			for _, pod := range config.Pods {
				newPods[pod.Uid] = struct{}{}
			}
			for podUid := range myPods {
				if _, ok := newPods[podUid]; !ok {
					unregisterPod(podUid)
					delete(myPods, podUid)
				}
			}
		default:
		}
		currentTime := time.Now().Unix()
		if int(currentTime-startTime) < config.StateProbe.Interval {
			duration := time.Duration(((int64)(config.StateProbe.Interval) - (currentTime - startTime)) * 1000_000_000)
			time.Sleep(duration)
			continue
		}
		// reset start time for next loop
		startTime = currentTime

		for _, pod := range config.Pods {
			if pod.Deleted {
				continue
			}

			resp, err := http.Get(fmt.Sprintf("http://%s:%d/%s", pod.Ip, config.StateProbe.Port, strings.TrimPrefix(config.StateProbe.Path, "/")))
			if err != nil {
				log.Log.Info("http get error: ", "error", err)
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
				log.Log.Info("http get returned: ", "code", resp.StatusCode, "status", resp.Status, "body", body)
			}

			pod.CurrentState = body
			pod.ContentType = resp.Header.Get("Content-Type")
			sendStateToHelpers(pod)
			//log.Log.Info("state is now", "state", string(body[:]), "ip", pod.ip)
		}
	}
}

func sendStateToHelpers(pod *util.PsPod) {
	helpers := helperclient.GetHelpers()
	for _, helper := range helpers {
		_, err := helper.NewState(&generated.State{
			Id:          pod.Uid,
			State:       pod.CurrentState,
			ContentType: pod.ContentType,
		})
		if err != nil {
			log.Log.Error(err, "could not send new state to helper", "helper", helper, "pod", pod.Name)
		}
	}
}

func unregisterPod(podUid string) {
	log.Log.Info("Unregister Pod: ", "pod", podUid)
	helpers := helperclient.GetHelpers()
	for _, helper := range helpers {
		unregisterPodWithHelper(podUid, helper)
	}
}

func unregisterPodWithHelper(podUid string, helper helperclient.Helper) {
	_, err := helper.DeletePod(&generated.PodId{Id: podUid})
	if err != nil {
		log.Log.Error(err, "Failed to delete Pod")
	}
}
