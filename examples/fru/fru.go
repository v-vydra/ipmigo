package main

import (
	"fmt"
	"github.com/v-vydra/ipmigo"
	"time"
)

func main() {
	c, err := ipmigo.NewClient(ipmigo.Arguments{
		Version:       ipmigo.V2_0,
		Address:       "172.30.1.241:623",
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

	c.SetSDRReadingBytes(32)
	c.SetFRUReadingBytes(63)

	if err := c.Open(); err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	// Get sensor records
	records, err := ipmigo.SDRGetRecordsRepo(c, func(id uint16, t ipmigo.SDRType) bool {
		return t == ipmigo.SDRTypeFRUDeviceLocator
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, r := range records {
		switch s := r.(type) {
		case *ipmigo.SDRFRUDeviceLocator:
			fmt.Printf("*** Device ID: %d\n", s.DeviceID)
			fru, err := ipmigo.FRUGetDeviceData(c, s.DeviceID, s.AccessLUN)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("%s\n", fru.String())
		}
	}

}
