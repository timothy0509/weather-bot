package hko

import "fmt"

// ReadingByPlace returns the first reading matching a place.
func ReadingByPlace(rg ReadingGroup, place string) (Reading, bool) {
	for _, r := range rg.Data {
		if r.Place == place {
			return r, true
		}
	}
	return Reading{}, false
}

// FormatTemperature formats a temperature reading.
func (r Reading) FormatTemperature() string {
	if r.Value == 0 && r.Place == "" {
		return "N/A"
	}
	unit := r.Unit
	if unit == "C" {
		unit = "°C"
	}
	return fmt.Sprintf("%.1f%s", r.Value, unit)
}

// FormatRainfall formats a rainfall reading.
func (r RainfallReading) FormatRainfall() string {
	unit := r.Unit
	if unit == "mm" && r.Max == 0 && r.Min == 0 {
		return "0 mm"
	}
	if r.Max == r.Min {
		return fmt.Sprintf("%.1f %s", r.Max, unit)
	}
	return fmt.Sprintf("%.1f–%.1f %s", r.Min, r.Max, unit)
}

// MaxRainfallReading returns the reading with the maximum rainfall.
func MaxRainfallReading(g RainfallGroup) (RainfallReading, bool) {
	if len(g.Data) == 0 {
		return RainfallReading{}, false
	}
	max := g.Data[0]
	for _, r := range g.Data[1:] {
		if r.Max > max.Max {
			max = r
		}
	}
	return max, true
}
