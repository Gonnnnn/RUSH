package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"rush/attendance"
	"rush/golang/env"
	rushHttp "rush/http"
	"rush/server"
	"rush/session"
	"rush/user"
)

func main() {
	env.Load("ENV_FILE")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongodbHost := env.GetRequiredStringVariable("MONGODB_HOST")
	mongodbPort := env.GetRequiredStringVariable("MONGODB_PORT")
	username := env.GetRequiredStringVariable("MONGODB_USERNAME")
	password := env.GetRequiredStringVariable("MONGODB_PASSWORD")
	mongoDbEndpoint := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, mongodbHost, mongodbPort)
	clientOptions := options.Client().ApplyURI(mongoDbEndpoint)
	log.Println("Connecting to MongoDB")
	mongodbClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = mongodbClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	mongodbDatabaseName := env.GetRequiredStringVariable("MONGODB_DB_NAME")
	mongodbCollectionName := env.GetRequiredStringVariable("MONGODB_SESSION_COLLECTION_NAME")
	sessionCollection := mongodbClient.Database(mongodbDatabaseName).Collection(mongodbCollectionName)

	log.Println("Connecting to SQLite")
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

	server := server.New(user.NewRepo(db), session.NewMongoDBRepo(sessionCollection), attendance.NewRepo(db))

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	router.Use(cors.New(corsConfig))

	rushHttp.SetUpRouter(router, server)

	log.Println("Starting server")
	router.Run(":8080")
}
