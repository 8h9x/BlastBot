package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	b, _ := os.ReadFile("./test/ClientSettings.Sav")

	fmt.Printf("file size: %d\n", len(b))

	if len(b) < 16 {
		log.Fatal("file too short to be zlib")
	}

	magic := string(b[:4])
	uncompressedSize := binary.LittleEndian.Uint32(b[4:8])

	fmt.Printf("magic      : %s\n", magic)
	fmt.Printf("target size: %d bytes\n", uncompressedSize)
	fmt.Printf("header dump: %s\n", hex.EncodeToString(b[:16]))

	buff := bytes.NewReader(b[16:])

	r, err := zlib.NewReader(buff)
	if err != nil {
		log.Fatalf("failed to create zlib reader: %v", err)
	}
	defer r.Close()

	decompressedData, err := io.ReadAll(r)
	if err != nil {
		log.Fatalf("decompression failed: %v", err)
	}

	fmt.Printf("full uncompressed size: %d\n", len(decompressedData))

	startParsing(decompressedData)

	// os.WriteFile("./test/dump2", decompressedData, 0666)
}

func readUEString(r *bytes.Reader) string {
	var length int32
	binary.Read(r, binary.LittleEndian, &length)

	if length == 0 {
		return ""
	}

	isUnicode := length < 0
	if isUnicode {
		length = -length * 2
	}

	buf := make([]byte, length)
	r.Read(buf)
	return strings.TrimRight(string(buf), "\x00")
}

func startParsing(data []byte) {
	reader := bytes.NewReader(data)

	reader.Seek(50, io.SeekStart)

	fmt.Println("--- Starting Property Dump ---")
	for i := 0; i < 10; i++ {
		propName := readUEString(reader)
		if propName == "None" || propName == "" {
			break
		}

		propType := readUEString(reader)

		var propSize int32
		binary.Read(reader, binary.LittleEndian, &propSize)

		var index int32
		binary.Read(reader, binary.LittleEndian, &index)

		fmt.Printf("Property: %s | Type: %s | Data Size: %d bytes\n", propName, propType, propSize)

		reader.Seek(int64(propSize)+1, io.SeekCurrent)
	}
}
