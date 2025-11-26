package valueobjects

type CapillaryMetrics struct {
	Density       float64
	Diameter      float64
	Tortuosity    float64
	Regularity    float64
	Visibility    float64
	Abnormalities []string
}

func NewCapillaryMetrics(density, diameter, tortuosity, regularity, visibility float64, abnormalities []string) *CapillaryMetrics {
	return &CapillaryMetrics{
		Density:       density,
		Diameter:      diameter,
		Tortuosity:    tortuosity,
		Regularity:    regularity,
		Visibility:    visibility,
		Abnormalities: abnormalities,
	}
}

func (cm *CapillaryMetrics) IsNormal() bool {
	return cm.Density >= 7.0 && cm.Density <= 10.0 &&
		cm.Diameter >= 10.0 && cm.Diameter <= 15.0 &&
		cm.Tortuosity <= 2.0 &&
		len(cm.Abnormalities) == 0
}
