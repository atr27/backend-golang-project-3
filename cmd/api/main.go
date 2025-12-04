package main

import (
	"flag"
	"log"

	"hr-backend/internal/config"
	"hr-backend/internal/database"
	"hr-backend/internal/handlers"
	"hr-backend/internal/middleware"
	"hr-backend/internal/repositories"
	"hr-backend/internal/services"
	"hr-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Parse command-line flags
	seedFlag := flag.Bool("seed", false, "Seed the database with initial data")
	resetFlag := flag.Bool("reset", false, "Reset the database (drop all tables)")
	flag.Parse()

	// Load configuration
	cfg := config.Load()

	// Initialize JWT
	utils.InitJWT(cfg.JWT.Secret)

	// Connect to database
	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Reset database if flag is provided
	if *resetFlag {
		if err := database.Reset(); err != nil {
			log.Fatalf("Failed to reset database: %v", err)
		}
		log.Println("Database reset completed!")
	}

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Seed database if flag is provided
	if *seedFlag {
		if err := database.Seed(); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
	}

	db := database.GetDB()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	employeeRepo := repositories.NewEmployeeRepository(db)
	deptRepo := repositories.NewDepartmentRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)
	leaveRepo := repositories.NewLeaveRepository(db)
	payrollRepo := repositories.NewPayrollRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	employeeService := services.NewEmployeeService(employeeRepo, userRepo, db)
	deptService := services.NewDepartmentService(deptRepo)
	attendanceService := services.NewAttendanceService(attendanceRepo, employeeRepo)
	leaveService := services.NewLeaveService(leaveRepo, employeeRepo)
	payrollService := services.NewPayrollService(payrollRepo, employeeRepo, db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	deptHandler := handlers.NewDepartmentHandler(deptService)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService)
	leaveHandler := handlers.NewLeaveHandler(leaveService)
	payrollHandler := handlers.NewPayrollHandler(payrollService)

	// Setup Gin router
	gin.SetMode(cfg.Server.GinMode)
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		utils.SuccessResponse(c, 200, "Service is healthy", gin.H{
			"status": "ok",
		})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/otentikasi")
		{
			auth.POST("/masuk", authHandler.Login)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Auth routes
			protected.POST("/otentikasi/keluar", authHandler.Logout)
			protected.POST("/otentikasi/ubah-kata-sandi", authHandler.ChangePassword)
			protected.GET("/otentikasi/profil", authHandler.GetProfile)

			// Dashboard
			protected.GET("/dasbor/statistik", employeeHandler.GetDashboardStats)

			// Department routes
			departments := protected.Group("/departemen")
			{
				departments.POST("", middleware.RoleMiddleware("admin", "hr_manager"), deptHandler.CreateDepartment)
				departments.GET("", deptHandler.GetDepartments)
				departments.GET("/:id", deptHandler.GetDepartmentByID)
				departments.PUT("/:id", middleware.RoleMiddleware("admin", "hr_manager"), deptHandler.UpdateDepartment)
				departments.DELETE("/:id", middleware.RoleMiddleware("admin"), deptHandler.DeleteDepartment)
			}

			// Employee routes
			employees := protected.Group("/karyawan")
			{
				employees.POST("", middleware.RoleMiddleware("admin", "hr_manager"), employeeHandler.CreateEmployee)
				employees.GET("", employeeHandler.GetEmployees)
				employees.GET("/buat-kode", middleware.RoleMiddleware("admin", "hr_manager"), employeeHandler.GenerateEmployeeCode)
				employees.GET("/:id", employeeHandler.GetEmployeeByID)
				employees.PUT("/:id", middleware.RoleMiddleware("admin", "hr_manager"), employeeHandler.UpdateEmployee)
				employees.DELETE("/:id", middleware.RoleMiddleware("admin"), employeeHandler.DeleteEmployee)
			}

			// Attendance routes
			attendance := protected.Group("/kehadiran")
			{
				attendance.POST("/absen-masuk", attendanceHandler.ClockIn)
				attendance.POST("/absen-keluar", attendanceHandler.ClockOut)
				attendance.GET("", attendanceHandler.GetAttendance)
				attendance.GET("/laporan", middleware.RoleMiddleware("admin", "hr_manager", "department_manager"), attendanceHandler.GetAttendanceReport)
				attendance.POST("/manual", middleware.RoleMiddleware("admin", "hr_manager"), attendanceHandler.CreateManualAttendance)
			}

			// Leave routes
			leaves := protected.Group("/cuti")
			{
				leaves.POST("", leaveHandler.CreateLeave)
				leaves.GET("", leaveHandler.GetLeaves)
				leaves.GET("/:id", leaveHandler.GetLeaveByID)
				leaves.PUT("/:id/setujui", middleware.RoleMiddleware("admin", "hr_manager", "department_manager"), leaveHandler.ApproveLeave)
				leaves.GET("/saldo/:employee_id", leaveHandler.GetLeaveBalance)
			}

			// Payroll routes
			payroll := protected.Group("/penggajian")
			payroll.Use(middleware.RoleMiddleware("admin", "hr_manager"))
			{
				payroll.POST("/buat", payrollHandler.GeneratePayroll)
				payroll.GET("", payrollHandler.GetPayrolls)
				payroll.GET("/ringkasan", payrollHandler.GetPayrollSummary)
				payroll.GET("/:id", payrollHandler.GetPayrollByID)
				payroll.GET("/:id/unduh", payrollHandler.DownloadPayrollSlip)
				payroll.PUT("/:id", payrollHandler.UpdatePayroll)
				payroll.POST("/:id/proses-pembayaran", payrollHandler.ProcessPayment)
			}
		}
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
