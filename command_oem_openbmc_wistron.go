package ipmigo

// GetSetOEMOpenBmcWistronI2CCommand I2C Read Write on Wistron's OpenBMC
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
