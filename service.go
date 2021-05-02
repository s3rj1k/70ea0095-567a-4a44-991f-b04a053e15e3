package main

import "time"

// Service is a Translator user.
type Service struct {
	translator Translator
}

func NewService() *Service {
	t := newRandomTranslator(
		100*time.Millisecond,
		500*time.Millisecond,
		0.1,
	)

	return &Service{
		translator: &TranslatorWithDeduplication{
			Translator: &TranslatorWithRetry{
				MaxTries:     5,
				InitialDelay: 50 * time.Millisecond,
				MaxDelay:     500 * time.Millisecond,

				Translator: &TranslatorWithCache{
					Translator: t,
				},
			},
		},
	}
}
