package minioClient

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Conn *minio.Client

func InitConnection() {
	endpoint := "127.0.0.1:9000"
	accessKey := "LOI19XZ652eFhm1n"
	secretKey := "gQT7hsb2vIIZ5eUeYndlBRkfY2SA20g1"
	var err error
	Conn, err = minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
	})

	if err != nil {
		log.Fatalln(err)
	}
}
