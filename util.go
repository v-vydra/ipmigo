package ipmigo

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func toJSON(s interface{}) string {
	r, _ := json.Marshal(s)
	return string(r)
}

func retry(retries int, f func() error) (err error) {
	for i := 0; i <= retries; i++ {
		err = f()
		switch e := err.(type) {
		case net.Error:
			if e.Timeout() {
				continue
			}
		}
		return
	}
	return
}

// ConvertBoardMfgDate converts a 3-byte Board Mfg Date from an IPMI FRU response to time.Time.
// The input `data` should contain at least 3 bytes starting from the Board Mfg Date field.
func ConvertBoardMfgDate(data []byte) (time.Time, error) {
	// Check if input slice has at least 3 bytes for the date
	if len(data) < 3 {
		return time.Now(), fmt.Errorf("data slice is too short, expected at least 3 bytes")
	}

	// IPMI FRU dates are offset in minutes from 1/1/1996
	baseDate := time.Date(1996, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Read the 3-byte value as a little-endian unsigned integer
	offsetMinutes := uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16

	// Calculate the actual date by adding the offset
	manufactureDate := baseDate.Add(time.Duration(offsetMinutes) * time.Minute)

	return manufactureDate, nil
}
