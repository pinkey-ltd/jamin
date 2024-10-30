package migrations

import "gofr.dev/pkg/gofr/migration"

const createTableConstructionRunningSQL = `CREATE TABLE mqtt_construction_running (
 	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(255) UNIQUE NOT NULL,
	container_id VARCHAR(64) NOT NULL,
	construction_id BIGINT NOT NULL,
	status VARCHAR(20) NOT NULL,
	comments TEXT,
	CONSTRAINT fk_construction_id 
		FOREIGN KEY(construction_id) 
		REFERENCES mqtt_construction(id) 
		ON DELETE CASCADE
);`

func createTableConstructionRunning() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createTableConstructionRunningSQL)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
