package main

import (
	"fmt"
	"log"
	"stream-app/pkg/api"
	"stream-app/pkg/config"
	"stream-app/pkg/logger"
	"stream-app/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	appLogger, err := logger.NewLogger(cfg.Log.Level)
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}
	defer appLogger.SugaredLogger.Sync()

	dsn := "host=" + cfg.Database.Host +
		" user=" + cfg.Database.User +
		" password=" + cfg.Database.Password +
		" dbname=" + cfg.Database.Name +
		" port=" + fmt.Sprintf("%d", cfg.Database.Port) +
		" sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		appLogger.SugaredLogger.Fatalf("cannot connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		appLogger.SugaredLogger.Fatalf("Ошибка миграции базы данных: %v", err)
	}

	server := api.NewServer(cfg, db, appLogger)
	if err := server.Run(); err != nil {
		appLogger.SugaredLogger.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
