package main

import (
	"gofr.dev/pkg/gofr"
	"io"
)

// ENV
//
// APP_NAME=mqtt_alarm
//
// DB_HOST=10.0.32.144
// DB_USER=<USER>
// DB_PASSWORD=<PASSWORD>
// DB_NAME=mqtt
// DB_PORT=5432
// DB_DIALECT=postgres
// DB_TABLE_NAME=<NAME>
//
// METRICS_PORT=2124
//
// REDIS_HOST=10.0.32.144
// REDIS_PORT=6379

var tableName string

type request struct {
	MqttDeviceId  string
	YwThresholdId float64
	MqttDataId    float64
}

func main() {
	app := gofr.New()
	// ENVs
	tableName = app.Config.GetOrDefault("DB_TABLE_NAME", "test_data")
	// register alarm post service
	app.AddHTTPService("alarm_receiver", "http://10.0.0.30/")

	app.AddCronJob("* * * * *", "check_alarm", checkAlarm)
	app.Run()
}

func checkAlarm(c *gofr.Context) {
	// get alarm weight values

	// get last 1-minute telemetries data
	rows, err := c.SQL.QueryContext(c, "SELECT id,create_at FROM "+tableName+" WHERE time > NOW() - INTERVAL '%s';")
	if err != nil {
		c.Logger.Error("get construction runner by construction id failed", err)
		return
	}
	for rows.Next() {
		var re usecase.ConstructionRunnerEntity
		if err := rows.Scan(&re.ID, &re.Name, &re.ConstructionID, &re.ContainerID); err != nil {
			return
		}
	}
	return
}

func sendAlarm(c *gofr.Context) (interface{}, error) {
	alarmSvc := c.GetHTTPService("alarm_receiver")
	res, err := alarmSvc.Post(c, "dev-api/zhjc/yjlist", nil, nil)
	if err != nil {
		c.Logger.Error(err, "send alarm to receiver error")
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.Logger.Error(err, "post to receiver services body close error")
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return string(body), nil
}
