package database
import (
	"fmt"
	"log"

	"keuangan_app/config"
	"keuangan_app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)
	log.Printf("DEBUG pw_len=%d host=%q user=%q db=%q port=%q", len(cfg.DBPassword), cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}


	if err = DB.AutoMigrate(&models.Transaksi{}); err != nil {
		log.Fatal("Failed to migrate: ", err)
	}

	log.Println("Database connected & migrated")
}