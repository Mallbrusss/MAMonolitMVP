package services

import (
	"mamonolitmvp/internal/models"
)

type InstrumentRepository interface {
	CreateInstruments(instruments []models.PlacementPrice) error
}

type InstrumentService struct {
	instrumentRepository InstrumentRepository
}

func NewInstrumentService(instrumentRepository InstrumentRepository) *InstrumentService {
	return &InstrumentService{
		instrumentRepository: instrumentRepository,
	}
}

func (is *InstrumentService) CreateInstruments(instruments []models.PlacementPrice) error {
	err := is.instrumentRepository.CreateInstruments(instruments)
	if err != nil {
		return err
	}
	return nil
}

func (is *InstrumentService) CreateCandles(candles []models.HistoricCandle) error {
	err := is.CreateCandles(candles)
	if err != nil {
		return err
	}
	return nil
}
