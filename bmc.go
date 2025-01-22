package ipmigo

import "fmt"

// ** Self Test
type SelfTestStatus uint8

const (
	SelfTestStatusAllPassed          SelfTestStatus = 0x55
	SelfTestStatusNotImplemented     SelfTestStatus = 0x56
	SelfTestStatusInaccessible       SelfTestStatus = 0x57
	SelfTestStatusFatalHardwareError SelfTestStatus = 0x58
)

func (c SelfTestStatus) String() string {
	switch c {
	case SelfTestStatusAllPassed:
		return "No error, All Self Tests Passed"
	case SelfTestStatusNotImplemented:
		return "Self Test function not implemented in this controller"
	case SelfTestStatusInaccessible:
		return "Corrupted or inaccessible data or devices"
	case SelfTestStatusFatalHardwareError:
		return "Fatal hardware error (BMC inoperative)"
	case 0xFF:
		return "Reserved"
	default:
		return fmt.Sprintf("Unknown internal failure, code 0x%00x", uint8(c))
	}
}

// GetSelfTestResultsCommand Get Self Test Results Command (section 20.4)
type GetSelfTestResultsCommand struct {
	// Request Data

	// Response Data
	Status       uint8 // byte 1
	TestsResults uint8 // byte 2
}

func (c *GetSelfTestResultsCommand) Name() string { return "Get Self Test Results" }
func (c *GetSelfTestResultsCommand) Code() uint8  { return 0x04 }

func (c *GetSelfTestResultsCommand) NetFnRsLUN() NetFnRsLUN {
	return NewNetFnRsLUN(NetFnAppReq, 0)
}

func (c *GetSelfTestResultsCommand) String() string           { return cmdToJSON(c) }
func (c *GetSelfTestResultsCommand) Marshal() ([]byte, error) { return []byte{}, nil }

func (c *GetSelfTestResultsCommand) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 2); err != nil {
		return nil, err
	}
	c.Status = buf[0]
	c.TestsResults = buf[1]

	return nil, nil
}

func (c *GetSelfTestResultsCommand) GetStatus() uint8      { return c.Status }
func (c *GetSelfTestResultsCommand) GetTestResults() uint8 { return c.TestsResults }
func (c *GetSelfTestResultsCommand) GetTestResultsAsString() string {
	result := ""
	switch SelfTestStatus(c.Status) {
	case SelfTestStatusAllPassed, SelfTestStatusNotImplemented, 0xFF:
		if c.TestsResults != 0 {
			result = fmt.Sprintf("Unexpected TestResult Data 0x%00x", c.TestsResults)
		}
	case SelfTestStatusFatalHardwareError:
		result = fmt.Sprintf("Device specific information 0x%00x", c.TestsResults)
	case SelfTestStatusInaccessible:
		result = c.GetTestsResultsBitfieldAsString()
	}
	return result
}

// GetTestsResultsBitfieldAsString when Status == 0x57h interpret as bitfield
func (c *GetSelfTestResultsCommand) GetTestsResultsBitfieldAsString() string {
	result := ""
	for i := 0; i < 8; i++ {
		if c.TestsResults&(1<<i) != 0 {
			if len(result) > 0 {
				result += ", "
			}
			switch i {
			case 7:
				result += "Cannot access SEL device"
			case 6:
				result += "Cannot access SDR Repository"
			case 5:
				result += "Cannot access BMC FRU device"
			case 4:
				result += "IPMB signal lines do not respond"
			case 3:
				result += "SDR Repository empty"
			case 2:
				result += "Internal Use Area of BMC FRU corrupted"
			case 1:
				result += "Controller update ‘boot block’ firmware corrupted"
			case 0:
				result += "Controller operational firmware corrupted"
			}
		}
	}
	return result
}

// ** Cold Reset

// SetColdResetCommand Cold Reset Command (section 20.2)
type SetColdResetCommand struct {
	// Request Data

	// Response Data

}

func (c *SetColdResetCommand) Name() string { return "Set Cold Reset" }
func (c *SetColdResetCommand) Code() uint8  { return 0x02 }

func (c *SetColdResetCommand) NetFnRsLUN() NetFnRsLUN {
	return NewNetFnRsLUN(NetFnAppReq, 0)
}

func (c *SetColdResetCommand) String() string           { return cmdToJSON(c) }
func (c *SetColdResetCommand) Marshal() ([]byte, error) { return []byte{}, nil }

func (c *SetColdResetCommand) Unmarshal(buf []byte) ([]byte, error) {
	if err := cmdValidateLength(c, buf, 0); err != nil {
		return nil, err
	}

	return nil, nil
}
