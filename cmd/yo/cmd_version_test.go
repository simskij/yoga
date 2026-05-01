//go:build integration

package main

import (
	"testing"
)

func TestVersion_DefaultIsDev(t *testing.T) {
	if version != "dev" {
		t.Errorf("expected default version to be 'dev', got %q", version)
	}
}
