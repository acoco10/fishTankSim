package entities

import (
	"github.com/acoco10/fishTankWebGame/game/movement"
	"image"
	"math"
)

type Corner uint8

const (
	FrontRight Corner = iota
	FrontLeft
	RearRight
	RearLeft
	MidPoint
)

type SurfaceType int

const (
	TopSurface SurfaceType = iota
	BottomSurface
	LeftSurface
	RightSurface
)

type CharState uint8

const (
	Starting CharState = iota
	Mowing
	Crashed
)

type CharDirection uint8

const (
	CharRight Direction = iota
	CharLeft
	Down
	Up
)

// State holds current movement state

// TankCharacter represents a character with position and rendering info
type TankCharacter struct {
	Direction       Direction
	Corners         *TankCorners
	Width, Height   int
	Collisions      []Collision
	WorldBoundaries image.Rectangle
	Colliders       []image.Rectangle

	state CharState
}

type Collision struct {
	Angle    float64
	Corner   Corner
	Velocity float64
	Object   image.Rectangle
	Surface  SurfaceType
}

func CharMovementUpdate(character *Entity) {
	// Get input acceleration
	character.MovementState.AccelX, character.MovementState.AccelY = character.MovementSystem.Input.GetAcceleration(&character.MovementSystem.Params)

	// CharUpdate velocity based on acceleration

	character.MovementState.VelX += character.MovementState.AccelX
	character.MovementState.VelY += character.MovementState.AccelY

	if character.MovementState.AccelX == 0 {
		character.MovementState.VelX *= 0.89
	}
	if character.MovementState.AccelY == 0 {
		character.MovementState.VelY *= 0.89
	}

	speed := math.Sqrt(character.MovementState.VelX*character.MovementState.VelX + character.MovementState.VelY*character.MovementState.VelY)
	if speed > character.MovementSystem.Params.MaxSpeed {
		character.MovementState.VelX = (character.MovementState.VelX / speed) * character.MovementSystem.Params.MaxSpeed
		character.MovementState.VelY = (character.MovementState.VelY / speed) * character.MovementSystem.Params.MaxSpeed
	}

	if character.MovementSystem.Params.Acceleration < 0 {
		if speed > character.MovementSystem.Params.MaxSpeed/3 {
			character.MovementState.VelX = (character.MovementState.VelX / speed) * character.MovementSystem.Params.MaxSpeed / 3
			character.MovementState.VelY = (character.MovementState.VelY / speed) * character.MovementSystem.Params.MaxSpeed / 3
		}
	}

	// CharUpdate position based on velocity
	character.Sprite.X += float32(character.MovementState.VelX)
	character.Sprite.Y += float32(character.MovementState.VelY)
}

func HandleCollision(character *Entity, collisions []Collision) {

	var frontRight bool
	var frontLeft bool
	var frontMid bool

	for _, collision := range collisions {
		if collision.Corner == FrontRight {
			frontRight = true
			character.MovementSystem.Params.Acceleration = math.Max(character.MovementSystem.Params.Acceleration-0.2, 0)
		}
		if collision.Corner == FrontLeft {
			frontLeft = true
			character.MovementSystem.Params.Acceleration = math.Max(character.MovementSystem.Params.Acceleration-0.2, 0)

		}
		if collision.Corner == RearRight || collision.Corner == RearLeft {
			character.MovementState.VelX = 0
			character.MovementState.VelY = 0
			character.MovementSystem.Params.Acceleration = 0
		}
		if collision.Corner == MidPoint {
			frontMid = true
		}
	}

	if frontRight {
		HandleGlance(character, FrontRight)
	}

	if frontLeft {
		HandleGlance(character, FrontLeft)
	}

	if frontMid {
		CrashFunc(character)
		return
	}

	// Get the direction the player is facing
}

func CompareSurfaceOfCrash(tankAngle float64, collision Collision, cor Corner) bool {

	switch cor {
	case FrontRight:
		return handleSurfaceCollision(collision.Surface, tankAngle)

	case FrontLeft:
		return handleSurfaceCollision(collision.Surface, tankAngle)

	case RearLeft:
		return handleSurfaceCollision(collision.Surface, tankAngle)

	case RearRight:

		return handleSurfaceCollision(collision.Surface, tankAngle)
	default:
		return false
	}
}

func handleSurfaceCollision(surface SurfaceType, tankAngle float64) bool {
	switch surface {
	case TopSurface:
		// Shallow if tank moving roughly horizontal (parallel to surface)
		return isVerticalMovement(tankAngle)

	case BottomSurface:
		// Shallow if tank moving roughly horizontal
		return isVerticalMovement(tankAngle)

	case LeftSurface:
		// Shallow if tank moving roughly vertical or parallel to surface
		return isHorizontalMovement(tankAngle)

	case RightSurface:
		return isHorizontalMovement(tankAngle)
	default:
		return false
	}
}

func isHorizontalMovement(angle float64) bool {
	return (angle >= -20 && angle <= 10) || (angle >= 170 && angle <= 200)
}

