package domain

import (
	"reflect"
	"testing"
)

func TestNewQuestRegistry(t *testing.T) {
	tests := []struct {
		name string
		want *QuestRegistry
	}{
		{
			name: "initialize empty quest registry",
			want: &QuestRegistry{quests: make(map[string]*Quest)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQuestRegistry(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQuestRegistry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestRegistry_AddQuest(t *testing.T) {
	q1 := &Quest{ID: "q1", Title: "Quest 1"}

	type fields struct {
		quests map[string]*Quest
	}
	type args struct {
		quest *Quest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "add new quest successfully",
			fields: fields{quests: make(map[string]*Quest)},
			args:   args{quest: q1},
			want:   true,
		},
		{
			name:   "fail to add quest with duplicate ID",
			fields: fields{quests: map[string]*Quest{"q1": q1}},
			args:   args{quest: q1},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qr := &QuestRegistry{
				quests: tt.fields.quests,
			}
			if got := qr.AddQuest(tt.args.quest); got != tt.want {
				t.Errorf("QuestRegistry.AddQuest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestRegistry_RemoveQuest(t *testing.T) {
	q1 := &Quest{ID: "q1", Title: "Quest 1"}

	type fields struct {
		quests map[string]*Quest
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "remove existing quest",
			fields: fields{quests: map[string]*Quest{"q1": q1}},
			args:   args{id: "q1"},
			want:   true,
		},
		{
			name:   "remove non-existing quest",
			fields: fields{quests: map[string]*Quest{"q1": q1}},
			args:   args{id: "q2"},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qr := &QuestRegistry{
				quests: tt.fields.quests,
			}
			if got := qr.RemoveQuest(tt.args.id); got != tt.want {
				t.Errorf("QuestRegistry.RemoveQuest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestRegistry_GetQuestById(t *testing.T) {
	q1 := &Quest{ID: "q1", Title: "Quest 1"}

	type fields struct {
		quests map[string]*Quest
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Quest
		want1  bool
	}{
		{
			name:   "quest found",
			fields: fields{quests: map[string]*Quest{"q1": q1}},
			args:   args{id: "q1"},
			want:   q1,
			want1:  true,
		},
		{
			name:   "quest not found",
			fields: fields{quests: map[string]*Quest{"q1": q1}},
			args:   args{id: "q2"},
			want:   nil,
			want1:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qr := &QuestRegistry{
				quests: tt.fields.quests,
			}
			got, got1 := qr.GetQuestById(tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QuestRegistry.GetQuestById() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("QuestRegistry.GetQuestById() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestQuestRegistry_ListQuests(t *testing.T) {
	q1 := &Quest{ID: "q1", Title: "Quest 1"}

	type fields struct {
		quests map[string]*Quest
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Quest
	}{
		{
			name:   "empty registry returns empty slice",
			fields: fields{quests: make(map[string]*Quest)},
			want:   []*Quest{},
		},
		{
			name:   "registry with single quest",
			fields: fields{quests: map[string]*Quest{"q1": q1}},
			want:   []*Quest{q1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qr := &QuestRegistry{
				quests: tt.fields.quests,
			}
			if got := qr.ListQuests(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QuestRegistry.ListQuests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestRegistry_CountQuests(t *testing.T) {
	q1 := &Quest{ID: "q1"}
	q2 := &Quest{ID: "q2"}

	type fields struct {
		quests map[string]*Quest
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "zero quests",
			fields: fields{quests: make(map[string]*Quest)},
			want:   0,
		},
		{
			name:   "two quests",
			fields: fields{quests: map[string]*Quest{"q1": q1, "q2": q2}},
			want:   2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qr := &QuestRegistry{
				quests: tt.fields.quests,
			}
			if got := qr.CountQuests(); got != tt.want {
				t.Errorf("QuestRegistry.CountQuests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestRegistry_FindByDifficulty(t *testing.T) {
	qEasy := &Quest{ID: "q1", Difficulty: 2}
	qHard := &Quest{ID: "q2", Difficulty: 8}

	type fields struct {
		quests map[string]*Quest
	}
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Quest
	}{
		{
			name:   "no quests match difficulty range",
			fields: fields{quests: map[string]*Quest{"q1": qEasy, "q2": qHard}},
			args:   args{min: 4, max: 6},
			want:   []*Quest{},
		},
		{
			name:   "one quest matches difficulty range",
			fields: fields{quests: map[string]*Quest{"q1": qEasy, "q2": qHard}},
			args:   args{min: 1, max: 3},
			want:   []*Quest{qEasy},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qr := &QuestRegistry{
				quests: tt.fields.quests,
			}
			if got := qr.FindByDifficulty(tt.args.min, tt.args.max); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QuestRegistry.FindByDifficulty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestRegistry_SortedByDifficulty(t *testing.T) {
	qLow := &Quest{ID: "q1", Difficulty: 2}
	qMid := &Quest{ID: "q2", Difficulty: 5}
	qHigh := &Quest{ID: "q3", Difficulty: 9}

	type fields struct {
		quests map[string]*Quest
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Quest
	}{
		{
			name:   "empty registry sorted",
			fields: fields{quests: make(map[string]*Quest)},
			want:   []*Quest{},
		},
		{
			name:   "quests are properly sorted by unique difficulties",
			fields: fields{quests: map[string]*Quest{"q2": qMid, "q3": qHigh, "q1": qLow}},
			want:   []*Quest{qLow, qMid, qHigh},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qr := &QuestRegistry{
				quests: tt.fields.quests,
			}
			if got := qr.SortedByDifficulty(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QuestRegistry.SortedByDifficulty() = %v, want %v", got, tt.want)
			}
		})
	}
}
