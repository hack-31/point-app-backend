package model

import (
	"testing"
	"time"

	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	type input struct {
		time time.Time
	}
	type want struct {
		time string
	}

	tests := map[string]struct {
		input input
		want  want
	}{
		"Format関数はyyyy/MM/dd hh:mm:ssに変換する": {
			input: input{
				time: clock.FixedClocker{}.Now(),
			},
			want: want{
				time: "2022/05/10 12:34:56",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			got := NewTime(tt.input.time).Format()
			assert.Equal(t, tt.want.time, got)
		})
	}
}
