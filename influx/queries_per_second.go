package influx

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

var (
	queriesPerSecondAggregate = `SELECT count(value) FROM deis_router_response_time_seconds WHERE time > now() - 1m GROUP BY time(1m)`
	queriesPerSecondPerApp    = `SELECT count(value) FROM deis_router_response_time_seconds WHERE time > now() - 1m AND app = '%s' GROUP BY time(1m)`
	appQuery                  = `SHOW TAG VALUES FROM deis_router_response_time_seconds WITH KEY = app`
)

// This code isnt being used right now since we can perform the calculations in grafana. However, we are keeping this around just in case we need it in a future release.

//QueriesPerSecond finds the qps value for the aggregate and individual endpoints
func QueriesPerSecond(influxClient client.Client) error {
	err := writeAggregateQueriesPerSecond(influxClient)
	if err != nil {
		return err
	}
	err = writePerAppQueriesPerSecond(influxClient)
	if err != nil {
		return err
	}
	return nil
}

func writeAggregateQueriesPerSecond(influxClient client.Client) error {
	results, err := Query(influxClient, queriesPerSecondAggregate)
	if err != nil {
		return err
	}
	tags := map[string]string{"app": "aggregate"}
	err = writeQueriesPerSecond(influxClient, results, tags)
	if err != nil {
		return err
	}
	return nil
}

func writePerAppQueriesPerSecond(influxClient client.Client) error {
	results, err := Query(influxClient, appQuery)
	if err != nil {
		return err
	}

	if len(results[0].Series) > 0 && len(results[0].Series[0].Values) > 0 {
		for _, value := range results[0].Series[0].Values {
			app := value[1].(string)
			appResults, err := Query(influxClient, fmt.Sprintf(queriesPerSecondPerApp, app))
			if err != nil {
				return err
			}
			tags := map[string]string{"app": app}
			err = writeQueriesPerSecond(influxClient, appResults, tags)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func writeQueriesPerSecond(influxClient client.Client, results []client.Result, tags map[string]string) error {
	if len(results[0].Series) > 0 {
		v, _ := results[0].Series[0].Values[1][1].(json.Number).Float64()
		err := WriteFloat64(influxClient, v/60, "deis_router_request_count_seconds", tags, time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}
