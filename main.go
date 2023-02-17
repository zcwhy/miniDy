package main

import (
	"miniDy/router"
)

func main() {
	r := router.InitRouter()

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
