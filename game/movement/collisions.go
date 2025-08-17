package movement

import (
	"math"
)

const TILE_SIZE = 16 // Adjust to your tile size

// Player represents the tank
type Player struct {
	X, Y      float64
	Direction float64 // In radians
}

func GetAngleInDegrees(radians float64) float64 {
	angle := math.Mod(radians-math.Pi/2, 2*math.Pi)
	if angle < 0 {
		angle += 2 * math.Pi
	}

	// Convert to degrees for easier calculation
	degrees := angle * 180 / math.Pi
	return degrees
}

// GetCardinalDirection determines which cardinal direction the player is facing
func GetCardinalDirection(radians float64) (dx, dy int) {
	// Normalize angle to 0-2π range
	degrees := GetAngleInDegrees(radians)

	// Determine cardinal direction based on angle ranges
	if degrees >= 315 || degrees < 45 {
		// Right (0°)
		return 2, 0
	} else if degrees >= 45 && degrees < 135 {
		// Down (90°)
		return 0, 2
	} else if degrees >= 135 && degrees < 225 {
		// Left (180°)
		return -2, 0
	} else {
		// Up (270°)
		return 0, -2
	}
}

// HandleCollision moves the player back one tile in the opposite direction

// Alternative version that pushes back based on exact direction vector
