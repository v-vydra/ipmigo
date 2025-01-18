package ipmigo

import (
	"encoding/binary"
)

// GetFRUInventoryAreaInfoCommand Get FRU Inventory Area Info Command (Section 34.1)
type GetFRUInventoryAreaInfoCommand struct {
	// Request Data
	DeviceID uint8
	Lun      uint8 // most of the cases 0

	// Response Data
	FruSize       uint16 // FRU Inventory area size in bytes
	AccessByWords bool   // 0b = Device is accessed by bytes, 1b = Device is accessed by words
}

func (c *GetFRUInventoryAreaInfoCommand) Name() string { return "Get FRU Inventory Area Info" }
func (c *GetFRUInventoryAreaInfoCommand) Code() uint8  { return 0x10 }

func (c *GetFRUInventoryAreaInfoCommand) NetFnRsLUN() NetFnRsLUN {
	return NewNetFnRsLUN(NetFnStorageReq, c.Lun)
}

func (c *GetFRUInventoryAreaInfoCommand) String() string           { return cmdToJSON(c) }
func (c *GetFRUInventoryAreaInfoCommand) Marshal() ([]byte, error) { return []byte{c.DeviceID}, nil }

func (c *GetFRUInventoryAreaInfoCommand) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 3); err != nil {
		return nil, err
	}

	//c.FruSize = binary.LittleEndian.Uint16(buf[0:1])
	c.FruSize = binary.LittleEndian.Uint16(buf)
	c.AccessByWords = buf[2]&0x01 == 0x01

	return nil, nil
}

// GetFRUDataCommand Read FRU Data Command (Section 34.2)
type GetFRUDataCommand struct {
	// Request Data
	DeviceID     uint8  // FRU Device ID. FFh = reserved.
	Lun          uint8  // should be 0
	Offset       uint16 // Offset is in bytes or words per device access type returned in the Get FRU Inventory Area Info command.
	CountRequest uint8  // Count to read (16?)

	// Response Data
	CountResponse uint8  // Count returned
	Data          []byte //  Requested data
}

func (c *GetFRUDataCommand) Name() string { return "Read FRU Data Command" }
func (c *GetFRUDataCommand) Code() uint8  { return 0x11 }

func (c *GetFRUDataCommand) NetFnRsLUN() NetFnRsLUN {
	return NewNetFnRsLUN(NetFnStorageReq, c.Lun)
}

func (c *GetFRUDataCommand) String() string { return cmdToJSON(c) }

func (c *GetFRUDataCommand) Marshal() ([]byte, error) {
	return []byte{c.DeviceID, byte(c.Offset), byte(c.Offset >> 8), c.CountRequest}, nil
}

func (c *GetFRUDataCommand) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 1); err != nil {
		return nil, err
	}

	c.CountResponse = buf[1]

	buf = buf[1:]
	if l := len(buf); l <= int(c.CountRequest) {
		c.Data = make([]byte, l)
		copy(c.Data, buf)
		return nil, nil
	} else {
		c.Data = make([]byte, c.CountResponse)
		copy(c.Data, buf)
		return buf[c.CountRequest:], nil
	}
}
