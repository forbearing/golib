package config

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	assert.NoError(t, Init())
	assert.NoError(t, Save("../testdata/config/config.ini"))
}

func TestMinio(t *testing.T) {
	ctx := context.Background()
	assert.NoError(t, Init())
	fmt.Println(App.MinioConfig)

	cli, err := minio.New(App.MinioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(App.MinioConfig.AccessKey, App.MinioConfig.SecretKey, ""),
		Secure: App.MinioConfig.UseSsl,
	})
	assert.NoError(t, err)
	_, _ = ctx, cli

	if len(App.MinioConfig.Region) > 0 {
		err = cli.MakeBucket(ctx, App.MinioConfig.Bucket, minio.MakeBucketOptions{Region: App.MinioConfig.Region})
	} else {
		err = cli.MakeBucket(ctx, App.MinioConfig.Bucket, minio.MakeBucketOptions{})
	}
	if err != nil {
		exists, errBucketExists := cli.BucketExists(ctx, App.MinioConfig.Bucket)
		if errBucketExists == nil && exists {
			goto CONTINUE
		}
		t.Fatal(err)
	}
CONTINUE:
	fileContent := make([]byte, 1024*1024)
	uploadInfo, err := cli.PutObject(ctx, App.MinioConfig.Bucket, "samplefile", bytes.NewBuffer(fileContent), int64(len(fileContent)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", uploadInfo)
	fmt.Println("===== successfully PutObject")

	object, err := cli.GetObject(ctx, App.MinioConfig.Bucket, "samplefile", minio.GetObjectOptions{})
	assert.NoError(t, err)
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, object)
	assert.NoError(t, err)
	fmt.Println("===== successfully GetObject: ", len(buf.Bytes()))
}
