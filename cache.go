package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"

	"golang.org/x/text/language"
)

type TranslatorWithCache struct {
	Data    map[string]string
	cacheMu sync.Mutex

	Translator
}

func GenerateDataKey(from, to language.Tag, data string) string {
	h := sha256.New()

	h.Write([]byte(from.String()))
	h.Write([]byte(to.String()))
	h.Write([]byte(data))

	return hex.EncodeToString(h.Sum(nil))
}

func (t *TranslatorWithCache) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	t.cacheMu.Lock()
	defer t.cacheMu.Unlock()

	if t.Data == nil {
		t.Data = make(map[string]string)
	}

	key := GenerateDataKey(from, to, data)

	if val, ok := t.Data[key]; ok {
		return val, nil
	}

	translated, err := t.Translator.Translate(ctx, from, to, data)
	if err == nil {
		t.Data[key] = translated

		return translated, nil
	}

	return "", nil
}
