package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func (opt *Opt) getPasswd() ([]User, error) {
	users := make([]User, 0)
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
		// kevin:x:1005:1006::/home/kevin:/usr/bin/zsh
		parts := strings.SplitN(line, ":", 7)
		if len(parts) < 6 || parts[0] == "" ||
			parts[0][0] == '+' || parts[0][0] == '-' {
			continue
		}
		shell := ""
		if len(parts) >= 7 {
			shell = parts[6]
		}
		uid, err := strconv.Atoi(parts[2])
		if err != nil {
			return users, err
		}
		ll, ok := lastLog[uid]
		if !ok {
			ll = 0
		}
		u := User{
			UID:      uid,
			UserName: parts[0],
			Shell:    shell,
			LastLog:  ll,
		}
		users = append(users, u)
	}
	if err := s.Err(); err != nil {
		return users, err
	}
	return users, nil
}
