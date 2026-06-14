package domain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewReward_Validation(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		title     string
		xp        int
		rarity    string
		wantErr   bool
		wantField string
	}{
		{
			name:    "valid reward",
			id:      "r1",
			title:   "Sword",
			xp:      100,
			rarity:  "rare",
			wantErr: false,
		},
		{
			name:      "empty title",
			id:        "r1",
			title:     "",
			xp:        100,
			rarity:    "rare",
			wantErr:   true,
			wantField: "title",
		},
		{
			name:      "xp must be > 0",
			id:        "r1",
			title:     "Sword",
			xp:        0,
			rarity:    "rare",
			wantErr:   true,
			wantField: "XPAMount",
		},
		{
			name:    "invalid rarity defaults to common",
			id:      "r1",
			title:   "Sword",
			xp:      100,
			rarity:  "invalid",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reward, err := NewReward(tt.id, tt.title, tt.xp, tt.rarity)

			if tt.wantErr {
				require.Error(t, err)

				var ve *ValidationError
				require.True(t, errors.As(err, &ve))
				require.Equal(t, tt.wantField, ve.Field)

				require.Nil(t, reward)
			} else {
				require.NoError(t, err)
				require.NotNil(t, reward)
			}
		})
	}
}

func TestNewReward(t *testing.T) {
	type args struct {
		id       string
		title    string
		XPAmount int
		rarity   string
	}
	tests := []struct {
		name    string
		args    args
		want    *Reward
		wantErr bool
	}{
		{
			name: "successfully create common reward",
			args: args{id: "rew-1", title: "Copper Coins", XPAmount: 50, rarity: "common"},
			want: &Reward{ID: "rew-1", Title: "Copper Coins", XPAmount: 50, Rarity: "common"},
			wantErr: false,
		},
		{
			name: "successfully create epic reward",
			args: args{id: "rew-2", title: "Dragon Scale", XPAmount: 500, rarity: "epic"},
			want: &Reward{ID: "rew-2", Title: "Dragon Scale", XPAmount: 500, Rarity: "epic"},
			wantErr: false,
		},
		{
			name: "invalid rarity falls back to common",
			args: args{id: "rew-3", title: "Rusty Key", XPAmount: 10, rarity: "super-rare"},
			want: &Reward{ID: "rew-3", Title: "Rusty Key", XPAmount: 10, Rarity: "common"},
			wantErr: false,
		},
		{
			name: "fail with zero XP amount",
			args: args{id: "rew-4", title: "Ghost Gift", XPAmount: 0, rarity: "legendary"},
			want: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewReward(tt.args.id, tt.args.title, tt.args.XPAmount, tt.args.rarity)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReward() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReward() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReward_String(t *testing.T) {
	type fields struct {
		ID       string
		Title    string
		XPAmount int
		Rarity   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "common reward string",
			fields: fields{ID: "1", Title: "Bread", XPAmount: 15, Rarity: "common"},
			want:   "[C] Bread - XP: 15 (common)",
		},
		{
			name:   "rare reward string",
			fields: fields{ID: "2", Title: "Silver Ring", XPAmount: 150, Rarity: "rare"},
			want:   "[R] Silver Ring - XP: 150 (rare)",
		},
		{
			name:   "epic reward string",
			fields: fields{ID: "3", Title: "Crystal Ball", XPAmount: 400, Rarity: "epic"},
			want:   "[E] Crystal Ball - XP: 400 (epic)",
		},
		{
			name:   "legendary reward string",
			fields: fields{ID: "4", Title: "Excalibur", XPAmount: 1000, Rarity: "legendary"},
			want:   "[L] Excalibur - XP: 1000 (legendary)",
		},
		{
			name:   "fallback placeholder for unknown rarity",
			fields: fields{ID: "5", Title: "Glitch Item", XPAmount: 1, Rarity: "unknown"},
			want:   "[?] Glitch Item - XP: 1 (unknown)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Reward{
				ID:       tt.fields.ID,
				Title:    tt.fields.Title,
				XPAmount: tt.fields.XPAmount,
				Rarity:   tt.fields.Rarity,
			}
			if got := r.String(); got != tt.want {
				t.Errorf("Reward.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReward_Value(t *testing.T) {
	type fields struct {
		ID       string
		Title    string
		XPAmount int
		Rarity   string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "common multiplier is 1x",
			fields: fields{XPAmount: 100, Rarity: "common"},
			want:   100,
		},
		{
			name:   "rare multiplier is 2x",
			fields: fields{XPAmount: 100, Rarity: "rare"},
			want:   200,
		},
		{
			name:   "epic multiplier is 5x",
			fields: fields{XPAmount: 100, Rarity: "epic"},
			want:   500,
		},
		{
			name:   "legendary multiplier is 10x",
			fields: fields{XPAmount: 100, Rarity: "legendary"},
			want:   1000,
		},
		{
			name:   "unknown rarity fallback multiplier is 1x",
			fields: fields{XPAmount: 100, Rarity: "mythic"},
			want:   100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Reward{
				ID:       tt.fields.ID,
				Title:    tt.fields.Title,
				XPAmount: tt.fields.XPAmount,
				Rarity:   tt.fields.Rarity,
			}
			if got := r.Value(); got != tt.want {
				t.Errorf("Reward.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}