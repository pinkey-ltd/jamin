package repo

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
