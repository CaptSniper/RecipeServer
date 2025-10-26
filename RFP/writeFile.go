// rfp3writer.go
package rfp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// writeChunk writes a chunk to the buffer with 8-byte alignment
func writeChunk(buf *bytes.Buffer, chunkType string, payload []byte) error {
	if len(chunkType) != 4 {
		return fmt.Errorf("chunk type must be 4 characters")
	}
	buf.WriteString(chunkType)
	binary.Write(buf, binary.LittleEndian, uint32(len(payload)))
	buf.Write(payload)

	// Add padding to next 8-byte boundary
	padding := (8 - (len(payload) % 8)) % 8
	if padding > 0 {
		buf.Write(make([]byte, padding))
	}
	return nil
}

// WriteRecipe writes a Recipe struct into a binary RFP3 file
func WriteRecipe(filename string, r Recipe) error {
	buf := &bytes.Buffer{}

	// --- HEADER ---
	buf.WriteString("RFP3")                            // Magic
	binary.Write(buf, binary.LittleEndian, uint16(1))  // Version
	binary.Write(buf, binary.LittleEndian, uint16(18)) // Header size
	binary.Write(buf, binary.LittleEndian, uint32(0))  // Placeholder: chunk count
	binary.Write(buf, binary.LittleEndian, uint16(0))  // Flags
	binary.Write(buf, binary.LittleEndian, uint32(0))  // Reserved

	chunkCount := 0

	// --- CORE CHUNK ---
	corePayload := &bytes.Buffer{}
	binary.Write(corePayload, binary.LittleEndian, r.PrepTimeMin)
	binary.Write(corePayload, binary.LittleEndian, r.CookTimeMin)
	binary.Write(corePayload, binary.LittleEndian, r.AdditionalTime)
	binary.Write(corePayload, binary.LittleEndian, r.TotalTimeMin)
	binary.Write(corePayload, binary.LittleEndian, r.Servings)

	binary.Write(corePayload, binary.LittleEndian, uint16(len(r.ImagePath)))
	corePayload.WriteString(r.ImagePath)

	writeChunk(buf, "CORE", corePayload.Bytes())
	chunkCount++

	// --- INGREDIENT CHUNKS ---
	for _, ing := range r.Ingredients {
		ingPayload := &bytes.Buffer{}
		binary.Write(ingPayload, binary.LittleEndian, uint16(len(ing)))
		ingPayload.WriteString(ing)
		writeChunk(buf, "INGR", ingPayload.Bytes())
		chunkCount++
	}

	// --- STEP CHUNKS ---
	for i, step := range r.Steps {
		stepPayload := &bytes.Buffer{}
		binary.Write(stepPayload, binary.LittleEndian, uint16(i+1))
		binary.Write(stepPayload, binary.LittleEndian, uint16(len(step)))
		stepPayload.WriteString(step)
		writeChunk(buf, "STEP", stepPayload.Bytes())
		chunkCount++
	}

	// --- PATCH CHUNK COUNT IN HEADER ---
	data := buf.Bytes()
	binary.LittleEndian.PutUint32(data[0x08:], uint32(chunkCount))

	// --- WRITE FILE ---
	return os.WriteFile(filename, data, 0644)
}
