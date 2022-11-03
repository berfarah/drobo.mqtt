package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	droboStatusTopic   = "drobo_status"
	droboCapacityTopic = "drobo_capacity"
	droboUsedTopic     = "drobo_used"

	gigabyteInBytes = 1_000_000_000
)

func NewPublisher(c *Config) *Publisher {
	return &Publisher{
		config:    c,
		connected: make(chan bool, 1),
	}
}

type Publisher struct {
	config    *Config
	client    mqtt.Client
	connected chan bool
}

func (p *Publisher) OnConnect(client mqtt.Client) {
	p.client = client
	log.Println("Connected!")
	p.discoverySetup()
	close(p.connected)
}

func (p *Publisher) discoverySetup() error {
	drobo, err := getDroboInfo(p.config.Drobo)
	if err != nil {
		return err
	}

	p.publishSensorDiscovery(drobo, map[string]interface{}{
		"name":      "Status",
		"object_id": droboStatusTopic,
		"icon":      "mdi:nas",
	})

	p.publishSensorDiscovery(drobo, map[string]interface{}{
		"name":                "Capacity",
		"object_id":           droboCapacityTopic,
		"icon":                "mdi:gauge",
		"unit_of_measurement": "GB",
	})

	p.publishSensorDiscovery(drobo, map[string]interface{}{
		"name":                "% Used",
		"object_id":           droboUsedTopic,
		"icon":                "mdi:percent",
		"unit_of_measurement": "%",
	})

	for i := range drobo.Slots.Nodes {
		p.publishSensorDiscovery(drobo, map[string]interface{}{
			"name":      fmt.Sprintf("Disk %d Status", i+1),
			"object_id": diskStatusTopic(i + 1),
			"icon":      "mdi:harddisk",
		})
	}

	return nil
}

func (p *Publisher) publishSensorDiscovery(drobo Drobo, config map[string]interface{}) error {
	config["device"] = map[string]interface{}{
		"manufacturer": "Drobo",
		"model":        drobo.Model,
		"name":         drobo.Name,
		"sw_version":   drobo.Version,
		"connections":  [][]string{{"host", p.config.Drobo}},
		"identifiers":  drobo.Serial,
	}
	objectID := config["object_id"].(string)
	config["unique_id"] = objectID + "_" + drobo.Serial
	config["state_topic"] = p.buildTopic(objectID)
	config["entity_category"] = "diagnostic"

	b, err := json.Marshal(config)
	if err != nil {
		return err
	}

	p.client.Publish(p.buildTopic(objectID+"/config"), 0, true, b)
	return nil
}

func (p *Publisher) publishSensorData(drobo Drobo) {
	var tokens []mqtt.Token
	tokens = append(tokens, p.client.Publish(p.buildTopic(droboStatusTopic), 0, true, strings.Join(drobo.Statuses(), ", ")))
	tokens = append(tokens, p.client.Publish(p.buildTopic(droboCapacityTopic), 0, true, strconv.Itoa(drobo.TotalCapacityProtected/gigabyteInBytes)))
	tokens = append(tokens, p.client.Publish(p.buildTopic(droboUsedTopic), 0, true, strconv.FormatFloat(
		float64(drobo.UsedCapacityProtected*100)/float64(drobo.TotalCapacityProtected),
		'f', 2, 32,
	)))

	for i, slot := range drobo.Slots.Nodes {
		tokens = append(tokens, p.client.Publish(p.buildTopic(diskStatusTopic(i+1)), 0, true, slot.StatusString()))
	}

	for _, token := range tokens {
		if token.Wait(); token.Error() != nil {
			log.Println("Error publishing data:", token.Error())
		}
	}
}

func (p *Publisher) Run(ctx context.Context) {
	log.Println("Connecting...")
	<-p.connected

	ticker := time.NewTicker(time.Duration(p.config.PollSeconds) * time.Second)

	log.Println("Starting polling cycle")
	for {
		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			drobo, err := getDroboInfo(p.config.Drobo)
			if err != nil {
				log.Println("Error fetching drobo data:", err)
			}
			p.publishSensorData(drobo)
		}
	}

}

func (p *Publisher) buildTopic(topic string) string {
	return strings.Join([]string{
		p.config.TopicPrefix,
		"sensor",
		p.config.TopicNodeID,
		topic,
	}, "/")
}

func diskStatusTopic(i int) string {
	return fmt.Sprintf("disk%d_status", i)
}
