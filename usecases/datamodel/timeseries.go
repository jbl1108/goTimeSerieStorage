package datamodel

import "time"

type TimeSeries struct {
	Measurement string            `json:"measurement"`
	Tags        map[string]string `json:"tags"`
	Fields      map[string]any    `json:"fields"`
	Timestamp   time.Time         `json:"timestamp"`
}

type BucketTimeSeries struct {
	Bucket string `json:"bucket"`
	TimeSeries
}

// Message a struct
type Message struct {
	Topic string `json:"topic"`
	Data  any    `json:"data"`
}
