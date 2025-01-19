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

	// temporary fix for IPMI devices which don't list main device id 0 in SDR, ex.: HPE iLO Builtin FRU Device 0
	foundMainDevice := false

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
			// if well defined:
			//fmt.Printf("%s\n", fru.ToString())

			// list all
			fmt.Printf("Board Info Area\n%sProduct Info Area\n%s\n",
				fru.GetBoardInfoAreaAsString(),
				fru.GetProductInfoAreaFieldsAsString(),
			)

			if s.DeviceID == 0 {
				foundMainDevice = true
			}
		}
	}

	// temporary fix for IPMI devices which don't list main device id 0 in SDR, ex.: HPE iLO Builtin FRU Device 0
	if !foundMainDevice {
		fmt.Printf("Warning: No device with ID 0 found in SDR - attempt to get directly\n")
		fru, err := ipmigo.FRUGetDeviceData(c, 0, 0)
		if err != nil {
			fmt.Printf("Warning: Device ID 0 is really missing on this system: %v\n", err)
		} else {
			fmt.Printf("%s\n", fru.CommonHeader.ToString())
			fmt.Printf("Board Info Area\n%sProduct Info Area\n%s\n",
				fru.GetBoardInfoAreaAsString(),
				fru.GetProductInfoAreaFieldsAsString(),
			)

		}
	}
}
