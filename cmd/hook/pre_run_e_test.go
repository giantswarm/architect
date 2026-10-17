package hook

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// testGit runs the git CLI in dir with a reproducible identity and no user or
// system configuration, so the layouts under test are the ones git creates.
func testGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...) // #nosec G204 -- the test's own git arguments
	cmd.Dir = dir
	cmd.Env = append(cmd.Environ(),
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

func runPreRunE(t *testing.T, workingDirectory string) *cobra.Command {
	t.Helper()
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("working-directory", workingDirectory, "")
	if err := PreRunE(cmd, nil); err != nil {
		t.Fatalf("PreRunE(%q): %v", workingDirectory, err)
	}
	return cmd
}

// TestPreRunE resolves the git defaults in a plain clone, in a subdirectory of
// it and in a linked worktree, whose .git is a file pointing at the main
// repository's git directory.
func TestPreRunE(t *testing.T) {
	t.Parallel()

	mainDir := t.TempDir()
	testGit(t, mainDir, "init", "--initial-branch=main")
	testGit(t, mainDir, "remote", "add", "origin", "https://example.com/org/repo.git")
	testGit(t, mainDir, "commit", "--allow-empty", "-m", "first")
	testGit(t, mainDir, "tag", "v1.2.3")
	mainSha := testGit(t, mainDir, "rev-parse", "HEAD")

	sub := filepath.Join(mainDir, "sub")
	if err := os.Mkdir(sub, 0o750); err != nil {
		t.Fatal(err)
	}

	wtDir := filepath.Join(t.TempDir(), "wt")
	testGit(t, mainDir, "worktree", "add", "-b", "feature", wtDir)

	testCases := []struct {
		name   string
		dir    string
		branch string
	}{
		{name: "plain clone", dir: mainDir, branch: "main"},
		{name: "subdirectory of a plain clone", dir: sub, branch: "main"},
		{name: "linked worktree", dir: wtDir, branch: "feature"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := runPreRunE(t, tc.dir)

			want := map[string]string{
				"sha":     mainSha,
				"tag":     "v1.2.3",
				"branch":  tc.branch,
				"version": "1.2.3",
			}
			for flag, value := range want {
				if got := cmd.Flag(flag).Value.String(); got != value {
					t.Errorf("flag %q = %q, want %q", flag, got, value)
				}
			}
		})
	}
}
