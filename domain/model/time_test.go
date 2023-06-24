package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/hack-31/point-app-backend/utils/clock"
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
			// アサーション
			if d := cmp.Diff(got, tt.want.time); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
		})
	}
}
