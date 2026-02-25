package usecases

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	datamodel "github.com/jbl1108/goTimeSeriesStorage/usecases/datamodel"
)

var msg datamodel.Message = datamodel.Message{
	Topic: "timeseries/metrics",
	Data: datamodel.TimeSeries{
		Measurement: "cpu",
		Tags: map[string]string{
			"host":   "server1",
			"region": "us-west",
		},
		Fields: map[string]any{
			"usage": 0.75,
		},
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	},
}

var expectedTimeSeries = datamodel.TimeSeries{
	Measurement: "cpu",
	Tags: map[string]string{
		"host":   "server1",
		"region": "us-west",
	},
	Fields: map[string]any{
		"usage": 0.75,
	},
	Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
}

type mockOutput struct {
	received datamodel.BucketTimeSeries
	err      error
	called   bool
}

func (m *mockOutput) StoreTimeSeries(data datamodel.BucketTimeSeries) error {
	m.received = data
	m.called = true
	return m.err
}

func TestHandleTimeSeries_Success(t *testing.T) {
	mock := &mockOutput{}
	uc := NewStoreTimeSeriesUseCase(mock)

	jsonString, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal time series: %v", err)
	}

	var myMsg datamodel.Message
	err = json.Unmarshal(jsonString, &myMsg)
	if err != nil {
		t.Fatalf("failed to unmarshal time series into message data: %v", err)
	}

	if err := uc.HandleTimeSeries(myMsg); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !mock.called {
		t.Fatalf("expected output port to be called")
	}
	if mock.received.Bucket != "metrics" {
		t.Fatalf("expected bucket 'metrics', got %q", mock.received.Bucket)
	}

	if !reflect.DeepEqual(mock.received.TimeSeries, expectedTimeSeries) {
		t.Fatalf("time series mismatch: expected %+v, got %+v", expectedTimeSeries, mock.received.TimeSeries)
	}
}

func TestHandleTimeSeries_InvalidTopicParts(t *testing.T) {
	mock := &mockOutput{}
	uc := NewStoreTimeSeriesUseCase(mock)

	msg := datamodel.Message{
		Topic: "timeseries", // too few parts
		Data: datamodel.TimeSeries{
			Measurement: "m",
		},
	}

	if err := uc.HandleTimeSeries(msg); err == nil {
		t.Fatalf("expected error for invalid topic parts, got nil")
	}
}

func TestHandleTimeSeries_InvalidPrefix(t *testing.T) {
	mock := &mockOutput{}
	uc := NewStoreTimeSeriesUseCase(mock)

	msg := datamodel.Message{
		Topic: "notimeseries/bucket",
		Data: datamodel.TimeSeries{
			Measurement: "m",
		},
	}

	if err := uc.HandleTimeSeries(msg); err == nil {
		t.Fatalf("expected error for invalid prefix, got nil")
	}
}

// -- additional tests below --

func TestHandleTimeSeries_OutputError(t *testing.T) {
	wantErr := errors.New("storage failure")
	mock := &mockOutput{err: wantErr}
	uc := NewStoreTimeSeriesUseCase(mock)

	jsonString, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal time series: %v", err)
	}

	var myMsg datamodel.Message
	if err := json.Unmarshal(jsonString, &myMsg); err != nil {
		t.Fatalf("failed to unmarshal time series into message data: %v", err)
	}

	err = uc.HandleTimeSeries(myMsg)
	if err == nil || err.Error() != wantErr.Error() {
		t.Fatalf("expected storage error, got %v", err)
	}
}

func TestHandleTimeSeries_NilData(t *testing.T) {
	mock := &mockOutput{}
	uc := NewStoreTimeSeriesUseCase(mock)

	msg := datamodel.Message{Topic: "timeseries/bucket", Data: nil}
	if err := uc.HandleTimeSeries(msg); err == nil {
		t.Fatalf("expected error for nil data, got nil")
	}
}

func TestHandleTimeSeries_DataNotMap(t *testing.T) {
	mock := &mockOutput{}
	uc := NewStoreTimeSeriesUseCase(mock)

	msg := datamodel.Message{Topic: "timeseries/bucket", Data: "not a map"}
	if err := uc.HandleTimeSeries(msg); err == nil {
		t.Fatalf("expected error for non-map data, got nil")
	}
}

func TestHandleTimeSeries_BadTimestamp(t *testing.T) {
	mock := &mockOutput{}
	uc := NewStoreTimeSeriesUseCase(mock)

	bad := map[string]any{
		"measurement": "cpu",
		"timestamp":   "not-a-time",
	}
	msg := datamodel.Message{Topic: "timeseries/metrics", Data: bad}

	before := time.Now()
	if err := uc.HandleTimeSeries(msg); err != nil {
		t.Fatalf("expected no error with bad timestamp, got %v", err)
	}
	if !mock.called {
		t.Fatalf("expected output port to be called")
	}
	if mock.received.TimeSeries.Timestamp.Before(before) {
		t.Fatalf("expected parsed timestamp to be >= %v, got %v", before, mock.received.TimeSeries.Timestamp)
	}
}
