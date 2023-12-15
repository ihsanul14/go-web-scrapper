package app

import (
	"go-web-scrapper/entity"
	databaseEntity "go-web-scrapper/entity/database"
	"go-web-scrapper/framework/database"
	"go-web-scrapper/usecase"
	"log"

	"go-web-scrapper/framework/logger"

	"github.com/subosito/gotenv"
)

func Run() {
	err := gotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	baseLogger := logger.InitLogger()
	dbConn, err := database.Connect()
	if err != nil {
		baseLogger.Logger.Fatal(err)
	}

	postgresEntity := databaseEntity.NewPostgres(dbConn)
	entity := entity.NewEntity(postgresEntity)
	usecase := usecase.NewUsecase(entity, baseLogger)
	usecase.Get()
	// time.Sleep(5 * time.Minute)
}
