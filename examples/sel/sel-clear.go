package main

import (
	"encoding/hex"
	"fmt"
	"github.com/v-vydra/ipmigo"
	"time"
)

func getSelEntries(ip string) {
	c, err := ipmigo.NewClient(ipmigo.Arguments{
		Version:       ipmigo.V2_0,
		Address:       fmt.Sprintf("%s:623", ip),
		Timeout:       3 * time.Second,
		Retries:       3,
		Username:      "root",
		Password:      "0penBmc",
		CipherSuiteID: 3,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if err := c.Open(); err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	// Get total count
	_, total, err := ipmigo.SELGetEntries(c, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Total SEL entries: %d\n", total)
	fmt.Printf("RecordID | Timestamp | SensorType | SensorNumber | Direction | Description\n")

	// Get latest 50 events
	count := 50
	offset := 0
	if total > count {
		offset = total - count
	}
	records, total, err := ipmigo.SELGetEntries(c, offset, count)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Output events
	for _, r := range records {
		switch s := r.(type) {
		case *ipmigo.SELEventRecord:
			dir := "Asserted"
			if !s.IsAssertionEvent() {
				dir = "De-asserted"
			}
			fmt.Printf("%-4d | %-25s | %-25s(0x%02x) | %-10s | %s\n",
				s.RecordID, &s.Timestamp, s.SensorType, s.SensorNumber, dir, s.Description())

		case *ipmigo.SELTimestampedOEMRecord:
			fmt.Printf("%-4d | %-25s | 0x%08x | 0x%s\n",
				s.RecordID, &s.Timestamp, s.ManufacturerID, hex.EncodeToString(s.OEMDefined))

		case *ipmigo.SELNonTimestampedOEMRecord:
			fmt.Printf("%-4d | 0x%s\n", s.RecordID, hex.EncodeToString(s.OEM))
		default:
			fmt.Printf("unknown record type: %+v", s)
		}
	}
}
func clearSel(ip string) {
	c, err := ipmigo.NewClient(ipmigo.Arguments{
		Version:       ipmigo.V2_0,
		Address:       fmt.Sprintf("%s:623", ip),
		Timeout:       3 * time.Second,
		Retries:       3,
		Username:      "root",
		Password:      "0penBmc",
		CipherSuiteID: 3,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if err := c.Open(); err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	err = ipmigo.ClearSELWaitFinish(c, 10)
	if err != nil {
		fmt.Printf("SEL Erasure error: %s", err)
	} else {
		fmt.Println("SEL Erasure success")
	}
}
func main() {
	ip := "172.30.1.241"
	getSelEntries(ip)
	clearSel(ip)
	getSelEntries(ip)
}
