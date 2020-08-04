package managers

import (
	"encoding/json"
	"fmt"
	"github.com/r4stl1n/algo-benchmark-discord-bot/pkg/dto"
	"net/http"
	"strconv"
)

type HttpManager struct {
	ServiceManager *ServiceManager
}

func CreateHttpManager(serviceManager *ServiceManager) *HttpManager {

	return &HttpManager{
		ServiceManager: serviceManager,
	}
}

func (httpManager *HttpManager) healthCheckCallback(w http.ResponseWriter, req *http.Request) {

	healthCheckStruct := dto.HealthCheckStruct{
		Status:  "",
		Healthy: true,
	}

	healthCheckMarshal, healthCheckMarshalError := json.Marshal(healthCheckStruct)

	if healthCheckMarshalError != nil {
		http.Error(w, healthCheckMarshalError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(healthCheckMarshal)
}

func (httpManager *HttpManager) ListenAndServe() {

	http.HandleFunc("/healthCheck", httpManager.healthCheckCallback)

	listenError := http.ListenAndServe(":"+strconv.Itoa(httpManager.ServiceManager.Config.HttpListenPort), nil)

	if listenError != nil {
		fmt.Println(listenError)
	}
}
