package repository

import (
	"fmt"
	"log"
	"mamonolitmvp/internal/models"

	"gorm.io/gorm"
)

type InstrumentRepository struct {
	db *gorm.DB
}

func NewInstrumentRepository(db *gorm.DB) *InstrumentRepository {
	return &InstrumentRepository{
		db: db,
	}
}

func (ir *InstrumentRepository) CreateInstruments(instruments []models.PlacementPrice) error {
	batchSize := 100

	for i := 0; i < len(instruments); i += batchSize {
		end := i + batchSize
		if end > len(instruments) {
			end = len(instruments)
		}

		batch := make([]*models.PlacementPrice, end-i)
		for j := 0; j < len(batch); j++ {
			batch[j] = &instruments[i+j]
		}

		tx := ir.db.Begin()
		if tx.Error != nil {
			log.Printf("failed to start transaction: %v", tx.Error)
			return tx.Error
		}

		batchNumber := (i / batchSize) + 1
		totalBatches := (len(instruments) + batchSize - 1) / batchSize

		log.Printf("Create batch %d of %d", batchNumber, totalBatches)

		err := tx.Create(&batch).Error
		if err != nil {
			tx.Rollback()
			log.Printf("failed to insert batch: %v", err)
			return err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			log.Printf("failed to commit transaction: %v", err)
			return err
		}
	}

	log.Println("Instument create success")
	return nil
}

func (ir *InstrumentRepository) CreateCandles(candles []models.HistoricCandle) error {
	batchSize := 100

	for i := 0; i < len(candles); i += batchSize {
		end := i + batchSize
		if end > len(candles) {
			end = len(candles)
		}

		batch := make([]*models.HistoricCandle, end-i)
		for j := 0; j < len(batch); j++ {
			batch[j] = &candles[i+j]
		}

		tx := ir.db.Begin()
		if tx.Error != nil {
			log.Printf("failed to start transaction: %v", tx.Error)
			return tx.Error
		}

		batchNumber := (i / batchSize) + 1
		totalBatches := (len(candles) + batchSize - 1) / batchSize

		log.Printf("Create batch %d of %d", batchNumber, totalBatches)

		err := tx.Create(&batch).Error
		if err != nil {
			tx.Rollback()
			log.Printf("failed to insert batch: %v", err)
			return err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			log.Printf("failed to commit transaction: %v", err)
			return err
		}
	}

	log.Println("Candles create success")
	return nil
}

func (ir *InstrumentRepository) GetInstrumentUIDAndFigi(ticker string) (models.Ids, error) {
	var ids models.Ids
	err := ir.db.Model(&models.PlacementPrice{}).Select("uid", "figi").Where("ticker=?", ticker).Scan(&ids).Error
	if err != nil || ids.Uid == "" && ids.Figi == "" {
		err = gorm.ErrRecordNotFound
		log.Printf("failed to Get InstrumentID: %v", err)
		return ids, err
	}
	return ids, nil
}

func (ir *InstrumentRepository) GetCandles(instrumentUID string) ([]models.HistoricCandle, error) {
	var candles []models.HistoricCandle
	err := ir.db.Model(&models.HistoricCandle{}).Where("instrument_id=?", instrumentUID).Find(&candles).Error
	fmt.Println("err2: ", err)
	if err != nil || len(candles) == 0 {
		err = gorm.ErrRecordNotFound
		log.Printf("no candles found for instrument UID: %s", instrumentUID)
		log.Printf("failed to Get Candles: %v", err)
		return nil, err
	}
	return candles, nil
}

func (ir *InstrumentRepository) GetTicker(instrumentUID string) (string, error) {
	var name string
	err := ir.db.Model(&models.PlacementPrice{}).Select("ticker").Where("uid=?", instrumentUID).Scan(&name).Error
	if err != nil {
		log.Printf("failed to Get Name: %v", err)
		return "", err
	}
	return name, nil
}
