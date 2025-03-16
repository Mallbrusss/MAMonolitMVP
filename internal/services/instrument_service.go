package services

import (
	"mamonolitmvp/internal/models"
)

type InstrumentRepository interface {
	CreateInstruments(instruments []models.PlacementPrice) error
	GetTicker(instrumentUID string) (string, error)
	//CreateCandles(candles []models.HistoricCandle) error
}

type InstrumentService struct {
	instrumentRepository InstrumentRepository
}

func NewInstrumentService(instrumentRepository InstrumentRepository) *InstrumentService {
	return &InstrumentService{
		instrumentRepository: instrumentRepository,
	}
}

func (s *InstrumentService) CreateInstruments(instruments []models.PlacementPrice) error {
	return s.instrumentRepository.CreateInstruments(instruments)
}

//func (s *InstrumentService) CreateCandles(candles []models.HistoricCandle) error {
//	return s.instrumentRepository.CreateCandles(candles)
//}
