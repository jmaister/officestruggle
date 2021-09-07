package state

import "math"

func CalculatePlayerStats(level int) StatsValues {
	return StatsValues{
		Health:    20 * pow(2, level-1),
		MaxHealth: 20 * pow(2, level-1),
		Defense:   4 * pow(2, level-1),
		Power:     12 * pow(2, level-1),
		Fov:       8 + (level - 1),
	}
}

func pow(x int, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}
