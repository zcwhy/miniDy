package util

import "fmt"

func GetFileURL(fileName string) string {
	fileURL := fmt.Sprintf("http://192.168.138.51:8080/static/%s", fileName)
	return fileURL
}
