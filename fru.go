package ipmigo

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

const (
	fruCommonHeaderSize = 8

	fruDefaultReadBytes = 16 // nowadays 63 should be ok
)

// FRUCommonHeader Table 8-1, COMMON HEADER
type FRUCommonHeader struct {
	HeaderFormatVersion        uint8 // Byte 0 bits 3:0 : should be 0b001
	InternalUseAreaStartOffset uint8 // Byte 1: Internal Use Area offset in multiples of 8 bytes [0x00h == not used]
	ChassisInfoAreaStartOffset uint8 // Byte 2: Chassis Info Area offset in multiples of 8 bytes
	BoardInfoAreaStartOffset   uint8 // Byte 3: Board Info Area offset in multiples of 8 bytes
	ProductInfoAreaStartOffset uint8 // Byte 4: Product Info Area offset in multiples of 8 bytes
	MultiRecordAreaStartOffset uint8 // Byte 5: MultiRecord Area offset in multiples of 8 bytes
	PadByte                    uint8 // Byte 6: PAD
	CheckSum                   uint8 // Byte 7: Checksum of bytes [0..6]
}

func (f *FRUCommonHeader) Unmarshal(buf []byte) ([]byte, error) {
	if l := len(buf); l < 1 {
		return nil, fmt.Errorf("no data in FRU Common Header")
	}
	f.HeaderFormatVersion = buf[0] & 0x0F // bits 3:0

	if l := len(buf); l < fruCommonHeaderSize {
		return nil, fmt.Errorf("invalid FRU Common Header size : %d/%d, header version: %d [Raw Hex Data: %s]", l, fruCommonHeaderSize, f.HeaderFormatVersion, hex.EncodeToString(buf))
	}

	f.InternalUseAreaStartOffset = buf[1]
	f.ChassisInfoAreaStartOffset = buf[2]
	f.BoardInfoAreaStartOffset = buf[3]
	f.ProductInfoAreaStartOffset = buf[4]
	f.MultiRecordAreaStartOffset = buf[5]
	f.PadByte = buf[6]
	f.CheckSum = buf[7]

	if f.HeaderFormatVersion != 0x01 {
		return nil, fmt.Errorf("invalid FRU Common Header Format Version : %000b", f.HeaderFormatVersion)
	}
	return nil, nil
}

func (f *FRUCommonHeader) ToString() string {
	return fmt.Sprintf(
		"  Internal Use Area offset : %d\n"+
			"  Chassis Info Area offset : %d\n"+
			"  Board Info Area offset   : %d\n"+
			"  Product Info Area offset : %d\n"+
			"  MultiRecord Area offset  : %d\n"+
			"  Pad                      : %d\n",
		f.InternalUseAreaStartOffset,
		f.ChassisInfoAreaStartOffset,
		f.BoardInfoAreaStartOffset,
		f.ProductInfoAreaStartOffset,
		f.MultiRecordAreaStartOffset,
		f.PadByte,
	)
}

type FRUDeviceData struct {
	DeviceID uint8
	Lun      uint8

	DataSize uint16 // Raw FRU Data size (at least 8)
	Data     []byte // Raw FRU Data bytes - full area

	CommonHeader *FRUCommonHeader
	BoardInfo    *FRUBoardInfoArea
	ProductInfo  *FRUProductInfoArea
}

func (f *FRUDeviceData) ToStringDebug() string {
	var res = f.CommonHeader.ToString()
	if f.BoardInfo != nil {
		res += "\n" + f.BoardInfo.ToString()
	}
	if f.ProductInfo != nil {
		res += "\n" + f.ProductInfo.ToString()
	}
	return res
}
func (f *FRUDeviceData) ToString() string {
	var res string
	if f.BoardInfo != nil {
		res += f.BoardInfo.ToString()
	}
	if f.ProductInfo != nil {
		res += "\n" + f.ProductInfo.ToString()
	}
	return res
}
func (f *FRUDeviceData) String() string {
	res := "Board Info Area:\n" + f.GetBoardInfoAreaFieldsAsString()
	res += "Product Info Area:\n" + f.GetProductInfoAreaFieldsAsString()
	return res
}
func (f *FRUDeviceData) DebugString() string {
	return fmt.Sprintf("DeviceID: %d, Lun: %d, DataSize: %d, Data: %s", f.DeviceID, f.Lun, f.DataSize, hex.EncodeToString(f.Data))
}

