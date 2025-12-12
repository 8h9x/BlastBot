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

	os.WriteFile("./test/dump", decompressedData, 0666)
}
