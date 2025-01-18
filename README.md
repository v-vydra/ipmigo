ipmigo
======

**Work In Progress**

About
-----

__Forked from [github.com/k-sone/ipmigo](github.com/k-sone/ipmigo)__

ipmigo is a golang implementation for IPMI client.

ChangeLog
---------
2025-01-18
* added SetSDRReadingBytes and SetFRUReadingBytes for faster SDR and FRU areas reading (SDR/FRU defaults: 32/16 - works for my 255/63)
* added FRU Board and Product Areas handling - ToDo: Chassis and MultiRecord Areas
* added example how to retrieve all devices FRU Board and Product area information [examples/fru/fru.go](https://github.com/v-vydra/ipmigo/examples/fru.go])

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
