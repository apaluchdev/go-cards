package session

import (
	"time"

	"example.com/go_cards_server/jwthelper"
	"github.com/google/uuid"
)

var Tickets map[string]time.Time = make(map[string]time.Time)
var ClaimsForTickets map[string]*jwthelper.Claims = make(map[string]*jwthelper.Claims)
var expirationTime = time.Minute

func GetTicket() string {
	ticketId := uuid.New().String()
	Tickets[ticketId] = time.Now()
	return ticketId
}

// TODO - Call this on a timer, currently mem leaking tickets that are never used
func CleanExpiredTickets() {
	for ticketId, ticketTime := range Tickets {
		if time.Since(ticketTime) > expirationTime {
			delete(Tickets, ticketId)
		}
	}
}

func DeleteClaimsForTicket(ticketId string) {
	delete(ClaimsForTickets, ticketId)
}

func IsTicketExpired(ticketId string) bool {
	ticketTime, exists := Tickets[ticketId]
	if !exists {
		return true
	}

	return time.Since(ticketTime) > expirationTime
}
