package main

import (
	"github.com/makarellav/codecapsule/internal/assert"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 4, 1, 9, 13, 0, 0, time.UTC),
			want: "01 Apr 2024 at 09:13",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 4, 1, 9, 13, 0, 0, time.FixedZone("CET", 60*60)),
			want: "01 Apr 2024 at 08:13",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
	}
}
