package main

import "fmt"

type Reward struct {
	ID       string
	Title    string
	XPAmount int
	Rarity   string
}

func NewReward(id, title string, XPAmount int, rarity string) *Reward {
	validRarities := map[string]bool{
		"common":    true,
		"rare":      true,
		"epic":      true,
		"legendary": true,
	}
	if !validRarities[rarity] {
		rarity = "common"
	}
	return &Reward{
		ID:       id,
		Title:    title,
		XPAmount: XPAmount,
		Rarity:   rarity,
	}
}

func (r *Reward) String() string {
	letter := "?"

	switch r.Rarity {
	case "common":
		letter = "C"
	case "rare":
		letter = "R"
	case "epic":
		letter = "E"
	case "legendary":
		letter = "L"
	}

	return fmt.Sprintf("[%s] %s - XP: %d (%s)", letter, r.Title, r.XPAmount, r.Rarity)
}

func (r *Reward) Value() int {
	multiplier := 1
	switch r.Rarity {
	case "common":
		multiplier = 1
	case "rare":
		multiplier = 2
	case "epic":
		multiplier = 5
	case "legendary":
		multiplier = 10
	}
	return r.XPAmount * multiplier
}
