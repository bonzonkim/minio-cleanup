package minioUtils

import (
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go"
)

func ConnectMinio(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) *minio.Client {
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v\n", err)
		return nil
	}
	fmt.Println("Succefully connected MiniO")
	return minioClient
}



func removeObject(bucketName string, objectName string, minioClient *minio.Client) (bool, error) {
	if err := minioClient.RemoveObject(bucketName, objectName); err != nil {
		return false, err
	}
	return true, nil
}

func listObjects(bucketName string, recursive bool, minioClient *minio.Client) <-chan minio.ObjectInfo {
	doneCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(doneCh)
		objectCh := minioClient.ListObjects(bucketName, "", recursive, nil)
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
	daysAgo := time.Now().AddDate(0, 0, -5)
	objectCh := listObjects(bucketName, recursive, minioClient)

	for object := range objectCh {
		if object.Key == "loki_cluster_seed.json" {
			log.Printf("Skipping for seed file %s :\n", object.Key)
			continue
		}

		if !object.LastModified.Before(daysAgo) {
			continue
		}

		ok, err := removeObject(bucketName, object.Key, minioClient)
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
