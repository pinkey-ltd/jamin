package adapter

import (
	"gofr.dev/pkg/gofr"
	"os/exec"
)

func ImagesHandler(c *gofr.Context) (interface{}, error) {
	cmd := exec.Command("docker", "images")
	res, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	c.Logger.Info(res)
	return res, nil
}
