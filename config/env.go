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

// Load required variables
func Load() *AccessKeys {
	requiredVars := []string{"ENDPOINT", "ACCESSKEYID", "SECRETACCESSKEY", "BUCKETNAME", "RETENTIONPERIOD"}
	// map to store all the envs
	env := make(map[string]string)

	for _, key := range requiredVars {
		value := os.Getenv(key)
		if value == "" {
			log.Fatalf("%s is not set or empty string", value)
		}
		env[key] = value
	}


	retenteionPeriod, err := strconv.Atoi(env["RETENTIONPERIOD"])
	if err != nil {
		log.Fatalf("Invalid RETENTIONPERIOD value %v", err)
	}
	

	fmt.Println("Successfully loaded env")

	return &AccessKeys {
		Endpoint:			env["ENDPOINT"],
		AccessKeyId:		env["ACCESSKEYID"],
		SecretAccessKey:	env["SECRETACCESSKEY"],
		BucketName:			env["BUCKETNAME"],	
		RetentionPeriod:	retenteionPeriod,
	}
}
