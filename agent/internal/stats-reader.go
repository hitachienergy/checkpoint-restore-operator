package internal

import (
	"archive/tar"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// CRIU magic values

var IMG_SERVICE_MAGIC uint32 = 0x55105940
var STATS_MAGIC uint32 = 0x57093306

func ReadStats(checkpointTar string) []byte {
	checkpointFile, err := os.Open(checkpointTar)
	defer checkpointFile.Close()
	if err != nil {
		return nil
	}
	reader := tar.NewReader(checkpointFile)

	for {
		header, err := reader.Next()
		if err != nil {
			return nil
		}
		if header.Name == "stats-dump" {
			data := make([]byte, header.Size)
			_, err := io.ReadFull(reader, data)
			if err != nil {
				return nil
			}
			magic := binary.LittleEndian.Uint32(data[:4])
			subMagic := binary.LittleEndian.Uint32(data[4:8])
			if magic != IMG_SERVICE_MAGIC {
				fmt.Printf("magic value mismatch, got: %X wanted: %X\n", magic, IMG_SERVICE_MAGIC)
				return nil
			}
			if subMagic != STATS_MAGIC {
				fmt.Printf("sub magic value mismatch, got: %X wanted: %X\n", subMagic, STATS_MAGIC)
				return nil
			}

			return data[12:]
		}
	}
}
