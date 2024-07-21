package main

import (
	"database/sql"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"rush/attendance"
	rushHttp "rush/http"
	"rush/server"
	"rush/session"
	"rush/user"
)

func main() {
	db, err := sql.Open("sqlite3", "./sqlite/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createTables(db); err != nil {
		log.Fatal(err)
	}
	if err := createDummyData(db); err != nil {
		log.Fatal(err)
	}

	server := server.New(user.NewRepo(db), session.NewRepo(db), attendance.NewRepo(db))

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(corsConfig))

	rushHttp.SetUpRouter(router, server)

	router.Run(":8080")
}
