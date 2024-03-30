package main

import (
	router "example.com/server/api"
	"example.com/server/session_manager"
)

func main() {
	session_manager.InitSessionEngine()
	router.InitializeRouter()
}
