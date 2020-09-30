package models

import (
	"encoding/json"
	"testing"
)

func TestValidUsername(t *testing.T) {
	t.Run("with Valid Name 1", func(t *testing.T) {
		got := validUsername("Mathias d'Arras")
		want := true

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with Valid Name 2", func(t *testing.T) {
		got := validUsername("Martin Luther King, Jr.")
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
		got := validUserPic("https://lh3.googleusercontent.com/a-/AOh14Gh7K0TarVPkr5pzGM4zq7sCEY-WDjemYKjkPlkI=s96-c")
		want := true

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with invalid long filename", func(t *testing.T) {
		got := validUserPic("zxcvbnmasdfghasldkjaskpascvblroqasdasdasdadascxzlxcasklcasdkljjklqwertyuiop1234567890zxcvbnmasdfghjkl")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

}

func TestUserSanitizeJson(t *testing.T) {
	t.Run("with Valid Json", func(t *testing.T) {
		mockData := User{Name: "Carlos Maceira", Email: "carlosmgc2003@yahoo.com.ar", DisplayPic: "foto_carlos.jpg"}
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
	t.Run("with Invalid user name", func(t *testing.T) {
		mockData := User{Name: "nahuel__salazar", Email: "carlosmgc2003@yahoo.com.ar", DisplayPic: "foto_carlos.jpg"}
		mockJson, _ := json.Marshal(mockData)
		err := mockData.FromJson(mockJson)
		var got error
		if err != nil {
			got = err
		}
		want := InvalidUserName

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with Invalid Email", func(t *testing.T) {
		mockData := User{Name: "nahuelsalazar", Email: "carlosmgc2003", DisplayPic: "foto_carlos.jpg"}
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
		mockData := User{Name: "Carlos Maceira", Email: "carlosmgc2003@gmail.com", DisplayPic: "zxcvbnmasdfghasldkjaskpascvblroqasdasdasdadascxzlxcasklcasdkljjklqwertyuiop1234567890zxcvbnmasdfghjkl"}
		mockJson, _ := json.Marshal(mockData)
		err := mockData.FromJson(mockJson)
		var got error
		if err != nil {
			got = err
		}
		want := InvalidUserPic

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
}
