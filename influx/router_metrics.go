package influx

import (
	"time"

	"github.com/deis/stdout-metrics/util"
	"github.com/influxdata/influxdb/client/v2"
)

//WriteRouterMetrics converts json to influx datatype and sends
func WriteRouterMetrics(influxClient client.Client, json map[string]interface{}) error {
	parsedMessage, err := util.ParseNginxLog(json["log"].(string))
	if err != nil {
		return err
	}
	err = writeData(influxClient, util.GetHost(json), parsedMessage)
	if err != nil {
		return err
	}
	return nil
}

func writeData(influxClient client.Client, host string, parsedMessage map[string]interface{}) error {
	tags := map[string]string{"host": host,
		"app":         parsedMessage["app"].(string),
		"status_code": parsedMessage["status_code"].(string),
	}

	// We cant use the nginx timestamp because its not precise enough
	//parsedMessage["time"].(time.Time)
	err := WriteFloat64(influxClient, parsedMessage["request_time"].(float64), "deis_router_response_time_seconds", tags, time.Now())
	if err != nil {
		return err
	}
	return nil
}
