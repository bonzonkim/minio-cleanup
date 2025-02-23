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

// Global MinIO client (singleton)
var (
	minioClient *minio.Client
	once        sync.Once
)

// ConnectMinio initializes the MinIO connection
func ConnectMinio(endpoint, accessKeyID, secretAccessKey string, useSSL bool) {
	once.Do(func() {
		var err error
		minioClient, err = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Fatalf("Failed to connect to MinIO: %v", err)
		}
		fmt.Println("✅ Successfully connected to MinIO")
	})
}

// removeObjectsBatch deletes a batch of objects concurrently
func removeObjectsBatch(ctx context.Context, bucketName string, objects []string) error {
	objectCh := make(chan minio.ObjectInfo, len(objects))

	// Fill objectCh with objects to delete
	go func() {
		defer close(objectCh)
		for _, obj := range objects {
			objectCh <- minio.ObjectInfo{Key: obj}
		}
	}()

	// Perform batch deletion
	opts := minio.RemoveObjectsOptions{GovernanceBypass: true}
	for err := range minioClient.RemoveObjects(ctx, bucketName, objectCh, opts) {
		if err.Err != nil {
			log.Printf("❌ Failed to delete %s: %v", err.ObjectName, err.Err)
		} else {
			log.Printf("✅ Successfully delete %s\n", err.ObjectName)
		}
	}
	return nil
}

// RemoveObjectsBeforeHour deletes objects older than the retention period
func RemoveObjectsBeforeHour(bucketName string, retentionPeriod int, parallelism int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Define the deletion threshold time
	hoursAgoFromNow := time.Now().Add(-time.Duration(retentionPeriod) * time.Hour)

	// Worker pool setup
	var wg sync.WaitGroup
	sem := make(chan struct{}, parallelism) // Worker limit
	objectCh := make(chan string, 10000)    // Channel for objects to delete

	// List and filter objects concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(objectCh)

		objectStream := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
			Recursive: true,
		})

		// Files to not delete
		excludedFiles := map[string]bool{
			"loki_cluster_seed.json": true,
			"tempo_cluster_seed.json": true,
		}

		for {
			select {
			case <-ctx.Done():
				log.Println("Context Canceled. Stopping object listing.")
				return
			case object, ok := <-objectStream:
				if !ok {
					return
				}
				if object.Err != nil {
					log.Printf("⚠️ Error listing object %s: %v", object.Key, object.Err)
					continue
				}
				if excludedFiles[object.Key] {
					log.Printf("Skpping %s", object.Key)
					continue
				}
				if object.LastModified.Before(hoursAgoFromNow) {
					objectCh <- object.Key
				}
			}
		}
	}()

	// Worker pool for batch deletions
	const batchSize = 5000
	batch := make([]string, 0, batchSize)

	for obj := range objectCh {
		batch = append(batch, obj)

		// Process batch when it reaches the batch size
		if len(batch) >= batchSize {
			wg.Add(1)
			sem <- struct{}{} // Acquire worker slot

			go func(objects []string) {
				defer wg.Done()
				defer func() { <-sem }() // Release worker slot

				if err := removeObjectsBatch(ctx, bucketName, objects); err != nil {
					log.Printf("⚠️ Error deleting batch: %v", err)
				} else {
					log.Printf("✅ Successfully deleted %d files", len(objects))
				}
			}(batch)

			batch = nil // Reset batch
		}
	}

	// Process any remaining objects in the final batch
	if len(batch) > 0 {
		wg.Add(1)
		sem <- struct{}{}

		go func(objects []string) {
			defer wg.Done()
			defer func() { <-sem }()

			if err := removeObjectsBatch(ctx, bucketName, objects); err != nil {
				log.Printf("⚠️ Error deleting batch: %v", err)
			} else {
				log.Printf("✅ Successfully deleted %d files", len(objects))
			}
		}(batch)
	}

	wg.Wait()
	fmt.Printf("🎯 SUCCESSFULLY DELETED FILES OLDER THAN %d HOURS\n", retentionPeriod)
}
