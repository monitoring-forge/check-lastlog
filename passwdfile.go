package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func (opt *Opt) parsePasswdUser(line string, lastLog map[int]int64) *User {
	// Skip empty lines and comments
	if line == "" || strings.HasPrefix(line, "#") {
		return nil
	}
	// kevin:x:1005:1006::/home/kevin:/usr/bin/zsh
	parts := strings.SplitN(line, ":", 7)
	// Need at least username, placeholder, and UID
	if len(parts) < 3 || parts[0] == "" ||
		parts[0][0] == '+' || parts[0][0] == '-' {
		return nil
	}
	// Parse UID
	uid, err := strconv.Atoi(parts[2])
	if err != nil {
		// Skip invalid UID entries
		return nil
	}
	// Get shell (if present)
	shell := ""
	if len(parts) >= 7 {
		shell = parts[6]
	}
	// Get last login time from lastlog
	ll, ok := lastLog[uid]
	if !ok {
		ll = 0 // Never logged in
	}
	return &User{
		UID:      uid,
		UserName: parts[0],
		Shell:    shell,
		LastLog:  ll,
	}
}

func (opt *Opt) getPasswd() ([]*User, error) {
	users := make([]*User, 0)
	lastLog, err := opt.getLastLog()
	if err != nil {
		return users, err
	}

	f, err := os.Open(opt.PasswdFile)
	if err != nil {
		return users, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		u := opt.parsePasswdUser(line, lastLog)
		if u != nil {
			users = append(users, u)
		}
	}
	if err := s.Err(); err != nil {
		return users, err
	}
	return users, nil
}
