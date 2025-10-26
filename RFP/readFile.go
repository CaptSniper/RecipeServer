package rfp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// readChunk reads a single chunk from the buffer
func readChunk(buf *bytes.Reader) (chunkType string, payload []byte, err error) {
	chunkHeader := make([]byte, 8)
	if _, err = buf.Read(chunkHeader); err != nil {
		return
	}
	chunkType = string(chunkHeader[:4])
	chunkSize := binary.LittleEndian.Uint32(chunkHeader[4:8])
	payload = make([]byte, chunkSize)
	if _, err = buf.Read(payload); err != nil {
		return
	}
	// Skip padding to next 8-byte boundary
	padding := (8 - (chunkSize % 8)) % 8
	if padding > 0 {
		buf.Seek(int64(padding), 1)
	}
	return
}

// ReadRecipeFile reads an RFP3 file into a Recipe struct
func ReadRecipeFile(filename string) (*Recipe, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(data)

	// Check magic
	magic := make([]byte, 4)
	if _, err = buf.Read(magic); err != nil {
		return nil, err
	}
	if string(magic) != "RFP3" {
		return nil, fmt.Errorf("not an RFP3 file")
	}

	// Skip version, header size
	buf.Seek(4, 1)
	var chunkCount uint32
	binary.Read(buf, binary.LittleEndian, &chunkCount)
	buf.Seek(6, 1) // flags + reserved

	recipe := &Recipe{}
	for i := uint32(0); i < chunkCount; i++ {
		chunkType, payload, err := readChunk(buf)
		if err != nil {
			return nil, err
		}
		rdr := bytes.NewReader(payload)

		switch chunkType {
		case "CORE":
			binary.Read(rdr, binary.LittleEndian, &recipe.PrepTime)
			binary.Read(rdr, binary.LittleEndian, &recipe.CookTime)
			binary.Read(rdr, binary.LittleEndian, &recipe.AdditionalTime)
			binary.Read(rdr, binary.LittleEndian, &recipe.TotalTime)

			var servings uint16
			binary.Read(rdr, binary.LittleEndian, &servings)
			servingsBytes := make([]byte, servings)
			rdr.Read(servingsBytes)
			recipe.Servings = string(servingsBytes)

			var pathLen uint16
			binary.Read(rdr, binary.LittleEndian, &pathLen)
			pathBytes := make([]byte, pathLen)
			rdr.Read(pathBytes)
			recipe.ImagePath = string(pathBytes)

		case "INGR":
			var strLen uint16
			binary.Read(rdr, binary.LittleEndian, &strLen)
			strBytes := make([]byte, strLen)
			rdr.Read(strBytes)
			recipe.Ingredients = append(recipe.Ingredients, string(strBytes))

		case "STEP":
			var stepNum uint16
			var textLen uint16
			binary.Read(rdr, binary.LittleEndian, &stepNum)
			binary.Read(rdr, binary.LittleEndian, &textLen)
			textBytes := make([]byte, textLen)
			rdr.Read(textBytes)
			recipe.Steps = append(recipe.Steps, string(textBytes))
		}
	}

	return recipe, nil
}
