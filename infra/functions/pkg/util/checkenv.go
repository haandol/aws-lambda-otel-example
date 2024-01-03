package util

import "os"

func IsLocalEnv() bool {
	return os.Getenv("AWS_LAMBDA_RUNTIME_API") == ""
}
