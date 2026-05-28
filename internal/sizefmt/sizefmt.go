package sizefmt

import "fmt"

const unitBase int64 = 1024

var iecUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

func BytesIEC(size int64) string {
	if size < unitBase {
		return fmt.Sprintf("%dB", size)
	}

	value := float64(size)
	unitIdx := 0
	for value >= float64(unitBase) && unitIdx < len(iecUnits)-1 {
		value /= float64(unitBase)
		unitIdx++
	}

	return fmt.Sprintf("%.1f%s", value, iecUnits[unitIdx])
}

func BytesRaw(size int64) string {
	return fmt.Sprintf("%dB", size)
}