// GetInternalUseArea return Internal Use Area byte buffer
//
//	ToDo: test confirmation needed
func (f *FRUDeviceData) GetInternalUseArea() []byte {
	if f.CommonHeader.InternalUseAreaStartOffset > 0 {
		var stopOffset uint8 = 0
		if f.CommonHeader.ChassisInfoAreaStartOffset > 0 {
			stopOffset = f.CommonHeader.ChassisInfoAreaStartOffset
		} else if f.CommonHeader.BoardInfoAreaStartOffset > 0 {
			stopOffset = f.CommonHeader.BoardInfoAreaStartOffset
		} else if f.CommonHeader.ProductInfoAreaStartOffset > 0 {
			stopOffset = f.CommonHeader.ProductInfoAreaStartOffset
		} else if f.CommonHeader.MultiRecordAreaStartOffset > 0 {
			stopOffset = f.CommonHeader.MultiRecordAreaStartOffset
		} else {
			return []byte{}
		}
		return f.Data[f.CommonHeader.InternalUseAreaStartOffset*8 : stopOffset*8]
	}
	return []byte{}
}
func (f *FRUDeviceData) ParseBoardInfoArea() error {
	if f.CommonHeader.BoardInfoAreaStartOffset > 0 {
		area := &FRUBoardInfoArea{}
		_, err := area.Unmarshal(f.Data[f.CommonHeader.BoardInfoAreaStartOffset*8:])
		if err != nil {
			return err
		} else {
			f.BoardInfo = area
			return nil
		}
	}
	f.BoardInfo = nil
	return nil
}
func (f *FRUDeviceData) ParseProductInfoArea() error {
	if f.CommonHeader.ProductInfoAreaStartOffset > 0 {
		area := &FRUProductInfoArea{}
		_, err := area.Unmarshal(f.Data[f.CommonHeader.ProductInfoAreaStartOffset*8:])
		if err != nil {
			return err
		} else {
			f.ProductInfo = area
			return nil
		}
	}
	f.ProductInfo = nil
	return nil
}
func (f *FRUDeviceData) GetBoardInfoAreaAsString() string {
	if f.BoardInfo != nil {
		return fmt.Sprintf("Manufacture Date   : %s\n%s\n",
			f.BoardInfo.ManufactureDateTime.Format(time.DateTime),
			f.BoardInfo.String())
	} else {
		return f.GetBoardInfoAreaFieldsAsString()
	}
}
func (f *FRUDeviceData) GetBoardInfoAreaFieldsAsString() string {
	if f.BoardInfo != nil {
		return f.BoardInfo.String()
	} else {
		return "  no Board Info Area not found\n"
	}
}
func (f *FRUDeviceData) GetProductInfoAreaFieldsAsString() string {
	if f.ProductInfo != nil {
		return f.ProductInfo.String()
	} else {
		return "  no Product Info Area not found\n"
	}
}

type FRUBoardInfoArea struct {
	AreaFormatVersion   uint8 // must be 0x01
	AreaLength          uint8 // in multiples of 8
	LanguageCode        uint8
	ManufactureDateTime time.Time

	Fields []FRUAreaType

	CheckSum uint8
}

type FRUAreaType struct {
	Type   uint8
	Length uint8
	Value  []byte
}

// GetValue decodes the textual field from the raw bytes, depending on type.
//
//	currently only language code english supported (0 and 20)
func (t FRUAreaType) GetValue(languageCode uint8) interface{} {
	switch t.Type {
	case 0b00: // binary or unspecified
		return t.Value
	case 0b01: // BCD Plus - see section 13.1
		// ToDo: make BDC Plus conversion
		return string(t.Value)
	case 0b10: // 6-bit ASCII packed (ignores language code), section 13.2
		// ToDo: make 6-bit ASCII conversion
		return t.Value
	case 0b11: // 8bit ASCII, depends on language code, we suppose it's that English
		if languageCode == 0 || languageCode == 25 {
			return strings.TrimSpace(string(t.Value))
		} else {
			// ToDo: unicode conversion
			return t.Value
		}
	default:
		return t.Value
	}
}

