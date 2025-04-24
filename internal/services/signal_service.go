package services

import (
	"fmt"
	"mamonolitmvp/internal/math/price_analysis"
	"mamonolitmvp/internal/models"
	"strconv"
)

func (s *TinkoffService) GetTotalSignal(instrumentInfo map[string]any) (string, price_analysis.Signal, []price_analysis.Fdi, []float64, error) {
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

	respBody, _ := s.Client.Post(url, headers, reqBody)

	response, _, _ := s.fixeRespBody(respBody, reqBody.InstrumentId)

	prices := make([]float64, 0, len(response.Candles))

	for _, v := range response.Candles {
		locUnit := v.Close.Units
		locNano := v.Close.Nano

		locPrice := fmt.Sprintf("%s.%d", locUnit, locNano)
		priceValue, err := strconv.ParseFloat(locPrice, 64)
		if err != nil {
			return "", price_analysis.Signal{}, nil, nil, err
		}

		prices = append(prices, priceValue)
	}

	sig, err := s.pa.TotalSignal(prices)
	if err != nil {
		return "", price_analysis.Signal{}, nil, nil, err
	}

	fdi, hurstSeries, err := s.pa.SlidingWindowAnalysis(prices, 100)
	if err != nil {
		return "", price_analysis.Signal{}, nil, nil, err
	}

	ticker, err := s.is.instrumentRepository.GetTicker(instrumentInfo["instrumentId"].(string))
	if err != nil {
		return "", price_analysis.Signal{}, nil, nil, err
	}

	return ticker, sig, fdi, hurstSeries, err
}
