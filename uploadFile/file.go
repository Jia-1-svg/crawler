package uploadFile

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func UploadFile(file *multipart.FileHeader) (string, error) {
	ctx := context.Background()
	endpoint := "115.190.43.4:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called testbucket.
	files, err := file.Open()
	if err != nil {
		return "", nil
	}
	//延迟关闭文件
	defer files.Close()

	bucketName := "2308a"

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// Upload the test file
	// Change the value of filePath if the file is in another location
	objectName := fmt.Sprintf("%s", uuid.New().String())
	contentType := "image/png"

	// Upload the test file with FPutObject
	info, err := minioClient.PutObject(ctx, bucketName, objectName, files, file.Size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}
	object, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, time.Hour, nil)
	if err != nil {
		return "", nil
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	return object.String(), err
}
