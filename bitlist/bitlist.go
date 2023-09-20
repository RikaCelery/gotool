package bitlist

import "fmt"

type BitList struct {
	data   []byte
	Length int
}

func NewBitList() *BitList {
	return &BitList{
		data:   []byte{},
		Length: 0,
	}
}

func (bl *BitList) appendBit(bit byte) {
	if bl.Length%8 == 0 {
		bl.data = append(bl.data, 0)
	}
	if bit == 1 {
		bl.data[bl.Length/8] |= 1 << (7 - uint(bl.Length%8))
	}
	bl.Length++
}

func (bl *BitList) Append(bits ...byte) {
	for _, bit := range bits {
		bl.appendBit(bit)
	}
}

func (bl *BitList) Concatenate(other *BitList) *BitList {
	newBitList := NewBitList()

	// Append bits from the first BitList
	for i := 0; i < bl.Length; i++ {
		bit := (bl.data[i/8] >> (7 - uint(i%8))) & 1
		newBitList.appendBit(byte(bit))
	}

	// Append bits from the second BitList
	for i := 0; i < other.Length; i++ {
		bit := (other.data[i/8] >> (7 - uint(i%8))) & 1
		newBitList.appendBit(byte(bit))
	}

	return newBitList
}

func (bl *BitList) Join(other *BitList) {
	// Append bits from the other BitList
	for i := 0; i < other.Length; i++ {
		bit := (other.data[i/8] >> (7 - uint(i%8))) & 1
		bl.appendBit(byte(bit))
	}
}

func (bl *BitList) GetBit(index int) byte {
	if index >= 0 && index < bl.Length {
		byteIndex := index / 8
		bitIndex := 7 - uint(index%8)
		return (bl.data[byteIndex] >> bitIndex) & 1
	}
	return 0 // Return 0 for out-of-bounds index
}
func (bl *BitList) SubList(startIndex, endIndex int) *BitList {
	endIndex--
	if startIndex < 0 {
		startIndex = 0
	}
	if endIndex >= bl.Length {
		endIndex = bl.Length - 1
	}

	if endIndex < startIndex {
		return nil
	}

	subBitList := NewBitList()

	for i := startIndex; i <= endIndex; i++ {
		bit := bl.GetBit(i)
		subBitList.appendBit(bit)
	}

	return subBitList
}
func (bl *BitList) String() string {
	result := ""
	for i := 0; i < bl.Length; i++ {
		byteIndex := i / 8
		bitIndex := 7 - uint(i%8)
		bit := (bl.data[byteIndex] >> bitIndex) & 1
		result += fmt.Sprintf("%d", bit)
	}
	return result
}

func BytesToBitList(byteSlice []byte) *BitList {
	bitList := NewBitList()

	for _, b := range byteSlice {
		// Append each bit of the byte to the BitList
		for i := 7; i >= 0; i-- {
			bit := (b >> uint(i)) & 1
			bitList.appendBit(byte(bit))
		}
	}

	return bitList
}
func BitStringToBitList(bitString string) *BitList {
	bitList := NewBitList()

	for _, char := range bitString {
		if char == '0' {
			bitList.appendBit(0)
		} else if char == '1' {
			bitList.appendBit(1)
		} else {
			// Skip invalid characters
			continue
		}
	}

	return bitList
}
func (bl *BitList) ToBytes() []byte {
	data := bl.data
	if bl.Length%8 != 0 {
		b := bl.Length % 8
		data[len(data)-1] >>= (8 - b)
	}
	return data
}
