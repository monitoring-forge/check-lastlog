package main

import (
	"encoding/binary"
	"io"
	"os"
)

// linux/amd64: struct lastlog { time_t(8) + line[32] + host[256] }
const (
	ttsize = 8
	llsize = 296
)

func (opt *Opt) getLastLog() (map[int]int64, error) {
	lastlog := make(map[int]int64)
	f, err := os.Open(opt.LastLogFile)
	if err != nil {
		return lastlog, err
	}
	defer f.Close()
	buf := make([]byte, llsize)
	pos := 0
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return lastlog, err
		}
		if n == 0 {
			break
		}
		unixTime := int64(binary.LittleEndian.Uint32(buf[:ttsize]))
		lastlog[pos] = unixTime
		pos++
	}
	return lastlog, nil
}
