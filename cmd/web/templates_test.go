package main

import (
	"testing"
	"time"

	"snippetbox.tmgasek.net/internal/assert"
)

// func TestExample(t *testing.T) {
// 	t.Run("Example sub-test 1", func(t *testing.T) {
// 		// Do a test.
// 	})
// 	t.Run("Example sub-test 2", func(t *testing.T) {
// 		// Do another test.
// 	})
// 	t.Run("Example sub-test 3", func(t *testing.T) {
// 		// And another...
// 	})
// }

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2022 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2022 at 09:15",
		},
	}

	for _, tt := range tests {
		// Run a sub-test for each test case
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want)
		})
	}

}
