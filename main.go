package main

import "proof-program-registry/src"
import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	fmt.Print("hallo")
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&src.Entry{})
	r := src.SetupRouter()
	r.Run(":8080")
}
