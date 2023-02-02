package main

import (
	"miniDy/router"
)

func main() {
	r := router.InitRouter()

	r.Run(":8080")
}
