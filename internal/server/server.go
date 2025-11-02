package server

import (
	"fmt"
	"log"

	"wallet_service/config"
	"wallet_service/handler"
	"wallet_service/internal/repository"
	"wallet_service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Server struct {
	DB            *sqlx.DB
	WalletRepo    repository.WalletRepositoryInterface
	WalletService *service.WalletService
	WalletHandler *handler.WalletHandler
	Router        *gin.Engine
}

func NewServer(cfg *config.Config) (*Server, error) {
	// Connect to database
	dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)

	db, err := sqlx.Connect("postgres", dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations with goose
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize layers
	walletRepo := repository.NewWalletRepository(db)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := handler.NewWalletHandler(walletService)

	// Setup Gin router
	r := gin.Default()
	api := r.Group("/api/v1")
	{
		api.POST("/wallets", walletHandler.CreateWallet)
		api.POST("/wallet", walletHandler.PerformWalletOperation)
		api.GET("/wallets/:wallet_uuid", walletHandler.GetWalletBalance)
	}

	server := &Server{
		DB:            db,
		WalletRepo:    walletRepo,
		WalletService: walletService,
		WalletHandler: walletHandler,
		Router:        r,
	}

	return server, nil
}

func runMigrations(db *sqlx.DB) error {
	goose.SetDialect("postgres")
	if err := goose.Up(db.DB, "./migrations"); err != nil {
		return fmt.Errorf("failed to run goose migrations: %w", err)
	}
	log.Println("Migrations completed successfully")
	return nil
}
