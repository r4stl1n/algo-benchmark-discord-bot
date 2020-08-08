package managers

import (
	"encoding/json"
	"fmt"
	"github.com/r4stl1n/algo-benchmark-discord-bot/pkg/dto"
	"github.com/sirupsen/logrus"
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

func (httpManager *HttpManager) getDailyBmCallback(w http.ResponseWriter, req *http.Request) {

	participantID := req.Header.Get("ABDB-PARTICIPANT-ID")
	restApiKey := req.Header.Get("ABDB-REST-API-KEY")

	participantExist := httpManager.ServiceManager.DatabaseClient.CheckIfParticipantExistByUUID(participantID)

	if participantExist != true {
		http.Error(w, "Participant does not exist", http.StatusForbidden)
		return
	}

	participantModel, participantModelError := httpManager.ServiceManager.DatabaseClient.GetParticipantByUUID(participantID)

	if participantModelError != nil {
		logrus.Error(participantModelError)
		http.Error(w, participantModelError.Error(), http.StatusInternalServerError)
		return
	}

	if participantModel.ApiKey != restApiKey {
		http.Error(w, "Rest api key is invalid", http.StatusForbidden)
		return
	}

	if participantModel.Approved == false {
		http.Error(w, "Need to be approved by other user", http.StatusForbidden)
		return
	}

	dailyBmModel := dto.DailyBmEntryModel{}

	dailyBmModel, _ = httpManager.ServiceManager.DatabaseClient.GetDailyBmForToday()

	dailyBmModelMarshal, dailyBmModelMarshalError := json.Marshal(dailyBmModel)

	if dailyBmModelMarshalError != nil {
		http.Error(w, dailyBmModelMarshalError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(dailyBmModelMarshal)
}

func (httpManager *HttpManager) getHistoricDailyBmCallback(w http.ResponseWriter, req *http.Request) {

	participantID := req.Header.Get("ABDB-PARTICIPANT-ID")
	restApiKey := req.Header.Get("ABDB-REST-API-KEY")

	participantExist := httpManager.ServiceManager.DatabaseClient.CheckIfParticipantExistByUUID(participantID)

	if participantExist != true {
		http.Error(w, "Participant does not exist", http.StatusForbidden)
		return
	}

	participantModel, participantModelError := httpManager.ServiceManager.DatabaseClient.GetParticipantByUUID(participantID)

	if participantModelError != nil {
		logrus.Error(participantModelError)
		http.Error(w, participantModelError.Error(), http.StatusInternalServerError)
		return
	}

	if participantModel.ApiKey != restApiKey {
		http.Error(w, "Rest api key is invalid", http.StatusForbidden)
		return
	}

	if participantModel.Approved == false {
		http.Error(w, "Need to be approved by other user", http.StatusForbidden)
		return
	}

	dailyBmModelEntries := []dto.DailyBmEntryModel{}

	dailyBmModelEntries, _ = httpManager.ServiceManager.DatabaseClient.GetAllDailyBmEntries()

	dailyBmModelMarshal, dailyBmModelMarshalError := json.Marshal(dailyBmModelEntries)

	if dailyBmModelMarshalError != nil {
		http.Error(w, dailyBmModelMarshalError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(dailyBmModelMarshal)
}

func (httpManager *HttpManager) ListenAndServe() {

	http.HandleFunc("/api/healthCheck", httpManager.healthCheckCallback)
	http.HandleFunc("/api/getDailyBm", httpManager.getDailyBmCallback)
	http.HandleFunc("/api/getHistoricDailyBm", httpManager.getHistoricDailyBmCallback)

	listenError := http.ListenAndServe(":"+strconv.Itoa(httpManager.ServiceManager.Config.HttpListenPort), nil)

	if listenError != nil {
		fmt.Println(listenError)
	}
}
