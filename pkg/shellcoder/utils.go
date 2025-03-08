package shellcoder

import (
	"debug/elf"
	"encoding/binary"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func deComment(asm string) string {
	var final string
	for _, each := range strings.Split(asm, "\n") {
		instruction := strings.SplitN(each, ";", 2)[0]
		if strings.TrimSpace(instruction) == "" {
			continue
		}
		final += instruction + "\n"
	}
	return final
}

func PrependPayloadSize(payload []byte) []byte {
	payloadSize := uint32(len(payload))
	lenBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBuf, payloadSize)
	return append(lenBuf, payload...)
}

func GetEntryAddr(f *os.File) int {
	elfFile, err := elf.NewFile(f)
	if err != nil {
		log.Fatal(err)
	}
	entry := elfFile.Entry - 0x40000
	log.Infof("Found Entry: %v %x", entry, entry)
	return int(entry)
}

func MergeBytes(bs ...[]byte) []byte {
	var res []byte
	for _, b := range bs {
		res = append(res, b...)
	}
	return res
}
