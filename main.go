package main

import (
	loadEnv "minio-test/config"
	minioUtils "minio-test/minio"
)

var (
	useSSL     = false
	recursive  = true
)

func main() {
	keys := loadEnv.Load()
	minioClient := minioUtils.ConnectMinio(keys.Endpoint, keys.AccessKeyId, keys.SecretAccessKey, useSSL)

	minioUtils.RemoveObjectsBeforeHour(keys.BucketName, keys.RetentionPeriod, recursive, minioClient)
}
