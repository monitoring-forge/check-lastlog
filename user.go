package main

import (
	"fmt"
	"time"
)

// User :
type User struct {
	UID      int
	UserName string
	Shell    string
	LastLog  int64
}

// LastLogTime : user.LastLog as time.Time
func (u *User) LastLogTime() time.Time {
	return time.Unix(u.LastLog, 0)
}

var noLoginShell = map[string]struct{}{
	"/bin/sync":         {},
	"/sbin/halt":        {},
	"/sbin/nologin":     {},
	"/sbin/shutdown":    {},
	"/usr/sbin/nologin": {},
	"/usr/bin/false":    {},
	"/bin/false":        {},
}

// NoLogin : User has nologin shell
func (u *User) NoLogin() bool {
	_, ok := noLoginShell[u.Shell]
	return ok
}

// LastLoginDays :
func (u *User) LastLoginDays() string {
	if u.LastLog == 0 {
		return "*Never logged in*"
	}
	t := time.Now().Unix() - u.LastLog
	if t < 0 {
		t = 0
	}
	return fmt.Sprintf("%d days", int(t/86400))
}
