package pokemon

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Client(t *testing.T) {

	t.Run("get pokemon - internal server error", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))

		client := NewPokemonClient(s.URL)

		_, err := client.GetPokemon(context.Background(), "pikachu")

		if err == nil {
			t.Fatal("expected error, got none")
		}

		if errors.Is(err, ErrPokemonNotFound) {
			t.Fatal("expected generic error, got not found error")
		}

	})
	t.Run("get pokemon - not found", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("request for %s", r.URL)
			w.WriteHeader(http.StatusNotFound)
		}))

		client := NewPokemonClient(s.URL)

		_, err := client.GetPokemon(context.Background(), "pikachu")

		if err == nil {

			t.Fatal("expected error, got none")
		}
		if !errors.Is(err, ErrPokemonNotFound) {
			t.Fatalf("expected not found error, got different error %v", err)
		}
	})
}
