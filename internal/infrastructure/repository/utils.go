package repository

import "math"

func CalculateTotalPages(total int64, limit int) int {
	if total == 0 || limit == 0 {
		return 1
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}
