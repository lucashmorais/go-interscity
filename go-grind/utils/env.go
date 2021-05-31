package utils

import "os"

func GetServerURL() string {
	return "http://" + os.Getenv("SERVER_LOCATION") + ":" + os.Getenv("SERVER_PORT")
}
