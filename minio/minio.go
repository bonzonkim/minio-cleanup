package minioUtils

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Connect to MiniO
func ConnectMinio(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) *minio.Client {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
		return nil
	}
	fmt.Println("Succefully connected MiniO")
	return minioClient
}



// remove object in the given bucket
func removeObject(ctx context.Context, bucketName string, objectName string, minioClient *minio.Client) (bool, error) {
	if err := minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{}); err != nil {
		return false, err
	}
	return true, nil
}

// listing all objects in the storage
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

// remove objects older than 'retentionPeriod' in the given bucket 
func RemoveObjectsBeforeHour(bucketName string, retentionPeriod int, recursive bool, minioClient *minio.Client) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	hoursAgoFromNow := time.Now().Add(-time.Duration(retentionPeriod) * time.Hour)
	objectCh := listObjects(ctx, bucketName, recursive, minioClient)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var sum int

	skippedFiles := map[string]bool {
		"loki_cluster_seed.json":		true,
		"tempo_cluster_seed.json":		true,
	}

	for object := range objectCh {
		goObj := object
		if skippedFiles[object.Key] {
			log.Printf("Skipping seed file %s :\n", object.Key)
			continue
		}

		if !object.LastModified.Before(hoursAgoFromNow) {
			continue
		}

		wg.Add(1)
		go func(obj minio.ObjectInfo) {
			defer wg.Done()
			ok, err := removeObject(ctx, bucketName, obj.Key, minioClient)
			if err != nil {
				log.Printf("Error occurred while deleting %s : %v\n", obj.Key, err)
				return
			}
			if ok {
				mu.Lock()
				fmt.Printf("Successfully deleted %s | %v\n", obj.Key, obj.LastModified)
				sum++
				mu.Unlock()
			}
		}(goObj)

	}
	wg.Wait()
	fmt.Printf("SUCCESSFULLY DELETED ALL OBJECTS OLDER THAN %d DAYS EXCEPT SEED FILE | DELETED %d files", retentionPeriod, sum)
}
