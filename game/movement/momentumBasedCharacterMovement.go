package movement

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

// Params holds all movement-related parameters
type Params struct {
	MaxSpeed     float64 // Maximum speed the character can reach
	Acceleration float64 // How fast the character accelerates
	Friction     float64 // How
	Direction    float64
}

// DefaultMovementParams returns sensible default movement parameters
func DefaultMovementParams() Params {
	return Params{
		MaxSpeed:     200.0, // pixels per second
		Acceleration: 800.0, // pixels per second squared
		Friction:     0.85,  // reduces velocity by 15% each frame
	}
}

type State struct {
	VelX, VelY     float64 // Current velocity
	AccelX, AccelY float64 // Current acceleration
}

func (h *WASDInputHandler) GetAccel(params *Params) (accelX, accelY float64) {
	// Track if any movement keys are pressed
	movementInput := false

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		params.Direction -= 0.05
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		params.Direction += 0.05
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		params.Acceleration += 0.01
		movementInput = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		params.Acceleration -= 0.01
		movementInput = true
	}

	// Apply drag when no movement input is detected
	if !movementInput {
		dragFactor := 0.7 // Adjust this value (0.9-0.98 range works well)
		params.Acceleration *= dragFactor

		// Stop very small accelerations to prevent endless tiny movements
		if math.Abs(params.Acceleration) < 0.001 {
			params.Acceleration = 0
		}
	}

	// Cap maximum acceleration to prevent runaway speeds
	maxAccel := 0.5
	if params.Acceleration > maxAccel {
		params.Acceleration = maxAccel
	} else if params.Acceleration < -maxAccel {
		params.Acceleration = -maxAccel
	}

	accelX = math.Cos(params.Direction-math.Pi/2) * params.Acceleration
	accelY = math.Sin(params.Direction-math.Pi/2) * params.Acceleration

	return accelX, accelY
}

type WASDInputHandler struct{}

func (h *WASDInputHandler) GetAcceleration(params *Params) (accelX, accelY float64) {

	input := false

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		params.Direction -= 0.05
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		params.Direction += 0.05
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		params.Acceleration += 0.01
		input = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		params.Acceleration -= 0.005
		input = true
	}

	if !input {
		params.Acceleration = 0
	}

	accelX = math.Cos(params.Direction-math.Pi/2) * params.Acceleration
	accelY = math.Sin(params.Direction-math.Pi/2) * params.Acceleration

	return accelX, accelY
}

// System handles all movement logic

type System struct {
	Params Params
	Input  InputHandler
}

// NewMovementSystem creates a new movement system
func NewMovementSystem(params Params, input InputHandler) *System {
	return &System{
		Params: params,
		Input:  input,
	}
}

type InputHandler interface {
	GetAcceleration(params *Params) (accelX, accelY float64)
}
