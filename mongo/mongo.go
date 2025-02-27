package mongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
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

// Init initializes the global MongoDB client.
// It reads MongoDB configuration from config.App.MongoConfig.
// If MongoDB is not enabled, it returns nil.
// The function is thread-safe and ensures the client is initialized only once.
func Init() (err error) {
	cfg := config.App.MongoConfig
	if !cfg.Enable {
		return nil
	}

	mu.Lock()
	defer mu.Unlock()
	if initialized {
		return nil
	}

	uri := buildURI(cfg)
	client, err = mongo.Connect(options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetConnectTimeout(cfg.ConnectTimeout),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return errors.Wrap(err, "failed to connect to mongodb")
	}
	zap.S().Infow("successfully connect to mongodb", "host", cfg.Host, "port", cfg.Port, "database", cfg.Database)

	initialized = true
	return err
}

// New returns a new MongoDB client instance with given configuration.
// It's the caller's responsibility to close the client,
// caller should always call Close() when it's no longer needed.
func New(cfg config.MongoConfig) (*mongo.Client, error) {
	uri := buildURI(cfg)
	return mongo.Connect(options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(cfg.MaxPoolSize).
		SetMinPoolSize(cfg.MinPoolSize).
		SetConnectTimeout(cfg.ConnectTimeout),
	)
}

func buildURI(cfg config.MongoConfig) string {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port,
		cfg.Database, cfg.AuthSource,
	)
	if len(cfg.Username) == 0 && len(cfg.Password) == 0 {
		uri = fmt.Sprintf("mongodb://%s:%d/%s", cfg.Host, cfg.Port, cfg.Database)
	}
	return uri
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
