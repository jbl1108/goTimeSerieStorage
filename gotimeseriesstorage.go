package main

import (
	"github.com/jbl1108/goTimeSeriesStorage/config"
)

func main() {
	config := config.NewApplication()
	config.Delivery.Connect()
	defer config.Delivery.Disconnect()

	// Keep the application running to listen for MQTT messages
	select {}
}
