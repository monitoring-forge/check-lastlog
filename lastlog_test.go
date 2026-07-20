package main

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

// TestLastLog : test lastlog.go
func TestLastLog(t *testing.T) {
	// create temporary lastlog file
	tmpdir := t.TempDir()
	file, err := os.Create(filepath.Join(tmpdir, "lastlog"))
	if err != nil {
		t.Fatalf("Failed to create temporary lastlog file: %v", err)
	}
	defer file.Close()

	// Write a sample lastlog entry for UID 0 (root) with a timestamp of 1 day ago
	unixTime := int64(0) // 0 means never logged in
	buf := make([]byte, llsize)
	binary.LittleEndian.PutUint32(buf[:ttsize], uint32(unixTime))
	_, err = file.Write(buf)
	if err != nil {
		t.Fatalf("Failed to write to temporary lastlog file: %v", err)
	}

	opt := Opt{
		LastLogFile: filepath.Join(tmpdir, "lastlog"),
	}
	lastlog, err := opt.getLastLog()
	if err != nil {
		t.Fatalf("Failed to get lastlog: %v", err)
	}
	if lastlog[0] != unixTime {
		t.Errorf("Expected lastlog for UID 0 to be %d, got %d", unixTime, lastlog[0])
	}
	// Write a sample lastlog entry for UID 1 (user) with a timestamp of 2 days ago
	unixTime = int64(2 * 86400) // 2 days ago
	binary.LittleEndian.PutUint32(buf[:ttsize], uint32(unixTime))
	_, err = file.Write(buf)
	if err != nil {
		t.Fatalf("Failed to write to temporary lastlog file: %v", err)
	}

	lastlog, err = opt.getLastLog()
	if err != nil {
		t.Fatalf("Failed to get lastlog: %v", err)
	}
	if lastlog[1] != unixTime {
		t.Errorf("Expected lastlog for UID 1 to be %d, got %d", unixTime, lastlog[1])
	}
}
