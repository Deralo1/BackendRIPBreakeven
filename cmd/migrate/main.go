package main

import (
	"Backeven/internal/app/ds"
	"Backeven/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	//migrate the schema
	err = db.AutoMigrate(
		&ds.User{},
		&ds.Expense{},
		&ds.ExpenseForRequest{},
		&ds.BreakevenRequest{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
