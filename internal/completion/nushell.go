package completion

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// GenerateNushell writes a nushell completion module for root to w.
func GenerateNushell(root *cobra.Command, w io.Writer) error {
	var b strings.Builder
	b.WriteString("module completions {\n")
	writeNushellExterns(&b, root, root.Name())
	b.WriteString("}\n\nuse completions *\n")
	_, err := io.WriteString(w, b.String())
	return err
}

func writeNushellExterns(b *strings.Builder, cmd *cobra.Command, path string) {
	visible := visibleSubcmds(cmd)

	if len(visible) > 0 {
		fmt.Fprintf(b, "\n  def \"nu-complete %s\" [] {\n    [\n", path)
		for _, sub := range visible {
			fmt.Fprintf(b, "      [\"%s\", \"%s\"]\n", sub.Name(), sub.Short)
		}
		b.WriteString("    ]\n  }\n")
	}

	fmt.Fprintf(b, "\n  export extern \"%s\" [\n", path)
	if len(visible) > 0 {
		fmt.Fprintf(b, "    subcommand?: string@\"nu-complete %s\"\n", path)
	}

	seen := map[string]bool{}
	writeFlag := func(f *pflag.Flag) {
		if f.Name == "help" || seen[f.Name] {
			return
		}
		seen[f.Name] = true
		if f.Value.Type() == "bool" {
			if f.Shorthand != "" {
				fmt.Fprintf(b, "    --%s(-%s)  # %s\n", f.Name, f.Shorthand, f.Usage)
			} else {
				fmt.Fprintf(b, "    --%s  # %s\n", f.Name, f.Usage)
			}
		} else {
			if f.Shorthand != "" {
				fmt.Fprintf(b, "    --%s(-%s): %s  # %s\n", f.Name, f.Shorthand, nushellType(f.Value.Type()), f.Usage)
			} else {
				fmt.Fprintf(b, "    --%s: %s  # %s\n", f.Name, nushellType(f.Value.Type()), f.Usage)
			}
		}
	}
	cmd.Flags().VisitAll(writeFlag)
	cmd.InheritedFlags().VisitAll(writeFlag)

	b.WriteString("    --help(-h)\n  ]\n")

	for _, sub := range visible {
		writeNushellExterns(b, sub, path+" "+sub.Name())
	}
}

func visibleSubcmds(cmd *cobra.Command) []*cobra.Command {
	var out []*cobra.Command
	for _, sub := range cmd.Commands() {
		if !sub.Hidden && sub.Name() != "help" && sub.Name() != "completion" {
			out = append(out, sub)
		}
	}
	return out
}

func nushellType(goType string) string {
	switch goType {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		return "int"
	default:
		return "string"
	}
}
