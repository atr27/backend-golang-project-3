package database

import (
	"fmt"
	"log"

	"hr-backend/internal/config"
	"hr-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.DatabaseConfig) error {
	var err error

	dsn := cfg.DSN()

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func Migrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Department{},
		&models.Employee{},
		&models.Attendance{},
		&models.Leave{},
		&models.LeaveBalance{},
		&models.Payroll{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrated successfully")
	return nil
}

// Reset drops all tables and recreates them
func Reset() error {
	log.Println("Resetting database...")
	
	// Drop tables in reverse order to respect foreign key constraints
	tables := []interface{}{
		&models.Payroll{},
		&models.LeaveBalance{},
		&models.Leave{},
		&models.Attendance{},
		&models.Employee{},
		&models.Department{},
		&models.User{},
	}
	
	for _, table := range tables {
		if err := DB.Migrator().DropTable(table); err != nil {
			log.Printf("Warning: Failed to drop table: %v", err)
		}
	}
	
	log.Println("All tables dropped successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
