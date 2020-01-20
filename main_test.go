package main

import (
	"os"
	"testing"
)

func TestHash(t *testing.T) {
	f, _ := createFakeFile()
	defer f.Close()
	got, err := hash(f.Name())
	// empty string hash is alwasy the same
	want := emptyStringHash()
	if err != nil {
		t.Error("error getting hash of file")
	}
	if got != want {
		t.Errorf("has() = %q want %q", got, want)
	}
}

func TestLoadHashes(t *testing.T) {
	f, _ := createFakeFile()
	defer f.Close()
	os.Setenv(configFilesEnvVar, f.Name())
	hashes, err := loadHashes()
	if err != nil {
		t.Fatal("error", err)
	}
	want := emptyStringHash()
	got := hashes["empty_file_test"]
	if want != got {
		t.Errorf("loadHashes() got: %q and want %q", got, want)
	}

	// when unset env var
	os.Unsetenv(configFilesEnvVar)
	_, err = loadHashes()
	if err == nil {
		t.Error("should return an error when no env var exist")
	}
}

func emptyStringHash() string {
	return "d41d8cd98f00b204e9800998ecf8427e"
}

func createFakeFile() (*os.File, error) {
	return os.Create("empty_file_test")
}
