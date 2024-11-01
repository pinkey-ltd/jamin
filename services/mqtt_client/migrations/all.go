package migrations

import "gofr.dev/pkg/gofr/migration"

func All(tableName string) map[int64]migration.Migrate {
	return map[int64]migration.Migrate{
		20240226153000: createTableTelemetries(tableName),
	}
}
