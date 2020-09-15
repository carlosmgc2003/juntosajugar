package models

import (
	"testing"
	"time"
)

func TestValidMeetingPlace(t *testing.T) {
	t.Run("with Valid meeting place", func(t *testing.T) {
		got := validMeetingPlace("Av. San Martin 2020")
		want := true

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with invalid meeting place", func(t *testing.T) {
		got := validMeetingPlace("Foo")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with Empty meetingplace", func(t *testing.T) {
		got := validMeetingPlace("")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
}

func TestValidScheduledTime(t *testing.T) {
	t.Run("with Valid scheduledtime", func(t *testing.T) {
		start := time.Now()
		inSevenDays := start.Add(time.Hour * 24 * 7)
		got := time.Now().Before(inSevenDays)
		want := true

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
}
