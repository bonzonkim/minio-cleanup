package loadEnv

import (
	"fmt"
	"os"
)


type AccessKeys struct {
	AccessKeyId     string
	SecretAccessKey string
}
func Load() *AccessKeys {
	Keys := &AccessKeys {
		AccessKeyId:     os.Getenv("ACCESSKEYID"),
		SecretAccessKey: os.Getenv("SECRETACCESSKEY"),
	}

	fmt.Println("Successfully loaded env")
	return Keys
}
