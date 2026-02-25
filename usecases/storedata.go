package usecases

import (
	"errors"
	"log"
	"strings"
	"time"

	datamodel "github.com/jbl1108/goTimeSeriesStorage/usecases/datamodel"
	"github.com/jbl1108/goTimeSeriesStorage/usecases/ports/output"
)

type StoreTimeSeriesUseCase struct {
	outputPort output.TimeSeriesOutputPort
}

func NewStoreTimeSeriesUseCase(port output.TimeSeriesOutputPort) *StoreTimeSeriesUseCase {
	return &StoreTimeSeriesUseCase{outputPort: port}
}

func (s *StoreTimeSeriesUseCase) HandleTimeSeries(message datamodel.Message) error {
	log.Printf("Handling time series message: %+v", message)
	if message.Data == nil {
		return errors.New("Message data is nil")
	}
	if _, ok := message.Data.(map[string]any); !ok {
		return errors.New("Message data is not a map[string]any")
	}
	bucket, err := s.getBucketName(message.Topic)
	if err != nil {
		return err
	}
	t := s.getTimeSeriesData(message)
	timeseries := datamodel.BucketTimeSeries{
		Bucket:     bucket,
		TimeSeries: *t,
	}
	return s.outputPort.StoreTimeSeries(timeseries)
}

func (s *StoreTimeSeriesUseCase) getTimeSeriesData(message datamodel.Message) *datamodel.TimeSeries {
	m := message.Data.(map[string]any)
	t := new(datamodel.TimeSeries)
	for k, v := range m {
		switch k {
		case "timestamp":
			t.Timestamp = s.getTimeOrNow(v)
		case "tags":
			t.Tags = s.getTag(v)
		case "fields":
			t.Fields = v.(map[string]any)
		case "measurement":
			t.Measurement = v.(string)
		default:
			log.Printf("Unknown field in message data: %s", k)
		}
	}
	return t
}

func (*StoreTimeSeriesUseCase) getBucketName(topic string) (string, error) {
	parts := strings.Split(topic, "/")
	if len(parts) < 2 {
		return "", errors.New("Invalid topic format, expected 'timeseries/{bucket}' got: " + topic)
	}
	if parts[0] != "timeseries" {
		return "", errors.New("Invalid topic format, expected 'timeseries/{bucket}' got: " + topic)
	}
	bucket := parts[1]
	return bucket, nil
}

func (*StoreTimeSeriesUseCase) getTimeOrNow(v any) time.Time {
	// start with current time, but only replace it if parsing succeeds
	timestamp := time.Now()
	if str, ok := v.(string); ok {
		if parsed, err := time.Parse(time.RFC3339, str); err == nil {
			timestamp = parsed
		} else {
			log.Printf("Error parsing timestamp: %v - using now()", err)
		}
	} else {
		log.Printf("timestamp value is not a string: %v", v)
	}
	return timestamp
}

func (*StoreTimeSeriesUseCase) getTag(v any) map[string]string {
	tags := make(map[string]string)
	for key, value := range v.(map[string]any) {
		if str, ok := value.(string); ok {
			tags[key] = str
		} else {
			log.Printf("Tag value for key '%s' is not a string: %v", key, value)
		}
	}
	return tags
}
