package analyzer

import (
	"github.com/labstack/echo/v4"
	"mamonolitmvp/internal/math/price_analysis"
	"mamonolitmvp/internal/models"
	"net/http"
)

type StockExchange interface {
	GetTotalSignal(instrumentInfo map[string]any) (string, price_analysis.Signal, error)
}

type Signal struct {
	Service StockExchange
}

func NewSignalHandler(service StockExchange) *Signal {
	return &Signal{
		Service: service,
	}
}

func (h *Signal) GetSignals(c echo.Context) error {
	c.Request().Header.Set("Content-Type", "application/json")
	var req models.GetCandlesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request format",
		})
	}

	instrumentInfo := map[string]any{
		"figi":         req.Figi,
		"from":         req.From,
		"to":           req.To,
		"interval":     req.Interval,
		"instrumentId": req.InstrumentId,
	}

	ticker, signal, err := h.Service.GetTotalSignal(instrumentInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch all candles",
			"err":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"ticker": ticker,
		"Hurst":  signal.Hurst,
		"MDFA": map[string]any{
			"LogFq": signal.LogFq,
			"Hq":    signal.Hq,
			"LogS":  signal.LogS,
		},
		"MFSpectrum": map[string]any{
			"Qsorted": signal.Qsorted,
			"Tau":     signal.Tau,
			"Alpha":   signal.Alpha,
			"FAlpha":  signal.FAlpha,
		},
	})
}
