package main

import (
	loadEnv "minio-test/config"
	minioUtils "minio-test/minio"
)

var (
    endpoint   = "<your minio endpoint>"
	useSSL     = false
	lokiChunks = "chunks"
	recursive  = true
)

func main() {
	keys := loadEnv.Load()
	minioClient := minioUtils.ConnectMinio(keys.Endpoint, keys.AccessKeyId, keys.SecretAccessKey, useSSL)

	minioUtils.RemoveObjectsBeforeWeek(keys.BucketName, recursive, minioClient)
}
