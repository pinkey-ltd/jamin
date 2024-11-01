package adapter

import (
	"context"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/specgen"
	"gofr.dev/pkg/gofr"
	"os/user"
)

type PodmanAdapter struct {
	ctx context.Context
}

func NewPodmanAdapter() (*PodmanAdapter, error) {
	// get uid
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	conn, err := bindings.NewConnection(context.Background(), "unix:///run/user/"+currentUser.Uid+"/podman/podman.sock")
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

func (p *PodmanAdapter) StartContainer(c *gofr.Context, containerId string) (bool, error) {
	if err := containers.Start(p.ctx, containerId, nil); err != nil {
		c.Logger.Error("container start failed", err)
		return false, err
	}
	return true, nil
}

func (p *PodmanAdapter) StopContainer(c *gofr.Context, containerId string) (bool, error) {
	if err := containers.Stop(p.ctx, containerId, nil); err != nil {
		c.Logger.Error("container stop failed", err)
		return false, err
	}
	return true, nil
}

func (p *PodmanAdapter) RemoveContainer(c *gofr.Context, containerId string) error {
	options := &containers.RemoveOptions{}
	options.WithForce(true)
	options.WithVolumes(true)

	_, err := containers.Remove(p.ctx, containerId, options)
	if err != nil {
		c.Logger.Error("container removal failed", err)
		return err
	}
	return nil
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
