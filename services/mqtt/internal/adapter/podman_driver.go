package adapter

import (
	"context"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/specgen"
	"gofr.dev/pkg/gofr"
)

type PodmanAdapter struct {
	ctx context.Context
}

func NewPodmanAdapter() (*PodmanAdapter, error) {
	conn, err := bindings.NewConnection(context.Background(), "unix:///run/user/1000/podman/podman.sock")
	if err != nil {
		return nil, err
	}
	return &PodmanAdapter{
		ctx: conn,
	}, nil
}

func (p *PodmanAdapter) RemoveImage(c *gofr.Context, imageId string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PodmanAdapter) RunContainer(c *gofr.Context, imageName string, name string) (string, error) {
	if p.ctx == nil {
		c.Logger.Error()
	}
	s := specgen.NewSpecGenerator(imageName, false)
	s.Name = name
	createResponse, err := containers.CreateWithSpec(p.ctx, s, nil)
	if err != nil {
		c.Logger.Error("container creation failed", err)
		return "", err
	}
	if err := containers.Start(p.ctx, createResponse.ID, nil); err != nil {
		c.Logger.Error("container start failed", err)
		return "", err
	}
	c.Logger.Debug("Container ID:", createResponse.ID)
	return createResponse.ID, nil
}

func (p *PodmanAdapter) StopContainer(c *gofr.Context, containerId string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PodmanAdapter) RemoveContainer(c *gofr.Context, containerId string) error {
	//TODO implement me
	panic("implement me")
}

func (p *PodmanAdapter) ShowContainerLogs(c *gofr.Context, containerId string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PodmanAdapter) ShowContainerStatus(c *gofr.Context, containerId string) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

var _ Adapter = (*PodmanAdapter)(nil)
