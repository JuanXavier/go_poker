package p2p

type GameStatus int32

const (
	GameStatusWaitingForCards GameStatus = iota
	GameStatusReceivingCards
	GameStatusDealing
	GameStatusPreFlop
	GameStatusFlop
	GameStatusTurn
	GameStatusRiver
)

func (g GameStatus) String() string {
	switch g {
	case GameStatusWaitingForCards:
		return "WAITING FOR CARDS"
	case GameStatusReceivingCards:
		return "RECEIVING CARDS"
	case GameStatusDealing:
		return "DEALING"
	case GameStatusPreFlop:
		return "Pre-flop"
	case GameStatusFlop:
		return "Flop"
	case GameStatusTurn:
		return "Turn"
	case GameStatusRiver:
		return "River"
	default:
		return "Unknown"
	}
}
