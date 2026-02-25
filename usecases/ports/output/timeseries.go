package output

import "github.com/jbl1108/goTimeSeriesStorage/usecases/datamodel"

type TimeSeriesOutputPort interface {
	StoreTimeSeries(data datamodel.BucketTimeSeries) error
}
