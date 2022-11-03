package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	// Drobo is host:port for the drobo instance
	Drobo string `json:"drobo"`
	// PollSeconds is the frequency at which updates are sent
	PollSeconds int `json:"poll_seconds"`
	// Broker is host:port for the MQTT broker
	Broker      string `json:"broker"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	TopicPrefix string `json:"topic_prefix"`
	TopicNodeID string `json:"topic_node_id"`
}

var configPath string

const defaultConfigPath = "/usr/local/etc/drobomqtt.conf"

func init() {
	flag.StringVar(&configPath, "c", defaultConfigPath, "(required) Specifiy config file path")
}

func LoadConfig() *Config {

	var config Config

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Couldn't open config file: %v\n", err)
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	if err := json.Unmarshal([]byte(bytes), &config); err != nil {
		log.Fatalf("Couldn't open config file: %v\n", err)
	}

	if config.TopicPrefix == "" {
		config.TopicPrefix = "homeassistant"
	}

	if config.TopicNodeID == "" {
		config.TopicNodeID = "drobo"
	}

	if config.PollSeconds == 0 {
		config.PollSeconds = 5
	}

	return &config
}
