package main

import (
	"github.com/ByteArena/box2d"
	"image"
	"math"
)

type TankCurrent struct {
	centerX, centerY float64 // Center of the circular current
	strength         float64 // How strong the current is
	radius           float64 // How far the current extends
	enabled          bool
	variableStrength bool    // Whether strength varies over time
	baseStrength     float64 // Base strength for variation
}

func NewTankCurrent(centerX, centerY, strength, radius float64) *TankCurrent {
	return &TankCurrent{
		centerX:          centerX,
		centerY:          centerY,
		strength:         strength,
		radius:           radius,
		enabled:          true,
		variableStrength: true,
		baseStrength:     strength,
	}
}

func (tc *TankCurrent) ApplyCurrentForce(body *box2d.B2Body, time float64) {
	if !tc.enabled {
		return
	}

	pos := body.GetPosition()

	// Calculate distance from current center
	dx := pos.X - tc.centerX
	dy := pos.Y - tc.centerY
	distance := math.Sqrt(dx*dx + dy*dy)

	// Only apply current if within radius
	if distance > tc.radius || distance < 0.1 { // Avoid division by zero
		return
	}

	// Calculate current strength based on distance (stronger near center)
	distanceStrength := (tc.radius - distance) / tc.radius

	// Variable strength over time (makes it feel more natural)
	currentStrength := tc.strength
	if tc.variableStrength {
		currentStrength = tc.baseStrength * (0.7 + 0.3*math.Sin(time*2))
	}

	finalStrength := currentStrength * distanceStrength

	// For counter-clockwise circular motion:
	// The velocity should be perpendicular to the radius vector
	// Counter-clockwise means: if radius points right, velocity points up
	forceX := -dy * finalStrength // Perpendicular to dy
	forceY := dx * finalStrength  // Perpendicular to dx

	currentForce := box2d.MakeB2Vec2(forceX, forceY)
	body.ApplyForceToCenter(currentForce, true)
}

// Create multiple current zones for more interesting flow
type MultiZoneCurrent struct {
	currents []TankCurrent
	time     float64
}

func NewMultiZoneCurrent(tankRect image.Rectangle, pixelsPerMeter float32) *MultiZoneCurrent {
	// Convert tank to world coordinates
	leftWorld := float64(tankRect.Min.X) / float64(pixelsPerMeter)
	rightWorld := float64(tankRect.Max.X) / float64(pixelsPerMeter)
	topWorld := float64(tankRect.Min.Y) / float64(pixelsPerMeter)
	bottomWorld := float64(tankRect.Max.Y) / float64(pixelsPerMeter)

	tankWidth := rightWorld - leftWorld
	tankHeight := bottomWorld - topWorld

	currents := []TankCurrent{
		// Main central current
		{
			centerX:          leftWorld + tankWidth*0.5,
			centerY:          topWorld + tankHeight*0.6,
			strength:         0.005,
			radius:           tankWidth * 0.4,
			enabled:          true,
			variableStrength: true,
			baseStrength:     0.005,
		},
		// Smaller secondary current (clockwise to create turbulence)
		{
			centerX:          leftWorld + tankWidth*0.7,
			centerY:          topWorld + tankHeight*0.3,
			strength:         -0.002, // Negative = clockwise
			radius:           tankWidth * 0.15,
			enabled:          true,
			variableStrength: true,
			baseStrength:     0.002,
		},
	}

	return &MultiZoneCurrent{
		currents: currents,
		time:     0,
	}
}

func (mzc *MultiZoneCurrent) Update(dt float64) {
	mzc.time += dt
}

func (mzc *MultiZoneCurrent) ApplyCurrentForces(body *box2d.B2Body) {
	for i := range mzc.currents {
		mzc.currents[i].ApplyCurrentForce(body, mzc.time)
	}
}

type WaterSimulation struct {
	waterLevel float64
	buoyancy   float64
	drag       float64
	enabled    bool
	current    *MultiZoneCurrent // Add current system
}

func NewWaterSimulation(waterLevel, buoyancy, drag float64, tankRect image.Rectangle, pixelsPerMeter float32) *WaterSimulation {
	return &WaterSimulation{
		waterLevel: waterLevel,
		buoyancy:   buoyancy,
		drag:       drag,
		enabled:    true,
		current:    NewMultiZoneCurrent(tankRect, pixelsPerMeter),
	}
}

func (ws *WaterSimulation) ApplyWaterForces(body *box2d.B2Body) {
	if !ws.enabled {
		return
	}

	pos := body.GetPosition()

	// Only apply water forces if object is below water level
	if pos.Y > ws.waterLevel {
		// Apply buoyancy (your existing code)

		dist := pos.Y - ws.waterLevel
		bforce := dist * 0.4 * ws.buoyancy

		buoyancyForce := box2d.MakeB2Vec2(0, -bforce)
		body.ApplyForceToCenter(buoyancyForce, true)

		// Apply drag (your existing code)
		velocity := body.GetLinearVelocity()
		dragForce := box2d.B2Vec2MulScalar(-ws.drag*0.1, velocity)
		body.ApplyForceToCenter(dragForce, true)

		// Apply current forces (NEW!)
		ws.current.ApplyCurrentForces(body)
	}
}

