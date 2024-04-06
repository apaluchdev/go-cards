package main

import (
	router "example.com/go_cards_server/api"
	"example.com/go_cards_server/sessionmgr"
)

func main() {
	sessionmgr.InitSessionEngine()
	router.InitializeRouter()
}
