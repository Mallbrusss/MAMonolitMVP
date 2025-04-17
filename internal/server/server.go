package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"log"
	"mamonolitmvp/config"
	"mamonolitmvp/internal/handlers/analyzer"
	"mamonolitmvp/internal/handlers/etl"
	"mamonolitmvp/internal/repository"
	"mamonolitmvp/internal/services"
	"mamonolitmvp/internal/storage/timescale"
	"net/http"
)

type Server struct {
	cfg *config.Config
	e   *echo.Echo
	db  *gorm.DB
}

func NewServer() *Server {
	return &Server{
		cfg: config.LoadConfig(),
		e:   echo.New(),
	}
}

func (s *Server) initializeMiddleware() {
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())
	s.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) initializeDatabase() error {
	db, err := timescale.InitDB(s.cfg.PostgresHost, s.cfg.PostgresUser,
		s.cfg.PostgresPassword, s.cfg.PostgresDatabase, s.cfg.PostgresPort)
	if err != nil {
		log.Fatal(err)
		return err
	}

	s.db = db
	return nil
}

func (s *Server) initializeRepository(db *gorm.DB) *repository.InstrumentRepository {
	return repository.NewInstrumentRepository(db)
}

func (s *Server) registerRoutes(instrumentRepository *repository.InstrumentRepository) {
	service := services.NewTinkoffService(s.cfg, instrumentRepository)
	etlHandler := etl.NewETLHandler(service)
	signalHandler := analyzer.NewSignalHandler(service)

	s.e.GET("/api/v1/ti/getClosePrices", etlHandler.GetClosePricesHandler)
	s.e.GET("/api/v1/ti/getBonds", etlHandler.GetAllBonds)
	s.e.GET("/api/v1/ti/getCandles", etlHandler.GetCandles)
	s.e.GET("/api/v1/sig/getSignals", signalHandler.GetSignals)

	dbHandler := etl.NewDBHandler(instrumentRepository)
	s.e.GET("/api/v1/db/getInstrumentIDs", dbHandler.GetInstrumentUIDAndFigi)
	s.e.GET("/api/v1/db/getCandles", dbHandler.GetCandles)

	log.Printf("Server is running on port %s...", s.cfg.ServerPort)
}

func (s *Server) Run() error {
	err := s.initializeDatabase()
	if err != nil {
		log.Fatal(err)
		return err
	}
	s.initializeMiddleware()
	initializeRepository := s.initializeRepository(s.db)
	s.registerRoutes(initializeRepository)

	address := fmt.Sprintf(":%s", s.cfg.ServerPort)
	return s.e.Start(address)
}
