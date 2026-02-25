package repositories

import (
	"context"
	"log"
	"slices"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/jbl1108/goTimeSeriesStorage/usecases/datamodel"
)

type InfluxRepository struct {
	client       influxdb2.Client
	org          string
	knownBuckets []string
}

func NewInfluxRepository(url, token, org string) *InfluxRepository {
	client := influxdb2.NewClient(url, token)
	if client == nil {
		panic("Could not create InfluxDB client")
	}

	defer client.Close()
	ctx := context.Background()
	bucketsAPI := client.BucketsAPI()
	buckets, err := bucketsAPI.GetBuckets(ctx)
	if err != nil {
		log.Printf("Error getting buckets: %v", err)
	}
	var bucketNames []string
	for _, b := range *buckets {
		bucketNames = append(bucketNames, b.Name)
	}
	log.Printf("Known buckets: %v", bucketNames)

	return &InfluxRepository{
		client:       client,
		org:          org,
		knownBuckets: bucketNames,
	}
}

func (r *InfluxRepository) StoreTimeSeries(ts datamodel.BucketTimeSeries) error {
	ctx := context.Background()
	err := r.createBucketIfNeeded(ctx, ts)
	if err != nil {
		return err
	}
	writeAPI := r.client.WriteAPIBlocking(r.org, ts.Bucket)
	p := write.NewPointWithMeasurement(ts.Measurement)
	for key, value := range ts.Tags {
		p = p.AddTag(key, value)
	}
	for key, value := range ts.Fields {
		p = p.AddField(key, value)
	}
	if ts.Timestamp.IsZero() == false {
		p = p.SetTime(ts.Timestamp)
	}
	return writeAPI.WritePoint(ctx, p)
}
func (r *InfluxRepository) Close() {
	r.client.Close()
}

func (r *InfluxRepository) createBucketIfNeeded(ctx context.Context, ts datamodel.BucketTimeSeries) error {
	if !slices.Contains(r.knownBuckets, ts.Bucket) {
		orgsAPI := r.client.OrganizationsAPI()
		org, err := orgsAPI.FindOrganizationByName(context.Background(), r.org)
		if err != nil {
			log.Fatalf("Failed to find organization: %v", err)
		}
		bucketsAPI := r.client.BucketsAPI()
		_, err = bucketsAPI.CreateBucketWithName(ctx, org, ts.Bucket)
		if err != nil {
			log.Printf("Error creating bucket %s: %v", ts.Bucket, err)
			return err
		}
		r.knownBuckets = append(r.knownBuckets, ts.Bucket)
	}
	return nil
}
