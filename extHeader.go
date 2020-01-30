// Copyright 2020 Johanna Amélie Schander <git@mimoja.de>
package intelfsp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type ExtendedFSPHeader struct {
	BinaryExtendedFSPHeader
	FspProducerData []byte
}

type BinaryExtendedFSPHeader struct {
	Signature           [4]byte
	HeaderLength        uint32
	Revision            uint8
	Reserved1           uint8
	FspProducerId       [6]byte
	FspProducerRevision uint32
	FspProducerDataSize uint32
}

func ParseExtendedHeader(b []byte) (*ExtendedFSPHeader, error) {
	binHeader := BinaryExtendedFSPHeader{}
	if len(b) < FixedExtendedInfoHeaderLength {
		return nil, fmt.Errorf("short FSP Info Header length %d; want at least %d", len(b), FixedInfoHeaderLength)
	}

	reader := bytes.NewReader(b)
	if err := binary.Read(reader, binary.LittleEndian, &b); err != nil {
		return nil, err
	}

	if !bytes.Equal(binHeader.Signature[:], []byte(FSPExtHeaderSignature)) {
		return nil, fmt.Errorf("invalid signature %v (%s); want %s", binHeader.Signature, string(binHeader.Signature[:]), FSPHeaderSignature)
	}

	if binHeader.Reserved1 != 0x0 {
		return nil, fmt.Errorf("reserved bytes must be zero")
	}

	if uint32(len(b)) < FixedExtendedInfoHeaderLength+binHeader.FspProducerDataSize {
		return nil, fmt.Errorf("extended info header size requested 0x%08x bytes", binHeader.FspProducerDataSize)
	}
	extHeader := ExtendedFSPHeader{
		BinaryExtendedFSPHeader: binHeader,
		FspProducerData:         b[FixedExtendedInfoHeaderLength : FixedExtendedInfoHeaderLength+binHeader.FspProducerDataSize],
	}

	return &extHeader, nil
}

func (eh ExtendedFSPHeader) Summary() string {
	s := fmt.Sprintf("Signature                   : %s\n", eh.Signature)
	s += fmt.Sprintf("Header Length               : %d\n", eh.HeaderLength)
	s += fmt.Sprintf("Revision                    : %#02x\n", eh.Revision)
	s += fmt.Sprintf("Reserved1                   : %#04x\n", eh.Reserved1)
	s += fmt.Sprintf("FSP Producer ID             : %v\n", eh.FspProducerId)
	s += fmt.Sprintf("FSP Producer Revision       : %#08x\n", eh.FspProducerRevision)
	s += fmt.Sprintf("FSP Producer Data Size      : %#08x %d\n", eh.FspProducerDataSize, eh.FspProducerDataSize)
	s += fmt.Sprintf("FSP Producer Data           : %v\n", eh.FspProducerData)

	return s
}
