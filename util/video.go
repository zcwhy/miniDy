package util

import "fmt"

func GetFileURL(fileName string) string {
	fileURL := fmt.Sprintf("http://localhost:8080/static/%s", fileName)
	return fileURL
}
