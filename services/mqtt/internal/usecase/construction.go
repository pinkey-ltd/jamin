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
	DataAdapter     string
	DataResource    ConstructionData
	DataDestination ConstructionData
	Args            string
	Schema          string
	SchemaAdapter   string
	Runners         []*ConstructionRunnerEntity
}

type ConstructionCURDEntity struct {
	ID              int64
	Name            string
	Status          ConstructionStatus
	DataAdapter     string
	DataResource    ConstructionData
	DataDestination ConstructionData
	Args            string
	Schema          string
	SchemaAdapter   string
}

func (c *ConstructionCURDEntity) RestPath() string {
	return "construction"
}

func (c *ConstructionCURDEntity) TableName() string {
	return "mqtt_construction"
}

type ConstructionRunnerEntity struct {
	ID             int64
	Name           string
	ContainerID    string
	ConstructionID int64
	Status         string
	Comments       string
}
