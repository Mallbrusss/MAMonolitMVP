package services

import (
	"encoding/json"
	"fmt"
	"mamonolitmvp/config"
	"mamonolitmvp/internal/math/price_analysis"
	"mamonolitmvp/internal/models"
	"mamonolitmvp/internal/repository"
	"mamonolitmvp/pkg/http_client"
)

const ChunkSize = 524288

type TinkoffService struct {
	Client *http_client.HTTPClient
	Config *config.Config
	is     *InstrumentService
	pa     *price_analysis.PriceAnalysis
}

func NewTinkoffService(cfg *config.Config, repo *repository.InstrumentRepository) *TinkoffService {
	return &TinkoffService{
		Client: http_client.NewHTTPClient(),
		Config: cfg,
		is:     NewInstrumentService(repo),
		pa:     price_analysis.NewPriceAnalysis(),
	}
}

func (s *TinkoffService) GetClosePrices(instruments []string) ([]models.ClosePrice, error) {
	var InstrumentRequests []models.InstrumentRequest
	for _, instrument := range instruments {
		InstrumentRequests = append(InstrumentRequests, models.InstrumentRequest{InstrumentID: instrument})
	}

	reqBody := models.GetClosePricesRequest{Instruments: InstrumentRequests}
	url := fmt.Sprintf("%s/tinkoff.public.invest.api.contract.v1.MarketDataService/GetClosePrices", s.Config.APIBaseURL)

	headers := map[string]string{
		"Authorization": "Bearer " + s.Config.APIToken,
		"Content-Type":  "application/json",
	}

	respBody, err := s.Client.Post(url, headers, reqBody)
	if err != nil {
		return nil, err
	}

	var response models.ClosePricesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return response.ClosePrices, nil
}

func (s *TinkoffService) GetAllInstruments(instrumentStatus string) ([]models.PlacementPrice, error) {
	reqBody := models.BondsRequest{InstrumentStatus: instrumentStatus}
	url := fmt.Sprintf("%s/tinkoff.public.invest.api.contract.v1.InstrumentsService/Shares", s.Config.APIBaseURL)

	headers := map[string]string{
		"Authorization": "Bearer " + s.Config.APIToken,
		"Content-Type":  "application/json",
	}

	respBody, err := s.Client.Post(url, headers, reqBody)
	if err != nil {
		return nil, err
	}

	var response models.GetBondsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	err = s.is.CreateInstruments(response.Instruments)
	if err != nil {
		return nil, err
	}

	return response.Instruments, nil
}

func (s *TinkoffService) GetCandles(instrumentInfo map[string]any) ([]models.HistoricCandle, error) {
	reqBody := models.GetCandlesRequest{
		Figi:         instrumentInfo["figi"].(string),
		From:         instrumentInfo["from"].(string),
		To:           instrumentInfo["to"].(string),
		Interval:     instrumentInfo["interval"].(string),
		InstrumentId: instrumentInfo["instrumentId"].(string),
	}

	url := fmt.Sprintf("%s/tinkoff.public.invest.api.contract.v1.MarketDataService/GetCandles", s.Config.APIBaseURL)

	headers := map[string]string{
		"Authorization": "Bearer " + s.Config.APIToken,
		"Content-Type":  "application/json",
	}

	respBody, err := s.Client.Post(url, headers, reqBody)
	if err != nil {
		return nil, err
	}

	response, _, err := s.fixeRespBody(respBody, reqBody.InstrumentId)
	if err != nil {
		return nil, err
	}

	err = s.is.CreateCandles(response.Candles)
	if err != nil {
		return nil, err
	}

	return response.Candles, nil
}

func (s *TinkoffService) fixeRespBody(respBody []byte, instrumentID string) (models.GetCandlesResponse, []byte, error) {
	var responce models.GetCandlesResponse
	err := json.Unmarshal(respBody, &responce)
	if err != nil {
		return models.GetCandlesResponse{}, nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	for i := range responce.Candles {
		responce.Candles[i].InstrumentId = instrumentID
	}

	data, err := json.Marshal(responce)
	if err != nil {
		return models.GetCandlesResponse{}, nil, err
	}

	return responce, data, nil
}