func (f *FRUBoardInfoArea) GetFieldValueStringById(id uint8) string {
	if int(id) > len(f.Fields) {
		return ""
	}
	if f.Fields[id-1].Type == 0b11 {
		return f.Fields[id-1].GetValue(f.LanguageCode).(string)
	} else {
		return fmt.Sprintf("Type:%00b, Hex: %s, Raw: %+v", f.Fields[id-1].Type, hex.EncodeToString(f.Fields[id-1].Value), f.Fields[id-1])
	}
}
func (f *FRUBoardInfoArea) ToString() string {
	var res string
	if !f.ManufactureDateTime.IsZero() {
		res = fmt.Sprintf(
			"  Manufacture Date/Time : %s\n",
			f.ManufactureDateTime.Format(time.RFC1123Z),
		)
	} else {
		res = ""
	}

	if len(f.Fields) >= 4 {
		res += fmt.Sprintf(
			"  Board Manufacturer    : %s\n"+
				"  Board Product Name    : %s\n"+
				"  Board Serial Name     : %s\n"+
				"  Board Part Number     : %s\n",
			f.GetFieldValueStringById(1),
			f.GetFieldValueStringById(2),
			f.GetFieldValueStringById(3),
			f.GetFieldValueStringById(4),
		)
	}
	if len(f.Fields) > 4 {
		for i := 5; i <= len(f.Fields); i++ {
			res += fmt.Sprintf("    OEM Field #%d : %s\n", i-4, f.GetFieldValueStringById(uint8(i)))
		}
	}

	return res
}
func (f *FRUBoardInfoArea) String() string {
	var res string
	if len(f.Fields) < 1 {
		res = "Board Info Area Fields array is empty\n"
	} else {
		for i := 1; i <= len(f.Fields); i++ {
			res += fmt.Sprintf(" Board Field #%d : %s\n", i-1, f.GetFieldValueStringById(uint8(i)))
		}
	}

	return res
}

func (f *FRUBoardInfoArea) Unmarshal(buf []byte) ([]byte, error) {
	buffLen := len(buf)
	if buffLen < 2 {
		return nil, fmt.Errorf("invalid Board Info area size : %d/%d", buffLen, 2)
	}
	f.AreaFormatVersion = buf[0] & 0x0F // bits 3:0
	f.AreaLength = buf[1]

	areaSize := int(f.AreaLength * 8)

	if areaSize > buffLen {
		return nil, fmt.Errorf("invalid Board Info area size : need %d , FRU Data Size %d", areaSize, buffLen)
	}
	if f.AreaFormatVersion != 0x01 {
		return nil, fmt.Errorf("invalid Board Info area Format Version : %000b", f.AreaFormatVersion)
	}

	f.LanguageCode = buf[2]

	convertDate, err := ConvertBoardMfgDate(buf[3:6])
	if err != nil {
		return nil, fmt.Errorf("invalid Board Info area Manufacture Date : %s", err)
	}
	f.ManufactureDateTime = convertDate

	// The variable-length fields start at byte (boardOffset + 6)
	fieldsStart := 6
	if fieldsStart >= areaSize {
		return nil, fmt.Errorf("board area length %d too small for variable fields - starting at offset %d", areaSize, fieldsStart)
	}

	// We'll parse fields in order until we see a type-length byte = 0xC1 (end marker)
	f.Fields = []FRUAreaType{}

	idx := fieldsStart
	fieldIndex := 0

	for idx < areaSize {
		fieldIndex++
		fieldLenByte := buf[idx]
		if fieldLenByte == 0xC1 { // valid fo English - non Unicode
			// End of fields
			idx++
			break
		}

		length := int(fieldLenByte & 0x3F)      // bits 7:6
		fieldType := (fieldLenByte & 0xC0) >> 6 // bits 5:0 Might be ASCII, BCD, 6-bit ASCII, etc.

		//fmt.Printf(" Field #%d: Type/Len %00000000b Type %d [%00b], Length %d [%0000b]\n", fieldIndex, fieldLenByte, fieldType, fieldType, length, length)

		idx++
		if fieldType != 0b000000 && length != 0 { // not an empty field

			dataEnd := idx + length
			if dataEnd > areaSize {
				// We might be hitting the padding/checksum
				break
			}

			rawFieldData := buf[idx:dataEnd]
			idx = dataEnd

			f.Fields = append(f.Fields, FRUAreaType{
				Type:   fieldType,
				Length: uint8(length),
				Value:  rawFieldData,
			})
		}
	}
	f.CheckSum = buf[len(buf)-1]

	return nil, nil
}

