package domain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewQuest_Validation(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		title       string
		description string
		difficulty  int
		wantErr     bool
		wantField   string
	}{
		{
			name:        "valid quest",
			id:          1,
			title:       "Quest",
			description: "Valid description",
			difficulty:  5,
			wantErr:     false,
		},
		{
			name:        "invalid title",
			id:          2,
			title:       "ab",
			description: "Valid description",
			difficulty:  5,
			wantErr:     true,
			wantField:   "title",
		},
		{
			name:        "invalid description",
			id:          3,
			title:       "Valid title",
			description: "",
			difficulty:  5,
			wantErr:     true,
			wantField:   "description",
		},
		{
			name:        "difficulty too low",
			id:          4,
			title:       "Valid title",
			description: "Valid description",
			difficulty:  0,
			wantErr:     true,
			wantField:   "difficulty",
		},
		{
			name:        "difficulty too high",
			id:          5,
			title:       "Valid title",
			description: "Valid description",
			difficulty:  11,
			wantErr:     true,
			wantField:   "difficulty",
		},
		{
			name:        "title and description invalid",
			id:          6,
			title:       "ab",
			description: "",
			difficulty:  5,
			wantErr:     true,
			wantField:   "title",
		},
		{
			name:        "title and difficulty invalid",
			id:          7,
			title:       "ab",
			description: "Valid description",
			difficulty:  0,
			wantErr:     true,
			wantField:   "title",
		},
		{
			name:        "all invalid",
			id:          8,
			title:       "a",
			description: "",
			difficulty:  0,
			wantErr:     true,
			wantField:   "title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewQuest(tt.id, tt.title, tt.description, tt.difficulty)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, q)

				var ve *ValidationError
				require.True(t, errors.As(err, &ve))
				require.Equal(t, tt.wantField, ve.Field)
			} else {
				require.NoError(t, err)
				require.NotNil(t, q)
			}
		})
	}
}

