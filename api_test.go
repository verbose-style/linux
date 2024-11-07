package linux_test

import (
	"encoding/json"
	"linux"
	"os"
	"testing"
)

func TestLinux(t *testing.T) {
	var Linux = linux.Native()

	header, err := Linux.Stat("./api_test.go")
	if err != nil {
		t.Fatal(err)
	}
	json.NewEncoder(os.Stdout).Encode(header)

	if _, err := Linux.Stat("./nothing"); err != new(linux.StatError).Types().DoesNotExist {
		t.Fatal(err)
	}
}
