package main

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/text/language"
)

type countingTranslator struct {
	Count *uint64
}

func (f *countingTranslator) Translate(_ context.Context, _, _ language.Tag, _ string) (string, error) {
	time.Sleep(2500 * time.Millisecond) // simulate work
	atomic.AddUint64(f.Count, 1)
	return "", nil
}

func TestTranslatorWithDedup(t *testing.T) {
	var cnt uint64

	x := TranslatorWithDeduplication{
		Translator: &countingTranslator{
			Count: &cnt,
		},
	}

	var (
		queriesCnt int = 2048
		wg         sync.WaitGroup
	)

	for i := 0; i < queriesCnt; i++ {
		wg.Add(1)
		go func() {
			_, _ = x.Translate(context.TODO(), language.English, language.Russian, text)
			wg.Done()
		}()
	}

	wg.Wait()

	if int(cnt) >= queriesCnt {
		t.Errorf("expecting that proccessed queries %d be less than %d", queriesCnt, cnt)
	}
}
