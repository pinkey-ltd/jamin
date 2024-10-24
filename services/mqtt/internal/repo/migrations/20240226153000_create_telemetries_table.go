package migrations

import "gofr.dev/pkg/gofr/migration"

const createTableTelemetriesSQL = `CREATE TABLE mqtt_telemetries (
  time TIMESTAMPTZ NOT NULL,
  data JSONB NOT NULL
);`

func createTableTelemetries() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createTableTelemetriesSQL)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
