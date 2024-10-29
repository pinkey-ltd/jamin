package repo

import (
	"github.com/pinkey-ltd/jamin/services/mqtt/internal/usecase"
	"gofr.dev/pkg/gofr"
)

type ConstructionPGRepository struct{}

func (c *ConstructionPGRepository) GetAll(ctx *gofr.Context) (*[]usecase.ConstructionEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConstructionPGRepository) Create(ctx *gofr.Context, data *usecase.ConstructionEntity) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConstructionPGRepository) Update(ctx *gofr.Context, data *usecase.ConstructionEntity) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConstructionPGRepository) Delete(ctx *gofr.Context, id int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConstructionPGRepository) GetById(ctx *gofr.Context, id int64) (*usecase.ConstructionEntity, error) {
	var res usecase.ConstructionEntity
	row := ctx.SQL.QueryRowContext(ctx, "SELECT id,name,status,data_adapter,data_destination,data_destination,args,schema,schema_adapter FROM mqtt_construction WHERE id = $1", id)
	err := row.Scan(&res.ID, &res.Name, &res.Status, &res.DataAdapter, &res.DataResource, &res.DataDestination, &res.Args, &res.Schema, &res.SchemaAdapter)
	if err != nil {
		ctx.Logger.Error("test handle", err)
		return nil, err
	}
	// Get Runners
	var rs []*usecase.ConstructionRunnerEntity
	rows, err := ctx.SQL.QueryContext(ctx, "SELECT id,name,construction_id,container_id,comments FROM mqtt_construction_running WHERE construction_id = $1", id)
	if err != nil {
		ctx.Logger.Error("get construction runner by construction id failed", err)
		return nil, err
	}
	for rows.Next() {
		var re usecase.ConstructionRunnerEntity
		if err := rows.Scan(&re.ID, &re.Name, &re.ConstructionID, &re.ContainerID, &re.Comments); err != nil {
			return nil, err
		}
		rs = append(rs, &re)
	}
	res.Runners = rs
	return &res, nil
}

var _ ConstructionRepo = (*ConstructionPGRepository)(nil)
