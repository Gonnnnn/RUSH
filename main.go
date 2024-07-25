package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ridge/must/v2"
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
	mongodbClient := must.OK1(mongo.Connect(ctx, clientOptions))
	must.OK(mongodbClient.Ping(ctx, nil))

	mongodbDatabaseName := env.GetRequiredStringVariable("MONGODB_DB_NAME")
	mongodbSessionColName := env.GetRequiredStringVariable("MONGODB_SESSION_COLLECTION_NAME")
	mongodbUserColName := env.GetRequiredStringVariable("MONGODB_USER_COLLECTION_NAME")
	mongodbAttendanceReportColName := env.GetRequiredStringVariable("MONGODB_ATTENDANCE_REPORT_COLLECTION_NAME")
	sessionCollection := mongodbClient.Database(mongodbDatabaseName).Collection(mongodbSessionColName)
	userCollection := mongodbClient.Database(mongodbDatabaseName).Collection(mongodbUserColName)
	attendanceReportCollection := mongodbClient.Database(mongodbDatabaseName).Collection(mongodbAttendanceReportColName)

	server := server.New(user.NewMongoDbRepo(userCollection), session.NewMongoDbRepo(sessionCollection), attendance.NewMongoDbRepo(attendanceReportCollection))

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
