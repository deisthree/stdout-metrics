package influx

import (
	"fmt"
	"os"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

//Query returns the result and an error
func Query(influxClient client.Client, queryString string) ([]client.Result, error) {
	query := client.NewQuery(queryString, "kubernetes", "s")
	response, err := influxClient.Query(query)
	if err != nil {
		return nil, err
	}

	if response.Error() != nil {
		return nil, response.Error()
	}
	return response.Results, nil
}

//WriteInt64 takes a int64 value and sends it to influx
func WriteInt64(influxClient client.Client, value int64, name string, tags map[string]string, timestamp time.Time) error {
	fields := map[string]interface{}{"value": value}
	return Write(influxClient, name, tags, fields, timestamp)
}

//WriteFloat64 takes a float64 value and sends it to influx
func WriteFloat64(influxClient client.Client, value float64, name string, tags map[string]string, timestamp time.Time) error {
	fields := map[string]interface{}{"value": value}
	return Write(influxClient, name, tags, fields, timestamp)
}

//Write sends data to influx
func Write(influxClient client.Client, name string, tags map[string]string, fields map[string]interface{}, timestamp time.Time) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: "kubernetes",
	})
	if err != nil {
		return err
	}

	pt, err := client.NewPoint(name, tags, fields, timestamp)
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	err = influxClient.Write(bp)
	if err != nil {
		return err
	}
	return nil
}

//Connect to influx
func Connect() (client.Client, error) {
	influxClient, influxError := client.NewHTTPClient(client.HTTPConfig{
		Timeout: time.Second * 5,
		Addr:    fmt.Sprintf("http://%s", os.Getenv("DEIS_MONITOR_INFLUXAPI_SERVICE_HOST")),
	})
	if influxError != nil {
		return nil, influxError
	}
	return influxClient, nil
}
