package http

import (
	"github.com/gin-gonic/gin"
	"github.com/project-capillary/backend/internal/infrastructure/http/handlers"
	"github.com/project-capillary/backend/internal/infrastructure/http/middleware"
	"github.com/project-capillary/backend/internal/infrastructure/ws"
)

type Router struct {
	patientHandler     *handlers.PatientHandler
	examinationHandler *handlers.ExaminationHandler
	userHandler        *handlers.UserHandler
	reportHandler      *handlers.ReportHandler
	deviceManager      *ws.DeviceManager
}

func NewRouter(
	patientHandler *handlers.PatientHandler,
	examinationHandler *handlers.ExaminationHandler,
	userHandler *handlers.UserHandler,
	reportHandler *handlers.ReportHandler,
	deviceManager *ws.DeviceManager,
) *Router {
	return &Router{
		patientHandler:     patientHandler,
		examinationHandler: examinationHandler,
		userHandler:        userHandler,
		reportHandler:      reportHandler,
		deviceManager:      deviceManager,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", r.userHandler.Login)
			auth.POST("/register", r.userHandler.CreateUser)
		}

		users := api.Group("/users")
		{
			users.GET("/:id", r.userHandler.GetUser)
		}

		patients := api.Group("/patients")
		{
			patients.POST("", r.patientHandler.CreatePatient)
			patients.GET("/:id", r.patientHandler.GetPatient)
			patients.PUT("/:id", r.patientHandler.UpdatePatient)
			patients.DELETE("/:id", r.patientHandler.DeletePatient)
			patients.GET("", r.patientHandler.ListPatients)
		}

		examinations := api.Group("/examinations")
		{
			examinations.POST("", r.examinationHandler.CreateExamination)
			examinations.GET("", r.examinationHandler.ListExaminations)
			examinations.GET("/:id", r.examinationHandler.GetExamination)
			examinations.POST("/:id/photos", r.examinationHandler.AttachPhotos)
			examinations.POST("/:id/analyze", r.examinationHandler.StartAnalysis)
			examinations.GET("/patient/:patientId", r.examinationHandler.GetPatientExaminations)
		}

		reports := api.Group("/reports")
		{
			reports.GET("", r.reportHandler.ListReports)
			reports.POST("", r.reportHandler.CreateReport)
			reports.GET("/:id", r.reportHandler.GetReport)
			reports.GET("/examination/:examinationId", r.reportHandler.GetExaminationReport)
		}
	}

	api.GET("/photos/:filename", func(c *gin.Context) {
		r.deviceManager.ServePhoto(c.Writer, c.Request, c.Param("filename"))
	})

	router.GET("/ws", func(c *gin.Context) {
		r.deviceManager.HandleWebSocket(c.Writer, c.Request)
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
