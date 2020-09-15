package models

import "testing"

func TestValidGameName(t *testing.T) {
	t.Run("with Valid Boardgame Name", func(t *testing.T) {
		got := validGameName("Estanciero")
		want := true

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("with invalid Boardgame Name", func(t *testing.T) {
		got := validGameName("Oa")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
}

func TestValidGameClass(t *testing.T) {
	t.Run("with each valid Boardgame class", func(t *testing.T) {
		gameclasses := []string{"Ingenio", "Estrategia", "Clasico", "Dados", "Palabras", "Cartas"}
		for _, class := range gameclasses {
			got := validGameClass(class)
			want := true
			if got != want {
				t.Errorf("got %t want %t", got, want)
			}
		}
	})
	t.Run("with invalid Boardgame class", func(t *testing.T) {
		got := validGameClass("Foo")
		want := false
		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
}

func TestFind(t *testing.T) {
	t.Run("with item that exists in slice", func(t *testing.T) {
		gameclasses := []string{"Ingenio", "Estrategia", "Clasico", "Dados", "Palabras", "Cartas"}
		gotInt, gotBool := find(gameclasses, gameclasses[0])
		wantInt, wantBool := 0, true
		if gotInt != wantInt {
			t.Errorf("got %d want %d", gotInt, wantInt)
		}
		if gotBool != wantBool {
			t.Errorf("got %t want %t", gotBool, wantBool)
		}
	})
	t.Run("with item that does not exists in slice", func(t *testing.T) {
		gameclasses := []string{"Ingenio", "Estrategia", "Clasico", "Dados", "Palabras", "Cartas"}
		gotInt, gotBool := find(gameclasses, "Foo")
		wantInt, wantBool := -1, false
		if gotInt != wantInt {
			t.Errorf("got %d want %d", gotInt, wantInt)
		}
		if gotBool != wantBool {
			t.Errorf("got %t want %t", gotBool, wantBool)
		}
	})
}

func TestValidGamePic(t *testing.T) {
	t.Run("with invalid Boardgame Name", func(t *testing.T) {
		got := validGamePic("foto_de_carlos.jpg")
		want := true

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
	t.Run("with invalid (long) Boardgame Name", func(t *testing.T) {
		got := validGamePic("foto_de_carlos_en_vacaciones_de_invierno_esquiando.jpg")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
	t.Run("with invalid (short) Boardgame Name", func(t *testing.T) {
		got := validGamePic("foto.jpg")
		want := false

		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
}
