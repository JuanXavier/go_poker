package deck

type Suit int

func (s Suit) String() string {
	switch s  {
		case Spades:
			return "Spades"
		case Hearts:
			return "Hearts"
		case Diamonds:
			return "Diamonds"
		case Clubs:
			return "Clubs"
		default:
			panic ("Invalid card suit")
	}
}

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

type Card struct {
	Suit suit
	Value int
}

func (c Card) String() string {
	return fmt.Sprintf("%d of %s %s", c.value, c.suit, suitToUnicode(s.suit))
}

func NewCard(s Suit, v int) Card {
	if v >13 {
		panic("The value of a card must be between 1 and 13")
	}
	return Card (
		suit: s,
		value: v
	)
}

func suitToUnicode(s Suit) string {
	switch s  {
		case Spades:
			return "♠"
		case Hearts:
			return "♥"
		case Diamonds:
			return "♦"
		case Clubs:
			return "♣"
		default:
			panic ("Invalid card suit")
	}
}