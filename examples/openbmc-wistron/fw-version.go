package main

import (
	"fmt"
	"github.com/v-vydra/ipmigo"
	"time"
)

func getFirmwareInfo(ip string, devId uint8) {
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

	cmd := &ipmigo.GetOEMOpenBmcWistronFirmwareInfoCommand{
		DevId: devId,
	}
	err = c.Execute(cmd)
	if err != nil {
		fmt.Printf("Error when getting Firmware Information for FRU Device %d: %s\n", devId, err)
	} else {
		fmt.Printf("FRU Device %d Firmware information data : %s\n", devId, cmd.GetFirmwareString())
	}
}

func main() {
	fmt.Printf("Firmware Information for Primary BMC:\n")
	getFirmwareInfo("172.30.1.241", 1)
}
