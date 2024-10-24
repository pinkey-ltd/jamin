package migrations

import "gofr.dev/pkg/gofr/migration"

const createTableDeviceSQL = `CREATE TABLE mqtt_device (
 	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(255) UNIQUE NOT NULL,
	status VARCHAR(255) NOT NULL,
	health CHAR(1) NOT NULL DEFAULT '1',
	lat DOUBLE PRECISION NOT NULL,
	lng DOUBLE PRECISION NOT NULL,
	template_id BIGINT NOT NULL,
	parent_id BIGINT,
	CONSTRAINT fk_template_id 
		FOREIGN KEY(template_id) 
		REFERENCES mqtt_device_template(id) 
		ON DELETE CASCADE,
	CONSTRAINT fk_parent_id
		FOREIGN KEY (parent_id)
		REFERENCES mqtt_device(id)
		ON DELETE CASCADE
);`

func createTableDevice() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createTableDeviceSQL)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
