package inventory

import (
	"sync"
	"testing"
	"time"
)

func TestBug5_ConcurrentGetNoSharedPointers(t *testing.T) {
	s := New()
	now := time.Date(2026, 8, 14, 12, 0, 0, 0, time.UTC)
	if _, err := s.Create(CreateInput{SKU: "CG", Name: "conc", Stock: 5}, now); err != nil {
		t.Fatal(err)
	}
	if _, err := s.StockIn("CG", AmountInput{Amount: 1}, now.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	results := make([]*Product, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			p, _ := s.Get("CG")
			results[idx] = p
		}(i)
	}
	wg.Wait()
	if results[0].LastInAt == nil || results[1].LastInAt == nil {
		t.Fatal("LastInAt should not be nil")
	}
	*results[0].LastInAt = time.Time{}
	if results[1].LastInAt.IsZero() {
		t.Fatal("BUG: concurrent Get returns shared *time.Time pointers")
	}
}
