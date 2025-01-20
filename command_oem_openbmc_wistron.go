package ipmigo

import "fmt"

// *** OpenBMC Wistron I2C

// GetSetOEMOpenBmcWistronI2CCommand Wistron OEM OpenBMC I2C Read Write Command
type GetSetOEMOpenBmcWistronI2CCommand struct {
	// Request Data
	BusID        uint8  // 0-based (count from 0)
	SlaveAddress uint8  // bits 7:1 Slave Address (7-bit), bit 0 - reserved
	ReadCount    uint8  // Number of bytes to read, 1-based [max 255]
	Offset       uint8  // Data Offset
	DataWrite    []byte // if empty => read I2C, if not empty => write I2C

	// Response Data
	Data []byte // if read from I2C and non-zero read count => read data
}

func (c *GetSetOEMOpenBmcWistronI2CCommand) Name() string { return "Read Write on Wistron OpenBMC" }
func (c *GetSetOEMOpenBmcWistronI2CCommand) Code() uint8  { return 0x25 }

func (c *GetSetOEMOpenBmcWistronI2CCommand) Input() []byte  { return c.DataWrite }
func (c *GetSetOEMOpenBmcWistronI2CCommand) Output() []byte { return c.Data }

func (c *GetSetOEMOpenBmcWistronI2CCommand) NetFnRsLUN() NetFnRsLUN {
	return NewNetFnRsLUN(NetFnOemOne, 0)
}

func (c *GetSetOEMOpenBmcWistronI2CCommand) String() string { return cmdToJSON(c) }
func (c *GetSetOEMOpenBmcWistronI2CCommand) Marshal() ([]byte, error) {
	cmd := []byte{c.BusID, c.SlaveAddress, c.ReadCount, c.Offset}
	if len(c.DataWrite) > 0 {
		cmd = append(cmd, c.DataWrite...)
	}
	return cmd, nil
}

func (c *GetSetOEMOpenBmcWistronI2CCommand) Unmarshal(buf []byte) ([]byte, error) {
	c.Data = buf

	return nil, nil
}

// *** OpenBMC Wistron FAN Control

type ControlMode uint8

const (
	ControlModeAuto   ControlMode = 0x00
	ControlModeManual ControlMode = 0x01
)

func (c ControlMode) String() string {
	switch c {
	case ControlModeAuto:
		return "Auto"
	case ControlModeManual:
		return "Manual"
	default:
		return fmt.Sprintf("Unknown mode [%00x]", uint8(c))
	}
}

// ** Get

// GetOEMOpenBmcWistronFanControlCommand Wistron OEM OpenBMC Get Fan Speed Control Command
type GetOEMOpenBmcWistronFanControlCommand struct {
	// Response Data
	ControlMode ControlMode // Auto == 0, Manual == 1
	FanSpeed    uint8       // 0 - 100, in percents
}

func (c *GetOEMOpenBmcWistronFanControlCommand) Name() string {
	return "Get Fan Speed Control Command on Wistron OpenBMC"
}
func (c *GetOEMOpenBmcWistronFanControlCommand) Code() uint8 { return 0x22 }

func (c *GetOEMOpenBmcWistronFanControlCommand) NetFnRsLUN() NetFnRsLUN {
	return NewNetFnRsLUN(NetFnOemOne, 0)
}

func (c *GetOEMOpenBmcWistronFanControlCommand) String() string           { return cmdToJSON(c) }
func (c *GetOEMOpenBmcWistronFanControlCommand) Marshal() ([]byte, error) { return []byte{}, nil }

func (c *GetOEMOpenBmcWistronFanControlCommand) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 2); err != nil {
		return nil, err
	}
	c.ControlMode = ControlMode(buf[0])
	c.FanSpeed = buf[1]

	return nil, nil
}

// * Set

// SetOEMOpenBmcWistronFanControlCommand Wistron OEM OpenBMC Set Fan Speed Control Command
type SetOEMOpenBmcWistronFanControlCommand struct {
	// Response Data
	ControlMode ControlMode // Auto == 0, Manual == 1
	FanSpeed    uint8       // 0 - 100, in percents, has meaning only in ControlModeManual
}

func (c *SetOEMOpenBmcWistronFanControlCommand) Name() string {
	return "Get Fan Speed Control Command on Wistron OpenBMC"
}
func (c *SetOEMOpenBmcWistronFanControlCommand) Code() uint8 { return 0x21 }

func (c *SetOEMOpenBmcWistronFanControlCommand) NetFnRsLUN() NetFnRsLUN {
	return NewNetFnRsLUN(NetFnOemOne, 0)
}

func (c *SetOEMOpenBmcWistronFanControlCommand) String() string { return cmdToJSON(c) }
func (c *SetOEMOpenBmcWistronFanControlCommand) Marshal() ([]byte, error) {
	if c.ControlMode == ControlModeManual {
		if c.FanSpeed > 100 {
			return nil, fmt.Errorf("fan speed %d%% is out of range [0-100]", c.FanSpeed)
		}
		return []byte{uint8(c.ControlMode), c.FanSpeed}, nil
	} else {
		return []byte{uint8(c.ControlMode)}, nil
	}
}

func (c *SetOEMOpenBmcWistronFanControlCommand) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 0); err != nil {
		return nil, err
	}

	return nil, nil
}
