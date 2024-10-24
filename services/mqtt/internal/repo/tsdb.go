package repo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gofr.dev/pkg/gofr"
	"time"
)

// Timescale Repository

const schemaName = "telemetries"

type SchemaTest struct {
	Deformation []struct {
		Absv float64 `json:"absv"`
		Id   int     `json:"id"`
		Ts   int64   `json:"ts"`
		X    float64 `json:"x"`
		Y    float64 `json:"y"`
		Z    float64 `json:"z"`
	} `json:"deformation"`
	DeviceId    string `json:"deviceId"`
	Environment struct {
		Humidity  float64 `json:"humidity"`
		InnerTemp float64 `json:"innerTemp"`
		Temp      float64 `json:"temp"`
	} `json:"environment"`
	RealDis int `json:"realDis"`
	Slant   struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"slant"`
	Ts int64 `json:"ts"`
}

func (i *SchemaTest) Value() (driver.Value, error) {
	return json.Marshal(i)
}

func (i *SchemaTest) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}
	return json.Unmarshal(data, &i)
}

type Schema struct {
	Time time.Time   `json:"time"`
	Data *SchemaTest `json:"data"`
}

type TelemetryTSRepository struct {
	C *gofr.Context
}

func (t *TelemetryTSRepository) GetTelemetriesByTime(timeEnum TimeEnum) (*[]Schema, error) {
	var res []Schema
	sql := fmt.Sprintf("SELECT * FROM %s WHERE time > NOW() - INTERVAL '%s';", schemaName, timeEnum)
	rows, err := t.C.SQL.QueryContext(t.C, sql)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r Schema
		if err := rows.Scan(&r.Time, &r.Data); err != nil {
			return nil, err
		}

		res = append(res, r)
	}
	return &res, nil
}

func (t *TelemetryTSRepository) CreateTelemetry(data Schema) (bool, error) {
	_, err := t.C.SQL.ExecContext(t.C, "INSERT INTO "+schemaName+" (time,data) VALUES (NOW(),?)", data)
	if err != nil {
		return false, err
	}
	return true, nil
}

var _ TelemetryRepo[Schema] = (*TelemetryTSRepository)(nil)
