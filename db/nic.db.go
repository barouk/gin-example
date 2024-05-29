package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB
var err error

func ConnectDb() *gorm.DB {

	end := godotenv.Load("env")
	if end != nil {
		panic("Faild to load .env file")
	}

	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbHost := os.Getenv("DBHOST")
	dbPort := os.Getenv("DBPORT")
	dbName := os.Getenv("DBNAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
	})
	if err != nil {
		panic("Failed to connect to postgres database")
	}
	//err = db.AutoMigrate(&entity.User{})
	//if err != nil {
	//	panic("Failed : Unable to migrate your postgres database")
	//}
	return db
}

func CloseDb(db *gorm.DB) {
	dbPsql, err := db.DB()
	if err != nil {
		panic("Failed: postgres database connection")
	}
	err = dbPsql.Close()
	if err != nil {
		panic("Failed: unable to close postgre connection database")
	}
}
