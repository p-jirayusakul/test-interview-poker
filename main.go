package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Suit string
type Rank int

type Card struct {
	Suit Suit
	Rank Rank
}

type Player struct {
	ID   int
	Hand []Card
}

type HandRank int

const (
	HighCard HandRank = iota
	OnePair
	TwoPair
	ThreeOfKind
	Straight
	Flush
	FullHouse
	FourOfKind
	StraightFlush
	RoyalFlush
)

type Result struct {
	PlayerID int
	Rank     HandRank
	Score    []int // ใช้ tie-breaker
}

// ------------------ Deck ------------------

// NewDeck สร้างไพ่ 52 ใบ (4 ดอก x 13 หน้า)
func NewDeck() []Card {
	suits := []Suit{"Spades", "Hearts", "Diamonds", "Clubs"}
	var deck []Card

	for _, s := range suits {
		for r := 2; r <= 14; r++ {
			deck = append(deck, Card{Suit: s, Rank: Rank(r)})
		}
	}
	return deck
}

// Shuffle สับไพ่แบบสุ่มใน memory
func Shuffle(deck []Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

// Deal แจกไพ่ให้ผู้เล่นคนละ cardsPerPlayer ใบ
func Deal(deck []Card, players []Player, cardsPerPlayer int) ([]Player, []Card) {
	index := 0

	for i := range players {
		players[i].Hand = deck[index : index+cardsPerPlayer]
		index += cardsPerPlayer
	}

	return players, deck[index:]
}

// ------------------ Evaluate ------------------

type rankGroup struct {
	Rank  int
	Count int
}

// EvaluateHand วิเคราะห์มือ 5 ใบ -> (ประเภทไพ่, คะแนนสำหรับเทียบ)
func EvaluateHand(hand []Card) (HandRank, []int) {
	// เรียง rank จากน้อยไปมาก เพื่อเช็ก straight ง่าย
	sort.Slice(hand, func(i, j int) bool {
		return hand[i].Rank < hand[j].Rank
	})

	isFlush := checkFlush(hand)                     // ดอกเดียวกันหมดไหม
	isStraight, highStraight := checkStraight(hand) // เรียงติดไหม (รวม A-2-3-4-5)

	groups := buildGroups(hand) // นับจำนวนหน้าไพ่ซ้ำ

	// จัดอันดับจากสูงสุด -> ต่ำสุด
	if isFlush && isStraight && highStraight == 14 {
		return RoyalFlush, []int{14}
	}

	if isFlush && isStraight {
		return StraightFlush, []int{highStraight}
	}

	if groups[0].Count == 4 {
		return FourOfKind, buildScore(groups)
	}

	if groups[0].Count == 3 && groups[1].Count == 2 {
		return FullHouse, buildScore(groups)
	}

	if isFlush {
		return Flush, extractRanksDesc(hand)
	}

	if isStraight {
		return Straight, []int{highStraight}
	}

	if groups[0].Count == 3 {
		return ThreeOfKind, buildScore(groups)
	}

	if groups[0].Count == 2 && groups[1].Count == 2 {
		return TwoPair, buildScore(groups)
	}

	if groups[0].Count == 2 {
		return OnePair, buildScore(groups)
	}

	return HighCard, extractRanksDesc(hand)
}

// ------------------ Helpers ------------------

// checkFlush ไพ่ทั้ง 5 ใบ ดอกเดียวกันไหม
func checkFlush(hand []Card) bool {
	for i := 1; i < len(hand); i++ {
		if hand[i].Suit != hand[0].Suit {
			return false
		}
	}
	return true
}

// checkStraight rank เรียงติดกันไหม
// รองรับ A-2-3-4-5 โดยคืน high = 5
func checkStraight(hand []Card) (bool, int) {
	isStraight := true
	for i := 1; i < len(hand); i++ {
		if hand[i].Rank != hand[i-1].Rank+1 {
			isStraight = false
			break
		}
	}

	if isStraight {
		return true, int(hand[len(hand)-1].Rank)
	}

	var ranks []int
	for _, c := range hand {
		ranks = append(ranks, int(c.Rank))
	}
	sort.Ints(ranks)

	if ranks[0] == 2 &&
		ranks[1] == 3 &&
		ranks[2] == 4 &&
		ranks[3] == 5 &&
		ranks[4] == 14 {
		return true, 5
	}

	return false, 0
}

// buildGroups นับ rank ว่าซ้ำกี่ใบ แล้ว sort ตามความสำคัญ
// สำคัญกว่า = count เยอะกว่า, ถ้า count เท่ากัน rank สูงกว่า
func buildGroups(hand []Card) []rankGroup {
	m := map[int]int{}
	for _, c := range hand {
		m[int(c.Rank)]++
	}

	var groups []rankGroup
	for r, c := range m {
		groups = append(groups, rankGroup{Rank: r, Count: c})
	}

	sort.Slice(groups, func(i, j int) bool {
		if groups[i].Count != groups[j].Count {
			return groups[i].Count > groups[j].Count
		}
		return groups[i].Rank > groups[j].Rank
	})

	return groups
}

// buildScore แปลง groups เป็น score สำหรับ tie-breaker
// เช่น FullHouse 10 over 7 -> [10,10,10,7,7]
func buildScore(groups []rankGroup) []int {
	var score []int
	for _, g := range groups {
		for i := 0; i < g.Count; i++ {
			score = append(score, g.Rank)
		}
	}
	return score
}

// extractRanksDesc ใช้กับ Flush/HighCard
func extractRanksDesc(hand []Card) []int {
	var ranks []int
	for _, c := range hand {
		ranks = append(ranks, int(c.Rank))
	}
	sort.Sort(sort.Reverse(sort.IntSlice(ranks)))
	return ranks
}

// ------------------ Compare ------------------

// Compare >0 แปลว่า a ชนะ b, <0 แปลว่า b ชนะ a, =0 เสมอ
func Compare(a, b Result) int {
	if a.Rank != b.Rank {
		return int(a.Rank - b.Rank)
	}

	for i := range a.Score {
		if a.Score[i] != b.Score[i] {
			return a.Score[i] - b.Score[i]
		}
	}
	return 0
}

// FindWinners หาผู้ชนะสูงสุด (รองรับเสมอหลายคน)
func FindWinners(results []Result) []Result {
	best := results[0]
	winners := []Result{best}

	for i := 1; i < len(results); i++ {
		cmp := Compare(results[i], best)

		if cmp > 0 {
			best = results[i]
			winners = []Result{results[i]}
		} else if cmp == 0 {
			winners = append(winners, results[i])
		}
	}

	return winners
}

// ------------------ Utils ------------------

func rankToString(r HandRank) string {
	names := []string{
		"High Card", "One Pair", "Two Pair", "Three of a Kind",
		"Straight", "Flush", "Full House", "Four of a Kind",
		"Straight Flush", "Royal Flush",
	}
	return names[r]
}

// ------------------ Main ------------------

func main() {
	// STEP 1: สร้างและสับไพ่
	deck := NewDeck()
	Shuffle(deck)

	// STEP 2: สร้างผู้เล่น 4 คน
	players := []Player{
		{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4},
	}

	// STEP 3: แจกคนละ 5 ใบ
	players, deck = Deal(deck, players, 5)

	var results []Result

	for _, p := range players {
		rank, score := EvaluateHand(p.Hand)

		results = append(results, Result{
			PlayerID: p.ID,
			Rank:     rank,
			Score:    score,
		})

		fmt.Printf("Player %d: %v -> %s\n",
			p.ID, p.Hand, rankToString(rank))
	}

	// STEP 5: แสดงไพ่คงเหลือ
	fmt.Printf("Cards left in deck: %d\n", len(deck))

	// STEP 6: หา winner
	winners := FindWinners(results)

	// STEP 7: แสดงผู้ชนะ
	fmt.Print("*** Winner: ")
	for _, w := range winners {
		fmt.Printf("Player %d ", w.PlayerID)
	}
	fmt.Printf("with %s ***\n", rankToString(winners[0].Rank))
}
