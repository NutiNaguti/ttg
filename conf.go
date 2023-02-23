package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/charmbracelet/log"
)

type Config struct {
	ApiId              int32  `json:"api_id"`
	ApiHash            string `json:"api_hash"`
	DatabaseDirectory  string `json:"database_directory"`
	UseMessageDatabase bool   `json:"use_message_database"`
	UseSecretChats     bool   `json:"use_secret_chats"`
	SystemLanguageCode string `json:"system_language_code"`
	DeviceModel        string `json:"device_model"`
	AppVersion         string `json:"application_version"`
}

func GetConfig() *Config {
	var conf Config

	file, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Error(err)
	}

	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Error(err)
	}

	return &conf
}
