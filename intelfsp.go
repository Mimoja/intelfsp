package intelfsp

import (
	"bytes"
	"fmt"
)

// 1.0   https://www.intel.com/content/dam/www/public/us/en/documents/technical-specifications/fsp-architecture-spec.pdf
// 1.1   https://www.intel.com/content/dam/www/public/us/en/documents/technical-specifications/fsp-architecture-spec-v1-1.pdf
// 1.1a  https://www.intel.com/content/dam/www/public/us/en/documents/technical-specifications/fsp-architecture-spec-v1-1a.pdf
// 2.0   https://www.intel.com/content/dam/www/public/us/en/documents/technical-specifications/fsp-architecture-spec-v2.pdf
// 2.1   https://digitallibrary.intel.com/content/dam/ccl/public/intel-fsp-external-architecture-specification-v2-1.pdf

type IntelFSP struct {
	Info               *BinaryFSPHeader
	ImageAttributes    *ImageAttributes
	ComponentAttribute *ComponentAttributes
	ExtendedInfo       *ExtendedFSPHeader
	Raw                []byte
}

func FindAndParse(bts []byte) []IntelFSP {
	fsp := []IntelFSP{}
	index := 0
	offset := 0
	for index >= 0 {
		index = bytes.Index(bts[offset:], []byte(FSPHeaderSignature))

		if index >= 0 &&
			bts[offset+index+8] == 0 && bts[offset+index+9] == 0 {
			hdr, err := Parse(bts[offset+index:])
			if err != nil {
				offset += index + 1
				continue
			}
			if hdr != nil {
				fsp = append(fsp, *hdr)
			}
			offset += index + 1
		}
	}
	return fsp;
}

func Parse(b []byte) (*IntelFSP, error) {
	fsp := IntelFSP{}
	fspInfo, err := ParseHeader(b)

	if err != nil {
		return nil, err
	}
	fsp.Info = fspInfo

	fsp.ImageAttributes = (*fsp.Info).GetImageAttributes()
	fsp.ComponentAttribute = (*fsp.Info).GetComponentAttributes()

	extInfo, err := ParseExtendedHeader(b[72:])
	if err == nil {
		fsp.ExtendedInfo = extInfo
	}

	imageSize := (*fsp.Info).GetImageSize()
	if uint32(len(b)) < imageSize {
		fsp.Raw = b[:imageSize]
	}
	return &fsp, nil
}

func (fsp IntelFSP) Summary() string {
	s := ""
	if fsp.Info != nil {
		s += fmt.Sprintf("--- INFO HEADER ---\n")
		s += (*fsp.Info).Summary()
	}

	if fsp.ImageAttributes != nil {
		s += fmt.Sprintf("--- IMAGE ATTRIBUTES ---\n")
		s += fsp.ImageAttributes.String()
		s += "\n"
	}

	if fsp.ComponentAttribute != nil {
		s += fmt.Sprintf("--- COMPONENT ATTRIBUTES ---\n")
		s += fsp.ComponentAttribute.String()
		s += "\n"
	}
	if fsp.ExtendedInfo != nil {
		s += fmt.Sprintf("--- EXT INFO HEADER ---\n")
		s += fsp.ExtendedInfo.Summary()
	}
	return s
}
