package minio

import (
	"bytes"
	"context"
	"io"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/util"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

var (
	once sync.Once
	cli  *minio.Client
	ctx  = context.Background()
)

// Init initializes the MiniIO clients as a singleton.
// It reads Minio configuration from config.App.MinioConfig
// It returns nil if minio is not enabled.
func Init() (err error) {
	cfg := config.App.MinioConfig
	if !cfg.Enable {
		return nil
	}
	once.Do(func() {
		if cli, err = New(cfg); err != nil {
			err = errors.Wrap(err, "failed to create minio client")
			zap.S().Error(err)
		}
	})
	zap.S().Infow("successfully connect to minio", "endpoint", cfg.Endpoint, "bucket", cfg.Bucket)
	return
}

// New creates a new Minio client instance with the given configuration.
func New(cfg config.MinioConfig) (*minio.Client, error) {
	return minio.New(cfg.Endpoint, buildOptions(cfg))
}

func buildOptions(cfg config.MinioConfig) *minio.Options {
	return &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSsl,
	}
}

func Put(reader io.Reader, size int64) (filename string, err error) {
	region := config.App.MinioConfig.Region
	bucket := config.App.MinioConfig.Bucket
	if len(region) > 0 {
		err = cli.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: region})
	} else {
		err = cli.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	}
	if err != nil {
		exists, errBucketExists := cli.BucketExists(ctx, config.App.MinioConfig.Bucket)
		if errBucketExists == nil && exists {
			goto CONTINUE
		}
		return
	}
CONTINUE:
	filename = util.UUID()
	_, err = cli.PutObject(ctx, bucket, filename, reader, size, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	return
}

func Get(filename string) ([]byte, error) {
	object, err := cli.GetObject(ctx, config.App.MinioConfig.Bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, object); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
