package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

var version string
var commit string

type Opt struct {
	Before         int64  `long:"before" description:"[Deprecated] Check for users whose login is older than DAYS"`
	Warn           int64  `short:"w" long:"warning" default:"60" description:"warning if users whose login is older than DAYS"`
	Crit           int64  `short:"c" long:"critical" default:"85" description:"critical if users whose login is older than DAYS"`
	MinUID         int    `long:"min-uid" default:"500" description:"min uid to check lastlog"`
	MaxUID         int    `long:"max-uid" default:"60000" description:"max uid to check lastlog"`
	WhiteUserNames string `long:"white-user-names" default:"" description:"comma-separated user names to whitelist"`
	Version        bool   `short:"v" long:"version" description:"Show version"`
	LastLogFile    string `long:"lastlog-file" default:"/var/log/lastlog" description:"lastlog file path"`
	PasswdFile     string `long:"passwd-file" default:"/etc/passwd" description:"passwd file path"`
	Verbose        bool   `short:"V" long:"verbose" description:"Show verbose log"`
}

func (opt *Opt) run() *checkers.Checker {
	whiteUserNames := make(map[string]struct{})
	if opt.WhiteUserNames != "" {
		names := strings.Split(opt.WhiteUserNames, ",")
		for _, n := range names {
			whiteUserNames[n] = struct{}{}
		}
	}

	now := time.Now().Unix()
	warn := now - opt.Warn*86400
	crit := now - opt.Crit*86400

	// for compatibility
	if opt.Before != 0 {
		warn = now - opt.Before*86400
		crit = now - opt.Before*86400
	}

	hasCrit := false
	hasWarn := false
	msgs := make([]string, 0)

	users, err := opt.getPasswd()
	if err != nil {
		return checkers.Unknown(fmt.Sprintf("UNKNOWN: %v", err))
	}
	for _, u := range users {
		if opt.Verbose {
			fmt.Fprintf(os.Stderr, "DEBUG: user=%s, uid=%d, shell=%s, lastlog=%s, lastlogin=%s\n", u.UserName, u.UID, u.Shell, u.LastLoginDays(), u.LastLogTime().Format(time.RFC3339))
		}
		if u.UID <= opt.MinUID {
			continue
		}
		if u.UID >= opt.MaxUID {
			continue
		}
		if _, ok := whiteUserNames[u.UserName]; ok {
			continue
		}
		if u.NoLogin() {
			continue
		}

		if u.LastLog < crit {
			hasCrit = true
			msgs = append(msgs, fmt.Sprintf("%s(%s)", u.UserName, u.LastLoginDays()))
		} else if u.LastLog < warn {
			hasWarn = true
			msgs = append(msgs, fmt.Sprintf("%s(%s)", u.UserName, u.LastLoginDays()))
		}
	}

	if hasCrit {
		// crit
		return checkers.Critical(fmt.Sprintf("CRITICAL: Found users who have not logged in recently: %s", strings.Join(msgs, ", ")))
	} else if hasWarn {
		// warn
		return checkers.Warning(fmt.Sprintf("WARNING: Found users who have not logged in recently: %s", strings.Join(msgs, ", ")))
	}

	// ok
	return checkers.Ok("OK: No users were found who have not logged in recently")
}

func main() {
	opt := Opt{}
	psr := flags.NewParser(&opt, flags.HelpFlag|flags.PassDoubleDash)
	_, err := psr.Parse()
	if opt.Version {
		if commit == "" {
			commit = "dev"
		}
		fmt.Printf(
			"%s-%s\n%s/%s, %s, %s\n",
			filepath.Base(os.Args[0]),
			version,
			runtime.GOOS,
			runtime.GOARCH,
			runtime.Version(),
			commit)
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	ckr := opt.run()
	ckr.Name = "check-lastlog"
	ckr.Exit()
}
