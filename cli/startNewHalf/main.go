package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rush/user"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	csvPath := flag.String("csv", "./new_members.csv", "path to new members CSV")
	outDir := flag.String("out", ".", "output directory for backup JSON")
	mongoURI := flag.String("mongo-uri", "", "MongoDB URI (required)")
	dbName := flag.String("db", "rush", "database name")
	usersCol := flag.String("users-col", "users", "users collection name")
	sessionsCol := flag.String("sessions-col", "sessions", "sessions collection name")
	attendancesCol := flag.String("attendances-col", "attendances", "attendances collection name")
	flag.Parse()

	if *mongoURI == "" {
		log.Fatal("-mongo-uri is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(*mongoURI))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")

	db := client.Database(*dbName)
	usersCollection := db.Collection(*usersCol)
	sessionsCollection := db.Collection(*sessionsCol)
	attendancesCollection := db.Collection(*attendancesCol)

	// --- Export ---
	// Use raw bson.M to preserve all fields exactly as stored (including is_deleted, force_apply, etc.)
	userDocs, err := fetchAll(ctx, usersCollection)
	if err != nil {
		log.Fatalf("failed to fetch users: %v", err)
	}
	sessionDocs, err := fetchAll(ctx, sessionsCollection)
	if err != nil {
		log.Fatalf("failed to fetch sessions: %v", err)
	}
	attendanceDocs, err := fetchAll(ctx, attendancesCollection)
	if err != nil {
		log.Fatalf("failed to fetch attendances: %v", err)
	}

	backup := map[string]interface{}{
		"exported_at": time.Now().Format(time.RFC3339),
		"users":       userDocs,
		"sessions":    sessionDocs,
		"attendances": attendanceDocs,
	}

	filename := fmt.Sprintf("rush_backup_%s.json", time.Now().Format("20060102_150405"))
	outPath := filepath.Join(*outDir, filename)

	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal backup: %v", err)
	}
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		log.Fatalf("failed to write backup file: %v", err)
	}

	log.Printf("Exported %d users, %d sessions, %d attendances to %s",
		len(userDocs), len(sessionDocs), len(attendanceDocs), outPath)

	// --- Confirmation ---
	fmt.Print("This will DELETE all data. Type 'yes' to continue: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if strings.TrimSpace(scanner.Text()) != "yes" {
		log.Fatal("Aborted")
	}

	// --- Drop ---
	if err := usersCollection.Drop(ctx); err != nil {
		log.Fatalf("failed to drop users: %v", err)
	}
	if err := sessionsCollection.Drop(ctx); err != nil {
		log.Fatalf("failed to drop sessions: %v", err)
	}
	if err := attendancesCollection.Drop(ctx); err != nil {
		log.Fatalf("failed to drop attendances: %v", err)
	}
	log.Println("Dropped all collections")

	// --- Seed ---
	newUsers, err := user.ParseCSV(*csvPath)
	if err != nil {
		log.Fatalf("failed to parse CSV: %v", err)
	}

	repo := user.NewMongoDbRepo(usersCollection)
	count, err := repo.AddMany(newUsers)
	if err != nil {
		log.Fatalf("failed to insert users: %v", err)
	}

	log.Printf("%d users inserted", count)
}

func fetchAll(ctx context.Context, col *mongo.Collection) ([]bson.M, error) {
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	if docs == nil {
		docs = []bson.M{}
	}
	return docs, nil
}
