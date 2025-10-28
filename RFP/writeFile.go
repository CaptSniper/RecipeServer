package rfp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// writeChunk creates a chunk to the buffer with 8-byte alignment
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

// WriteRecipe writes a Recipe struct into a binary RFP file
func WriteRecipe(dir, filename string, r Recipe) error {
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

	// write property count
	binary.Write(corePayload, binary.LittleEndian, uint16(len(r.CoreProps)))

	// write properties (key/value pairs)
	for k, v := range r.CoreProps {
		binary.Write(corePayload, binary.LittleEndian, uint16(len(k)))
		corePayload.WriteString(k)
		binary.Write(corePayload, binary.LittleEndian, uint16(len(v)))
		corePayload.WriteString(v)
	}

	// write image path
	binary.Write(corePayload, binary.LittleEndian, uint16(len(r.ImagePath)))
	corePayload.WriteString(r.ImagePath)

	// write name last
	binary.Write(corePayload, binary.LittleEndian, uint16(len(r.Name)))
	corePayload.WriteString(r.Name)

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
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

	filename = nonAlphanumericRegex.ReplaceAllString(filename, "")
	return os.WriteFile(filepath.Join(dir, filename+".rfp"), data, 0644)
}
