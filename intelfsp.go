package intelfsp

import (
	"fmt"
)

// 1.0   https://www.intel.com/content/dam/www/public/us/en/documents/technical-specifications/fsp-architecture-spec.pdf
// 1.1   https://www.intel.com/content/dam/www/public/us/en/documents/technical-specifications/fsp-architecture-spec-v1-1.pdf
// 1.1a  https://www.intel.com/content/dam/www/public/us/en/documents/technical-specifications/fsp-architecture-spec-v1-1a.pdf
// 2.0   https://www.intel.com/content/dam/www/public/us/en/documents/technical-specifications/fsp-architecture-spec-v2.pdf
// 2.1   https://digitallibrary.intel.com/content/dam/ccl/public/intel-fsp-external-architecture-specification-v2-1.pdf


type IntelFSP struct {
	Info         *BinaryFSPHeader
	ExtendedInfo *ExtendedFSPHeader
}


func Parse(b []byte) (*IntelFSP, error) {
	fsp := IntelFSP{}
	fspInfo, err := ParseHeader(b)

	if err != nil {
		return nil, err
	}
	fsp.Info = fspInfo

	extInfo, err := ParseExtendedHeader(b[72:])
	if err == nil {
		fsp.ExtendedInfo = extInfo
	}
	return &fsp, nil
}


func (fsp IntelFSP) Summary() string {
	s := ""
	if fsp.Info != nil {
		s += fmt.Sprintf("--- INFO HEADER ---\n")
		s += (*fsp.Info).Summary()
	}
	if fsp.ExtendedInfo != nil {
		s += fmt.Sprintf("--- EXT INFO HEADER ---\n")
		s += fsp.ExtendedInfo.Summary()
	}
	return s
}
