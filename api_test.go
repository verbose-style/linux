package linux_test

import (
	"encoding/json"
	"os"
	"testing"

	"verbose.style/linux"
	"verbose.style/linux/internal"
)

func TestLinux(t *testing.T) {
	var Linux = linux.Native()

	header, err := Linux.Stat("./api_test.go")
	if err != nil {
		t.Fatal(err)
	}
	json.NewEncoder(os.Stdout).Encode(header)

	f, err := Linux.Open("./api_test.go", linux.FileAccessReadOnly, 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	mmap, err := Linux.MapFileIntoMemory(nil, int(header.Size), linux.MemoryAllowReads, linux.MapPrivate, 0, f.Descriptor, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := mmap.WriteAt([]byte{1}, 0); err == nil {
		t.Fatal("expected error")
	}
}

func TestValues(t *testing.T) {
	internal.Test(t)
}
