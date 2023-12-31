package database

import (
	"fmt"
	"go-web-scrapper/entity/model"

	"gorm.io/gorm"
)

type Postgres struct {
	Db *gorm.DB
}

type IPostgres interface {
	Insert(data []*model.Data) error
}

func NewPostgres(db *gorm.DB) IPostgres {
	return &Postgres{Db: db}
}

func (p *Postgres) Insert(data []*model.Data) error {
	err := p.Db.Create(&data).Error
	if err != nil {
		return fmt.Errorf("entity.database.postgres.Insert: %v", err)
	}
	return nil
}
