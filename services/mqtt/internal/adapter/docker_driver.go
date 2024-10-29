package adapter

import (
	"errors"
	"gofr.dev/pkg/gofr"
	"os/exec"
	"strings"
)

// QueryDocker

func ImagesHandler(c *gofr.Context) (interface{}, error) {
	query := c.Param("query")

	return queryDocker(c, []string{"images"}, query)
}

func ContainersHandler(c *gofr.Context) (interface{}, error) {
	query := c.Param("query")

	return queryDocker(c, []string{"ps", "-a"}, query)
}

func queryDocker(c *gofr.Context, args []string, query string) (interface{}, error) {
	cmd := exec.Command("podman", args...)
	raw, err := cmd.CombinedOutput()
	if err != nil {
		c.Logger.Error(err, "failed to exec docker ", args)
		return nil, err
	}
	c.Logger.Debug("run docker ", args, " result ", string(raw))
	// Formatter
	resStr := strings.Split(string(raw), "\n")
	res := make([]string, 0)
	for _, img := range resStr[1:] {
		if img == "" {
			continue
		}
		if strings.Contains(img, query) {
			res = append(res, img)
		}
	}
	return res, nil
}

// DockerAdapter provides an interface for interacting with Docker or Podman to manage containers and images.
type DockerAdapter struct{}

func (d *DockerAdapter) RemoveImage(c *gofr.Context, imageId string) (interface{}, error) {
	cmd := exec.Command("podman", "rmi", imageId)
	err := cmd.Run()
	if err != nil {
		c.Logger.Error(err, "failed to exec docker remove image ", imageId)
		return nil, err
	}
	return nil, nil
}

func (d *DockerAdapter) RunContainer(c *gofr.Context, imageName string, name string) (string, error) {
	if imageName == "" {
		return "", errors.New("image name is empty")
	}
	args := []string{"run", "--restart=always", "-d"}
	if name != "" {
		args = append(args, "--name "+name)
	}
	args = append(args, imageName)
	cmd := exec.Command("podman", args...)
	r, err := cmd.CombinedOutput()
	if err != nil {
		c.Logger.Error(err, "failed to exec docker run container ", imageName)
		return "", err
	}
	// Format
	res := string(r)[:12]
	return res, nil
}

func (d *DockerAdapter) StopContainer(c *gofr.Context, containerId string) (interface{}, error) {
	cmd := exec.Command("podman", "stop", containerId)
	err := cmd.Run()
	if err != nil {
		c.Logger.Error(err, "failed to exec docker stop container ", containerId)
		return nil, err
	}
	return nil, nil
}

func (d *DockerAdapter) RemoveContainer(c *gofr.Context, containerId string) error {
	cmd := exec.Command("podman", "rm", "-f", containerId)
	err := cmd.Run()
	if err != nil {
		c.Logger.Error(err, "failed to exec docker remove container ", containerId)
		return err
	}
	return nil
}

func (d *DockerAdapter) ShowContainerLogs(c *gofr.Context, containerId string) (interface{}, error) {
	if containerId == "" {
		return nil, nil
	}

	c.Logger.Debug("show container logs ", containerId)

	cmd := exec.Command("podman", "logs", containerId)
	raw, err := cmd.CombinedOutput()
	if err != nil {
		c.Logger.Error(err, "failed to get docker logs container ", containerId)
	}
	// Formatter
	res := strings.Split(string(raw), "\n")

	return res, nil
}

func (d *DockerAdapter) ShowContainerStatus(c *gofr.Context, containerId string) (interface{}, error) {
	if containerId == "" {
		return nil, nil
	}

	c.Logger.Debug("show container status ", containerId)

	cmd := exec.Command("podman", "ps", "-a")

	r, err := cmd.CombinedOutput()
	if err != nil {
		c.Logger.Error(err, "failed to exec docker status container ", containerId)
		return nil, err
	}
	// Formatter
	resStr := strings.Split(string(r), "\n")
	res := ""
	for _, c := range resStr[1:] {
		if c == "" {
			continue
		}
		if strings.Contains(c, containerId) {
			res = c
			break
		}
	}
	return res, nil
}
