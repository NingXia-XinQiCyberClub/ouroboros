package productcontract

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestRunRealRepo runs the embedded driver against this repository.
// It is the smoke check for the port; skipped when node is absent or the
// GUI dependencies are not installed.
func TestRunRealRepo(t *testing.T) {
	if _, err := exec.LookPath("node"); err != nil {
		t.Skip("node not on PATH")
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	guiRoot := filepath.Join(cwd, "..", "..", "..", "gui")
	// The driver imports postcss as a bare specifier, resolvable only from
	// gui/node_modules. Skip on checkouts without GUI dependencies installed
	// (bare `make test` on the Go side) instead of failing on the environment.
	if _, err := os.Stat(filepath.Join(guiRoot, "node_modules", "postcss")); err != nil {
		t.Skip("gui/node_modules/postcss not installed; run npm install under gui/ to enable")
	}
	repoRoot := filepath.Dir(guiRoot)
	if err := Run(context.Background(), Options{GUIRoot: guiRoot, RepoRoot: repoRoot, Stdout: os.Stdout, Stderr: os.Stderr}); err != nil {
		t.Fatalf("product contract failed: %v (exit code %d)", err, ExitCode(err))
	}
}
