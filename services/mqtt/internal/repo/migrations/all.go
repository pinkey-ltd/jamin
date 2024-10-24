package migrations

import "gofr.dev/pkg/gofr/migration"

func All() map[int64]migration.Migrate {
	return map[int64]migration.Migrate{
		20240226153000: createTableTelemetries(),
		20240227090000: createTableConstruction(),
		20240227140000: createTableDeviceTemplate(),
		20240227143000: createTableDevice(),
		20240227153000: alterTableDevicePk(),
	}
}
