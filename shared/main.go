package shared

import (
	"bytes"
	"encoding/binary"
	"log"
)

func Float64SliceToBytes(floats [][2]float64) []byte {
	var flat []float64
	for idx, slice := range floats {
		flat[idx] = slice[0]
	}

	buf := new(bytes.Buffer)

	for _, f := range flat {
		binary.Write(buf, binary.LittleEndian, f)
	}
	return buf.Bytes()
}

func BytesToFloat64Slice(b []byte) []float64 {
	floats := make([]float64, len(b)/8)
	buf := bytes.NewReader(b)
	for i := range floats {
		err := binary.Read(buf, binary.LittleEndian, &floats[i])
		if err != nil {
			log.Fatalf("binary.Read failed: %v", err)
		}
	}
	return floats
}
