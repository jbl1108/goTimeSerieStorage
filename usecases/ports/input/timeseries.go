package input

import "github.com/jbl1108/goTimeSeriesStorage/usecases/datamodel"

type TimeSeriesInputPort interface {
	HandleTimeSeries(message datamodel.Message) error
}
