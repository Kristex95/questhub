package main

type QuestRegistry struct {
	quests map[string]*Quest
}

func NewQuestRegistry() *QuestRegistry {
	return &QuestRegistry{
		quests: make(map[string]*Quest),
	}
}

func (qr *QuestRegistry) AddQuest(quest *Quest) bool {
	_, found := qr.GetQuestById(quest.ID)
	if !found {
		qr.quests[quest.ID] = quest
		return true
	}
	return false
}

func (qr *QuestRegistry) RemoveQuest(id string) bool {
	_, found := qr.GetQuestById(id)
	if found {
		delete(qr.quests, id)
		return true
	}
	return false
}

func (qr *QuestRegistry) GetQuestById(id string) (*Quest, bool) {
	quest, ok := qr.quests[id]
	return quest, ok
}

func (qr *QuestRegistry) ListQuests() []*Quest {
	quests := make([]*Quest, 0, len(qr.quests))

	for _, quest := range qr.quests {
		quests = append(quests, quest)
	}
	return quests
}

func (qr *QuestRegistry) CountQuests() int {
	return len(qr.quests)
}

func (qr *QuestRegistry) FindByDifficulty(min, max int) []*Quest {
	validQuests := make([]*Quest, 0, qr.CountQuests())
	for _, q := range qr.ListQuests() {
		if q.Difficulty >= min && q.Difficulty <= max {
			validQuests = append(validQuests, q)
		}
	}
	return validQuests
}

func (qr *QuestRegistry) SortedByDifficulty() []*Quest {
	quests := qr.ListQuests()
	n := len(quests)
	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-i-1; j++ {
			if quests[j].Difficulty > quests[j+1].Difficulty {
				quests[j], quests[j+1] = quests[j+1], quests[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
	return quests
}
