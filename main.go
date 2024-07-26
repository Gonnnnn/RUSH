package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ridge/must/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/forms/v1"
	"google.golang.org/api/option"

	"rush/attendance"
	"rush/golang/env"
	rushHttp "rush/http"
	"rush/server"
	"rush/session"
	rushUser "rush/user"
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

	googleCreds := getGoogleCredentials(ctx, env.GetRequiredStringVariable("ENVIRONMENT"))
	log.Printf("project id: %s", googleCreds.ProjectID)
	formsService := must.OK1(forms.NewService(ctx, option.WithCredentials(googleCreds)))

	server := server.New(rushUser.NewMongoDbRepo(userCollection), session.NewMongoDbRepo(sessionCollection),
		attendance.NewMongoDbRepo(attendanceReportCollection), session.NewFormHandler(formsService))

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

func getGoogleCredentials(ctx context.Context, environment string) *google.Credentials {
	if environment == "local" {
		googleCredsPath := must.OK1(getAbsolutePath((env.GetRequiredStringVariable("GOOGLE_CREDENTIALS_PATH"))))
		jsonCreds := must.OK1(os.ReadFile(googleCredsPath))
		return must.OK1(google.CredentialsFromJSON(ctx, jsonCreds, forms.FormsBodyScope))
	}

	base64UrlEncodedFile := env.GetRequiredStringVariable("GOOGLE_CREDENTIALS_JSON_BASE64URL_ENCODED")
	encoding := base64.RawURLEncoding
	decoded := must.OK1(encoding.DecodeString(base64UrlEncodedFile))
	jsonCreds := []byte(decoded)
	return must.OK1(google.CredentialsFromJSON(ctx, jsonCreds, forms.FormsBodyScope))
}

func getAbsolutePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	if path[:1] == "~" {
		return filepath.Join(usr.HomeDir, path[1:]), nil
	}

	return filepath.Abs(path)
}
