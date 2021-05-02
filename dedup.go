package main

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/text/language"
)

type TranslatorWithDeduplication struct {
	queue sync.Map

	Translator
}

var errorDedup = errors.New("query is currently being proccessed by another worker")

func (t *TranslatorWithDeduplication) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	key := GenerateDataKey(from, to, data)

	defer func() {
		t.queue.Delete(key)
	}()

	_, loaded := t.queue.LoadOrStore(key, struct{}{})
	if loaded {
		return "", errorDedup
	}

	return t.Translator.Translate(ctx, from, to, data)
}
