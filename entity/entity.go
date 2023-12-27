package entity

import (
	"go-web-scrapper/entity/database"
	"go-web-scrapper/entity/model"
)

type Entity struct {
	Postgres database.IPostgres
}

type IEntity interface {
	Insert(data []*model.Data) error
}

func NewEntity(postgres database.IPostgres) IEntity {
	return &Entity{Postgres: postgres}
}

func (e *Entity) Insert(data []*model.Data) error {
	return e.Postgres.Insert(data)
}
