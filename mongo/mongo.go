package mongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/forbearing/golib/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.uber.org/zap"
)

var (
	initialized bool
	client      *mongo.Client
	mu          sync.RWMutex
)

func Init() (err error) {
	if !config.App.MongoConfig.Enable {
		return nil
	}

	mu.Lock()
	defer mu.Unlock()
	if initialized {
		return nil
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s",
		config.App.MongoConfig.Username,
		config.App.MongoConfig.Password,
		config.App.MongoConfig.Host,
		config.App.MongoConfig.Port,
		config.App.MongoConfig.Database,
		config.App.MongoConfig.AuthSource,
	)
	if len(config.App.MongoConfig.Username) == 0 && len(config.App.MongoConfig.Password) == 0 {
		uri = fmt.Sprintf("mongodb://%s:%d/%s",
			config.App.MongoConfig.Host,
			config.App.MongoConfig.Port,
			config.App.MongoConfig.Database,
		)
	}
	client, err = mongo.Connect(options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(config.App.MongoConfig.MaxPoolSize).
		SetMinPoolSize(config.App.MongoConfig.MinPoolSize).
		SetConnectTimeout(config.App.MongoConfig.ConnectTimeout),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("ping mongodb failed: %w", err)
	}
	zap.S().Infow("successfully connect to mongodb",
		"host", config.App.MongoConfig.Host,
		"port", config.App.MongoConfig.Port,
		"database", config.App.MongoConfig.Database,
	)

	initialized = true
	return err
}

// Client returns the MongoDB client instance
func Client() (*mongo.Client, error) {
	mu.RLock()
	defer mu.RUnlock()
	if !initialized {
		return nil, fmt.Errorf("mongo client not initialized, call Init() first")
	}
	if client == nil {
		return nil, fmt.Errorf("mongo client is nil")
	}
	return client, nil
}

// Database returns a handle to the specified database
func Database(name string) (*mongo.Database, error) {
	c, err := Client()
	if err != nil {
		return nil, err
	}
	return c.Database(name), nil
}

// Collection returns a handle to the specified collection
func Collection(dbName, collName string) (*mongo.Collection, error) {
	db, err := Database(dbName)
	if err != nil {
		return nil, err
	}
	return db.Collection(collName), nil
}

// Close closes the MongoDB client connection
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect MongoDB client: %w", err)
		}
		client = nil
		initialized = false
	}
	return nil
}

// Health checks if the MongoDB connection is healthy
func Health() error {
	c, err := Client()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.Ping(ctx, readpref.Primary())
}
