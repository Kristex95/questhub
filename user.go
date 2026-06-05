package main

import "fmt"

type User struct {
	ID       string
	Username string
	Email    string
	Level    int
	XP       int
	CompletedQuests []*Quest
	TotalXPEarned   int
}

func NewUser(id, username, email string) *User {
	return &User{
		ID:       id,
		Username: username,
		Email:    email,
		Level:    1,
		XP:       0,
		TotalXPEarned: 0,
	}
}

func (u *User) AddXP(amount int) {
	if amount <= 0 {
		return
	}

	fmt.Printf("%s gained %d XP\n", u.Username, amount)
	u.TotalXPEarned += amount

	for amount > 0 {
		requiredXP := u.Level * 1000
		needed := requiredXP - u.XP

		if amount >= needed {
			u.LevelUp()
			amount -= needed

			fmt.Printf("%s leveled up to level %d!\n", u.Username, u.Level)
		} else {
			u.XP += amount
			amount = 0
		}
	}
}

func (u *User) LevelUp() {
	u.Level += 1
	u.XP = 0
}

func (u *User) CompleteQuest(quest *Quest) {
	completed := quest.IsCompleted()
	if !completed {
		fmt.Printf("Quest %s is not completed\n", quest.Title)
		return
	}
	u.CompletedQuests = append(u.CompletedQuests, quest)
	u.AddXP(quest.TotalXP())
}

func (u *User) PrintInfo() {
	fmt.Println("=== User Profile ===")
	fmt.Printf("ID: %s\n", u.ID)
	fmt.Printf("Username: %s\n", u.Username)
	fmt.Printf("Email: %s\n", u.Email)
	fmt.Printf("Level: %d\n", u.Level)
	fmt.Printf("XP: %d\n", u.XP)
}

func (u *User) Stats() string {
	return fmt.Sprintf(
		"=== User Stats ===\nCompleted Quests: %d\nTotal XP Earned: %d\n",
		len(u.CompletedQuests),
		u.TotalXPEarned,
	)
}
