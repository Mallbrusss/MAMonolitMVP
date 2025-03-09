package timescale

import (
	"fmt"
	"log"
	"mamonolitmvp/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const maxRetries = 10
const retryInterval = 5 * time.Second

func InitDB(host, user, password, dBname, port string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dBname, port)

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		log.Printf("Ошибка подключения к базе данных: %v. Попытка %d из %d", err, i+1, maxRetries)
		time.Sleep(retryInterval)
	}

	migrate(db)
	return db, err
}

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.HistoricCandle{}, &models.High{}, &models.Low{}, &models.Close{}, &models.Open{})
	if err != nil {
		log.Println("error migrate Candle table")
	}

	err = db.AutoMigrate(&models.PlacementPrice{})
	if err != nil {
		log.Println("error migrate placementPrice table")
	}

	err = db.AutoMigrate(&models.Nominal{})
	if err != nil {
		log.Println("error migrate nominal table")
	}

	err = db.AutoMigrate(&models.MinPriceIncrement{})
	if err != nil {
		log.Println("error migrate minPriceIncrement table")
	}

	log.Println("Success connect to Postgres")
}
