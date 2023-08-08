package minio

import (
	"context"
	"os"

	"github.com/Goboolean/shared/pkg/resolver"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)



const bucketName = "test"


type Storage struct {
	client *minio.Client
}


func New(c *resolver.ConfigMap) (*Storage, error) {

	accessKey, err := c.GetStringKey("ACCESS_KEY")
	if err != nil {
		return nil, err
	}

	secretKey, err := c.GetStringKey("SECRET_KEY")
	if err != nil {
		return nil, err
	}

	host, err := c.GetStringKey("HOST")
	if err != nil {
		return nil, err
	}

	client, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	return &Storage{client: client}, nil
}



func (m *Storage) GetFile(ctx context.Context, name string) (*os.File, error) {
	object, err :=  m.client.GetObject(ctx, bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	return file, nil
}



func (m *Storage) Close() {}