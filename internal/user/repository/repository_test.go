package repository_test

import (
	"testing"

	"remora/internal/user/repository"
)

func TestNew(t *testing.T) {
	t.Parallel()

	repo := repository.New(nil)
	if repo == nil {
		t.Fatal("New() returned nil")
	}
}
