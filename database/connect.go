package database

import (
	"fmt"
	"log"
	"os"

	"github.com/sankaungkin/fiber-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() {

	var err error

	Host := os.Getenv("DB_HOST")
	Port := os.Getenv("API_PORT")
	Password := os.Getenv("POSTGRES_PASSWORD")
	User := os.Getenv("POSTGRES_USER")
	DBName := os.Getenv("POSTGRES_DB")
	SSLMode := os.Getenv("SSLMODE")

	dsn := fmt.Sprintf(
		"host=%s port=%s password=%s user=%s dbname=%s sslmode=%s",
		Host, Port, Password, User, DBName, SSLMode)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Migrating database.......")
	err = DB.AutoMigrate(&models.Category{}, &models.Customer{}, &models.Supplier{}, &models.Product{}, &models.Inventory{}, &models.Sale{}, &models.SaleDetail{}, &models.Purchase{}, &models.PurchaseDetail{}, &models.User{})
	if err != nil {
		log.Fatal(err)
	}

}