type FRUProductInfoArea struct {
	AreaFormatVersion uint8 // must be 0x01
	AreaLength        uint8 // in multiples of 8
	LanguageCode      uint8

	Fields []FRUAreaType

	CheckSum uint8
}

func (f *FRUProductInfoArea) GetFieldValueStringById(id uint8) string {
	if int(id) > len(f.Fields) {
		return ""
	}
	if f.Fields[id-1].Type == 0b11 {
		return f.Fields[id-1].GetValue(f.LanguageCode).(string)
	} else {
		return fmt.Sprintf("Type:%00b, Hex: %s, Raw: %+v", f.Fields[id-1].Type, hex.EncodeToString(f.Fields[id-1].Value), f.Fields[id-1])
	}
}
func (f *FRUProductInfoArea) ToString() string {
	var res string

	if len(f.Fields) >= 5 {
		res += fmt.Sprintf(
			"  Product Manufacturer  : %s\n"+
				"  Product Name          : %s\n"+
				"  Product Part/Model    : %s\n"+
				"  Product Version       : %s\n"+
				"  Product Serial Number : %s\n",
			f.GetFieldValueStringById(1),
			f.GetFieldValueStringById(2),
			f.GetFieldValueStringById(3),
			f.GetFieldValueStringById(4),
			f.GetFieldValueStringById(5),
		)
	}
	if len(f.Fields) > 5 {
		for i := 6; i <= len(f.Fields); i++ {
			res += fmt.Sprintf("    OEM Field #%d : %s\n", i-6, f.GetFieldValueStringById(uint8(i)))
		}
	}

	return res
}
func (f *FRUProductInfoArea) String() string {
	var res string
	if len(f.Fields) < 1 {
		res = "Product Info Area Fields array is empty\n"
	} else {
		for i := 1; i <= len(f.Fields); i++ {
			res += fmt.Sprintf(" Product Field #%d : %s\n", i-1, f.GetFieldValueStringById(uint8(i)))
		}
	}

	return res
}

