package main

import (
	"context"
	"flag"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	flag.Parse()
	config := LoadConfig()
	pub := NewPublisher(config)

	opts := mqtt.NewClientOptions()
	if config.Username != "" {
		opts.SetUsername(config.Username)
		opts.SetPassword(config.Password)
	}

	opts.SetAutoReconnect(true)
	opts.AddBroker(config.Broker)
	opts.SetConnectionLostHandler(func(_ mqtt.Client, e error) {
		log.Println("Lost connection:", e)
	})
	opts.SetClientID("docker-mqtt")
	opts.SetOnConnectHandler(pub.OnConnect)
	client := mqtt.NewClient(opts)

	token := client.Connect()
	if token.Wait(); token.Error() != nil {
		log.Fatalf("Couldn't connect to MQTT broker: %v\n", token.Error())
	}

	pub.Run(context.Background())

	client.Disconnect(200)
}
