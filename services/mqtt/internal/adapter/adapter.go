package adapter

import "gofr.dev/pkg/gofr"

type Adapter interface {
	RemoveImage(c *gofr.Context, imageId string) (interface{}, error)
	RunContainer(c *gofr.Context, imageName string, name string) (string, error)
	StopContainer(c *gofr.Context, containerId string) (bool, error)
	RemoveContainer(c *gofr.Context, containerId string) error
	ShowContainerLogs(c *gofr.Context, containerId string) (interface{}, error)
	ShowContainerStatus(c *gofr.Context, containerId string) (interface{}, error)
	StartContainer(c *gofr.Context, containerId string) (bool, error)
}
