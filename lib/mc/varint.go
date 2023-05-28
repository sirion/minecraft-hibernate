package mc

import (
	"errors"
	"io"
)

const MASK_DATA int = 0b01111111
const MASK_MARK int = 0b10000000
const MASK_DATA_U uint32 = 0b01111111
const MASK_MARK_U uint32 = 0b10000000

/// ReadVarInt reads from the given reader byte by byte and interprets the data as VarInt.
/// Returns the value, the bytes read and an error in case the read data cannot be interpreted as VarInt
func ReadVarInt(cr io.ByteReader) (int, int, error) {
	value := 0
	position := 0

	readBytes := 0
	for {
		currentByte, err := cr.ReadByte()
		readBytes += 1
		if err != nil {
			return 0, 0, err
		}

		value = value | ((int(currentByte) & MASK_DATA) << position)

		if (int(currentByte) & MASK_MARK) == 0 {
			break
		}

		position += 7

		if position >= 32 {
			return 0, 0, errors.New("value too hight for VarInt")
		}
	}

	return value, readBytes, nil
}

func ToVarInt(signedValue int) []byte {
	value := uint32(signedValue)

	varInt := make([]byte, 0, 5)
	for {
		if (value & ^MASK_DATA_U) == 0 {
			varInt = append(varInt, byte(value))
			break
		}

		varInt = append(varInt, byte((value&MASK_DATA_U)|MASK_MARK_U))

		// Note: >>> means that the sign bit is shifted with the rest of the number rather than being left alone
		value = value >> 7
	}

	return varInt
}
