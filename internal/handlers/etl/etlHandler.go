package etl

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"mamonolitmvp/internal/models"
	"net/http"
)

type StockExchange interface {
	GetClosePrices(instruments []string) ([]models.ClosePrice, error)
	GetAllInstruments(instrumentStatus string) ([]models.PlacementPrice, error)
	GetCandles(instrumentInfo map[string]any) ([]models.HistoricCandle, error)
}

type ETLHandler struct {
	Service StockExchange
}

func NewETLHandler(service StockExchange) *ETLHandler {
	return &ETLHandler{
		Service: service,
	}
}

func (h *ETLHandler) GetClosePricesHandler(c echo.Context) error {
	c.Request().Header.Set("Content-Type", "application/json")

	var req models.GetClosePricesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request format",
		})
	}

	instruments := make([]string, 0, len(req.Instruments))
	for _, instrument := range req.Instruments {
		instruments = append(instruments, instrument.InstrumentID)
	}

	closePrices, err := h.Service.GetClosePrices(instruments)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch close prices",
		})
	}

	response := make(map[string]models.ClosePrice, len(closePrices))
	for _, v := range closePrices {
		response[v.InstrumentUid] = v
	}

	return c.JSON(http.StatusOK, response)
}

func (h *ETLHandler) GetAllBonds(c echo.Context) error {
	c.Request().Header.Set("Content-Type", "application/json")

	allBonds, err := h.Service.GetAllInstruments("INSTRUMENT_STATUS_BASE")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch all bonds",
		})
	}

	return c.JSON(http.StatusOK, allBonds)
}

func (h *ETLHandler) GetCandles(c echo.Context) error {
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

	candles, err := h.Service.GetCandles(instrumentInfo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Failed to fetch all candles",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"candles": candles,
	})
}

type Repository interface {
	GetInstrumentUIDAndFigi(ticker string) (models.Ids, error)
	GetCandles(instrumentUID string) ([]models.HistoricCandle, error)
}

type DBHandler struct {
	Repository Repository
}

func NewDBHandler(repository Repository) *DBHandler {
	return &DBHandler{
		Repository: repository,
	}
}

func (h *DBHandler) GetInstrumentUIDAndFigi(c echo.Context) error {
	c.Request().Header.Set("Content-Type", "application/json")
	ticker := c.QueryParam("ticker")

	instrumentIDS, err := h.Repository.GetInstrumentUIDAndFigi(ticker)
	if err != nil {
		log.Printf("Error getting instrument uid: %v", err)
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Invalid name, not found instrument uid",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"instrumentIDS": instrumentIDS,
	})
}

func (h *DBHandler) GetCandles(c echo.Context) error {
	c.Request().Header.Set("Content-Type", "application/json")

	historicCandle, err := h.Repository.GetCandles(c.QueryParam("instrument_id"))
	fmt.Println("err", err)
	if err != nil {
		fmt.Println("err:", err)
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "Invalid uid, not found candles",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"candles": historicCandle,
	})
}
