package sizefmt

import "fmt"

const unitBase int64 = 1024

var iecUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

func BytesIEC(size int64) string {
	value := float64(size)

	for idx, unit := range iecUnits {
		isLastUnit := idx == len(iecUnits)-1
		if value < float64(unitBase) || isLastUnit {
			if idx == 0 {
				return fmt.Sprintf("%d%s", int64(value), unit)
			}

			return fmt.Sprintf("%.1f%s", value, unit)
		}

		value /= float64(unitBase)
	}

	return "0B"
}

func BytesRaw(size int64) string {
	return fmt.Sprintf("%dB", size)
}
