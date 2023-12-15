package entity

import (
	"context"
	"go-web-scrapper/entity/database"
	"go-web-scrapper/entity/model"
)

type Entity struct {
	Postgres database.IPostgres
}

type IEntity interface {
	Insert(ctx context.Context, data []*model.Data) error
}

func NewEntity(postgres database.IPostgres) IEntity {
	return &Entity{Postgres: postgres}
}

func (e *Entity) Insert(ctx context.Context, data []*model.Data) error {
	return e.Postgres.Insert(ctx, data)
}
