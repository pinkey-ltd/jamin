package migrations

import "gofr.dev/pkg/gofr/migration"

const createTableConstructionSQL = `CREATE TABLE mqtt_construction (
 	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(255) UNIQUE NOT NULL,
	status VARCHAR(10) NOT NULL,
	adapter VARCHAR(255) NOT NULL,
	data_resource VARCHAR(10) NOT NULL,
	data_destination VARCHAR(10) NOT NULL,
	args TEXT NOT NULL,
	schema TEXT NOT NULL
);`

func createTableConstruction() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createTableConstructionSQL)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
