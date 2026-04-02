package main

import "testing"

func c(s Suit, r Rank) Card {
	return Card{Suit: s, Rank: r}
}

func TestEvaluateHand_RoyalFlush(t *testing.T) {
	hand := []Card{
		c("Hearts", 10),
		c("Hearts", 11),
		c("Hearts", 12),
		c("Hearts", 13),
		c("Hearts", 14),
	}
	rank, score := EvaluateHand(hand)
	if rank != RoyalFlush {
		t.Fatalf("expected RoyalFlush, got %v", rank)
	}
	if len(score) != 1 || score[0] != 14 {
		t.Fatalf("expected score [14], got %v", score)
	}
}

func TestEvaluateHand_StraightFlush_A2345(t *testing.T) {
	hand := []Card{
		c("Spades", 14),
		c("Spades", 2),
		c("Spades", 3),
		c("Spades", 4),
		c("Spades", 5),
	}
	rank, score := EvaluateHand(hand)
	if rank != StraightFlush {
		t.Fatalf("expected StraightFlush, got %v", rank)
	}
	if len(score) != 1 || score[0] != 5 {
		t.Fatalf("expected score [5], got %v", score)
	}
}

func TestEvaluateHand_FourOfKind(t *testing.T) {
	hand := []Card{
		c("Spades", 9),
		c("Hearts", 9),
		c("Diamonds", 9),
		c("Clubs", 9),
		c("Hearts", 2),
	}
	rank, score := EvaluateHand(hand)
	if rank != FourOfKind {
		t.Fatalf("expected FourOfKind, got %v", rank)
	}
	want := []int{9, 9, 9, 9, 2}
	for i := range want {
		if score[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, score)
		}
	}
}

func TestEvaluateHand_FullHouse(t *testing.T) {
	hand := []Card{
		c("Spades", 10),
		c("Hearts", 10),
		c("Diamonds", 10),
		c("Clubs", 7),
		c("Hearts", 7),
	}
	rank, score := EvaluateHand(hand)
	if rank != FullHouse {
		t.Fatalf("expected FullHouse, got %v", rank)
	}
	want := []int{10, 10, 10, 7, 7}
	for i := range want {
		if score[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, score)
		}
	}
}

func TestCompare_OnePair_Kicker(t *testing.T) {
	// P1: pair 8 + A kicker
	a := Result{
		PlayerID: 1,
		Rank:     OnePair,
		Score:    []int{8, 8, 14, 7, 3},
	}
	// P2: pair 8 + K kicker
	b := Result{
		PlayerID: 2,
		Rank:     OnePair,
		Score:    []int{8, 8, 13, 7, 3},
	}
	if Compare(a, b) <= 0 {
		t.Fatalf("expected player a to win by kicker, got Compare=%d", Compare(a, b))
	}
}

func TestCompare_TwoPair(t *testing.T) {
	// A: two pair Aces+2
	a := Result{
		PlayerID: 1,
		Rank:     TwoPair,
		Score:    []int{14, 14, 2, 2, 9},
	}
	// B: two pair Kings+Queens
	b := Result{
		PlayerID: 2,
		Rank:     TwoPair,
		Score:    []int{13, 13, 12, 12, 14},
	}
	if Compare(a, b) <= 0 {
		t.Fatalf("expected player a to win by higher top pair, got Compare=%d", Compare(a, b))
	}
}

func TestFindWinners_SingleWinner(t *testing.T) {
	results := []Result{
		{PlayerID: 1, Rank: Straight, Score: []int{9}},
		{PlayerID: 2, Rank: Flush, Score: []int{13, 11, 9, 6, 2}},
		{PlayerID: 3, Rank: OnePair, Score: []int{14, 14, 10, 8, 2}},
		{PlayerID: 4, Rank: HighCard, Score: []int{14, 13, 11, 9, 3}},
	}
	winners := FindWinners(results)
	if len(winners) != 1 || winners[0].PlayerID != 2 {
		t.Fatalf("expected player 2 as winner, got %+v", winners)
	}
}

func TestFindWinners_Tie(t *testing.T) {
	results := []Result{
		{PlayerID: 1, Rank: Straight, Score: []int{10}},
		{PlayerID: 2, Rank: Straight, Score: []int{10}},
		{PlayerID: 3, Rank: TwoPair, Score: []int{9, 9, 4, 4, 2}},
		{PlayerID: 4, Rank: OnePair, Score: []int{14, 14, 13, 5, 2}},
	}
	winners := FindWinners(results)
	if len(winners) != 2 {
		t.Fatalf("expected 2 winners, got %+v", winners)
	}
	if winners[0].PlayerID != 1 || winners[1].PlayerID != 2 {
		t.Fatalf("expected players 1 and 2 tie, got %+v", winners)
	}
}
