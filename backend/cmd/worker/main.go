package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/project-capillary/backend/internal/application/usecases"
	"github.com/project-capillary/backend/internal/domain/entities"
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
	examinationRepo := postgres.NewExaminationRepository(db.DB)
	reportRepo := postgres.NewReportRepository(db.DB)
	userRepo := postgres.NewUserRepository(db.DB)

	reportUseCase := usecases.NewReportUseCase(reportRepo, examinationRepo, analysisRepo)

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

		err = checkAndGenerateReport(ctx, analysis.ExaminationID, analysisRepo, examinationRepo, reportRepo, reportUseCase, userRepo)
		if err != nil {
			log.Printf("Worker: Failed to check/generate report: %v", err)
		}

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

func checkAndGenerateReport(
	ctx context.Context,
	examinationID string,
	analysisRepo *postgres.AnalysisRepositoryImpl,
	examinationRepo *postgres.ExaminationRepositoryImpl,
	reportRepo *postgres.ReportRepositoryImpl,
	reportUseCase *usecases.ReportUseCase,
	userRepo *postgres.UserRepositoryImpl,
) error {
	analyses, err := analysisRepo.GetByExaminationID(ctx, examinationID)
	if err != nil {
		return fmt.Errorf("failed to get analyses: %w", err)
	}

	if len(analyses) == 0 {
		return nil
	}
	allCompleted := true
	for _, analysis := range analyses {
		if analysis.Status != entities.AnalysisStatusCompleted {
			allCompleted = false
			break
		}
	}

	if !allCompleted {
		log.Printf("Worker: Not all analyses completed for examination %s", examinationID)
		return nil
	}

	existingReport, err := reportRepo.GetByExaminationID(ctx, examinationID)
	if err == nil && existingReport != nil {
		log.Printf("Worker: Report already exists for examination %s", examinationID)
		return nil
	}

	examination, err := examinationRepo.GetByID(ctx, examinationID)
	if err != nil || examination == nil {
		return fmt.Errorf("failed to get examination: %w", err)
	}

	summary, diagnosis, recommendations := generateReportContent(analyses)

	systemUser, err := userRepo.GetByUsername(ctx, "system")
	if err != nil || systemUser == nil {
		users, err := userRepo.List(ctx, 1, 0)
		if err != nil || len(users) == 0 {
			return fmt.Errorf("no users found to assign as report generator")
		}
		systemUser = users[0]
	}

	report := entities.NewReport(
		uuid.New().String(),
		examinationID,
		fmt.Sprintf("Отчёт по исследованию %s", examination.Description),
		fmt.Sprintf("Проанализировано изображений: %d\n\nДетали:\n%s", len(analyses), summary),
		summary,
		diagnosis,
		recommendations,
		systemUser.ID,
	)

	err = reportRepo.Create(ctx, report)
	if err != nil {
		return fmt.Errorf("failed to create report: %w", err)
	}

	examination.Complete()
	err = examinationRepo.Update(ctx, examination)
	if err != nil {
		log.Printf("Worker: Failed to update examination status: %v", err)
	}

	log.Printf("Worker: Successfully generated report %s for examination %s", report.ID, examinationID)
	return nil
}

func generateReportContent(analyses []*entities.Analysis) (summary, diagnosis, recommendations string) {
	totalDensity := 0.0
	totalDiameter := 0.0
	totalTortuosity := 0.0
	count := float64(len(analyses))

	for _, analysis := range analyses {
		if density, ok := analysis.Metrics["density"].(float64); ok {
			totalDensity += density
		}
		if diameter, ok := analysis.Metrics["diameter"].(float64); ok {
			totalDiameter += diameter
		}
		if tortuosity, ok := analysis.Metrics["tortuosity"].(float64); ok {
			totalTortuosity += tortuosity
		}
	}

	avgDensity := totalDensity / count
	avgDiameter := totalDiameter / count
	avgTortuosity := totalTortuosity / count

	summary = fmt.Sprintf(
		"Средние показатели капилляров:\n"+
			"- Плотность: %.2f сосудов/мм²\n"+
			"- Диаметр: %.2f мкм\n"+
			"- Извитость: %.2f",
		avgDensity, avgDiameter, avgTortuosity,
	)

	if avgDensity < 8.0 {
		diagnosis = "Наблюдается пониженная плотность капилляров"
		recommendations = "Рекомендуется дополнительное обследование для выявления причин снижения капиллярной плотности"
	} else if avgDensity > 9.5 {
		diagnosis = "Плотность капилляров в норме"
		recommendations = "Продолжать плановое наблюдение"
	} else {
		diagnosis = "Плотность капилляров в пределах нормы"
		recommendations = "Рекомендуется контрольное обследование через 6 месяцев"
	}

	if avgTortuosity > 2.0 {
		diagnosis += ". Повышенная извитость капилляров"
		recommendations += ". Требуется консультация специалиста"
	}

	return summary, diagnosis, recommendations
}
