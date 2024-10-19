package clients

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	pgOnce      sync.Once
	mongoOnce   sync.Once
	pgPool      *pgxpool.Pool
	mongoClient *mongo.Client
	pgErr       error
	mongoErr    error
)

func GetPostgresPool() (*pgxpool.Pool, error) {
	pgOnce.Do(func() {
		ctx := context.Background()
		connStr := os.Getenv("POSTGRES_CONNECTION_URL")
		pgPool, pgErr = pgxpool.New(ctx, connStr)
		if pgErr != nil {
			log.Fatalf("Unable to connect to PostgreSQL: %v\n", pgErr)
		}
		fmt.Println("Connected to PostgreSQL")
	})
	return pgPool, pgErr
}

func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {
		mongoURI := os.Getenv("MONGODB_CONNECTION_URL")
		clientOptions := options.Client().ApplyURI(mongoURI)

		mongoClient, mongoErr = mongo.Connect(context.TODO(), clientOptions)
		if mongoErr != nil {
			log.Fatalf("MongoDB connection error: %v", mongoErr)
		}

		fmt.Println("Connected to MongoDB")
	})
	return mongoClient, mongoErr
}

func CleanupDbConnections() {
	if pgPool != nil {
		pgPool.Close()
		fmt.Println("PostgreSQL connection closed.")
	}
	if mongoClient != nil {
		err := mongoClient.Disconnect(context.Background())
		if err != nil {
			log.Fatalf("MongoDB disconnect error: %v", err)
		}
		fmt.Println("MongoDB connection closed.")
	}
}
