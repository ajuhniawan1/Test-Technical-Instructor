package main

import (
	"log"

	"assignment-platform/internal/config"
	"assignment-platform/internal/database"
	"assignment-platform/internal/handler"
	"assignment-platform/internal/middleware"
	"assignment-platform/internal/repository"
	"assignment-platform/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load konfigurasi dari file .env atau environment variable.
	cfg := config.Load()

	// Membuka koneksi ke MySQL.
	db := database.ConnectMySQL(cfg)
	defer db.Close()

	// Inisialisasi repository.
	userRepo := repository.NewUserRepository(db)
	classRepo := repository.NewClassRepository(db)
	assignmentRepo := repository.NewAssignmentRepository(db)
	submissionRepo := repository.NewSubmissionRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	// Inisialisasi service.
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	classService := service.NewClassService(classRepo)
	assignmentService := service.NewAssignmentService(assignmentRepo, classRepo)
	submissionService := service.NewSubmissionService(db, submissionRepo, assignmentRepo, classRepo)
	reviewService := service.NewReviewService(db, reviewRepo, classRepo)

	// Inisialisasi handler.
	authHandler := handler.NewAuthHandler(authService)
	classHandler := handler.NewClassHandler(classService)
	assignmentHandler := handler.NewAssignmentHandler(assignmentService)
	submissionHandler := handler.NewSubmissionHandler(submissionService)
	reviewHandler := handler.NewReviewHandler(reviewService)

	r := gin.Default()

	// Endpoint health check untuk memastikan API hidup.
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Bootcamp Assignment API is running"})
	})

	api := r.Group("/api/v1")

	// Endpoint public.
	api.POST("/auth/login", authHandler.Login)

	// Endpoint protected, wajib JWT.
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Authentication.
		protected.GET("/auth/me", authHandler.Me)

		// Class management.
		protected.POST("/classes", middleware.RequireRole("admin"), classHandler.CreateClass)
		protected.GET("/classes", classHandler.ListClasses)
		protected.GET("/classes/:id", classHandler.GetClassByID)
		protected.POST("/classes/:id/trainers", middleware.RequireRole("admin"), classHandler.AssignTrainer)
		protected.POST("/classes/:id/talents", middleware.RequireRole("admin"), classHandler.AssignTalent)

		// Endpoint untuk melihat trainer/talent dalam class.
		protected.GET("/classes/:id/trainers", middleware.RequireRole("admin", "trainer"), classHandler.ListTrainersByClass)
		protected.GET("/classes/:id/talents", middleware.RequireRole("admin", "trainer"), classHandler.ListTalentsByClass)

		// Assignment management.
		protected.POST("/classes/:id/assignments", middleware.RequireRole("admin", "trainer"), assignmentHandler.CreateAssignment)
		protected.GET("/classes/:id/assignments", assignmentHandler.ListAssignmentsByClass)
		protected.GET("/assignments/:id", assignmentHandler.GetAssignmentByID)
		protected.PUT("/assignments/:id", middleware.RequireRole("admin", "trainer"), assignmentHandler.UpdateAssignment)
		protected.PATCH("/assignments/:id/close", middleware.RequireRole("admin", "trainer"), assignmentHandler.CloseAssignment)

		// Submission management.
		protected.POST("/assignments/:assignmentId/submissions", middleware.RequireRole("talent"), submissionHandler.SubmitAssignment)
		protected.PUT("/submissions/:id", middleware.RequireRole("talent"), submissionHandler.ResubmitSubmission)
		protected.GET("/submissions/me", middleware.RequireRole("talent"), submissionHandler.ListMySubmissions)
		protected.GET("/classes/:id/submissions", middleware.RequireRole("admin", "trainer"), submissionHandler.ListSubmissionsByClass)

		// Review management.
		protected.POST("/submissions/:id/review", middleware.RequireRole("admin", "trainer"), reviewHandler.ReviewSubmission)
		protected.POST("/submissions/:id/request-revision", middleware.RequireRole("admin", "trainer"), reviewHandler.RequestRevision)

		// Progress tracking.
		protected.GET("/classes/:id/progress", middleware.RequireRole("admin", "trainer"), submissionHandler.GetClassProgress)
		protected.GET("/talents/:talentId/progress", middleware.RequireRole("admin", "trainer", "talent"), submissionHandler.GetTalentProgress)
	}

	log.Println("Server running on port", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
