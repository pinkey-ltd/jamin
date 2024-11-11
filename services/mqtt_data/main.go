package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"gofr.dev/pkg/gofr"
	"strings"
	"time"
)

// ENV
//
// APP_NAME=mqtt_data
//
// DB_HOST=10.0.32.144
// DB_USER=<USER>
// DB_PASSWORD=<PASSWORD>
// DB_NAME=mqtt
// DB_PORT=5432
// DB_DIALECT=postgres
// DB_TABLE_NAME=<NAME>
//
// METRICS_PORT=2123
//
// REDIS_HOST=10.0.32.144
// REDIS_PORT=6379

var appName string
var initialled bool = false
var tableName string

type Schema struct {
	Time time.Time       `json:"time"`
	Data json.RawMessage `json:"data"`
}

func main() {
	app := gofr.New()

	// ENVs
	appName = app.Config.GetOrDefault("APP_NAME", "test_data")
	tableName = app.Config.GetOrDefault("DB_TABLE_NAME", "mqtt_telemetries")
	// initial db pre 30 seconds
	app.AddCronJob("*/30 * * * * *", "init", initDB)

	app.AddCronJob("* * * * *", "transform_data", transformData)

	app.Run()
}

func initDB(c *gofr.Context) {
	if initialled {
		return
	}
	c.Logger.Info("start data conversion")
	// table exits?
	exists := false
	query := "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1);"
	err := c.SQL.QueryRowContext(c, query, appName).Scan(&exists)
	if exists {
		initialled = true
		return
	}

	if err != nil {
		c.Logger.Error(err, "query table "+appName+" exists error")
		// go on
	}
	// Get Configure from redis
	jsonStr, err := c.Redis.Get(c, appName).Result()
	if err != nil {
		c.Logger.Error(err, "can't get app configure from redis")
		panic(err)
	}
	keys, _, err := transformJSON(c, jsonStr)
	keyPairs := make([]string, 0)
	for _, k := range keys {
		kp := "\"" + k + "\" text"
		kp = strings.Replace(kp, ".", "_", -1)
		keyPairs = append(keyPairs, kp)
	}
	sql := `CREATE TABLE ` + appName + ` (
id bigserial PRIMARY KEY NOT NULL,
create_at timestamptz NOT NULL,
` + strings.Join(keyPairs, ",") + `
);`
	_, err = c.SQL.Exec(sql)
	if err != nil {
		c.Logger.Error(err, "create table "+appName+" error")
	} else {
		initialled = true
	}
	return
}

func transformData(c *gofr.Context) {
	if !initialled {
		return
	}
	c.Logger.Info("transform string")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	// get raw data
	var lastTimeStamp time.Time
	lastTimestampRow := c.SQL.QueryRow("SELECT create_at FROM " + appName + " ORDER BY create_at DESC LIMIT 1;")
	err := lastTimestampRow.Scan(&lastTimeStamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			lastTimeStamp = time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)
		} else {
			c.Logger.Error(err, "query last recode from "+appName+" error")
			return
		}
	}
	var res []Schema
	sql := "SELECT * FROM " + tableName + " WHERE time > '" + lastTimeStamp.In(loc).Format("2006-01-02 15:04:05") + "';"
	rows, err := c.SQL.QueryContext(c, sql)
	if err != nil {
		c.Logger.Error(err, "get raw data failed")
		return
	}
	for rows.Next() {
		var r Schema
		if err := rows.Scan(&r.Time, &r.Data); err != nil {
			c.Logger.Error(err, "")
			return
		}
		res = append(res, r)
	}
	if len(res) == 0 {
		c.Logger.Debug("no more received data")
		return
	}
	c.Logger.Debug("hit:", string(res[0].Data))
	// trans
	insertSQLHeader := `INSERT INTO ` + appName + `("create_at",`
	for i, r := range res {
		if i == 0 {
			schema, v, err := transformJSON(c, "{\"data\":"+string(r.Data)+"}")
			if err != nil {
				c.Logger.Error(err, "insert data error ")
				return
			}
			schemaQuotation := make([]string, 0)
			for _, s := range schema {
				schemaQuotation = append(schemaQuotation, "\""+s+"\"")
			}
			insertSQLHeader = insertSQLHeader + strings.Join(schemaQuotation, ",") + ") VALUES ("
			insertSQLHeader = strings.Replace(insertSQLHeader, ".", "_", -1)
			localTime := r.Time.In(loc)
			insertSQLValue := "'" + localTime.Format("2006-01-02 15:04:05") + "'," + strings.Join(v, ",") + `);`
			insertSQL := insertSQLHeader + insertSQLValue
			_, err = c.SQL.ExecContext(c, insertSQL)
			if err != nil {
				c.Logger.Error(err, "insert data error ")
				return
			}
		} else {
			_, v, err := transformJSON(c, "{\"data\":"+string(r.Data)+"}")
			if err != nil {
				c.Logger.Error(err, "insert data error ")
				return
			}
			localTime := r.Time.In(loc)
			insertSQL := insertSQLHeader + "'" + localTime.Format("2006-01-02 15:04:05") + "'," + strings.Join(v, ",") + `);`
			_, err = c.SQL.ExecContext(c, insertSQL)
			if err != nil {
				c.Logger.Error(err, "insert data error and continue")
				continue
			}
		}
	}
	return
}

func transformJSON(c *gofr.Context, jsonStr string) (keys []string, values []string, err error) {
	result := gjson.Get(jsonStr, "data")
	if result.String() == "" {
		c.Logger.Error(err, "can't get data in json")
		return nil, nil, fmt.Errorf("json format error")
	}
	result.ForEach(func(key, value gjson.Result) bool {
		if !value.IsObject() && !value.IsArray() {
			c.Logger.Debug("the key is 1 ", key.String(), " and value is ", value.String())
			keys = append(keys, key.String())
			values = append(values, value.String())
		}

		// 如果值是数组或对象，可以进一步处理
		if value.IsObject() {
			value.ForEach(func(nestedKey, nestedValue gjson.Result) bool {
				if !nestedValue.IsObject() && !nestedValue.IsArray() {
					keys = append(keys, key.String()+"."+nestedKey.String())
					values = append(values, nestedValue.String())
					c.Logger.Debug("the key is 2 ", key.String(), " and value is ", nestedValue.String())
				}
				return true // 返回 true 继续遍历
			})
		}
		if value.IsArray() {
			value.ForEach(func(nestedKey, nestedValue gjson.Result) bool {
				if !nestedValue.IsObject() && !nestedValue.IsArray() {
					keys = append(keys, key.String()+"."+nestedKey.String())
					values = append(values, nestedValue.String())
					c.Logger.Debug("the key is 3 ", key.String()+"."+nestedKey.String(), " and value is ", nestedValue.String())
				}
				nestedValue.ForEach(func(innerKey, innerValue gjson.Result) bool {
					if !innerValue.IsObject() && !innerValue.IsArray() {
						keys = append(keys, key.String()+"."+nestedKey.String()+"."+innerKey.String())
						values = append(values, innerValue.String())
						c.Logger.Debug("the key is 4 ", key.String()+"."+nestedKey.String()+"."+innerKey.String(), " and value is ", innerValue.String())
					}
					return true
				})
				return true // 返回 true 继续遍历
			})
		}
		return true // 返回 true 继续遍历
	})
	return
}
