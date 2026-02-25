package delivery

import (
	"encoding/json"
	"log"

	"github.com/jbl1108/goTimeSeriesStorage/usecases/datamodel"
	"github.com/jbl1108/goTimeSeriesStorage/usecases/ports/input"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client       mqtt.Client
	topic        string
	inputUsecase input.TimeSeriesInputPort
}

func NewMQTTClient(broker string, username string, password string, topic string, inputUsecase input.TimeSeriesInputPort) *MQTTClient {
	opts := mqtt.NewClientOptions().AddBroker(broker)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(defaultMessageHandler)
	client := mqtt.NewClient(opts)
	return &MQTTClient{client: client, topic: topic, inputUsecase: inputUsecase}
}

func defaultMessageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message on topic: %s, payload: %s", msg.Topic(), string(msg.Payload()))
}

func (m *MQTTClient) Connect() {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	log.Printf("Connecting to topic: %v", m.topic)
	m.client.Subscribe(m.topic, 1, m.messageHandler)
}

func (m *MQTTClient) Disconnect() {
	m.client.Disconnect(250)
}

func (m *MQTTClient) messageHandler(client mqtt.Client, msg mqtt.Message) {
	var message datamodel.Message
	if err := json.Unmarshal(msg.Payload(), &message); err != nil {
		log.Printf("failed to parse message: %v", err)
		return
	}
	log.Printf("Received message: %+v", message)
	message.Topic = msg.Topic() // Override topic with the actual MQTT topic
	err := m.inputUsecase.HandleTimeSeries(message)
	if err != nil {
		log.Printf("Error handling time series: %v", err)
	}
}
