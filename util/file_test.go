package util

import (
	"code.google.com/p/go-uuid/uuid"

	"os"
	"path"
	"testing"
)

func TestAssertFileState(t *testing.T) {
	var (
		directory = os.TempDir()
		filename  = path.Join(directory, uuid.New())
	)
	if f, err := os.Create(filename); err == nil {
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
	if err := AssertFileState(filename, true); err != nil {
		t.Fatal("unexpected error")
	}
	if err := AssertFileState(filename, false); err == nil {
		t.Fatal("error expected")
	}
	if err := os.Remove(filename); err != nil {
		t.Fatal(err)
	}
	if err := AssertFileState(filename, true); err == nil {
		t.Fatal("error expected")
	}
	if err := AssertFileState(filename, false); err != nil {
		t.Fatal("unexpected error")
	}
}
