package intelfsp

type InfoHeaderV4 struct {
	InfoHeaderV3
}

func (ih InfoHeaderV4) GetImageAttributes() *ImageAttributes {
	graphicsSupport := ih.ImageAttribute&0b0001 != 0
	dispatchSupport := ih.ImageAttribute&0b0010 != 0
	return &ImageAttributes{
		GraphicsSupport:     &graphicsSupport,
		DispatchModeSupport: &dispatchSupport,
	}
}
