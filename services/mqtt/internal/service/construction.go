package service

import (
	"github.com/pinkey-ltd/jamin/services/mqtt/internal/adapter"
	"github.com/pinkey-ltd/jamin/services/mqtt/internal/usecase"
	"gofr.dev/pkg/gofr"
	"strconv"
)

type ConstructionGetter interface {
	GetById(ctx *gofr.Context, id int64) (*usecase.ConstructionEntity, error)
	Update(ctx *gofr.Context, data *usecase.ConstructionEntity) (bool, error)
	UpdateOrInsert(ctx *gofr.Context, data *usecase.ConstructionEntity) (bool, error)
}

type Construction struct {
	repo    ConstructionGetter //interface
	data    *usecase.ConstructionEntity
	adapter adapter.Adapter //interface
}

func NewConstruction(repo ConstructionGetter, data *usecase.ConstructionEntity, adapter adapter.Adapter) *Construction {
	return &Construction{
		repo:    repo,
		data:    data,
		adapter: adapter,
	}
}

// Common methods

func (c *Construction) getValueById(ctx *gofr.Context) error {
	idStr := ctx.PathParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.Logger.Error("incorrect ID format", err)
		return err
	}
	c.data, err = c.repo.GetById(ctx, id)
	return nil
}

// TestHandle retrieves a construction entity by ID for testing purposes.
// It parses the ID from the request context, queries the database, and logs the operation details.
// Returns the ConstructionEntity on success or an error if the operation fails.
func (c *Construction) TestHandle(ctx *gofr.Context) (interface{}, error) {
	ctx.Logger.Info("test handle")
	err := c.getValueById(ctx)
	if err != nil {
		return nil, err
	}
	// status
	if c.data.Status != usecase.ConstructionStatusDraft && c.data.Status != usecase.ConstructionStatusTestFailed {
		ctx.Logger.Error("invalid status, can't to test", err)
		return nil, err
	}

	// TODO test args & schema
	return c.data, nil
}

// StartHandle initiates the start process for a construction entity after validating its status.
// It fetches the entity by ID, checks if the status is either 'tested' or 'stopped',
// and logs an error if the status is invalid for starting.
func (c *Construction) StartHandle(ctx *gofr.Context) (interface{}, error) {
	ctx.Logger.Info("start handle")
	err := c.getValueById(ctx)
	if err != nil {
		return nil, err
	}

	if c.data.Status != usecase.ConstructionStatusTested && c.data.Status != usecase.ConstructionStatusStopped {
		ctx.Logger.Error("invalid status, can't to start", err)
		return nil, err
	}
	// step 1 get raw data
	if cid, err := c.adapter.RunContainer(ctx, c.data.DataAdapter, c.data.Name+"_data"); err != nil {
		ctx.Logger.Error("can't start "+c.data.DataAdapter, err)
		return nil, err
	} else {
		r := usecase.ConstructionRunnerEntity{
			Name:           c.data.Name + "_data",
			ContainerID:    cid,
			ConstructionID: c.data.ID,
		}
		c.data.Runners = append(c.data.Runners, &r)
	}
	// step 2  schema data
	if cid, err := c.adapter.RunContainer(ctx, c.data.SchemaAdapter, c.data.Name+"_schema"); err != nil {
		ctx.Logger.Error("can't start "+c.data.SchemaAdapter, err)
		if err = c.adapter.RemoveContainer(ctx, c.data.Runners[0].ContainerID); err != nil {
			ctx.Logger.Error("can't remove "+c.data.Runners[0].ContainerID, err)
			return nil, err
		}
		return nil, err
	} else {
		r := usecase.ConstructionRunnerEntity{
			Name:           c.data.Name + "_schema",
			ContainerID:    cid,
			ConstructionID: c.data.ID,
		}
		c.data.Runners = append(c.data.Runners, &r)
	}
	// step 3 alarm data
	if cid, err := c.adapter.RunContainer(ctx, c.data.SchemaAdapter, c.data.Name+"_alarm"); err != nil {
		ctx.Logger.Error("can't start "+c.data.AlarmAdapter, err)
		if err = c.adapter.RemoveContainer(ctx, c.data.Runners[0].ContainerID); err != nil {
			ctx.Logger.Error("can't remove "+c.data.Runners[0].ContainerID, err)
			return nil, err
		}
		if err = c.adapter.RemoveContainer(ctx, c.data.Runners[1].ContainerID); err != nil {
			ctx.Logger.Error("can't remove "+c.data.Runners[1].ContainerID, err)
			return nil, err
		}
		return nil, err
	} else {
		r := usecase.ConstructionRunnerEntity{
			Name:           c.data.Name + "_alarm",
			ContainerID:    cid,
			ConstructionID: c.data.ID,
		}
		c.data.Runners = append(c.data.Runners, &r)
	}

	c.data.Status = usecase.ConstructionStatusRunning
	_, err = c.repo.UpdateOrInsert(ctx, c.data)
	if err != nil {
		ctx.Logger.Error("can't update construction id="+c.data.Name, err)
		return nil, err
	}

	return nil, nil
}

func (c *Construction) StopHandle(ctx *gofr.Context) (interface{}, error) {
	ctx.Logger.Info("stop handle")
	err := c.getValueById(ctx)
	if err != nil {
		return nil, err
	}
	// TODO stop
	return nil, nil
}

func (c *Construction) RemoveHandle(ctx *gofr.Context) (interface{}, error) {
	ctx.Logger.Info("remove handle")
	err := c.getValueById(ctx)
	if err != nil {
		return nil, err
	}
	// TODO remove
	return nil, nil
}

func (c *Construction) StatusHandle(ctx *gofr.Context) (interface{}, error) {
	ctx.Logger.Info("status handle")
	err := c.getValueById(ctx)
	if err != nil {
		return nil, err
	}

	if c.data.Status != usecase.ConstructionStatusRunning && c.data.Status != usecase.ConstructionStatusStopped {
		return nil, nil
	}
	ss := make([]string, len(c.data.Runners))
	for i, r := range c.data.Runners {
		s, err := c.adapter.ShowContainerStatus(ctx, r.ContainerID)
		if err != nil {
			ctx.Logger.Error("can't get status", err)
			return nil, err
		}
		ss[i] = s.(string)
	}
	return ss, nil
}

func (c *Construction) LogsHandle(ctx *gofr.Context) (interface{}, error) {
	ctx.Logger.Info("logs handle")
	err := c.getValueById(ctx)
	if err != nil {
		return nil, err
	}
	if c.data.Status != usecase.ConstructionStatusRunning {
		return nil, nil
	}
	logs := make([][]string, len(c.data.Runners))
	for i, r := range c.data.Runners {
		log, err := c.adapter.ShowContainerLogs(ctx, r.ContainerID)
		if err != nil {
			ctx.Logger.Error("can't get logs", err)
			return nil, err
		}
		logs[i] = append(logs[i], log.([]string)...)
	}
	return logs, nil
}
