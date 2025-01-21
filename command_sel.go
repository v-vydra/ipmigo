package ipmigo

import (
	"encoding/binary"
	"fmt"
)

// GetSELInfoCommand Get SEL Info (Section 31.2)
type GetSELInfoCommand struct {
	// Response Data
	SELVersion        uint8
	Entries           uint16
	FreeSpace         uint16
	LastAddTime       uint32
	LastDelTime       uint32
	SupportAllocInfo  bool
	SupportReserve    bool
	SupportPartialAdd bool
	SupportDelete     bool
	Overflow          bool
}

func (c *GetSELInfoCommand) Name() string { return "Get SEL Info" }
func (c *GetSELInfoCommand) Code() uint8  { return 0x40 }

func (c *GetSELInfoCommand) NetFnRsLUN() NetFnRsLUN {
	return NewNetFnRsLUN(NetFnStorageReq, 0)
}

func (c *GetSELInfoCommand) String() string           { return cmdToJSON(c) }
func (c *GetSELInfoCommand) Marshal() ([]byte, error) { return []byte{}, nil }

func (c *GetSELInfoCommand) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 14); err != nil {
		return nil, err
	}

	c.SELVersion = buf[0]
	c.Entries = binary.LittleEndian.Uint16(buf[1:3])
	c.FreeSpace = binary.LittleEndian.Uint16(buf[3:5])
	c.LastAddTime = binary.LittleEndian.Uint32(buf[5:9])
	c.LastDelTime = binary.LittleEndian.Uint32(buf[9:13])
	c.SupportAllocInfo = buf[13]&0x01 != 0
	c.SupportReserve = buf[13]&0x02 != 0
	c.SupportPartialAdd = buf[13]&0x04 != 0
	c.SupportDelete = buf[13]&0x08 != 0
	c.Overflow = buf[13]&0x80 != 0

	return buf[14:], nil
}

// ReserveSELCommand Reserve SEL Command (Section 31.4)
type ReserveSELCommand struct {
	// Response Data
	ReservationID uint16
}

func (c *ReserveSELCommand) Name() string { return "Reserve SEL" }
func (c *ReserveSELCommand) Code() uint8  { return 0x42 }

func (c *ReserveSELCommand) NetFnRsLUN() NetFnRsLUN {
	return NewNetFnRsLUN(NetFnStorageReq, 0)
}

func (c *ReserveSELCommand) String() string           { return cmdToJSON(c) }
func (c *ReserveSELCommand) Marshal() ([]byte, error) { return []byte{}, nil }

func (c *ReserveSELCommand) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 2); err != nil {
		return nil, err
	}
	c.ReservationID = binary.LittleEndian.Uint16(buf)
	return buf[2:], nil
}

// GetSELEntryCommand Get SEL Entry Command (Section 31.5)
type GetSELEntryCommand struct {
	// Request Data
	ReservationID uint16
	RecordID      uint16
	RecordOffset  uint8
	ReadBytes     uint8

	// Response Data
	NextRecordID uint16
	RecordData   []byte
}

func (c *GetSELEntryCommand) Name() string           { return "Get SDR" }
func (c *GetSELEntryCommand) Code() uint8            { return 0x43 }
func (c *GetSELEntryCommand) NetFnRsLUN() NetFnRsLUN { return NewNetFnRsLUN(NetFnStorageReq, 0) }
func (c *GetSELEntryCommand) String() string         { return cmdToJSON(c) }

func (c *GetSELEntryCommand) Marshal() ([]byte, error) {
	return []byte{byte(c.ReservationID), byte(c.ReservationID >> 8), byte(c.RecordID), byte(c.RecordID >> 8),
		c.RecordOffset, c.ReadBytes}, nil
}

func (c *GetSELEntryCommand) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 2); err != nil {
		return nil, err
	}

	c.NextRecordID = binary.LittleEndian.Uint16(buf)
	buf = buf[2:]
	if l := len(buf); l <= int(c.ReadBytes) {
		c.RecordData = make([]byte, l)
		copy(c.RecordData, buf)
		return nil, nil
	} else {
		c.RecordData = make([]byte, c.ReadBytes)
		copy(c.RecordData, buf)
		return buf[c.ReadBytes:], nil
	}
}

// ** Clear SEL

type ClearSELAction uint8

const (
	ClearSELActionInitialErase     ClearSELAction = 0xAA
	ClearSELActionGetErasureStatus ClearSELAction = 0x00
	EraseAll
)

type ClearSELErasureProgress uint8

func (c ClearSELErasureProgress) String() string {
	switch c {
	case 0x00:

		return "erasure in progress"
	case 0x01:
		return "erase completed"
	default:
		return fmt.Sprintf("unknown progress [0x%00x]", uint8(c))
	}
}
func (c ClearSELErasureProgress) IsErasureInProgress() bool {
	return c == 0x00
}
func (c ClearSELErasureProgress) IsEraseCompleted() bool {
	return c == 0x01
}

// ClearSELCommand Clear SEL Command (Section 31.9)
type ClearSELCommand struct {
	// Request Data
	ReservationID uint16
	Action        ClearSELAction

	// Response Data
	ErasureProgress ClearSELErasureProgress // bits 3-0, 0h == erasure in progress, 1h == erase completed
}

func (c *ClearSELCommand) Name() string { return "Clear SEL Command" }
func (c *ClearSELCommand) Code() uint8  { return 0x47 }

func (c *ClearSELCommand) NetFnRsLUN() NetFnRsLUN { return NewNetFnRsLUN(NetFnStorageReq, 0) }

func (c *ClearSELCommand) String() string { return cmdToJSON(c) }
func (c *ClearSELCommand) Marshal() ([]byte, error) {
	return []byte{
		byte(c.ReservationID), byte(c.ReservationID >> 8), // LSB byte first
		0x43, 0x4c, 0x52, // C L R
		uint8(c.Action),
	}, nil
}

func (c *ClearSELCommand) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 1); err != nil {
		return nil, err
	}
	c.ErasureProgress = ClearSELErasureProgress(buf[0] & 0x0F) // bits 3:0
	return nil, nil
}
