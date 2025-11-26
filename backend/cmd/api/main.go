package main

import (
	"fmt"
	"log"

	"github.com/project-capillary/backend/internal/application/usecases"
	"github.com/project-capillary/backend/internal/infrastructure/config"
	"github.com/project-capillary/backend/internal/infrastructure/db/postgres"
	httpInfra "github.com/project-capillary/backend/internal/infrastructure/http"
	"github.com/project-capillary/backend/internal/infrastructure/http/handlers"
	"github.com/project-capillary/backend/internal/infrastructure/mq"
	"github.com/project-capillary/backend/internal/infrastructure/ws"
)

func main() {
	cfg := config.Load()

	db, err := postgres.NewPostgresDB(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection established")

	patientRepo := postgres.NewPatientRepository(db.DB)
	examinationRepo := postgres.NewExaminationRepository(db.DB)
	imageRepo := postgres.NewImageRepository(db.DB)
	analysisRepo := postgres.NewAnalysisRepository(db.DB)
	reportRepo := postgres.NewReportRepository(db.DB)
	userRepo := postgres.NewUserRepository(db.DB)
	_ = postgres.NewDeviceRepository(db.DB)

	mqPublisher, err := mq.NewRabbitMQ(cfg.RabbitMQ.URL, cfg.RabbitMQ.QueueName)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer mqPublisher.Close()

	log.Println("RabbitMQ connection established")

	patientUseCase := usecases.NewPatientUseCase(patientRepo)
	examinationUseCase := usecases.NewExaminationUseCase(examinationRepo, analysisRepo, imageRepo, mqPublisher)
	reportUseCase := usecases.NewReportUseCase(reportRepo, examinationRepo, analysisRepo, imageRepo)
	userUseCase := usecases.NewUserUseCase(userRepo)

	patientHandler := handlers.NewPatientHandler(patientUseCase)
	examinationHandler := handlers.NewExaminationHandler(examinationUseCase)
	reportHandler := handlers.NewReportHandler(reportUseCase)
	userHandler := handlers.NewUserHandler(userUseCase)

	deviceManager := ws.NewDeviceManager(cfg.Storage.PhotoPath)

	router := httpInfra.NewRouter(
		patientHandler,
		examinationHandler,
		userHandler,
		reportHandler,
		deviceManager,
	)

	engine := router.Setup()

	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)

	if err := engine.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
