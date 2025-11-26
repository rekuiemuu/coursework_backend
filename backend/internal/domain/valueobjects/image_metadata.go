package valueobjects

type ImageMetadata struct {
	Width      int
	Height     int
	Format     string
	SizeBytes  int64
	ColorSpace string
	DPI        int
}

func NewImageMetadata(width, height int, format string, sizeBytes int64, colorSpace string, dpi int) *ImageMetadata {
	return &ImageMetadata{
		Width:      width,
		Height:     height,
		Format:     format,
		SizeBytes:  sizeBytes,
		ColorSpace: colorSpace,
		DPI:        dpi,
	}
}

func (im *ImageMetadata) AspectRatio() float64 {
	if im.Height == 0 {
		return 0
	}
	return float64(im.Width) / float64(im.Height)
}

func (im *ImageMetadata) MegaPixels() float64 {
	return float64(im.Width*im.Height) / 1000000.0
}
