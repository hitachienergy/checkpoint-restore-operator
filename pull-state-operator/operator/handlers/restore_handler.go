package handlers

import (
	"hitachienergy.com/pull-state-operator/generated"
	helperclient "hitachienergy.com/pull-state-operator/helper-client"
	"hitachienergy.com/pull-state-operator/util"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func HandleRestore() {
	for {
		restore := <-util.Restores
		var helper helperclient.Helper
		for _, h := range helperclient.GetHelpers() {
			log.Log.Info("helper: ", "host", h.HostIp, "ip", h.Ip)
			if h.HostIp == restore.HostIp {
				helper = h
				break
			}
		}
		log.Log.Info("restoring ", "restore", generated.RestoreSpec{
			FromId: restore.FromPid,
			Ip:     restore.Ip,
			Mode:   restore.Mode,
			Path:   restore.Path,
			Port:   restore.Port,
		}, "helper", helper.HostIp, "restore", restore.HostIp)
		_, err := helper.Restore(&generated.RestoreSpec{
			FromId: restore.FromPid,
			Ip:     restore.Ip,
			Mode:   restore.Mode,
			Path:   restore.Path,
			Port:   restore.Port,
		})
		if err != nil {
			log.Log.Error(err, "Failed to restore state")
		}
	}
}
