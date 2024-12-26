package loadEnv

import (
	"fmt"
	"log"
	"os"
	"strconv"
)


type AccessKeys struct {
	Endpoint			string
	AccessKeyId			string
	SecretAccessKey		string
	BucketName			string
	RetentionPeriod		int
}

func Load() *AccessKeys {
	reteiontionStr := os.Getenv("RETENTIONPERIOD")
	retentionPeriod, err := strconv.ParseInt(reteiontionStr, 10, strconv.IntSize)
	if err != nil {
		log.Fatalln("failed to parse retentionPeriod to int")
	}
	Keys := &AccessKeys {
		Endpoint:			os.Getenv("ENDPOINT"),
		AccessKeyId:		os.Getenv("ACCESSKEYID"),
		SecretAccessKey: 	os.Getenv("SECRETACCESSKEY"),
		BucketName:		 	os.Getenv("BUCKETNAME"),
		RetentionPeriod: 	int(retentionPeriod),
	}

	fmt.Println("Successfully loaded env")
	return Keys
}
