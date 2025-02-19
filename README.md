ipmigo
======

_Work in progress_, but I use it personally as it is now in git _without any issues_. 

About
-----

**ipmigo** is a golang implementation of **IPMI** client. 

__Forked from [github.com/k-sone/ipmigo](github.com/k-sone/ipmigo)__

All new functionality is tested against Wistron OpenBMC OEM IPMI + for functionality not depending on OEM, most of the time, also validating against HPE iLO IPMI and Supermicro IPMI.

ChangeLog
---------
2025-01-21  
* added ClearSELCommand - Clear SEL Command  
* added helper function ClearSELWaitFinish - will initiate SEL Erase and wait maxWaitSeconds to complete  
  - example at [examples/sel/sel-clear.go](https://github.com/v-vydra/ipmigo/blob/master/examples/sel/sel-clear.go)  

2025-01-20
* added GetOEMOpenBmcWistronXcvrPortPageCommand/SetOEMOpenBmcWistronXcvrPortPageCommand - Get/Set XCVR port page RAW data on Wistron OEM OpenBMC  
    - example at [examples/openbmc-wistron/xcvr-eeprom.go](https://github.com/v-vydra/ipmigo/blob/master/examples/openbmc-wistron/xcvr-eeprom.go)  
* added GetOEMOpenBmcWistronFirmwareInfoCommand - Wistron OEM OpenBMC Get Firmware Information Command  
    - example at [examples/openbmc-wistron/fw-version.go](https://github.com/v-vydra/ipmigo/blob/master/examples/openbmc-wistron/fw-version.go)  
* added GetOEMOpenBmcWistronFanControlCommand/SetOEMOpenBmcWistronFanControlCommand - Wistron OEM OpenBMC Set Fan Speed Control Command  
  - example at [examples/openbmc-wistron/fan.go](https://github.com/v-vydra/ipmigo/blob/master/examples/openbmc-wistron/fan.go)  

2025-01-19  
* added GetSetOEMOpenBmcWistronI2CCommand - Read/Write via I2C over iPMI on Wistron's OpenBMC
* added NetFnOemOne 0x30
* added chassis power control example [examples/power/power.go](https://github.com/v-vydra/ipmigo/blob/master/examples/power/power.go)
* new SetChassisControlCommand - power control
* FRU example works on OpenBMC Wistron and HPE iLO iPMI  

2025-01-18  
* added SetSDRReadingBytes and SetFRUReadingBytes for faster SDR and FRU areas reading (SDR/FRU defaults: 32/16 - works for my 255/63)
* added FRU Board and Product Areas handling - ToDo: Chassis and MultiRecord Areas
* added example how to retrieve all devices FRU Board and Product area information [examples/fru/fru.go](https://github.com/v-vydra/ipmigo/blob/master/examples/fru/fru.go)

2025-01-17  
* added IPMI FRU InventoryAreaInfo and Data command

Supported Version
-----------------

* IPMI v2.0(lanplus)

Examples
--------

```go
package main

import (
    "fmt"

    "github.com/v-vydra/ipmigo"
)

func main() {
    c, err := ipmigo.NewClient(ipmigo.Arguments{
        Version:       ipmigo.V2_0,
        Address:       "192.168.1.1:623",
        Username:      "myuser",
        Password:      "mypass",
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

    cmd := &ipmigo.GetPOHCounterCommand{}
    if err := c.Execute(cmd); err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Power On Hours", cmd.PowerOnHours())
}
```

License
-------

MIT
