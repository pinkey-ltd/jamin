package migrations

import "gofr.dev/pkg/gofr/migration"

const alterTableDevicePkSQL = `ALTER TABLE mqtt_device ADD jcd_id bigint NULL;
ALTER TABLE mqtt_device 
	ADD CONSTRAINT mqtt_device_yw_jcd_fk FOREIGN KEY (id) 
	REFERENCES yw_jcd(id) 
	ON UPDATE SET NULL;
`

func alterTableDevicePk() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(alterTableDevicePkSQL)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
