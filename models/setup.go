package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Connection() {
	dsn := "root:@tcp(localhost:3306)/cars_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// Auto-migrate the Car model
	if err := db.AutoMigrate(&Car{}); err != nil {
		log.Fatalf("Error during auto-migration: %v", err)
	}
	db = DB
}
func InitDB() {
	dsn := "root:@tcp(localhost:3306)/cars_db?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database!")
	}
}
