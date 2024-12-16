package minioUtils

import (
	"context"
	"fmt"
	"log"
	"time"

	//"github.com/minio/minio-go"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectMinio(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) *minio.Client {
	//minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v\n", err)
		return nil
	}
	fmt.Println("Succefully connected MiniO")
	return minioClient
}



func removeObject(ctx context.Context, bucketName string, objectName string, minioClient *minio.Client) (bool, error) {
	if err := minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{}); err != nil {
		return false, err
	}
	return true, nil
}

func listObjects(ctx context.Context, bucketName string, recursive bool, minioClient *minio.Client) <-chan minio.ObjectInfo {
	doneCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(doneCh)
		objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
			Prefix: "",
			Recursive: true,
		})
		for object := range objectCh {
			if object.Err != nil {
				log.Printf("Error occurred for %s: %v\n:", object.Key, object.Err)
				continue
			}
			doneCh <- object
		}
	}()
	return doneCh
}

func RemoveObjectsBeforeWeek(bucketName string, recursive bool, minioClient *minio.Client) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	daysAgoFromNow := time.Now().AddDate(0, 0, -5)
	objectCh := listObjects(ctx, bucketName, recursive, minioClient)

	for object := range objectCh {
		if object.Key == "loki_cluster_seed.json" {
			log.Printf("Skipping for seed file %s :\n", object.Key)
			continue
		}

		if !object.LastModified.Before(daysAgoFromNow) {
			continue
		}

		ok, err := removeObject(ctx, bucketName, object.Key, minioClient)
		if err != nil {
			log.Printf("Error occcurred while deleting %s : %v\n", object.Key, err)
			continue
		}

		if ok {
		fmt.Printf("Successfully deleted %s | %v\n", object.Key, object.LastModified)
		}
	}
	fmt.Println("SUCCESSFULLY DELETED ALL OBJECTS OLDER THAN 5 DAYS EXCEPT SEED FILE")
}
