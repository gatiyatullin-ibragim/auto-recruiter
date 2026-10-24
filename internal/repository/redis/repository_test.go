package redis

import (
	"context"
	"testing"
	"time"

	"github.com/gatiyatullin-ibragim/auto-recruiter/internal/domain"
)

func TestStoreGetSet(t *testing.T) {
	t.Parallel()
	store := NewStore()
	if err := store.Set(context.Background(), "key", "value", 0); err != nil {
		t.Fatal(err)
	}
	got, err := store.Get(context.Background(), "key")
	if err != nil || got != "value" {
		t.Fatalf("Get() = %q, %v", got, err)
	}
	if _, err := store.Get(context.Background(), "missing"); err != domain.ErrCacheMiss {
		t.Fatalf("missing Get() error = %v, want ErrCacheMiss", err)
	}
	if err := store.Set(context.Background(), "", "value", time.Minute); err != domain.ErrInvalidInput {
		t.Fatalf("empty key Set() error = %v, want ErrInvalidInput", err)
	}
}
