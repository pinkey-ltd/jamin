package migrations

import "gofr.dev/pkg/gofr/migration"

const createTableTelemetriesSQL = `CREATE TABLE $1 (
  time TIMESTAMPTZ NOT NULL,
  data JSONB NOT NULL
);`

func createTableTelemetries(name string) migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createTableTelemetriesSQL, name)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
