package domain

type Progress struct {
	QuestTitle     string
	CompletedTasks int
	TotalTasks     int
	Percentage     float64 // 0.0 - 100.0
}
