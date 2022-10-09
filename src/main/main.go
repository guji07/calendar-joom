package main

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"

	"cryptoColony/src/config"
	"cryptoColony/src/controller"
	"cryptoColony/src/service"
	"cryptoColony/src/storage"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func main() {
	appConfig := getAppConfig()
	dbConfig := appConfig.Database
	connectionUrl := fmt.Sprintf("postgresql://%s:%d/%s?user=%s&password=%s&sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.User, dbConfig.Password)
	upMigration(connectionUrl)

	startServer(appConfig, connectionUrl)
}

func upMigration(connectionUrl string) {
	m, err := migrate.New(
		"file://db/migrations",
		connectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func getAppConfig() *config.AppConfig {
	cfgPath, err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	appConfig, err := config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	return appConfig
}

func startServer(appConfig *config.AppConfig, connectionUrl string) {
	srv := gin.Default()
	err := srv.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal(err)
	}

	logger := logrus.New()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", connectionUrl)
	goquDb := goqu.New("postgres", db)
	repository := storage.NewRepository(goquDb)

	validatorStruct := validator.New()

	calendarController := controller.NewCalendarController(service.NewUserService(&repository),
		service.NewEventService(&repository), logger, validatorStruct)
	srv.POST("/users", calendarController.CreateUser)
	srv.GET("/users/:user_id/events/", calendarController.GetUserEvents)
	srv.POST("/events", calendarController.CreateEvent)
	srv.GET("/events/:id", calendarController.GetEvent)
	srv.GET("/events/:id/respond", calendarController.RespondOnEvent)
	srv.GET("/events/window_by_users", calendarController.FindWindowForEvent)
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
