package main

import (
	"encoding/json"
	"github.com/r4stl1n/algo-benchmark-discord-bot/pkg/dto"
	"github.com/r4stl1n/algo-benchmark-discord-bot/pkg/managers"
	"github.com/r4stl1n/algo-benchmark-discord-bot/pkg/util"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func main() {

	configStruct := dto.ConfigStruct{}

	util.PrintBanner()

	logrus.SetLevel(logrus.DebugLevel)

	logrus.Info("Loading configuration file")

	configFile, configLoadError := ioutil.ReadFile("config.json")

	if configLoadError != nil {

		logrus.Warn("Could not load the config file")
		logrus.Info("Creating new config file, modify and run again")

		marshaledStruct, marshaledStructError := json.MarshalIndent(configStruct, "", "")

		if marshaledStructError != nil {
			logrus.Fatal("Failed to marshal blank structure")
		}

		writeFileError := ioutil.WriteFile("config.json", marshaledStruct, 0644)

		if writeFileError != nil {
			logrus.Fatal("Could not write empty config file")
		}

		return
	}

	logrus.Info("Config file loaded")
	configUnmarshalError := json.Unmarshal(configFile, &configStruct)

	if configUnmarshalError != nil {
		logrus.Fatal("Could not unmarshal config file")
	}

	logrus.Info("Connecting to database")

	databaseManager, databaseManagerError := managers.CreateDatabaseManager("data.db")

	if databaseManagerError != nil {
		logrus.Panic(databaseManagerError)
	}

	serviceManager := managers.CreateServiceManager(&configStruct, databaseManager)

	serviceInitError := serviceManager.Initalize()

	if serviceInitError != nil {
		logrus.Fatal(serviceInitError)
	}

	logrus.Info("Discord service active")

	httpManager := managers.CreateHttpManager(serviceManager)

	logrus.Info("Starting the http service")
	httpManager.ListenAndServe()
}
