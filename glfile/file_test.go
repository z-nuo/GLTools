package glfile

import (
	"path/filepath"
	"testing"
)

func TestEnsureDirCreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "dir")
	if Exists(dir) {
		t.Fatal("Exists() = true before directory creation")
	}
	if err := EnsureDir(dir); err != nil {
		t.Fatal(err)
	}
	if !Exists(dir) {
		t.Fatal("Exists() = false, want true")
	}
	if !IsDir(dir) {
		t.Fatal("IsDir() = false, want true")
	}
	if IsFile(dir) {
		t.Fatal("IsFile() = true, want false")
	}
}

func TestWriteTextAndReadText(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "note.txt")
	if err := WriteText(path, "hello"); err != nil {
		t.Fatal(err)
	}
	if !IsFile(path) {
		t.Fatal("IsFile() = false, want true")
	}
	got, err := ReadText(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello" {
		t.Fatalf("ReadText() = %q, want hello", got)
	}
}

func TestExt(t *testing.T) {
	if got := Ext("/tmp/archive.tar.gz"); got != ".gz" {
		t.Fatalf("Ext() = %q, want .gz", got)
	}
}

func TestJoin(t *testing.T) {
	if got := Join("a", "b", "c.txt"); got != filepath.Join("a", "b", "c.txt") {
		t.Fatalf("Join() = %q", got)
	}
}

func TestMissingPathChecks(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.txt")
	if Exists(path) {
		t.Fatal("Exists() = true, want false")
	}
	if IsFile(path) {
		t.Fatal("IsFile() = true, want false")
	}
	if IsDir(path) {
		t.Fatal("IsDir() = true, want false")
	}
}
