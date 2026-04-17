package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeImportFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "import.env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestImportFromFile_AddsNewKeys(t *testing.T) {
	p := writeImportFile(t, "FOO=bar\nBAZ=qux\n")
	dst := map[string]string{}
	res, err := ImportFromFile(p, dst, ImportOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Imported != 2 {
		t.Errorf("expected 2 imported, got %d", res.Imported)
	}
	if dst["FOO"] != "bar" || dst["BAZ"] != "qux" {
		t.Errorf("unexpected dst: %v", dst)
	}
}

func TestImportFromFile_SkipsExistingWithoutOverwrite(t *testing.T) {
	p := writeImportFile(t, "FOO=new\n")
	dst := map[string]string{"FOO": "old"}
	res, err := ImportFromFile(p, dst, ImportOptions{Overwrite: false})
	if err != nil {
		t.Fatal(err)
	}
	if res.Skipped != 1 || dst["FOO"] != "old" {
		t.Errorf("expected skip, got %v / %s", res, dst["FOO"])
	}
}

func TestImportFromFile_OverwriteReplacesExisting(t *testing.T) {
	p := writeImportFile(t, "FOO=new\n")
	dst := map[string]string{"FOO": "old"}
	_, err := ImportFromFile(p, dst, ImportOptions{Overwrite: true})
	if err != nil {
		t.Fatal(err)
	}
	if dst["FOO"] != "new" {
		t.Errorf("expected new, got %s", dst["FOO"])
	}
}

func TestImportFromFile_DryRunDoesNotMutate(t *testing.T) {
	p := writeImportFile(t, "FOO=bar\n")
	dst := map[string]string{}
	res, err := ImportFromFile(p, dst, ImportOptions{DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(dst) != 0 {
		t.Errorf("dst should be empty in dry-run")
	}
	if res.Imported != 1 {
		t.Errorf("expected imported count 1, got %d", res.Imported)
	}
}

func TestImportFromFile_IgnoresCommentsAndBlanks(t *testing.T) {
	p := writeImportFile(t, "# comment\n\nKEY=val\n")
	dst := map[string]string{}
	res, _ := ImportFromFile(p, dst, ImportOptions{})
	if res.Imported != 1 {
		t.Errorf("expected 1 imported, got %d", res.Imported)
	}
}

func TestImportFromFile_MissingFile(t *testing.T) {
	_, err := ImportFromFile("/nonexistent/file.env", map[string]string{}, ImportOptions{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