func isVerticalMovement(angle float64) bool {
	return (angle >= 80 && angle <= 110) || (angle >= 260 && angle <= 290)
}

func CrashFunc(character *Entity) {

	dx, dy := movement.GetCardinalDirection(character.MovementSystem.Params.Direction)
	// Move back one tile in the opposite direction
	character.Sprite.X -= float32(dx * 8)
	character.Sprite.Y -= float32(dy * 8)
	collisions := CheckCollision(*character.TankMovement.Corners, character.TankMovement.Colliders)
	x, y := findClosestValidPosition(collisions, *character, 20, 1)
	character.Sprite.X = x
	character.Sprite.Y = y
	character.UpdateCorners()
	character.MovementState.VelX = 0
	character.MovementState.VelY = 0
	character.MovementState.AccelX = 0
	character.MovementState.AccelY = 0
	character.StateMachine.Transition(character)
}

type TankCorners struct {
	FrontLeft     image.Point //0
	FrontRight    image.Point //1
	RearLeft      image.Point //2
	RearRight     image.Point //3
	FrontMidPoint image.Point
}

func GetCharCorners(character *Entity) *TankCorners {

	// Get half dimensions
	halfWidth := float32(character.Sprite.SpriteWidth())/2 - 4
	halfHeight := float32(character.Sprite.SpriteHeight())/2 - 2

	// Center position
	centerX := character.Sprite.X
	centerY := character.Sprite.Y

	// Rotation angle (adjusted for sprite orientation)
	angle := character.MovementSystem.Params.Direction
	cos := float32(math.Cos(angle))
	sin := float32(math.Sin(angle))

	// Local corner coordinates (relative to center)
	localCorners := struct {
		frontLeft, frontRight, rearLeft, rearRight [2]float32
	}{
		frontLeft:  [2]float32{-halfWidth, -halfHeight}, // Front-left
		frontRight: [2]float32{halfWidth, -halfHeight},  // Front-right
		rearRight:  [2]float32{halfWidth, halfHeight},   // Rear-right
		rearLeft:   [2]float32{-halfWidth, halfHeight},  // Rear-left
	}

	// Transform each corner
	transformPoint := func(localX, localY float32) image.Point {
		rotatedX := localX*cos - localY*sin
		rotatedY := localX*sin + localY*cos
		return image.Point{
			X: int(centerX + rotatedX),
			Y: int(centerY + rotatedY),
		}
	}

	// Helper function to transform a local point to world coordinates

	corners := &TankCorners{
		FrontLeft:     transformPoint(localCorners.frontLeft[0], localCorners.frontLeft[1]),
		FrontRight:    transformPoint(localCorners.frontRight[0], localCorners.frontRight[1]),
		RearLeft:      transformPoint(localCorners.rearLeft[0], localCorners.rearLeft[1]),
		RearRight:     transformPoint(localCorners.rearRight[0], localCorners.rearRight[1]),
		FrontMidPoint: transformPoint(0, -halfHeight),
	}

	return corners

}

func (c *Entity) UpdateCorners() {
	corners := GetCharCorners(c)
	c.TankMovement.Corners = corners
}

func (c *Entity) Update(worldBounds image.Rectangle) {
	CalcDirection(c)
	c.UpdateCorners()

	//update to current location

	CharMovementUpdate(c)

	collisions := CheckCollision(*c.TankMovement.Corners, c.TankMovement.Colliders)

	if len(collisions) > 0 {
		HandleCollision(c, collisions) //handle collisions and glances
	}
	collisions = CheckCollision(*c.TankMovement.Corners, c.TankMovement.Colliders)
	if len(collisions) > 0 {
		x, y := findClosestValidPosition(collisions, *c, 10, 1)
		c.Sprite.X = x
		c.Sprite.Y = y
		c.UpdateCorners()
	}

	c.Sprite.Update()
}

func findClosestValidPosition(collisions []Collision, character Entity, maxRadius int, step float32) (float32, float32) {
	if len(collisions) == 0 {
		return character.Sprite.X, character.Sprite.Y
	}

	startX := character.Sprite.X
	startY := character.Sprite.Y

	// Simple approach: try moving in all 8 directions
	directions := [][2]int{
		{0, -1},  // up
		{0, 1},   // down
		{-1, 0},  // left
		{1, 0},   // right
		{-1, -1}, // up-left
		{1, -1},  // up-right
		{-1, 1},  // down-left
		{1, 1},   // down-right
	}

	for r := 1; r <= maxRadius; r++ {
		for _, dir := range directions {
			testX := startX + float32(dir[0]*r)*step
			testY := startY + float32(dir[1]*r)*step

			character.Sprite.X = testX
			character.Sprite.Y = testY
			character.UpdateCorners()
			if CheckPosition(*character.TankMovement.Corners, character.TankMovement.Colliders) {
				return testX + float32(dir[0]*r), testY + float32(dir[1]*r)
			}
		}
	}
	println("No Valid Position")
	return startX, startY
}

