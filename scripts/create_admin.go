package main

import (
	"fmt"
	"log"

	"hr-backend/internal/config"
	"hr-backend/internal/database"
	"hr-backend/internal/models"
	"hr-backend/internal/utils"
)

func main() {
	fmt.Println("Creating admin user...")

	// Load configuration
	cfg := config.Load()

	// Initialize JWT (required for utils)
	utils.InitJWT(cfg.JWT.Secret)

	// Connect to database
	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Ensure tables exist
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	db := database.GetDB()

	// Check if admin already exists
	var existingUser models.User
	if err := db.Where("email = ?", "admin@hrms.com").First(&existingUser).Error; err == nil {
		fmt.Println("Admin user already exists!")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Create admin user
	user := &models.User{
		Email:        "admin@hrms.com",
		PasswordHash: hashedPassword,
		Role:         "admin",
		IsActive:     true,
	}

	if err := db.Create(user).Error; err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	fmt.Println("✓ Admin user created successfully!")
	fmt.Println("\nCredentials:")
	fmt.Println("  Email:    admin@hrms.com")
	fmt.Println("  Password: admin123")
	fmt.Println("\n⚠️  Please change the password after first login!")
}
