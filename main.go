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
	minioClient := minioUtils.ConnectMinio(endpoint, keys.AccessKeyId, keys.SecretAccessKey, useSSL)

	minioUtils.RemoveObjectsBeforeWeek(lokiChunks, recursive, minioClient)
}
