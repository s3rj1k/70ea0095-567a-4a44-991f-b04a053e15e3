package main

import (
	"context"
	"fmt"
	"testing"

	"golang.org/x/text/language"
)

type failingTranslator struct{}

func (f *failingTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	return "", fmt.Errorf("translate failed")
}

func TestTranslatorWithCache(t *testing.T) {
	data := map[string]string{
		GenerateDataKey(language.English, language.Russian, "Twinkle, twinkle, little star"): "Мерцай, мерцай, маленькая звездочка",
		GenerateDataKey(language.English, language.Russian, "How I wonder what you are"):     "Как я жажду узнать, кто ты",
	}

	tt := []struct {
		in      string
		out     string
		isError bool
	}{
		{"Twinkle, twinkle, little star", "Мерцай, мерцай, маленькая звездочка", false},
		{"How I wonder what you are", "Как я жажду узнать, кто ты", false},
		{"Up above the world so high", "", true},
		{"Like a diamond in the sky", "", true},
	}

	x := TranslatorWithCache{
		Data:       data,
		Translator: &failingTranslator{},
	}

	for i := range tt {
		s, err := x.Translate(context.TODO(), language.English, language.Russian, tt[i].in)
		if (err != nil) && !tt[i].isError {
			t.Fatalf("do not expect error: %v", err)
		}

		if s != tt[i].out {
			t.Fatalf("expecting output to contan: %s, but contains: %s, input: %s", tt[i].out, s, tt[i].in)
		}
	}
}
