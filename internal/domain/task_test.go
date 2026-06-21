package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTask_Validation(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		title    string
		xp       int
		questID  int
		wantErr  bool
		wantField string
	}{
		{
			name:    "valid task",
			id:      1,
			title:   "Task 1",
			xp:      100,
			questID: 1,
			wantErr: false,
		},
		{
			name:     "empty title",
			id:       1,
			title:    "",
			xp:       100,
			questID:  1,
			wantErr:  true,
			wantField: "title",
		},
		{
			name:    "negative xp normalized",
			id:      1,
			title:   "Task",
			xp:      -10,
			questID: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewTask(tt.id, tt.title, tt.xp, tt.questID)

			if tt.wantErr {
				require.Error(t, err)

				var ve *ValidationError
				require.True(t, errors.As(err, &ve))
				require.Equal(t, tt.wantField, ve.Field)

				require.Nil(t, task)
			} else {
				require.NoError(t, err)
				require.NotNil(t, task)
			}
		})
	}
}