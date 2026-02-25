package config

import (
	"github.com/jbl1108/goTimeSeriesStorage/delivery"
	"github.com/jbl1108/goTimeSeriesStorage/repositories"
	"github.com/jbl1108/goTimeSeriesStorage/usecases"
	"github.com/jbl1108/goTimeSeriesStorage/usecases/ports/input"
	"github.com/jbl1108/goTimeSeriesStorage/usecases/ports/output"
)

type Application struct {
	usecases               input.TimeSeriesInputPort
	outputPort             output.TimeSeriesOutputPort
	Delivery               *delivery.MQTTClient
	storeTimeSeriesUseCase *usecases.StoreTimeSeriesUseCase
}

func NewApplication() Application {
	c := NewConfig()
	outputPort := repositories.NewInfluxRepository(c.InfluxDBURL(), c.InfluxDBToken(), c.InfluxDBOrg())
	usecases := usecases.NewStoreTimeSeriesUseCase(outputPort)
	mqttClient := delivery.NewMQTTClient(c.MQTTAddress(), c.MQTTUsername(), c.MQTTPassword(), "timeseries/#", usecases)

	return Application{
		usecases:               usecases,
		outputPort:             outputPort,
		storeTimeSeriesUseCase: usecases,
		Delivery:               mqttClient,
	}
}
