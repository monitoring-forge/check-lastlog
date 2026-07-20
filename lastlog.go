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

	for pos := 0; pos <= opt.MaxUID; pos++ {
		_, err := io.ReadFull(f, buf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			// End of file reached, which is normal
			break
		}
		if err != nil {
			return lastlog, err
		}

		// Read the Unix time (8 bytes for time_t)
		unixTime := int64(binary.LittleEndian.Uint64(buf[:ttsize]))
		lastlog[pos] = unixTime
	}
	return lastlog, nil
}
