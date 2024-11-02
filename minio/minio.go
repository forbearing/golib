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
)

var (
	once sync.Once
	cli  *minio.Client
	ctx  = context.Background()
)

func Init() (err error) {
	if !config.App.MinioConfig.Enable {
		return nil
	}

	once.Do(func() {
		endpoint := config.App.MinioConfig.Endpoint
		accessKey := config.App.MinioConfig.AccessKey
		secretKey := config.App.MinioConfig.SecretKey
		useSsl := config.App.MinioConfig.UseSsl
		cli, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: useSsl,
		})
	})
	return
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
