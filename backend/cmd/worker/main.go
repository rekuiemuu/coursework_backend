package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/project-capillary/backend/internal/infrastructure/config"
	"github.com/project-capillary/backend/internal/infrastructure/db/postgres"
	"github.com/project-capillary/backend/internal/infrastructure/mq"
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

	log.Println("Worker: Database connection established")

	analysisRepo := postgres.NewAnalysisRepository(db.DB)

	mqConsumer, err := mq.NewRabbitMQ(cfg.RabbitMQ.URL, cfg.RabbitMQ.QueueName)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer mqConsumer.Close()

	log.Println("Worker: RabbitMQ connection established")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Worker: Received shutdown signal")
		cancel()
	}()

	log.Println("Worker: Starting message consumer...")

	err = mqConsumer.Consume(ctx, func(body []byte) error {
		var taskMsg mq.AnalysisTaskMessage
		if err := json.Unmarshal(body, &taskMsg); err != nil {
			log.Printf("Worker: Failed to parse message: %v", err)
			return err
		}

		log.Printf("Worker: Processing analysis task %s for image %s", taskMsg.AnalysisID, taskMsg.ImageID)

		analysis, err := analysisRepo.GetByID(ctx, taskMsg.AnalysisID)
		if err != nil {
			log.Printf("Worker: Failed to get analysis: %v", err)
			return err
		}

		if analysis == nil {
			log.Printf("Worker: Analysis %s not found", taskMsg.AnalysisID)
			return nil
		}

		analysis.StartProcessing()
		analysisRepo.Update(ctx, analysis)

		time.Sleep(2 * time.Second)

		metrics := generateMockMetrics()

		analysis.Complete(metrics)
		err = analysisRepo.Update(ctx, analysis)
		if err != nil {
			log.Printf("Worker: Failed to update analysis: %v", err)
			return err
		}

		log.Printf("Worker: Completed analysis task %s", taskMsg.AnalysisID)
		return nil
	})

	if err != nil {
		log.Printf("Worker: Consumer error: %v", err)
	}

	log.Println("Worker: Shutdown complete")
}

func generateMockMetrics() map[string]interface{} {
	return map[string]interface{}{
		"density":       7.0 + rand.Float64()*3.0,
		"diameter":      10.0 + rand.Float64()*5.0,
		"tortuosity":    1.0 + rand.Float64()*1.5,
		"regularity":    0.7 + rand.Float64()*0.3,
		"visibility":    0.8 + rand.Float64()*0.2,
		"abnormalities": []string{},
	}
}
