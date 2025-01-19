package main

import (
	"fmt"
	"github.com/v-vydra/ipmigo"
	"time"
)

// String prompt.
func String(prompt string, args ...interface{}) string {
	var s string
	fmt.Printf(prompt+": ", args...)
	fmt.Scanln(&s)
	return s
}

// Confirm continues prompting until the input is boolean-ish.
func Confirm(prompt string, args ...interface{}) bool {
	for {
		switch String(prompt, args...) {
		case "YeS":
			return true
		default:
			return false
		}
	}
}

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

	if err := c.Open(); err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	fmt.Printf("Geting chassis board informations ...\n")
	fru, err := ipmigo.FRUGetDeviceData(c, 0, 0)
	if err != nil {
		fmt.Printf(" Device ID 0 is missing on this system: %v\n", err)
		return
	}

	fmt.Printf("Geting current chassis status ...\n")
	cmdGetChassisStatus := &ipmigo.GetChassisStatusCommand{}
	if err := c.Execute(cmdGetChassisStatus); err != nil {
		fmt.Printf("unable to get chassis status: %v", err)
		return
	}

	if fru.BoardInfo == nil || len(fru.BoardInfo.Fields) < 3 {
		fmt.Printf("*** No board info or serial found in FRU\n")
		return
	}

	boardManufacturer := fru.BoardInfo.GetFieldValueStringById(1)
	boardProduct := fru.BoardInfo.GetFieldValueStringById(2)
	boardSerial := fru.BoardInfo.GetFieldValueStringById(3)
	if len(boardManufacturer) < 1 {
		fmt.Printf("*** No Valid board manufacturer found in FRU\n")
		return
	}
	if len(boardProduct) < 1 {
		fmt.Printf("*** No Valid board product found in FRU\n")
		return
	}
	if len(boardSerial) < 1 {
		fmt.Printf("*** No Valid board serial found in FRU\n")
		return
	}

	fmt.Printf("Board manufacturer           : %s\n", boardManufacturer)
	fmt.Printf("Board product                : %s\n", boardProduct)
	fmt.Printf("Board serial                 : %s\n", boardSerial)

	fmt.Printf("Current chassis power status : ")
	if cmdGetChassisStatus.PowerIsOn {
		fmt.Printf(" ON\n")
	} else {
		fmt.Printf(" OFF\n")
	}

	fmt.Printf("\n\n")
	if cmdGetChassisStatus.PowerIsOn {
		powerChange := Confirm("Power is ON, do you really want to power off? [YeS/n]")
		if powerChange {
			fmt.Printf("Sending power DOWN command ...\n")

			cmdChassisControl := &ipmigo.SetChassisControlCommand{
				ChassisControl: ipmigo.ChassisControlPowerDown,
			}
			if err := c.Execute(cmdChassisControl); err != nil {
				fmt.Printf("unable to set chassis control: %v", err)
				return
			}
		} else {
			fmt.Printf("Nothing to do\n")
		}

	} else {
		powerChange := Confirm("Power is OFF, do you really want to power on? [YeS/n]")
		if powerChange {
			fmt.Printf("Sending power UP command ...\n")

			cmdChassisControl := &ipmigo.SetChassisControlCommand{
				ChassisControl: ipmigo.ChassisControlPowerUp,
			}
			if err := c.Execute(cmdChassisControl); err != nil {
				fmt.Printf("unable to set chassis control: %v", err)
				return
			}
		} else {
			fmt.Printf("Nothing to do\n")
		}
	}
}
