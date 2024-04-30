package cards

import "math/rand"

type Card struct {
	Suit  string
	Value string
}

type Deck struct {
	Cards []Card
}

var CardValues = map[int]string{
	0:  "2",
	1:  "3",
	2:  "4",
	3:  "5",
	4:  "6",
	5:  "7",
	6:  "8",
	7:  "9",
	8:  "10",
	9:  "J",
	10: "Q",
	11: "K",
	12: "A",
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

func (d *Deck) DrawNCards(n int) []Card {
	if n <= len(d.Cards) {
		var drawnCards []Card = d.Cards[:n]
		d.Cards = d.Cards[n:]
		return drawnCards
	} else {
		// TODO - handle not enough cards in deck when drawing
		return nil
	}
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

func (d *Deck) GetNumberOfCardsRemaining() int {
	return len(d.Cards)
}