func (ws *WaterSimulation) Update(dt float64) {
	ws.current.Update(dt)
}

type RealisticWaterSimulation struct {
	waterLevel   float64
	fluidDensity float64 // Density of water (kg/m³)
	gravity      float64 // Gravitational acceleration
	drag         float64 // Water resistance
	enabled      bool
	current      *MultiZoneCurrent
	maxDepth     float64 // Bottom of the water body
}

func NewRealisticWaterSimulation(waterLevel, fluidDensity, gravity, drag, maxDepth float64, tankRect image.Rectangle, pixelsPerMeter float32) *RealisticWaterSimulation {
	return &RealisticWaterSimulation{
		waterLevel:   waterLevel,
		fluidDensity: fluidDensity,
		gravity:      gravity,
		drag:         drag,
		enabled:      true,
		current:      NewMultiZoneCurrent(tankRect, pixelsPerMeter),
		maxDepth:     maxDepth,
	}
}

func (rws *RealisticWaterSimulation) ApplyWaterForces(body *box2d.B2Body) {
	if !rws.enabled {
		return
	}

	pos := body.GetPosition()

	// Only apply forces if object is below water level
	if pos.Y <= rws.waterLevel {
		return
	}

	// Calculate how deep the object is submerged
	depth := pos.Y - rws.waterLevel
	maxPossibleDepth := rws.maxDepth - rws.waterLevel

	if depth > maxPossibleDepth {
		depth = maxPossibleDepth // Cap at tank bottom
	}

	// Get object properties
	//mass := body.GetMass()

	// For a sphere/circle, calculate the submerged volume based on depth
	// Simplified: assume the object is small enough to be fully submerged quickly
	fixture := body.GetFixtureList()
	if fixture == nil {
		return
	}

	var objectRadius float64
	var submergedVolume float64

	// Try to get the radius from the shape
	if circleShape, ok := fixture.GetShape().(*box2d.B2CircleShape); ok {
		objectRadius = circleShape.GetRadius()

		// Calculate submerged volume for a sphere partially underwater
		if depth >= objectRadius*2 {
			// Fully submerged
			submergedVolume = (4.0 / 3.0) * math.Pi * objectRadius * objectRadius * objectRadius
		} else if depth > 0 {
			// Partially submerged - use spherical cap formula
			h := math.Min(depth, objectRadius*2) // Height of submerged portion
			submergedVolume = math.Pi * h * h * (3*objectRadius - h) / 3
		} else {
			submergedVolume = 0
		}
	} else {
		// For non-circle shapes, use approximation
		objectRadius = 0.05 // Default small radius
		if depth > 0 {
			// Simple approximation: linear increase with depth
			submergedVolume = math.Pi * objectRadius * objectRadius * math.Min(depth, objectRadius*2)
		}
	}

	if submergedVolume <= 0 {
		return
	}

	// Calculate buoyant force using Archimedes' principle
	// F_buoyant = ρ_fluid × g × V_displaced
	buoyantForce := rws.fluidDensity * rws.gravity * submergedVolume
	buoyancyVector := box2d.MakeB2Vec2(0, -buoyantForce) // Upward force

	body.ApplyForceToCenter(buoyancyVector, true)

	// Apply drag proportional to submerged volume and velocity
	velocity := body.GetLinearVelocity()
	speed := velocity.Length()

	if speed > 0 {
		// Drag increases with submerged volume
		dragCoeff := rws.drag * (submergedVolume / (math.Pi * objectRadius * objectRadius * objectRadius))
		dragMagnitude := dragCoeff * speed * speed * 0.5
		dragDirection := box2d.B2Vec2MulScalar(-1, velocity)
		dragForce := box2d.B2Vec2MulScalar(dragMagnitude, dragDirection)
		body.ApplyForceToCenter(dragForce, true)
	}

	// Apply current forces only when submerged
	rws.current.ApplyCurrentForces(body)

	// Optional: Add some pressure-based effects for realism
	// Deeper objects experience more pressure, affecting their behavior
	pressureEffect := depth / maxPossibleDepth
	if pressureEffect > 0.1 {
		// Slight downward force due to pressure (very small effect)
		pressureForce := box2d.MakeB2Vec2(0, pressureEffect*0.01)
		body.ApplyForceToCenter(pressureForce, true)
	}
}

func (rws *RealisticWaterSimulation) Update(dt float64) {
	rws.current.Update(dt)
}

// Helper function to set realistic material properties
func SetFoodParticleDensity(body *box2d.B2Body, materialDensity float64) {
	// Typical densities (kg/m³):
	// Water: 1000
	// Fish food flakes: 400-800 (less dense than water = floats)
	// Fish food pellets: 1100-1400 (denser than water = sinks)

	fixture := body.GetFixtureList()
	if fixture != nil {
		// Calculate volume and set mass based on density
		if circleShape, ok := fixture.GetShape().(*box2d.B2CircleShape); ok {
			radius := circleShape.GetRadius()
			volume := (4.0 / 3.0) * math.Pi * radius * radius * radius

			// Set the body's density to achieve desired material density
			desiredMass := materialDensity * volume
			fixture.SetDensity(desiredMass / volume)
			body.ResetMassData()
		}
	}
}
