package models

import "math/rand"

type Card struct {
	Suit  string
	Value string
}

type Deck struct {
	Cards []Card
}

func (d *Deck) Shuffle() {
	for i := range d.Cards {
		j := i + rand.Intn(len(d.Cards)-i)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
}

func (d *Deck) Draw() Card {
	var drawnCard Card = d.Cards[0]
	d.Cards = d.Cards[1:]
	return drawnCard
}

func (d *Deck) DrawNCards(n uint16) []Card {
	var drawnCards []Card = d.Cards[:n]
	d.Cards = d.Cards[n:]
	return drawnCards
}

func (d *Deck) AddCard(c Card) {
	d.Cards = append(d.Cards, c)
}

func (d *Deck) AddCards(cards []Card) {
	d.Cards = append(d.Cards, cards...)
}

func (d *Deck) GetCardCount() int {
	return len(d.Cards)
}

func (d *Deck) GetCards() []Card {
	return d.Cards
}
