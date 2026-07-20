package main

import (
	"testing"
	"time"
)

// Test Nologin shell
func TestNoLogin(t *testing.T) {
	shells := []string{
		"/bin/sync",
		"/sbin/halt",
		"/sbin/nologin",
		"/sbin/shutdown",
		"/usr/sbin/nologin",
		"/usr/bin/false",
		"/bin/false",
	}

	for _, shell := range shells {
		u := User{Shell: shell}
		if !u.NoLogin() {
			t.Errorf("User with shell %s should be NoLogin", shell)
		}
	}

	u := User{Shell: "/bin/bash"}
	if u.NoLogin() {
		t.Errorf("User with shell /bin/bash should not be NoLogin")
	}
}

// Test LastLoginDays
func TestLastLoginDays(t *testing.T) {
	u := User{LastLog: 0}
	if u.LastLoginDays() != "*Never logged in*" {
		t.Errorf("Expected '*Never logged in*', got %s", u.LastLoginDays())
	}

	// Set LastLog to 2 days ago
	twoDaysAgo := time.Now().Unix() - 2*86400
	u.LastLog = twoDaysAgo
	expected := "2 days"
	if u.LastLoginDays() != expected {
		t.Errorf("Expected '%s', got %s", expected, u.LastLoginDays())
	}
}