func CheckPosition(corners TankCorners, colliders []image.Rectangle) bool {
	minX := min(corners.FrontLeft.X, corners.FrontRight.X, corners.RearLeft.X, corners.RearRight.X)
	maxX := max(corners.FrontLeft.X, corners.FrontRight.X, corners.RearLeft.X, corners.RearRight.X)
	minY := min(corners.FrontLeft.Y, corners.FrontRight.Y, corners.RearLeft.Y, corners.RearRight.Y)
	maxY := max(corners.FrontLeft.Y, corners.FrontRight.Y, corners.RearLeft.Y, corners.RearRight.Y)

	rect := image.Rect(minX, minY, maxX, maxY)
	for _, col := range colliders {
		if rect.Overlaps(col) {
			return false
		}
	}
	return true
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func CheckCollision(corners TankCorners, colliders []image.Rectangle) []Collision {
	var collisions []Collision

	shrinkPercent := float32(0.10)
	dx := float32(corners.FrontRight.X - corners.FrontLeft.X)
	dy := float32(corners.FrontRight.Y - corners.FrontLeft.Y)

	xChange := int(dx * shrinkPercent)
	yChange := int(dy * shrinkPercent)

	x1 := corners.FrontRight.X
	x2 := corners.FrontLeft.X
	y1 := corners.FrontRight.Y
	y2 := corners.FrontLeft.Y

	lineMinX := min(x1, x2)
	lineMaxX := max(x1, x2)
	lineMinY := min(y1, y2)
	lineMaxY := max(y1, y2)

	lineBounds := image.Rect(lineMinX+xChange, lineMinY+yChange, lineMaxX-xChange, lineMaxY-yChange)

	// Check if bounding boxes overlap

	for _, col := range colliders {
		if corners.FrontRight.In(col) {
			collision := Collision{Corner: FrontRight, Object: col}
			collisions = append(collisions, collision)

		}
		if corners.FrontLeft.In(col) {
			collision := Collision{Corner: FrontLeft, Object: col}
			collisions = append(collisions, collision)
		}
		if corners.RearLeft.In(col) {
			collision := Collision{Corner: RearLeft, Object: col}
			collisions = append(collisions, collision)

		}
		if corners.RearRight.In(col) {
			collision := Collision{Corner: RearRight, Object: col}
			collisions = append(collisions, collision)

		}
		if lineBounds.Overlaps(col) {

			intersection := lineBounds.Intersect(col)
			penetrationX := min(intersection.Dx(), lineBounds.Dx())
			penetrationY := min(intersection.Dy(), lineBounds.Dy())
			maxPenetration := max(penetrationX, penetrationY)
			if maxPenetration > 7 {
				collision := Collision{Corner: MidPoint, Object: col}
				collisions = append(collisions, collision)
			}
		}
	}
	return collisions
}

func CalcDirection(c *Entity) {
	// Normalize angle to 0-2π range
	angle := math.Mod(c.MovementSystem.Params.Direction-math.Pi/2, 2*math.Pi)
	if angle < 0 {
		angle += 2 * math.Pi
	}

	// Convert to degrees for easier calculation
	degrees := angle * 180 / math.Pi

	// Determine cardinal direction based on angle ranges
	if degrees >= 315 || degrees < 45 {
		c.TankMovement.Direction = CharRight
	} else if degrees >= 45 && degrees < 135 {
		c.TankMovement.Direction = Down
	} else if degrees >= 135 && degrees < 225 {
		c.TankMovement.Direction = CharLeft
	} else {
		c.TankMovement.Direction = Up
	}

}

func HandleGlance(c *Entity, cor Corner) {

	if c.MovementState.VelY > 0.2 {
		dx, dy, directionModifier := Glance(cor)
		c.Sprite.X += dx * 16
		c.Sprite.Y += dy * 16
		c.MovementSystem.Params.Direction += directionModifier
	} else {
		dx, dy, directionModifier := Glance(cor)
		c.Sprite.X += dx * 16
		c.Sprite.Y += dy * 16
		c.MovementSystem.Params.Direction += directionModifier
		c.MovementState.VelX = 0.0
		c.MovementState.VelY = 0.0
	}

}

func Glance(cor Corner) (dx float32, dy float32, direction float64) {
	switch cor {
	case FrontLeft:
		// Hit front-left: back up and turn right
		return float32(0), float32(-0.2), 0.3
	case FrontRight:
		// Hit front-right: back up and turn left
		return float32(0), float32(-0.1), -0.3
	case RearLeft:
		// Hit rear-left: move forward and turn right
		return float32(0), float32(0.1), 0.3
	case RearRight:
		// Hit rear-right: move forward and turn left
		return float32(0), float32(0.1), -0.3
	default:
		return 0.0, 0.0, 0.0
	}
}

type StateData struct {
	name         CharState
	condition    func(character *TankCharacter) bool
	transitionTo CharState
}
