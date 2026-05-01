package dots

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func Diff(dotsPath, machine string, out io.Writer) error {
	entries, err := Collect(dotsPath, machine)
	if err != nil {
		return err
	}

	any := false
	for _, e := range entries {
		if e.Status == StatusLinked || e.Status == StatusMissing {
			continue
		}
		any = true
		srcData, err := os.ReadFile(e.Src)
		if err != nil {
			return err
		}
		dstData, err := os.ReadFile(e.Dst)
		if err != nil {
			return err
		}
		if string(srcData) == string(dstData) {
			continue
		}
		fmt.Fprintf(out, "--- %s\n+++ %s\n", e.Dst, e.Src)
		printUnifiedDiff(out, string(dstData), string(srcData))
	}
	if !any {
		fmt.Fprintln(out, "No differences.")
	}
	return nil
}

func printUnifiedDiff(out io.Writer, a, b string) {
	aLines := strings.Split(a, "\n")
	bLines := strings.Split(b, "\n")
	maxLen := len(aLines)
	if len(bLines) > maxLen {
		maxLen = len(bLines)
	}
	for i := 0; i < maxLen; i++ {
		aLine, bLine := "", ""
		if i < len(aLines) {
			aLine = aLines[i]
		}
		if i < len(bLines) {
			bLine = bLines[i]
		}
		if aLine == bLine {
			fmt.Fprintf(out, "  %s\n", aLine)
		} else {
			if i < len(aLines) {
				fmt.Fprintf(out, "- %s\n", aLine)
			}
			if i < len(bLines) {
				fmt.Fprintf(out, "+ %s\n", bLine)
			}
		}
	}
}
