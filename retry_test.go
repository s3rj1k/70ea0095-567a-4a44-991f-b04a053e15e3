package main

import (
	"context"
	"math"
	"strings"
	"testing"
	"time"

	"golang.org/x/text/language"
)

const (
	maxTries = 5
	minDelay = 100 * time.Millisecond
	maxDelay = 1 * time.Second

	text = "Hello world!"
)

func TestTranslatorWithRetryUnconditionalFail(t *testing.T) {
	x := TranslatorWithRetry{
		MaxTries:     maxTries,
		InitialDelay: minDelay,
		MaxDelay:     maxDelay,

		Translator: newRandomTranslator(minDelay, maxDelay, 1),
	}

	s, err := x.Translate(context.TODO(), language.English, language.Russian, text)
	if err == nil {
		t.Errorf("expecting error")
	}

	if s != "" {
		t.Errorf("expecting empty string")
	}
}

func TestTranslatorWithRetryUnconditionalSuccess(t *testing.T) {
	x := TranslatorWithRetry{
		MaxTries:     maxTries,
		InitialDelay: minDelay,
		MaxDelay:     maxDelay,

		Translator: newRandomTranslator(minDelay, maxDelay, 0),
	}

	s, err := x.Translate(context.TODO(), language.English, language.Russian, text)
	if err != nil {
		t.Fatalf("do not expect error: %v", err)
	}

	if !strings.Contains(s, text) {
		t.Fatalf("expecting output to contan: %s", text)
	}
}

func TestTranslatorWithRetryEventualSuccess(t *testing.T) {
	x := TranslatorWithRetry{
		MaxTries:     math.MaxInt32,
		InitialDelay: minDelay,
		MaxDelay:     maxDelay,

		Translator: newRandomTranslator(minDelay, maxDelay, 0.9),
	}

	s, err := x.Translate(context.TODO(), language.English, language.Russian, text)
	if err != nil {
		t.Fatalf("do not expect error: %v", err)
	}

	if !strings.Contains(s, text) {
		t.Fatalf("expecting output to contan: %s", text)
	}
}
