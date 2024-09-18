package main

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"github.com/benbjohnson/clock"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ridge/must/v2"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/forms/v1"
	"google.golang.org/api/option"

	"rush/attendance"
	"rush/auth"
	"rush/golang/array"
	"rush/golang/env"
	rushHttp "rush/http"
	"rush/job"
	"rush/server"
	"rush/session"
	rushUser "rush/user"
)

func main() {
	env.Load("ENV_FILE")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoDbEndpoint := env.GetRequiredStringVariable("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(mongoDbEndpoint)

	log.Println("Connecting to MongoDB")
	mongodbClient := must.OK1(mongo.Connect(ctx, clientOptions))
	must.OK(mongodbClient.Ping(ctx, nil))

	mongodbDatabaseName := env.GetRequiredStringVariable("MONGODB_DB_NAME")
	mongodbSessionColName := env.GetRequiredStringVariable("MONGODB_SESSION_COLLECTION_NAME")
	mongodbUserColName := env.GetRequiredStringVariable("MONGODB_USER_COLLECTION_NAME")
	mongodbAttendanceColName := env.GetRequiredStringVariable("MONGODB_ATTENDANCE_COLLECTION_NAME")
	sessionCollection := mongodbClient.Database(mongodbDatabaseName).Collection(mongodbSessionColName)
	userCollection := mongodbClient.Database(mongodbDatabaseName).Collection(mongodbUserColName)
	attendanceCollection := mongodbClient.Database(mongodbDatabaseName).Collection(mongodbAttendanceColName)

	googleCreds := getGoogleCredentials(ctx, env.GetRequiredStringVariable("ENVIRONMENT"))
	log.Printf("project id: %s", googleCreds.ProjectID)
	googleOption := option.WithCredentials(googleCreds)
	formsService := must.OK1(forms.NewService(ctx, googleOption))
	driveService := must.OK1(drive.NewService(ctx, googleOption))
	firebaseAuthClient := must.OK1(must.OK1(firebase.NewApp(ctx, nil, googleOption)).Auth(ctx))

	clock := clock.New()
	userRepo := rushUser.NewMongoDbRepo(userCollection)
	sessionRepo := session.NewMongoDbRepo(sessionCollection)
	server := server.New(
		auth.NewFbAuth(firebaseAuthClient),
		// https://learn.microsoft.com/en-us/dotnet/api/system.security.cryptography.hmacsha256.-ctor?view=net-8.0
		// The secret key is recommended to be 64 bytes long for HMACSHA256. RushAuth uses HMACSHA256 to sign the token.
		auth.NewRushAuth(env.GetRequiredStringVariable("JWT_SECRET_KEY"), clock),
		userRepo, rushUser.NewAdder(userRepo),
		sessionRepo,
		attendance.NewFormHandler(formsService, driveService),
		attendance.NewMongoDbRepo(attendanceCollection, clock),
		must.OK1(time.LoadLocation("Asia/Seoul")),
	)

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{env.GetRequiredStringVariable("CORS_ORIGIN")}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	rushHttp.SetUpRouter(router, server)

	jobExecutor := job.NewExecutor(sessionRepo, server)
	if env.GetRequiredStringVariable("ENVIRONMENT") != "local" {
		cron.New().AddFunc("30 * * * *", func() {
			result, err := jobExecutor.CloseSessions()
			if err != nil {
				// TODO(#150): Extract the logic and test logging logic.
				failedSessionsIds := result.FailedSessionIds
				errors := result.Errors
				stringifiedErrors := array.Map(errors, func(err error) string { return err.Error() })
				log.Printf("Failed to close sessions [%s]: %s", strings.Join(failedSessionsIds, ", "), strings.Join(stringifiedErrors, ", "))
			}
			log.Printf("Closed sessions: %s", strings.Join(result.SessionIdsSucceeded, ", "))
		})
	}

	log.Println("Starting server")
	router.Run(":8080")
}

func getGoogleCredentials(ctx context.Context, environment string) *google.Credentials {
	scopes := []string{forms.FormsBodyScope, drive.DriveScope}
	if environment == "local" {
		googleCredsPath := must.OK1(getAbsolutePath((env.GetRequiredStringVariable("GOOGLE_CREDENTIALS_PATH"))))
		jsonCreds := must.OK1(os.ReadFile(googleCredsPath))
		return must.OK1(google.CredentialsFromJSON(ctx, jsonCreds, scopes...))
	}

	base64UrlEncodedFile := env.GetRequiredStringVariable("GOOGLE_CREDENTIALS_JSON_BASE64URL_ENCODED")
	encoding := base64.RawURLEncoding
	decoded := must.OK1(encoding.DecodeString(base64UrlEncodedFile))
	jsonCreds := []byte(decoded)
	return must.OK1(google.CredentialsFromJSON(ctx, jsonCreds, scopes...))
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
