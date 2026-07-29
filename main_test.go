package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// test run
func TestRun(t *testing.T) {
	tmpdir := t.TempDir()
	file, err := os.Create(filepath.Join(tmpdir, "lastlog"))
	if err != nil {
		t.Fatalf("Failed to create temporary lastlog file: %v", err)
	}
	defer file.Close()
	now := time.Now().Unix()
	testUsers := []User{
		{UserName: "user0", UID: 0, Shell: "/bin/bash", LastLog: 0},
		{UserName: "user1", UID: 1, Shell: "/bin/bash", LastLog: now - 2*86400},      // 2 days ago
		{UserName: "user2", UID: 2, Shell: "/bin/bash", LastLog: now - 5*86400},      // 5 days ago
		{UserName: "user3", UID: 3, Shell: "/sbin/nologin", LastLog: now - 10*86400}, // 10 days ago
		{UserName: "user4", UID: 4, Shell: "/bin/bash", LastLog: now - 10*86400},     // 10 days ago
	}

	for _, u := range testUsers {
		buf := make([]byte, llsize)

		binary.LittleEndian.PutUint64(buf[:ttsize], uint64(u.LastLog))
		_, err = file.Write(buf)
		if err != nil {
			t.Fatalf("Failed to write to temporary lastlog file: %v", err)
		}
	}

	// Create a temporary passwd file
	passwdFile, err := os.Create(filepath.Join(tmpdir, "passwd"))
	if err != nil {
		t.Fatalf("Failed to create temporary passwd file: %v", err)
	}
	defer passwdFile.Close()

	for _, u := range testUsers {
		line := u.UserName + ":x:" + fmt.Sprintf("%d", u.UID) + ":1000::/home/" + u.UserName + ":" + u.Shell + "\n"
		_, err = passwdFile.WriteString(line)
		if err != nil {
			t.Fatalf("Failed to write to temporary passwd file: %v", err)
		}
	}

	opt := Opt{
		LastLogFile: filepath.Join(tmpdir, "lastlog"),
		PasswdFile:  filepath.Join(tmpdir, "passwd"),
		MaxUID:      100,
		MinUID:      1,
		Warn:        3, // 3 days
		Crit:        6, // 6 days
		Verbose:     true,
	}
	chk := opt.run()
	require.NotNil(t, chk, "Checker should not be nil")
	require.Equal(t, checkers.CRITICAL, chk.Status, "Expected status to be CRITICAL")
	assert.Equal(t, chk.Message, "Found users who have not logged in recently: user2(5 days), user4(10 days)", "Expected message to be CRITICAL")

	opt = Opt{
		LastLogFile: filepath.Join(tmpdir, "lastlog"),
		PasswdFile:  filepath.Join(tmpdir, "passwd"),
		MaxUID:      100,
		MinUID:      1,
		Warn:        4,  // 4 days
		Crit:        20, // 20 days
	}
	chk = opt.run()
	require.NotNil(t, chk, "Checker should not be nil")
	require.Equal(t, checkers.WARNING, chk.Status, "Expected status to be WARNING")
	assert.Equal(t, chk.Message, "Found users who have not logged in recently: user2(5 days), user4(10 days)", "Expected message to be WARNING")

	opt = Opt{
		LastLogFile: filepath.Join(tmpdir, "lastlog"),
		PasswdFile:  filepath.Join(tmpdir, "passwd"),
		MaxUID:      100,
		MinUID:      1,
		Warn:        20, // 20 days
		Crit:        30, // 30 days
	}
	chk = opt.run()
	require.NotNil(t, chk, "Checker should not be nil")
	require.Equal(t, checkers.OK, chk.Status, "Expected status to be OK")
	assert.Contains(t, chk.Message, "No users were found who have not logged in recently", "Expected message to indicate no users found")

	opt = Opt{
		LastLogFile:    filepath.Join(tmpdir, "lastlog"),
		PasswdFile:     filepath.Join(tmpdir, "passwd"),
		MaxUID:         100,
		MinUID:         0,
		Warn:           3, // 3 days
		Crit:           6, // 6 days
		WhiteUserNames: "user0,user1",
		Verbose:        true,
	}
	chk = opt.run()
	require.NotNil(t, chk, "Checker should not be nil")
	require.Equal(t, checkers.CRITICAL, chk.Status, "Expected status to be CRITICAL")
	assert.NotContains(t, chk.Message, "user0", "Expected message to not contain user0 (UID 0)")
}
