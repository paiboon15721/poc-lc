package main

import (
	"bufio"
	"errors"
	"os"
	"unicode/utf8"
)

func getRealHardwareID() (realHardwareID string, err error) {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "", err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanHardwareID := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		start := 0
		for width := 0; start < len(data); start += width {
			var r rune
			r, width = utf8.DecodeRune(data[start:])
			if r == 'S' {
				if string(data[start:start+6]) == "Serial" {
					start += 6
					break
				}
			}
		}
		for width, i := 0, start; i < len(data); i += width {
			var r rune
			r, width = utf8.DecodeRune(data[i:])
			if r == '\n' {
				return i + width, data[start:i], nil
			}
		}
		if atEOF && len(data) > start {
			return len(data), data[start:], nil
		}
		return start, nil, nil
	}
	scanner.Split(scanHardwareID)
	if scanner.Scan() {
		return "", errors.New("can not get hardwareID")
	}
	realHardwareID = scanner.Text()
	realHardwareID = realHardwareID[len(realHardwareID)-16:]
	return realHardwareID, nil
}
