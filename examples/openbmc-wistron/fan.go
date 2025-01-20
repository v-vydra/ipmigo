package main

import (
	"fmt"
	"github.com/v-vydra/ipmigo"
	"time"
)

func getFansSpeed(ip string) {
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

	cmd := &ipmigo.GetOEMOpenBmcWistronFanControlCommand{}
	err = c.Execute(cmd)
	if err != nil {
		fmt.Printf("Error when getting FAN Speed Setting: %s\n", err)
	} else {
		fmt.Printf("Actual FAN Control Mode: %s, PWM Duty Cycle = %d%%\n", cmd.ControlMode, cmd.FanSpeed)
	}
}

func setFansSpeed(ip string, controlMode ipmigo.ControlMode, fanSpeed uint8) {
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

	cmd := &ipmigo.SetOEMOpenBmcWistronFanControlCommand{
		ControlMode: controlMode,
		FanSpeed:    fanSpeed,
	}
	err = c.Execute(cmd)
	if err != nil {
		fmt.Printf("Error when getting FAN Speed Setting: %s\n", err)
	} else {
		fmt.Printf("New FAN Control Mode: %s", cmd.ControlMode)
		if controlMode == ipmigo.ControlModeManual {
			fmt.Printf(", PWM Duty Cycle = %d%%\n", fanSpeed)
		} else {
			fmt.Println()
		}
	}
}

func main() {
	openBmcIp := "172.30.1.241"
	getFansSpeed(openBmcIp)
	setFansSpeed(openBmcIp, ipmigo.ControlModeManual, 70)
	getFansSpeed(openBmcIp)
	setFansSpeed(openBmcIp, ipmigo.ControlModeAuto, 0) // fanSpeed not needed in this case
	getFansSpeed(openBmcIp)
}
