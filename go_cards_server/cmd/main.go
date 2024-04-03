package main

import (
	router "example.com/go_cards_server/api"
	"example.com/go_cards_server/session_manager"
)

func main() {
	session_manager.InitSessionEngine()
	router.InitializeRouter()
}
