package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUser_Validation(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		username  string
		email     string
		wantErr   bool
		wantField string
	}{
		{
			name:     "valid user",
			id:       "u1",
			username: "kirill",
			email:    "test@mail.com",
			wantErr:  false,
		},
		{
			name:      "username too short",
			id:        "u1",
			username:  "a",
			email:     "test@mail.com",
			wantErr:   true,
			wantField: "username",
		},
		{
			name:     "empty id allowed",
			id:       "",
			username: "kirill",
			email:    "test@mail.com",
			wantErr:  false,
		},
		{
			name:     "empty email allowed",
			id:       "u1",
			username: "kirill",
			email:    "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.id, tt.username, tt.email)

			if tt.wantErr {
				require.Error(t, err)

				var ve *ValidationError
				require.True(t, errors.As(err, &ve))
				require.Equal(t, tt.wantField, ve.Field)

				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
			}
		})
	}
}

func TestUser_AddXP(t *testing.T) {
	type fields struct {
		ID              string
		Username        string
		Email           string
		Level           int
		XP              int
		CompletedQuests []*Quest
		TotalXPEarned   int
	}
	type args struct {
		amount int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantLevel int
		wantXP    int
		wantTotal int
	}{
		{
			name: "add xp without level up",
			fields: fields{
				ID:       "u1",
				Username: "kirill",
				Level:    1,
				XP:       0,
			},
			args: args{amount: 500},
			wantLevel: 1,
			wantXP:    500,
			wantTotal: 500,
		},
		{
			name: "level up once",
			fields: fields{
				ID:       "u1",
				Username: "kirill",
				Level:    1,
				XP:       0,
			},
			args: args{amount: 1200},
			wantLevel: 2,
			wantXP:    200,
			wantTotal: 1200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:              tt.fields.ID,
				Username:        tt.fields.Username,
				Email:           tt.fields.Email,
				Level:           tt.fields.Level,
				XP:              tt.fields.XP,
				CompletedQuests: tt.fields.CompletedQuests,
				TotalXPEarned:   tt.fields.TotalXPEarned,
			}
			u.AddXP(tt.args.amount)
			if u.Level != tt.wantLevel {
				t.Errorf("Level = %d, want %d", u.Level, tt.wantLevel)
			}
			if u.XP != tt.wantXP {
				t.Errorf("XP = %d, want %d", u.XP, tt.wantXP)
			}
			if u.TotalXPEarned != tt.wantTotal {
				t.Errorf("TotalXPEarned = %d, want %d", u.TotalXPEarned, tt.wantTotal)
			}
		})
	}
}

func TestUser_LevelUp(t *testing.T) {
	type fields struct {
		ID              string
		Username        string
		Email           string
		Level           int
		XP              int
		CompletedQuests []*Quest
		TotalXPEarned   int
	}
	tests := []struct {
		name string
		fields fields
		wantLevel int
		wantXP int
	}{
		{
			name: "level up resets xp",
			fields: fields{
				Level: 3,
				XP:    500,
			},
			wantLevel: 4,
			wantXP: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:              tt.fields.ID,
				Username:        tt.fields.Username,
				Email:           tt.fields.Email,
				Level:           tt.fields.Level,
				XP:              tt.fields.XP,
				CompletedQuests: tt.fields.CompletedQuests,
				TotalXPEarned:   tt.fields.TotalXPEarned,
			}
			u.LevelUp()
			if u.Level != tt.wantLevel {
				t.Errorf("Level = %d, want %d", u.Level, tt.wantLevel)
			}
			if u.XP != tt.wantXP {
				t.Errorf("XP = %d, want %d", u.XP, tt.wantXP)
			}
		})
	}
}

func TestUser_CompleteQuest(t *testing.T) {
	type fields struct {
		ID              string
		Username        string
		Email           string
		Level           int
		XP              int
		CompletedQuests []*Quest
		TotalXPEarned   int
	}
	type args struct {
		quest *Quest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "quest not completed returns error",
			fields: fields{
				Username: "kirill",
			},
			args: args{
				quest: &Quest{
					Title: "q1",
					Tasks: []*Task{
						{
							Title: "task1",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "quest completed adds xp",
			fields: fields{
				Username: "kirill",
				Level:    1,
				XP:       0,
			},
			args: args{
				quest: &Quest{
					Title: "q1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:              tt.fields.ID,
				Username:        tt.fields.Username,
				Email:           tt.fields.Email,
				Level:           tt.fields.Level,
				XP:              tt.fields.XP,
				CompletedQuests: tt.fields.CompletedQuests,
				TotalXPEarned:   tt.fields.TotalXPEarned,
			}
			err := u.CompleteQuest(tt.args.quest);
			if (err != nil) != tt.wantErr {
				t.Errorf("User.CompleteQuest() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && len(u.CompletedQuests) == 0 && !tt.wantErr {
				t.Errorf("quest was not added to CompletedQuests")
			}
		})
	}
}
