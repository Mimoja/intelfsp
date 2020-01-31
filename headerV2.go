// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Copyright 2020 Johanna Amélie Schander <git@mimoja.de>
package intelfsp

import (
	"fmt"
)

type InfoHeaderV2 struct {
	CommonInfoHeader
	ImageRevision             uint32
	ImageID                   [8]byte
	ImageSize                 uint32
	ImageBase                 uint32
	ImageAttribute            uint32
	CfgRegionOffset           uint32
	CfgRegionSize             uint32
	ApiEntryNum               uint32
	TempRAMInitEntryOffset    uint32
	FspInitEntryOffset        uint32
	NotifyPhaseEntryOffset    uint32
	FSPMemoryInitEntryOffset  uint32
	TempRAMExitEntryOffset    uint32
	FSPSiliconInitEntryOffset uint32
}

func (ih InfoHeaderV2) Summary() string {
	s := fmt.Sprintf("Signature                   : %s\n", ih.Signature)
	s += fmt.Sprintf("Header Length               : %d\n", ih.HeaderLength)
	s += fmt.Sprintf("Reserved1                   : %#04x\n", ih.Reserved1)
	s += fmt.Sprintf("Spec Version                : %s\n", ih.SpecVersion)
	s += fmt.Sprintf("Header Revision             : %d\n", ih.HeaderRevision)
	s += fmt.Sprintf("Image Revision              : %#08x\n", ih.ImageRevision)
	s += fmt.Sprintf("Image ID                    : %#08x\n", ih.ImageID)
	s += fmt.Sprintf("Image Size                  : %#08x %d\n", ih.ImageSize, ih.ImageSize)
	s += fmt.Sprintf("Image Base                  : %#08x %d\n", ih.ImageBase, ih.ImageBase)
	s += fmt.Sprintf("Image Attribute             : %#08x\n", ih.ImageAttribute)
	s += fmt.Sprintf("Image Attribute             : %#08x\n", ih.ImageAttribute)

	s += fmt.Sprintf("Cfg Region Offset           : %#08x %d\n", ih.CfgRegionOffset, ih.CfgRegionOffset)
	s += fmt.Sprintf("Cfg Region Size             : %#08x %d\n", ih.CfgRegionSize, ih.CfgRegionSize)
	s += fmt.Sprintf("API Entry Num               : %#08x\n", ih.ApiEntryNum)
	s += fmt.Sprintf("TempRAMInit Entry Offset    : %#08x %d\n", ih.TempRAMInitEntryOffset, ih.TempRAMInitEntryOffset)
	s += fmt.Sprintf("FSP Init Entry Offset       : %#08x\n", ih.FspInitEntryOffset)
	s += fmt.Sprintf("NotifyPhase Entry Offset    : %#08x %d\n", ih.NotifyPhaseEntryOffset, ih.NotifyPhaseEntryOffset)
	s += fmt.Sprintf("FSPMemoryInit Entry Offset  : %#08x %d\n", ih.FSPMemoryInitEntryOffset, ih.FSPMemoryInitEntryOffset)
	s += fmt.Sprintf("TempRAMExit Entry Offset    : %#08x %d\n", ih.TempRAMExitEntryOffset, ih.TempRAMExitEntryOffset)
	s += fmt.Sprintf("FSPSiliconInit Entry Offset : %#08x %d\n", ih.FSPSiliconInitEntryOffset, ih.FSPSiliconInitEntryOffset)
	return s
}

func (ih InfoHeaderV2) GetImageSize() uint32 {
	return ih.ImageSize
}

func (ih InfoHeaderV2) GetImageAttributes() *ImageAttributes {
	graphicsSupport := ih.ImageAttribute&0b0001 != 0
	return &ImageAttributes{
		GraphicsSupport: &graphicsSupport,
	}
}
func (ih InfoHeaderV2) GetComponentAttributes() *ComponentAttributes { return nil }
