package main

import (
	"os"
	"path/filepath"
	"testing"
)

// Test getPasswd
func TestGetPasswd(t *testing.T) {
	tmpdir := t.TempDir()
	// create tmp lastlog
	lastlog, err := os.Create(filepath.Join(tmpdir, "lastlog"))
	if err != nil {
		t.Fatalf("Failed to create temporary lastlog file: %v", err)
	}
	lastlog.Close()
	// create test passwd
	file, err := os.Create(filepath.Join(tmpdir, "passwd"))
	if err != nil {
		t.Fatalf("Failed to create temporary passwd file: %v", err)
	}
	defer file.Close()
	_, err = file.WriteString(`root:x:0:0:root:/root:/bin/bash
bin:x:1:1:bin:/bin:/sbin/nologin
daemon:x:2:2:daemon:/sbin:/sbin/nologin
adm:x:3:4:adm:/var/adm:/sbin/nologin
lp:x:4:7:lp:/var/spool/lpd:/sbin/nologin
sync:x:5:0:sync:/sbin:/bin/sync
shutdown:x:6:0:shutdown:/sbin:/sbin/shutdown
halt:x:7:0:halt:/sbin:/sbin/halt
mail:x:8:12:mail:/var/spool/mail:/sbin/nologin
operator:x:11:0:operator:/root:/sbin/nologin
games:x:12:100:games:/usr/games:/sbin/nologin
ftp:x:14:50:FTP User:/var/ftp:/sbin/nologin
nobody:x:65534:65534:Kernel Overflow User:/:/sbin/nologin
kevin:x:1005:1006::/home/kevin:/usr/bin/zsh
`)
	if err != nil {
		t.Fatalf("Failed to write to temporary passwd file: %v", err)
	}

	opt := &Opt{
		PasswdFile:  filepath.Join(tmpdir, "passwd"),
		LastLogFile: filepath.Join(tmpdir, "lastlog"),
	}
	passwd, err := opt.getPasswd()
	if err != nil {
		t.Fatalf("Failed to get passwd: %v", err)
	}
	if len(passwd) != 14 {
		t.Fatalf("passwd entry mismtached: expected: %d, actual: %d", 14, len(passwd))
	}
	if passwd[13].UserName != "kevin" {
		t.Fatalf("last passwd entry mismtached: expected: %s, actual: %s", "kevin", passwd[13].UserName)
	}
}
