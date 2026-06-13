package main

import "fmt"

func main() {
	// --- PART 6 ---
	user := NewUser("12345678", "Ezio Auditore da Firenze", "ezio@assassin.com")
	user.PrintInfo()
	fmt.Println()

	// REWARD
	r1 := NewReward("r1", "Bronze Sword", 50, "common")
	r2 := NewReward("r2", "Silver Blade", 120, "rare")
	r3 := NewReward("r3", "Assassin Hood", 300, "epic")

	fmt.Println(r1.String())
	fmt.Println(r2.String())
	fmt.Println(r3.String())

	// QUEST REGISTRY
	reg := NewQuestRegistry()
	// QUEST 1
	q1 := NewQuest("q1", "Training", "Basic training", 2)
	q1.AddTask(NewTask("t1", "Run 1km", 50))
	q1.AddTask(NewTask("t2", "Climb tower", 70))
	// QUEST 2
	q2 := NewQuest("q2", "Stealth Mission", "Infiltration practice", 5)
	q2.AddTask(NewTask("t3", "Scout area", 80))
	q2.AddTask(NewTask("t4", "Avoid guards", 120))
	q2.AddTask(NewTask("t5", "Steal document", 150))
	// QUEST 3
	q3 := NewQuest("q3", "Assassination", "Eliminate target", 9)
	q3.AddTask(NewTask("t6", "Find target", 200))
	q3.AddTask(NewTask("t7", "Prepare weapons", 180))
	q3.AddTask(NewTask("t8", "Execute mission", 400))
	q3.AddTask(NewTask("t9", "Escape city", 250))

	// REGISTER QUESTS
	reg.AddQuest(q1)
	reg.AddQuest(q2)
	reg.AddQuest(q3)

	fmt.Println("\nTotal quests:", reg.CountQuests())
	for _, q := range reg.ListQuests() {
		fmt.Println(q.Summary())
	}

	// === PROGRESS SIMULATION ===
	fmt.Println("\n--- Completing tasks ---")
	// Quest 1 progress
	q1.CompleteTask("t1")
	fmt.Println("Q1 progress:", q1.GetProgressPercentage(), "%")
	q1.CompleteTask("t2")
	fmt.Println("Q1 completed tasks:")
	for _, task := range q1.GetCompletedTasks() {
		task.Print()
	}
	fmt.Println("Q1 remaining:", q1.GetRemainingTasks())
	// Quest 2 progress
	q2.CompleteTask("t3")
	q2.CompleteTask("t4")
	fmt.Println("Q2 progress:", q2.GetProgressPercentage(), "%")
	// Quest 3 partial
	q3.CompleteTask("t6")
	q3.CompleteTask("t7")
	fmt.Println("Q3 progress:", q3.GetProgressPercentage(), "%")

	// === COMPLETE QUESTS ===
	fmt.Println("\n--- Completing quests ---")
	user.CompleteQuest(q1)
	user.CompleteQuest(q2)
	user.CompleteQuest(q3)

	// === FINAL STATE ===
	fmt.Println("\n\n=== FINAL STATE ===")
	user.PrintInfo()
	fmt.Println()
	for _, q := range reg.ListQuests() {
		fmt.Println(q.Summary())
	}

	// --- PART 7 ---
	fmt.Println("\nPart 7. Additional tasks")
	fmt.Println("- q2 Deactivation")
	q2.Deactivate()
	fmt.Println(q2.Summary())
	fmt.Println("- q2 Activation")
	q2.Activate()
	fmt.Println(q2.Summary())
	q2.Summary()

	fmt.Println("\n Removing quest q2")
	reg.RemoveQuest("q2")
	for _, q := range reg.ListQuests() {
		fmt.Println(q.Summary())
	}
	_, ok := reg.GetQuestById("q2")
	fmt.Println("q2 in registry:", ok)

	// COMPLETING Quest2
	q3.CompleteTask("t8")
	q3.CompleteTask("t9")
	user.CompleteQuest(q3)

	// ADDING DUPLICATE
	falseQuest := NewQuest("q3", "False quest", "False", 1)
	fmt.Println("Trying to add duplicate quest:", reg.AddQuest(falseQuest))

	// --- PART 8 ---
	// FindByDifficulty
	fmt.Println("\nQuests by difficulty (2..6):")
	filtered := reg.FindByDifficulty(2, 6)
	for _, q := range filtered {
		fmt.Println(q.Summary())
	}

	// SortedByDifficulty
	fmt.Println("\nSorted by difficulty:")
	sorted := reg.SortedByDifficulty()
	for _, q := range sorted {
		fmt.Println(q.ID, q.Title, "diff:", q.Difficulty)
	}

	// User statistics
	fmt.Println("\nUser stats:")
	fmt.Println(user.Stats())

	// Reward values
	fmt.Println("\nReward values:")
	fmt.Printf("%s value: %d\n", r1.Title, r1.Value())
	fmt.Printf("%s value: %d\n", r2.Title, r2.Value())
	fmt.Printf("%s value: %d\n", r3.Title, r3.Value())

	fmt.Println()
	user.PrintInfo()
}
