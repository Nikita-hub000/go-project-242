package sizefmt

import (
	"fmt"
	"math"
)

const unitBase = 1024.0

var iecUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// FormatIEC formats size as a human-readable IEC-style value.
func FormatIEC(size int64) string {
	if size < 0 {
		return "0B"
	}

	value := float64(size)
	unitIndex := 0

	for unitIndex < len(iecUnits)-1 && value >= unitBase {
		value /= unitBase
		unitIndex++
	}

	if unitIndex == 0 {
		return fmt.Sprintf("%d%s", int64(value), iecUnits[unitIndex])
	}

	rounded := math.Round(value*10) / 10
	if rounded >= unitBase && unitIndex < len(iecUnits)-1 {
		unitIndex++
		rounded = 1.0
	}

	return fmt.Sprintf("%.1f%s", rounded, iecUnits[unitIndex])
}

// FormatRaw formats size as bytes without unit conversion.
func FormatRaw(size int64) string {
	return fmt.Sprintf("%dB", size)
}

// FormatLine formats one CLI output row with size and path.
func FormatLine(size, path string) string {
	return fmt.Sprintf("%s\t%s\n", size, path)
}
