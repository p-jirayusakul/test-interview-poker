package manual

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type Card struct {
	Suit  string // Spade, Heart, Diamond, Club
	Rank  string // 2, 3, 4, 5, 6, 7, 8, 9, 10, J, Q, K, A
	Value int    // 2-14
}

type Player struct {
	ID   int
	Hand []Card
}

const (
	HighCard      = 0
	OnePair       = 1
	TwoPair       = 2
	ThreeOfKind   = 3
	Straight      = 4
	Flush         = 5
	FullHouse     = 6
	FourOfKind    = 7
	StraightFlush = 8
	RoyalFlush    = 9
)

func newDeck() []Card {
	deck := make([]Card, 0, 52)
	for _, suit := range []string{"Spade", "Heart", "Diamond", "Club"} {
		for value := 2; value <= 14; value++ {
			var rank string
			if value >= 2 && value <= 10 {
				rank = strconv.Itoa(value)
				deck = append(deck, Card{Suit: suit, Rank: rank, Value: value})
				continue
			}

			switch value {
			case 11:
				rank = "J"
			case 12:
				rank = "Q"
			case 13:
				rank = "K"
			default:
				rank = "Ace"
			}
			deck = append(deck, Card{Suit: suit, Rank: rank, Value: value})
		}
	}
	return deck
}

func shuffle(cards []Card) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := len(cards) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

func Deal(deck []Card, players []Player, cardsPerPlayer int) ([]Player, []Card) {

	var i, j int
	for i = 0; i < len(deck); i++ {

		// ถ้าผู้เล่นทุกคนได้รับจำนวนการ์ดที่กำหนดแล้ว
		if (len(players) * cardsPerPlayer) == i {
			break
		}

		// ตรวจสอบว่าผู้เล่นมีการ์ดครบตามจำนวนที่กำหนดหรือไม่
		if len(players[j].Hand) != cardsPerPlayer {
			players[j].Hand = append(players[j].Hand, deck[i])
		}

		// แจกการ์ดวนจนครบรอบผู้เล่นแล้วหรือไม่ ถ้าครบรอบแล้ว ให้เริ่มใหม่ที่คนแรก (index 0)
		j++
		if j == len(players) {
			j = 0
		}

	}

	return players, deck[i:]
}

func evaluateHand(hand []Card) int {
	sort.Slice(hand, func(i, j int) bool {
		return hand[i].Value < hand[j].Value
	})

	isFlush := checkFlush(hand)
	isStraight, lastCardValue := checkStraight(hand)

	// 14 = Ace
	if isFlush && isStraight && lastCardValue == 14 {
		return RoyalFlush
	}

	if isFlush && isStraight {
		return StraightFlush
	}

	return 0
}

func checkFlush(hand []Card) bool {
	for i := 1; i < len(hand); i++ {
		if hand[i].Suit != hand[0].Suit {
			return false
		}
	}
	return true
}

func checkStraight(hand []Card) (bool, int) {
	for i := 0; i < len(hand); i++ {
		if i+1 < len(hand) {
			if hand[i].Value+1 != hand[i+1].Value {
				return false, hand[len(hand)-1].Value
			}
		}
	}
	return true, hand[len(hand)-1].Value
}

func printHand(hand []Card) string {
	result := "("
	for i, card := range hand {
		result += fmt.Sprintf("%s %s", card.Suit, card.Rank)
		if i != len(hand)-1 {
			result += ", "
		}
	}

	return result + ")"
}

func RunManual() {
	deck := newDeck()
	shuffle(deck)

	players := []Player{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4},
	}

	players, deck = Deal(deck, players, 5)
	for _, p := range players {

		result := evaluateHand(p.Hand)
		resultStr := ""
		switch result {
		case RoyalFlush:
			resultStr = "Royal Flush"
		case StraightFlush:
			resultStr = "Straight Flush"
		}

		fmt.Printf("Player %d: %s -> %s\n", p.ID, printHand(p.Hand), resultStr)
	}

	fmt.Println("Deck:", len(deck))
}
