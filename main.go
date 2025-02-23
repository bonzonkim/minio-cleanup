package main

import (
	loadEnv "minio-cleanup/config"
	minioUtils "minio-cleanup/minio"
)

var (
	useSSL     = false
	recursive  = true
)

func main() {
	keys := loadEnv.Load()
	minioUtils.ConnectMinio(keys.Endpoint, keys.AccessKeyId, keys.SecretAccessKey, useSSL)

	minioUtils.RemoveObjectsBeforeHour(keys.BucketName, keys.RetentionPeriod, 10)
}
