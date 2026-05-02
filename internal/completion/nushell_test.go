package completion_test

import (
	"strings"
	"testing"

	"github.com/simskij/yo/internal/completion"
	"github.com/spf13/cobra"
)

func buildTestCmd() *cobra.Command {
	root := &cobra.Command{Use: "yo", Short: "Personal CLI"}

	sub := &cobra.Command{Use: "dots", Short: "Manage dotfiles"}
	sub.PersistentFlags().String("machine", "", "override machine name")

	apply := &cobra.Command{Use: "apply", Short: "Symlink dotfiles into $HOME"}
	apply.Flags().Bool("dry-run", false, "preview without changes")
	apply.Flags().String("subpath", "", "limit to subpath")

	sub.AddCommand(apply)
	root.AddCommand(sub)

	init_ := &cobra.Command{Use: "init", Short: "Initialise yo"}
	init_.Flags().String("path", "", "storage path")
	root.AddCommand(init_)

	return root
}

func generate(t *testing.T) string {
	t.Helper()
	var buf strings.Builder
	if err := completion.GenerateNushell(buildTestCmd(), &buf); err != nil {
		t.Fatalf("GenerateNushell: %v", err)
	}
	return buf.String()
}

func TestGenerateNushell_ModuleWrapper(t *testing.T) {
	out := generate(t)
	if !strings.Contains(out, "module completions {") {
		t.Error("missing module completions block")
	}
	if !strings.Contains(out, "use completions *") {
		t.Error("missing use completions *")
	}
}

func TestGenerateNushell_RootExtern(t *testing.T) {
	out := generate(t)
	if !strings.Contains(out, `export extern "yo" [`) {
		t.Error("missing root extern")
	}
}

func TestGenerateNushell_SubcommandCompletionDef(t *testing.T) {
	out := generate(t)
	if !strings.Contains(out, `def "nu-complete yo" []`) {
		t.Error("missing nu-complete def for root")
	}
	if !strings.Contains(out, `"dots"`) {
		t.Error("missing dots entry in nu-complete def")
	}
}

func TestGenerateNushell_NestedExtern(t *testing.T) {
	out := generate(t)
	if !strings.Contains(out, `export extern "yo dots apply" [`) {
		t.Error("missing nested extern for dots apply")
	}
}

func TestGenerateNushell_BoolFlagIsSwitch(t *testing.T) {
	out := generate(t)
	if !strings.Contains(out, "--dry-run") {
		t.Error("missing --dry-run flag")
	}
	// bool flags must not have a type annotation
	if strings.Contains(out, "--dry-run: bool") {
		t.Error("bool flag should be a switch, not annotated with : bool")
	}
}

func TestGenerateNushell_StringFlagHasType(t *testing.T) {
	out := generate(t)
	if !strings.Contains(out, "--machine: string") {
		t.Error("missing --machine: string annotation")
	}
}

func TestGenerateNushell_HelpFlag(t *testing.T) {
	out := generate(t)
	if !strings.Contains(out, "--help(-h)") {
		t.Error("missing --help(-h) on at least one extern")
	}
}
