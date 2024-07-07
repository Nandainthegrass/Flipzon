package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Nandainthegrass/Flipzon/Authentication/cmd/api"
	"github.com/Nandainthegrass/Flipzon/Authentication/db"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading env files: %v", err)
	}
	fmt.Println(os.Getenv("User"))
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 os.Getenv("User"),
		Passwd:               os.Getenv("Passwd"),
		Addr:                 os.Getenv("Addr"),
		DBName:               os.Getenv("DBName"),
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	initStorage(db)

	//	raddr := os.Getenv("REDIS_ADDR")
	//	rdb := Redis.NewRedisClient(raddr) //type: *redis.Client

	server := api.NewAPIServer(":8001", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}

func initStorage(db *sql.DB) {
	err := db.Ping() //pings the database
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully connected!")
}
