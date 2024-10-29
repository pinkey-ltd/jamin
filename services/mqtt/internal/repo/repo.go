package repo

import (
	"github.com/pinkey-ltd/jamin/services/mqtt/internal/usecase"
	"gofr.dev/pkg/gofr"
)

type TimeEnum string

const (
	OneMinute      TimeEnum = "1 minutes"
	FifteenMinutes TimeEnum = "15 minutes"
	TwelveHours    TimeEnum = "12 hours"
	OneDay         TimeEnum = "1 day"
)

type TelemetryRepo[T any] interface {
	GetTelemetriesByTime(timeEnum TimeEnum) (*[]T, error)
	CreateTelemetry(data Schema) (bool, error)
}

type ConstructionRepo interface {
	GetAll(ctx *gofr.Context) (*[]usecase.ConstructionEntity, error)
	GetById(ctx *gofr.Context, id int64) (*usecase.ConstructionEntity, error)
	Create(ctx *gofr.Context, data *usecase.ConstructionEntity) (bool, error)
	Update(ctx *gofr.Context, data *usecase.ConstructionEntity) (bool, error)
	Delete(ctx *gofr.Context, id int64) (bool, error)
}
