package main

import (
	"context"
	"fmt"
	"time"

	"github.com/flowchartsman/retry"
	"golang.org/x/text/language"
)

type TranslatorWithRetry struct {
	MaxTries     int
	InitialDelay time.Duration
	MaxDelay     time.Duration

	Translator
}

func (t *TranslatorWithRetry) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	retrier := retry.NewRetrier(t.MaxTries, t.InitialDelay, t.MaxDelay)

	var translated string

	err := retrier.RunContext(ctx, func(ctx context.Context) error {
		var err error

		if translated, err = t.Translator.Translate(ctx, from, to, data); err == nil {
			return nil
		}

		return err
	})
	if err != nil {
		return "", fmt.Errorf("no mo retries: %w", err)
	}

	return translated, nil
}
