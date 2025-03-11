package minio

import (
	"bytes"
	"context"
	"io"
	"sync"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/util"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	once   sync.Once
	client *minio.Client
	ctx    = context.Background()
)

// Init initializes the MiniIO clients as a singleton.
// It reads Minio configuration from config.App.MinioConfig
// It returns nil if minio is not enabled.
func Init() (err error) {
	cfg := config.App.Minio
	if !cfg.Enable {
		return nil
	}
	once.Do(func() {
		client, err = New(cfg)
	})
	if err != nil {
		return errors.Wrap(err, "failed to create minio client")
	}
	zap.S().Infow("successfully connect to minio", "endpoint", cfg.Endpoint, "bucket", cfg.Bucket)
	return nil
}

// New creates a new Minio client instance with the given configuration.
func New(cfg config.Minio) (*minio.Client, error) {
	return minio.New(cfg.Endpoint, buildOptions(cfg))
}

func buildOptions(cfg config.Minio) *minio.Options {
	return &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSsl,
	}
}

func Put(reader io.Reader, size int64) (filename string, err error) {
	region := config.App.Minio.Region
	bucket := config.App.Minio.Bucket
	if len(region) > 0 {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: region})
	} else {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	}
	if err != nil {
		exists, errBucketExists := client.BucketExists(ctx, config.App.Minio.Bucket)
		if errBucketExists == nil && exists {
			goto CONTINUE
		}
		return
	}
CONTINUE:
	filename = util.UUID()
	_, err = client.PutObject(ctx, bucket, filename, reader, size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	return
}

func Get(filename string) ([]byte, error) {
	object, err := client.GetObject(ctx, config.App.Minio.Bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, object); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
