package rdbparser

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// opcodes
const (
	AUX          = 0xFA
	RESIZEDB     = 0xFB
	EXPIRETIMEMS = 0xFC
	EXPIRETIME   = 0xFD
	SELECTDB     = 0xFE
	EOF          = 0xFF
)

func handleLengthEncoding(data []byte, cursor int) (int, interface{}, int, error) {
	n := len(data)
	if cursor >= n {
		return 0, nil, cursor, errors.New("cursor out of bounds")
	}

	byteValue := data[cursor]
	lengthType := (byteValue & 0b11000000) >> 6 // First two bits represent length type

	if lengthType == 0 { 
		length := int(byteValue & 0b00111111)	// this byte, 6 bits length 
		cursor++
		return length, string(data[cursor : min(n, cursor + length)]), cursor + length, nil
	}

	if lengthType == 1 {
		if cursor + 1 >= len(data) {
			return 0, nil, cursor, errors.New("incomplete 14-bit length")
		} // (this 6 + next 8) 14 bits length BE
		length := int(byteValue & 0b00111111)<<8 | int(data[cursor+1]) 
		cursor += 2
		return length, string(data[cursor : min(n, cursor + length)]), cursor + length, nil
	}

	if lengthType == 2 {
		if cursor+4 >= len(data) {
			return 0, nil, cursor, errors.New("incomplete 32-bit length")
		} // Next 32 bits length BE
		length := int(binary.BigEndian.Uint32(data[cursor+1 : cursor+5]))
		cursor += 5
		return length, string(data[cursor : min(n, cursor + length)]), cursor + length, nil
	}

	if lengthType == 3 { // String encoding types
		stringType := byteValue & 0b00111111	//last 6 bits represent string type
		switch stringType {
			case 0: // 8-bit length
				length := int(data[cursor+1])
				return length, length, cursor + 2, nil

			case 1: // 16-bit length (Little-Endian)
				if cursor+2 >= len(data) {
					return 0, nil, cursor, errors.New("incomplete 16-bit string length")
				}
				length := int(binary.LittleEndian.Uint16(data[cursor+1 : cursor+3]))
				return length, length, cursor + 3, nil

			case 2: // 32-bit length (Little-Endian)
				if cursor+4 >= len(data) {
					return 0, nil, cursor, errors.New("incomplete 32-bit string length")
				}
				length := int(binary.LittleEndian.Uint32(data[cursor+1 : cursor+5]))
				return length, length, cursor + 5, nil

			default:
				return 0, nil, cursor, fmt.Errorf("invalid string type %d", stringType)
		}
	}

	return 0, nil, cursor, fmt.Errorf("invalid length encoding %d at %d", lengthType, cursor)
}


func ParseRDB(data []byte) []map[string]interface{} {
	var cursor int = 9 // Skip the Redis Magic String and Version
	fmt.Println("Header: ",string(data[:cursor]))

	var keyValuePairs []map[string]interface{}
	var expType string
	var expTime uint64
	for cursor < len(data) {
		if AUX == data[cursor] {
			cursor++
			// Key
			_,key, newCursor, err := handleLengthEncoding(data, cursor)
			if err != nil {
				fmt.Println("Error decoding AUX key:", err)
				return keyValuePairs
			}
			cursor = newCursor
			
			// Value
			_,value, newCursor, err := handleLengthEncoding(data, cursor)
			if err != nil {
				fmt.Println("Error decoding AUX value:", err)
				return keyValuePairs
			}
			cursor = newCursor
			fmt.Printf("MetaData: %v : %v\n", key, value)
		} else

		if SELECTDB == data[cursor] {
			cursor++
			dbIndex, _,_, err := handleLengthEncoding(data, cursor)
			if err != nil {
				fmt.Println("Error decoding database index:", err)
				return keyValuePairs
			}
			cursor++
			fmt.Println("Database Index:", dbIndex)
		} else

		if RESIZEDB == data[cursor] {
			cursor++
			// hash table size
			hashSize, _,_, err := handleLengthEncoding(data, cursor)
			if err != nil {
				fmt.Println("Error decoding hash table size: ", err)
				return keyValuePairs
			}
			cursor++
			
			// expiry table size
			expirySize, _,_, err := handleLengthEncoding(data, cursor)
			if err != nil {
				fmt.Println("Error decoding expiry table size:", err)
				return keyValuePairs
			}
			cursor++
			fmt.Printf("Hash Table Size: %v \nExpiry Table Size: %v\n", hashSize, expirySize)
		} else
	
		if EXPIRETIME == data[cursor] {
			cursor++
			expType = "PX"
			expTime = uint64(binary.LittleEndian.Uint32(data[cursor : cursor+4]))
			cursor += 4
		} else

		if EXPIRETIMEMS == data[cursor] {
			cursor++
			expType = "EX"
			expTime = uint64(binary.LittleEndian.Uint64(data[cursor : cursor+8]))
			cursor += 8
		} else

		if EOF == data[cursor] {
			return keyValuePairs
		} else

		if data[cursor] == 0x00{
			// Key and Value
			cursor++
			_,key, newCursor, err := handleLengthEncoding(data, cursor)
			if err != nil {
				fmt.Println("Error decoding key:", err)
				return keyValuePairs
			}
			cursor = newCursor
			
			_,value, newCursor, err := handleLengthEncoding(data, cursor)
			if err != nil {
				fmt.Println("Error decoding value:", err)
				return keyValuePairs
			}
			cursor = newCursor
			
			keyValuePairs = append(keyValuePairs, map[string]interface{}{"key": key, "value": value, "type": expType, "time": expTime})
			expType = ""
			expTime = 0
		} else

		{
			fmt.Println("Inc: ",cursor)
			cursor++
		}
		
	}

	return keyValuePairs
}