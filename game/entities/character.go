package entities

import (
	"github.com/acoco10/fishTankWebGame/game/movement"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"image"
	"math"
)

type Corner uint8

const (
	FrontRight Corner = iota
	FrontLeft
	RearRight
	RearLeft
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
	Direction      Direction
	Corners        *TankCorners
	MovementState  movement.State
	Sprite         *sprite.AnimatedSprite
	AnimationMap   map[string]*sprite.AnimatedSprite
	Width, Height  int
	MovementSystem *movement.System
	Collision      bool
	Moving         bool
	state          CharState
}

type Collision struct {
	Corner   Corner
	Velocity float64
	Object   image.Rectangle
}

func NewCharacter(x, y float64, sprite *sprite.AnimatedSprite) *TankCharacter {
	char := TankCharacter{
		MovementState: movement.State{}, // Zero values
		Sprite:        sprite,
		Width:         32,
		Height:        32,
	}

	return &char
}

func CharMovementUpdate(character *TankCharacter) {
	// Get input acceleration
	character.MovementState.AccelX, character.MovementState.AccelY = character.MovementSystem.Input.GetAcceleration(&character.MovementSystem.Params)

	// CharUpdate velocity based on acceleration

	character.MovementState.VelX += character.MovementState.AccelX
	character.MovementState.VelY += character.MovementState.AccelY

	if character.MovementState.AccelX == 0 {
		character.MovementState.VelX *= 0.95
	}
	if character.MovementState.AccelY == 0 {
		character.MovementState.VelY *= 0.95
	}

	speed := math.Sqrt(character.MovementState.VelX*character.MovementState.VelX + character.MovementState.VelY*character.MovementState.VelY)
	if speed > character.MovementSystem.Params.MaxSpeed {
		character.MovementState.VelX = (character.MovementState.VelX / speed) * character.MovementSystem.Params.MaxSpeed
		character.MovementState.VelY = (character.MovementState.VelY / speed) * character.MovementSystem.Params.MaxSpeed
	}

	if !character.Moving {
		return
	}

	// CharUpdate position based on velocity
	character.Sprite.X += float32(character.MovementState.VelX)
	character.Sprite.Y += float32(character.MovementState.VelY)
}

func HandleCollision(character *TankCharacter, collisions []Collision) {

	var frontRight bool
	var frontLeft bool

	for _, collision := range collisions {

		if collision.Corner == FrontRight {
			frontRight = true
		}
		if collision.Corner == FrontLeft {
			frontLeft = true
		}
	}

	if frontLeft && frontRight {
		character.Moving = false
		character.Sprite = character.AnimationMap["Crash"]
		character.Sprite.X = character.AnimationMap["Moving"].X
		character.Sprite.Y = character.AnimationMap["Moving"].Y
		dx, dy := movement.GetCardinalDirection(character.MovementSystem.Params.Direction)
		character.state = Crashed
		// Move back one tile in the opposite direction
		character.Sprite.X -= float32(dx * 16)
		character.Sprite.Y -= float32(dy * 16)
		return
	}

	if frontRight {
		HandleGlance(character, FrontRight)
	}

	if frontLeft {
		HandleGlance(character, FrontLeft)
	}

	// Get the direction the player is facing

}

func HandleCollisionSmooth(player *TankCharacter) {
	// Calculate the exact direction vector
	dirX := math.Cos(player.MovementSystem.Params.Direction)
	dirY := math.Sin(player.MovementSystem.Params.Direction)

	// Push back one tile distance in the opposite direction
	player.Sprite.X += float32(dirX * 32)
	player.Sprite.Y += float32(dirY * 32)
}

type TankCorners struct {
	FrontLeft  image.Point //0
	FrontRight image.Point //1
	RearLeft   image.Point //2
	RearRight  image.Point //3
}

func GetCharCorners(character *TankCharacter) *TankCorners {

	// Get half dimensions
	halfWidth := float32(character.Sprite.SpriteWidth)/2 - 2
	halfHeight := float32(character.Sprite.SpriteHeight)/2 - 1

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
		FrontLeft:  transformPoint(localCorners.frontLeft[0], localCorners.frontLeft[1]),
		FrontRight: transformPoint(localCorners.frontRight[0], localCorners.frontRight[1]),
		RearLeft:   transformPoint(localCorners.rearLeft[0], localCorners.rearLeft[1]),
		RearRight:  transformPoint(localCorners.rearRight[0], localCorners.rearRight[1]),
	}

	return corners

}

func (c *TankCharacter) UpdateCorners() {
	corners := GetCharCorners(c)
	c.Corners = corners
}

func (c *TankCharacter) Update(collisions []Collision) {
	CalcDirection(c)
	c.UpdateCorners()

	//update to current location
	c.Sprite.Update()
	CharMovementUpdate(c) //move
	if len(collisions) > 0 {
		HandleCollision(c, collisions) //handle collisions and glances
	}
	EnforceBoundaries(c, image.Rectangle{image.Point{0, 0}, image.Point{15 * 16, 10 * 16}})
}

func CalcDirection(c *TankCharacter) {
	// Normalize angle to 0-2π range
	angle := math.Mod(c.MovementSystem.Params.Direction-math.Pi/2, 2*math.Pi)
	if angle < 0 {
		angle += 2 * math.Pi
	}

	// Convert to degrees for easier calculation
	degrees := angle * 180 / math.Pi

	// Determine cardinal direction based on angle ranges
	if degrees >= 315 || degrees < 45 {
		c.Direction = CharRight
	} else if degrees >= 45 && degrees < 135 {
		c.Direction = Down
	} else if degrees >= 135 && degrees < 225 {
		c.Direction = CharLeft
	} else {
		c.Direction = Up
	}

}

func HandleGlance(c *TankCharacter, cor Corner) {

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

func EnforceBoundaries(c *TankCharacter, worldBounds image.Rectangle) {
	// Get current tank bounds
	tankBounds := image.Rect(
		int(c.Sprite.X)-c.Width/2,
		int(c.Sprite.Y)-c.Height/2,
		int(c.Sprite.X)+c.Width/2,
		int(c.Sprite.Y)+c.Height/2,
	)

	// Check and correct each boundary
	if tankBounds.Min.X < worldBounds.Min.X {
		c.Sprite.X = float32(worldBounds.Min.X + c.Width/2)
		c.MovementState.VelX = math.Max(0, c.MovementState.VelX) // Only allow positive velocity
	}
	if tankBounds.Max.X > worldBounds.Max.X {
		c.Sprite.X = float32(worldBounds.Max.X - c.Width/2)
		c.MovementState.VelX = math.Min(0, c.MovementState.VelX) // Only allow negative velocity
	}
	if tankBounds.Min.Y < worldBounds.Min.Y {
		c.Sprite.Y = float32(worldBounds.Min.Y + c.Height/2)
		c.MovementState.VelY = math.Max(0, c.MovementState.VelY) // Only allow positive velocity
	}
	if tankBounds.Max.Y > worldBounds.Max.Y {
		c.Sprite.Y = float32(worldBounds.Max.Y - c.Height/2)
		c.MovementState.VelY = math.Min(0, c.MovementState.VelY) // Only allow negative velocity
	}
}