func (f *FRUProductInfoArea) Unmarshal(buf []byte) ([]byte, error) {
	buffLen := len(buf)
	if buffLen < 2 {
		return nil, fmt.Errorf("invalid Product Info area size : %d/%d", buffLen, 2)
	}
	f.AreaFormatVersion = buf[0] & 0x0F // bits 3:0
	f.AreaLength = buf[1]

	areaSize := int(f.AreaLength * 8)

	if areaSize > buffLen {
		return nil, fmt.Errorf("invalid Product Info area size : need %d , FRU Data Size %d", areaSize, buffLen)
	}
	if f.AreaFormatVersion != 0x01 {
		return nil, fmt.Errorf("invalid Product Info area Format Version : %000b", f.AreaFormatVersion)
	}

	f.LanguageCode = buf[2]

	// The variable-length fields start at byte (boardOffset + 3)
	fieldsStart := 3
	if fieldsStart >= areaSize {
		return nil, fmt.Errorf("product area length %d too small for variable fields - starting at offset %d", areaSize, fieldsStart)
	}

	// We'll parse fields in order until we see a type-length byte = 0xC1 (end marker)
	f.Fields = []FRUAreaType{}

	idx := fieldsStart
	fieldIndex := 0

	for idx < areaSize {
		fieldIndex++
		fieldLenByte := buf[idx]
		if fieldLenByte == 0xC1 { // valid fo English - non Unicode
			// End of fields
			idx++
			break
		}

		length := int(fieldLenByte & 0x3F)      // bits 7:6
		fieldType := (fieldLenByte & 0xC0) >> 6 // bits 5:0 Might be ASCII, BCD, 6-bit ASCII, etc.

		//fmt.Printf(" Product Field #%d: Type/Len %00000000b [%00x] Type %d [%00b], Length %d [%0000b]\n", fieldIndex, fieldLenByte, fieldLenByte, fieldType, fieldType, length, length)

		idx++
		if fieldType != 0b000000 && length != 0 { // not an empty field

			dataEnd := idx + length
			if dataEnd > areaSize {
				// We might be hitting the padding/checksum
				break
			}

			rawFieldData := buf[idx:dataEnd]
			idx = dataEnd

			f.Fields = append(f.Fields, FRUAreaType{
				Type:   fieldType,
				Length: uint8(length),
				Value:  rawFieldData,
			})
			//fmt.Printf("   - Product Field #%d: value:%+v, str:%s\n", fieldIndex, rawFieldData, string(rawFieldData))
		} else {
			//fmt.Printf("   - Product Field #%d: ignoring\n", fieldIndex)
		}
	}
	f.CheckSum = buf[len(buf)-1]

	return nil, nil
}

func FRUGetDeviceData(c *Client, deviceId uint8, lun uint8) (*FRUDeviceData, error) {
	// Obtain the FRU size
	gfa := &GetFRUInventoryAreaInfoCommand{
		DeviceID: deviceId,
		Lun:      lun,
	}
	if err := c.Execute(gfa); err != nil {
		return nil, err
	}

	// Read Common Header data - first 8 bytes in raw FRU data
	cmdGetCommonHeader := &GetFRUDataCommand{
		DeviceID:     deviceId,
		Lun:          lun,
		Offset:       0,
		CountRequest: 8,
	}
	if err := c.Execute(cmdGetCommonHeader); err != nil {
		return nil, err
	}
	cah := &FRUCommonHeader{}
	if _, err := cah.Unmarshal(cmdGetCommonHeader.Data); err != nil {
		return nil, err
	}
	cahCheckSum := checksum(cmdGetCommonHeader.Data[:7])
	if cahCheckSum != cah.CheckSum {
		return nil, fmt.Errorf("device %d has incorrect FRU Common header checksum : %d != %d ", deviceId, cahCheckSum, cah.CheckSum)
	}

	// Read all FRU areas
	fdData := &FRUDeviceData{
		DeviceID:     deviceId,
		Lun:          lun,
		CommonHeader: cah,
		DataSize:     gfa.FruSize,
	}
	fdData.Data = append(fdData.Data, cmdGetCommonHeader.Data...)

	// load full FRU Data
	orgSize := fdData.DataSize - 8 // already loaded header
	var size uint8
	var id uint16 = 0

	for orgSize > 0 {
		id++
		if orgSize > uint16(c.fruReadingBytes) {
			size = c.fruReadingBytes
			orgSize -= uint16(c.fruReadingBytes)
		} else {
			size = uint8(orgSize)
			orgSize = 0
		}

		// data
		cmd1 := &GetFRUDataCommand{
			DeviceID:     deviceId,
			Lun:          lun,
			Offset:       8 + (id-1)*uint16(c.fruReadingBytes),
			CountRequest: size,
		}
		if err := c.Execute(cmd1); err != nil {
			fmt.Println(err)
			return fdData, err
		}
		fdData.Data = append(fdData.Data, cmd1.Data...)
	}

	err := fdData.ParseBoardInfoArea()
	if err != nil {
		return fdData, err
	}

	err = fdData.ParseProductInfoArea()
	if err != nil {
		return fdData, err
	}

	return fdData, nil
}
