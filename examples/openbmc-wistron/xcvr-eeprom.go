package main

import (
	"encoding/hex"
	"fmt"
	"github.com/v-vydra/ipmigo"
	"time"
)

func getXcvrPage(ip string, port uint8, page uint8) {
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

	cmd := &ipmigo.GetOEMOpenBmcWistronXcvrPortPageCommand{
		Port:     port,
		Function: 1,
		Page:     page,
		Offset:   0,
		Length:   0xFF,
	}
	err = c.Execute(cmd)
	if err != nil {
		fmt.Printf("Error when getting XCVR Port %d Full Page %d data: %s\n", port, page, err)
	} else {
		fmt.Printf("XCVR Port %d full Page %d Data : %s\n", port, page, hex.EncodeToString(cmd.Output()))
	}
}

func main() {
	getXcvrPage("172.30.1.241", 5, 0)
}
