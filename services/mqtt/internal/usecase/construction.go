package usecase

type ConstructionStatus string

const (
	ConstructionStatusDraft      ConstructionStatus = "draft"
	ConstructionStatusTestFailed ConstructionStatus = "test_failed"
	ConstructionStatusTested     ConstructionStatus = "tested"
	ConstructionStatusRunning    ConstructionStatus = "running"
	ConstructionStatusStopped    ConstructionStatus = "stopped"
	ConstructionStatusRemoved    ConstructionStatus = "removed"
)

type ConstructionData string

const (
	ConstructionDataInner ConstructionData = "inner"
	ConstructionDataOuter ConstructionData = "outer"
)

type ConstructionEntity struct {
	ID              int64
	Name            string
	Status          ConstructionStatus
	Adapter         string
	DataResource    ConstructionData
	DataDestination ConstructionData
	Args            string
	Schema          string
}

func (c *ConstructionEntity) RestPath() string {
	return "construction"
}

func (c *ConstructionEntity) TableName() string {
	return "mqtt_construction"
}
