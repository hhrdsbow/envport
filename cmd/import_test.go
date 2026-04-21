package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportShellFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVPORT_DIR", dir)

	content := "export FOO=bar\nexport BAZ=qux\n"
	filePath := filepath.Join(dir, "myenv.sh")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	out, err := executeCommand(rootCmd, "import", filePath, "--name", "testsnap")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}

	if !strings.Contains(out, "2 variable(s)") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "testsnap") {
		t.Errorf("expected snapshot name in output, got: %s", out)
	}
}

func TestImportDotenvFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVPORT_DIR", dir)

	content := "FOO=hello\nBAR=world\n"
	filePath := filepath.Join(dir, "prod.env")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	out, err := executeCommand(rootCmd, "import", filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}

	if !strings.Contains(out, "prod") {
		t.Errorf("expected auto-detected name 'prod' in output, got: %s", out)
	}
}

func TestImportJSONFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVPORT_DIR", dir)

	content := `{"KEY":"value","OTHER":"data"}`
	filePath := filepath.Join(dir, "vars.json")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	out, err := executeCommand(rootCmd, "import", filePath, "--name", "jsonsnap")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}

	if !strings.Contains(out, "2 variable(s)") {
		t.Errorf("expected 2 variables imported, got: %s", out)
	}
}

func TestImportMissingFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVPORT_DIR", dir)

	_, err := executeCommand(rootCmd, "import", "/nonexistent/path/file.sh")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
