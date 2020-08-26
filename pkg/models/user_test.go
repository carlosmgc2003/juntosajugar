package models

import (
	"encoding/json"
	"testing"
)

func TestValidUsername(t *testing.T) {
	t.Run("with Valid Username", func(t *testing.T) {
		got := validUsername("carlosmgc2003")
		want := true

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with Two Underscores", func(t *testing.T) {
		got := validUsername("nahuel__salazar")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with Empty Spaces", func(t *testing.T) {
		got := validUsername("matias kobold")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("Longer than allowed", func(t *testing.T) {
		got := validUsername("mnbvcxzasdfghjklqwertyuiop12345678909")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("Shorter than allowed", func(t *testing.T) {
		got := validUsername("toto")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
}

func TestValidEmail(t *testing.T) {
	t.Run("with Valid email", func(t *testing.T) {
		got := validEmail("carlosmgc2003@gmail.com")
		want := true

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with invalid utf8 character", func(t *testing.T) {
		got := validEmail("ñsalazar@fie.undef.edu.ar")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with No Domain", func(t *testing.T) {
		got := validEmail("cmaceria")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with no username", func(t *testing.T) {
		got := validEmail("@gmail.com.ar")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

}

func TestValidFilename(t *testing.T) {
	t.Run("with Valid filename", func(t *testing.T) {
		got := validFilename("foto_carlos.jpg")
		want := true

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with invalid long filename", func(t *testing.T) {
		got := validFilename("zxcvbnmasdfghjklqwertyuiop1234567890zxcvbnmasdfghjkl")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with invalid short filename", func(t *testing.T) {
		got := validFilename("zxcv")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
}

func TestUserSanitizeJson(t *testing.T) {
	t.Run("with Valid Json", func(t *testing.T) {
		mockData := User{Name: "carlosmgc2003", Email: "carlosmgc2003@yahoo.com.ar", Display_pic: "foto_carlos.jpg"}
		mockJson, _ := json.Marshal(mockData)
		err := mockData.FromJson(mockJson)
		var got bool
		if err != nil {
			got = true
		}
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
	t.Run("with Invalid Username", func(t *testing.T) {
		mockData := User{Name: "nahuel__salazar", Email: "carlosmgc2003@yahoo.com.ar", Display_pic: "foto_carlos.jpg"}
		mockJson, _ := json.Marshal(mockData)
		err := mockData.FromJson(mockJson)
		var got error
		if err != nil {
			got = err
		}
		want := InvalidUsername

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with Invalid Email", func(t *testing.T) {
		mockData := User{Name: "nahuelsalazar", Email: "carlosmgc2003", Display_pic: "foto_carlos.jpg"}
		mockJson, _ := json.Marshal(mockData)
		err := mockData.FromJson(mockJson)
		var got error
		if err != nil {
			got = err
		}
		want := InvalidEmail

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with Invalid filename", func(t *testing.T) {
		mockData := User{Name: "carlosmgc2003", Email: "carlosmgc2003@gmail.com", Display_pic: "foto"}
		mockJson, _ := json.Marshal(mockData)
		err := mockData.FromJson(mockJson)
		var got error
		if err != nil {
			got = err
		}
		want := InvalidFilename

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
}
