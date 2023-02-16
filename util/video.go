package util

import "fmt"

func GetFileURL(fileName string) string {
	fileURL := fmt.Sprintf("https://192.168.0.102:8080/static/%s", fileName)
	return fileURL
}
