// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Copyright 2020 Johanna Amélie Schander <git@mimoja.de>
package intelfsp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	FSPHeaderSignature            = "FSPH"
	FSPExtHeaderSignature         = "FSPE"
	FixedInfoHeaderLength         = 12
	HeaderFSP1Length              = 64
	HeaderFSP2Length              = 72
	HeaderFSP3Length              = 72
	HeaderFSP4Length              = 72
	FixedExtendedInfoHeaderLength = 20
)

type ImageAttributes struct {
	GraphicsSupport     *bool
	DispatchModeSupport *bool
}

func (ia ImageAttributes) String() string {
	var attrs []string

	if ia.GraphicsSupport != nil && *ia.GraphicsSupport{
		attrs = append(attrs,"Graphics Display Supported")
	} else if ia.GraphicsSupport != nil && !*ia.GraphicsSupport{
		attrs = append(attrs,"Graphics Display Not Supported")
	}

	if ia.DispatchModeSupport != nil && *ia.DispatchModeSupport{
		attrs = append(attrs,"Dispatch Mode Supported")
	} else if ia.GraphicsSupport != nil && !*ia.GraphicsSupport{
		attrs = append(attrs,"Dispatch Mode Not Supported")
	}

	return strings.Join(attrs, "|")
}


type ComponentAttributes struct {
	ReleaseBuild    bool
	OfficialRelease bool
	Type            Type
	TypeName        string
}

func (ca ComponentAttributes) String() string {
	var attrs []string
	if ca.ReleaseBuild {
		attrs = append(attrs, "Release Build")
	} else {
		attrs = append(attrs, "Debug Build")
	}
	if ca.OfficialRelease {
		attrs = append(attrs, "Official Release")
	} else {
		attrs = append(attrs, "Test Release")
	}

	attrs = append(attrs, ca.TypeName)

	return strings.Join(attrs, "|")
}

type BinaryFSPHeader interface {
	Summary() string
	GetImageSize() uint32
	GetImageAttributes() *ImageAttributes
	GetComponentAttributes() *ComponentAttributes
}

type CommonInfoHeader struct {
	Signature      [4]byte
	HeaderLength   uint32
	Reserved1      uint16
	SpecVersion    SpecVersion
	HeaderRevision HeaderRevision
}

// ImageRevision is the image revision field of the FSP info header.
type HeaderRevision uint8

func (hr HeaderRevision) GetHeader() (BinaryFSPHeader, error) {
	switch hr {
	case 0x01:
		return &InfoHeaderV1{}, nil
	case 0x02:
		return &InfoHeaderV2{}, nil
	case 0x03:
		return &InfoHeaderV3{}, nil
	case 0x04:
		return &InfoHeaderV4{}, nil
	default:
		return nil, fmt.Errorf("Unknown Header Revision: 0x%03X", hr)
	}
}

// SpecVersion represents the spec version as a packed BCD two-digit,
// dot-separated unsigned integer.
type SpecVersion uint8

func (sv SpecVersion) String() string {
	return fmt.Sprintf("%d.%d", (sv>>4)&0x0f, sv&0x0f)
}

type Type uint8

var (
	TypeT Type = 1
	TypeM Type = 2
	TypeS Type = 3
	TypeO Type = 8
	// TypeReserved is a fake type that represents a reserved FSP type.
	TypeReserved Type
)

var fspTypeNames = map[Type]string{
	TypeT:        "FSP-T",
	TypeM:        "FSP-M",
	TypeS:        "FSP-S",
	TypeO:        "FSP-O",
	TypeReserved: "FSP-ReservedType",
}

func ParseHeader(b []byte) (*BinaryFSPHeader, error) {
	if len(b) < FixedInfoHeaderLength {
		return nil, fmt.Errorf("short FSP Info Header length %d; want at least %d", len(b), FixedInfoHeaderLength)
	}
	var f CommonInfoHeader

	reader := bytes.NewReader(b)
	if err := binary.Read(reader, binary.LittleEndian, &f); err != nil {
		return nil, err
	}

	if !bytes.Equal(f.Signature[:], []byte(FSPHeaderSignature)) {
		return nil, fmt.Errorf("invalid signature %v (%s); want %s", f.Signature, string(f.Signature[:]), FSPHeaderSignature)
	}

	if f.Reserved1 != 0x0 {
		return nil, fmt.Errorf("reserved bytes must be zero")
	}

	var bh BinaryFSPHeader

	bh, err := f.HeaderRevision.GetHeader()

	if err != nil {
		return nil, fmt.Errorf("Could not determin header version: %v", err)
	}

	reader = bytes.NewReader(b)
	if err := binary.Read(reader, binary.LittleEndian, bh); err != nil {
		return nil, err
	}
	return &bh, nil
}
