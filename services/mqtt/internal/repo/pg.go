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
	sqlStr := "UPDATE mqtt_construction SET name = $1,status =$2,data_adapter = $3 ,data_resource =$4,data_destination=$5 ,args=$6,schema =$7 , schema_adapter=$8 WHERE id = $9;"
	subSqlStr := "UPDATE mqtt_construction_running SET name = $1,container_id = $2, construction_id = $3 WHERE id = $4;"
	tx, err := ctx.SQL.Begin()
	if err != nil {
		ctx.Logger.Error("update construction failed", err)
		return false, err
	}

	_, err = tx.ExecContext(ctx, sqlStr, data.Name, data.Status, data.DataAdapter, data.DataResource, data.DataDestination, data.Args, data.Schema, data.SchemaAdapter, data.ID)
	if err != nil {
		tx.Rollback()
		ctx.Logger.Error("update construction failed", err)
		return false, err
	}

	for _, r := range data.Runners {
		_, err = tx.ExecContext(ctx, subSqlStr, r.Name, r.ContainerID, data.ID, r.ID)
		if err != nil {
			tx.Rollback()
			ctx.Logger.Error("update construction runner failed", err)
			return false, err
		}
	}

	if err = tx.Commit(); err != nil {
		ctx.Logger.Error("update construction failed", err)
		return false, err
	}
	return true, nil
}

func (c *ConstructionPGRepository) UpdateOrInsert(ctx *gofr.Context, data *usecase.ConstructionEntity) (bool, error) {
	sqlStr := "UPDATE mqtt_construction SET name = $1,status =$2,data_adapter = $3 ,data_resource =$4,data_destination=$5 ,args=$6,schema =$7 , schema_adapter=$8 WHERE id = $9;"
	subSqlStr := "INSERT INTO mqtt_construction_running (name,container_id,construction_id) VALUES ($1,$2,$3);"
	tx, err := ctx.SQL.Begin()
	if err != nil {
		ctx.Logger.Error("update construction failed", err)
		return false, err
	}

	_, err = tx.ExecContext(ctx, sqlStr, data.Name, data.Status, data.DataAdapter, data.DataResource, data.DataDestination, data.Args, data.Schema, data.SchemaAdapter, data.ID)
	if err != nil {
		tx.Rollback()
		ctx.Logger.Error("update construction failed", err)
		return false, err
	}

	for _, r := range data.Runners {
		_, err = tx.ExecContext(ctx, subSqlStr, r.Name, r.ContainerID, r.ConstructionID)
		ctx.Logger.Debug("update construction runner", *r)
		if err != nil {
			tx.Rollback()
			ctx.Logger.Error("update construction runner failed", err)
			return false, err
		}
	}

	if err = tx.Commit(); err != nil {
		ctx.Logger.Error("update construction failed", err)
		return false, err
	}
	return true, nil
}

func (c *ConstructionPGRepository) UpdateOrDelete(ctx *gofr.Context, data *usecase.ConstructionEntity) (bool, error) {
	sqlStr := "UPDATE mqtt_construction SET name = $1,status =$2,data_adapter = $3 ,data_resource =$4,data_destination=$5 ,args=$6,schema =$7 , schema_adapter=$8 WHERE id = $9;"
	deleteSqlStr := "DELETE FROM mqtt_construction_running WHERE id = $1;"
	tx, err := ctx.SQL.Begin()
	if err != nil {
		ctx.Logger.Error("update construction failed", err)
		return false, err
	}

	_, err = tx.ExecContext(ctx, sqlStr, data.Name, data.Status, data.DataAdapter, data.DataResource, data.DataDestination, data.Args, data.Schema, data.SchemaAdapter, data.ID)
	if err != nil {
		tx.Rollback()
		ctx.Logger.Error("update construction failed", err)
		return false, err
	}

	for _, r := range data.Runners {
		_, err = tx.ExecContext(ctx, deleteSqlStr, r.ID)
		ctx.Logger.Debug("delete construction runner", *r)
		if err != nil {
			tx.Rollback()
			ctx.Logger.Error("delete construction runner failed", err)
			return false, err
		}
	}

	if err = tx.Commit(); err != nil {
		ctx.Logger.Error("update construction failed", err)
		return false, err
	}
	return true, nil
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
		ctx.Logger.Error("get construction error", err)
		return nil, err
	}
	// Get Runners
	var rs []*usecase.ConstructionRunnerEntity
	rows, err := ctx.SQL.QueryContext(ctx, "SELECT id,name,construction_id,container_id FROM mqtt_construction_running WHERE construction_id = $1", id)
	if err != nil {
		ctx.Logger.Error("get construction runner by construction id failed", err)
		return nil, err
	}
	for rows.Next() {
		var re usecase.ConstructionRunnerEntity
		if err := rows.Scan(&re.ID, &re.Name, &re.ConstructionID, &re.ContainerID); err != nil {
			return nil, err
		}
		rs = append(rs, &re)
	}
	res.Runners = rs
	return &res, nil
}

var _ ConstructionRepo = (*ConstructionPGRepository)(nil)
