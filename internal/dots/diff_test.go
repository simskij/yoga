package dots

import (
	"strings"
	"testing"
)

func TestPrintUnifiedDiff_Identical(t *testing.T) {
	var buf strings.Builder
	printUnifiedDiff(&buf, "line1\nline2\n", "line1\nline2\n")
	out := buf.String()
	if strings.Contains(out, "+") || strings.Contains(out, "-") {
		t.Errorf("expected no diff markers for identical content, got:\n%s", out)
	}
}

func TestPrintUnifiedDiff_Addition(t *testing.T) {
	var buf strings.Builder
	printUnifiedDiff(&buf, "line1\n", "line1\nline2\n")
	out := buf.String()
	if !strings.Contains(out, "+ line2") {
		t.Errorf("expected addition marker, got:\n%s", out)
	}
}

func TestPrintUnifiedDiff_Removal(t *testing.T) {
	var buf strings.Builder
	printUnifiedDiff(&buf, "line1\nline2\n", "line1\n")
	out := buf.String()
	if !strings.Contains(out, "- line2") {
		t.Errorf("expected removal marker, got:\n%s", out)
	}
}

func TestPrintUnifiedDiff_Changed(t *testing.T) {
	var buf strings.Builder
	printUnifiedDiff(&buf, "old\n", "new\n")
	out := buf.String()
	if !strings.Contains(out, "- old") || !strings.Contains(out, "+ new") {
		t.Errorf("expected both markers, got:\n%s", out)
	}
}
