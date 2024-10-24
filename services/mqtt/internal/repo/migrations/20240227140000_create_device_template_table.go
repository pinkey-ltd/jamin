package migrations

import "gofr.dev/pkg/gofr/migration"

const createTableDeviceTemplateSQL = `CREATE TABLE mqtt_device_template (
 	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(255) UNIQUE NOT NULL,
	model VARCHAR(255) NOT NULL,
	maker VARCHAR(255) NOT NULL,
	types VARCHAR(20) NOT NULL,
	construction_id BIGINT NOT NULL,
	CONSTRAINT fk_construction_id 
		FOREIGN KEY(construction_id) 
		REFERENCES mqtt_construction(id) 
		ON DELETE CASCADE
);`

func createTableDeviceTemplate() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createTableDeviceTemplateSQL)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