func TestNewQuest(t *testing.T) {
	type args struct {
		id          int
		title       string
		description string
		difficulty  int
	}
	tests := []struct {
		name    string
		args    args
		want    *Quest
		wantErr bool
	}{
		{
			name: "valid quest creation",
			args: args{
				id:          1,
				title:       "Quest Title",
				description: "Some description",
				difficulty:  5,
			},
			want: &Quest{
				ID:          1,
				Title:       "Quest Title",
				Description: "Some description",
				Difficulty:  5,
				XPReward:    500,
				isActive:    true,
				Tasks:       []*Task{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQuest(tt.args.id, tt.args.title, tt.args.description, tt.args.difficulty)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewQuest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQuest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuest_AddTask(t *testing.T) {
	type fields struct {
		ID          int
		Title       string
		Description string
		Difficulty  int
		isActive    bool
		XPReward    int
		Tasks       []*Task
	}
	type args struct {
		task *Task
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "add single task",
			fields: fields{
				ID:          1,
				Title:       "Quest",
				Description: "Desc",
				Difficulty:  5,
				isActive:    true,
				XPReward:    500,
				Tasks:       []*Task{},
			},
			args: args{
				task: &Task{
					ID:          1,
					isCompleted: false,
					XPReward:    100,
				},
			},
		},
		{
			name: "add multiple tasks",
			fields: fields{
				ID:          1,
				Title:       "Quest",
				Description: "Desc",
				Difficulty:  5,
				isActive:    true,
				XPReward:    500,
				Tasks:       []*Task{},
			},
			args: args{
				task: &Task{
					ID:          2,
					isCompleted: false,
					XPReward:    200,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quest{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Difficulty:  tt.fields.Difficulty,
				isActive:    tt.fields.isActive,
				XPReward:    tt.fields.XPReward,
				Tasks:       tt.fields.Tasks,
			}
			q.AddTask(tt.args.task)
			if len(q.Tasks) == 0 {
				t.Errorf("task was not added")
			}
			if q.Tasks[len(q.Tasks)-1].ID != tt.args.task.ID {
				t.Errorf("expected last task to be %d, got %d", tt.args.task.ID, q.Tasks[len(q.Tasks)-1].ID)
			}
		})
	}
}

func TestQuest_CompleteTask(t *testing.T) {
	type fields struct {
		ID          int
		Title       string
		Description string
		Difficulty  int
		isActive    bool
		XPReward    int
		Tasks       []*Task
	}
	type args struct {
		taskId int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "complete existing task",
			fields: fields{
				Tasks: []*Task{
					{ID: 1, isCompleted: false, XPReward: 100},
				},
			},
			args:    args{taskId: 1},
			wantErr: false,
		},
		{
			name: "task not found",
			fields: fields{
				Tasks: []*Task{},
			},
			args:    args{taskId: 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quest{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Difficulty:  tt.fields.Difficulty,
				isActive:    tt.fields.isActive,
				XPReward:    tt.fields.XPReward,
				Tasks:       tt.fields.Tasks,
			}
			if err := q.CompleteTask(tt.args.taskId); (err != nil) != tt.wantErr {
				t.Errorf("Quest.CompleteTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuest_Summary(t *testing.T) {
	type fields struct {
		ID          int
		Title       string
		Description string
		Difficulty  int
		isActive    bool
		XPReward    int
		Tasks       []*Task
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "summary with tasks",
			fields: fields{
				ID:          1,
				Title:       "Quest",
				Description: "Desc",
				Difficulty:  5,
				isActive:    true,
				XPReward:    500,
				Tasks: []*Task{
					{ID: 1, isCompleted: true},
					{ID: 2, isCompleted: false},
				},
			},
			want: "[q1] Quest | Difficulty: 5 | Progress : 1/2 | XP: 500 | active",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quest{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Difficulty:  tt.fields.Difficulty,
				isActive:    tt.fields.isActive,
				XPReward:    tt.fields.XPReward,
				Tasks:       tt.fields.Tasks,
			}
			if got := q.Summary(); got != tt.want {
				t.Errorf("Quest.Summary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuest_TotalXP(t *testing.T) {
	type fields struct {
		ID          int
		Title       string
		Description string
		Difficulty  int
		isActive    bool
		XPReward    int
		Tasks       []*Task
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "total xp includes quest and tasks",
			fields: fields{
				XPReward: 500,
				Tasks: []*Task{
					{XPReward: 100},
					{XPReward: 200},
				},
			},
			want: 800,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quest{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Difficulty:  tt.fields.Difficulty,
				isActive:    tt.fields.isActive,
				XPReward:    tt.fields.XPReward,
				Tasks:       tt.fields.Tasks,
			}
			if got := q.TotalXP(); got != tt.want {
				t.Errorf("Quest.TotalXP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuest_Activate(t *testing.T) {
	type fields struct {
		ID          int
		Title       string
		Description string
		Difficulty  int
		isActive    bool
		XPReward    int
		Tasks       []*Task
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "activate quest",
			fields: fields{
				isActive: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quest{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Difficulty:  tt.fields.Difficulty,
				isActive:    tt.fields.isActive,
				XPReward:    tt.fields.XPReward,
				Tasks:       tt.fields.Tasks,
			}
			q.Activate()
			if !q.isActive {
				t.Errorf("expected quest to be active")
			}
		})
	}
}

func TestQuest_Deactivate(t *testing.T) {
	type fields struct {
		ID          int
		Title       string
		Description string
		Difficulty  int
		isActive    bool
		XPReward    int
		Tasks       []*Task
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "deactivate quest",
			fields: fields{
				isActive: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quest{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Difficulty:  tt.fields.Difficulty,
				isActive:    tt.fields.isActive,
				XPReward:    tt.fields.XPReward,
				Tasks:       tt.fields.Tasks,
			}
			q.Deactivate()
			if q.isActive {
				t.Errorf("expected quest to be inactive")
			}
		})
	}
}

func TestQuest_GetCompletedTasks(t *testing.T) {
	type fields struct {
		ID          int
		Title       string
		Description string
		Difficulty  int
		isActive    bool
		XPReward    int
		Tasks       []*Task
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Task
	}{
		{
			name: "filter completed tasks",
			fields: fields{
				Tasks: []*Task{
					{ID: 1, isCompleted: true},
					{ID: 2, isCompleted: false},
				},
			},
			want: []*Task{
				{ID: 1, isCompleted: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quest{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Difficulty:  tt.fields.Difficulty,
				isActive:    tt.fields.isActive,
				XPReward:    tt.fields.XPReward,
				Tasks:       tt.fields.Tasks,
			}
			if got := q.GetCompletedTasks(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Quest.GetCompletedTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuest_GetRemainingTasks(t *testing.T) {
	type fields struct {
		ID          int
		Title       string
		Description string
		Difficulty  int
		isActive    bool
		XPReward    int
		Tasks       []*Task
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Task
	}{
		{
			name: "filter remaining tasks",
			fields: fields{
				Tasks: []*Task{
					{ID: 1, isCompleted: true},
					{ID: 2, isCompleted: false},
				},
			},
			want: []*Task{
				{ID: 2, isCompleted: false},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quest{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Difficulty:  tt.fields.Difficulty,
				isActive:    tt.fields.isActive,
				XPReward:    tt.fields.XPReward,
				Tasks:       tt.fields.Tasks,
			}
			if got := q.GetRemainingTasks(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Quest.GetRemainingTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuest_GetProgressPercentage(t *testing.T) {
	type fields struct {
		ID          int
		Title       string
		Description string
		Difficulty  int
		isActive    bool
		XPReward    int
		Tasks       []*Task
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "50 percent progress",
			fields: fields{
				Tasks: []*Task{
					{isCompleted: true},
					{isCompleted: false},
				},
			},
			want: 50.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quest{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Difficulty:  tt.fields.Difficulty,
				isActive:    tt.fields.isActive,
				XPReward:    tt.fields.XPReward,
				Tasks:       tt.fields.Tasks,
			}
			if got := q.GetProgressPercentage(); got != tt.want {
				t.Errorf("Quest.GetProgressPercentage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuest_IsCompleted(t *testing.T) {
	type fields struct {
		ID          int
		Title       string
		Description string
		Difficulty  int
		isActive    bool
		XPReward    int
		Tasks       []*Task
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "quest completed",
			fields: fields{
				Tasks: []*Task{
					{isCompleted: true},
					{isCompleted: true},
				},
			},
			want: true,
		},
		{
			name: "quest not completed",
			fields: fields{
				Tasks: []*Task{
					{isCompleted: true},
					{isCompleted: false},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quest{
				ID:          tt.fields.ID,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Difficulty:  tt.fields.Difficulty,
				isActive:    tt.fields.isActive,
				XPReward:    tt.fields.XPReward,
				Tasks:       tt.fields.Tasks,
			}
			if got := q.IsCompleted(); got != tt.want {
				t.Errorf("Quest.IsCompleted() = %v, want %v", got, tt.want)
			}
		})
	}
}
