package game

type PaperBall struct {
	X, Y, Z          float64 // 3D position
	VelX, VelY, VelZ float64 // 3D velocity
	Mass             float64
	Drag             float64 // Air resistance coefficient
}

const (
	Gravity    = 9.8 // Adjust for game feel
	AirDensity = 0.1 // Adjust for paper-like physics
)

func (p *PaperBall) Update(deltaTime float64) {
	// Apply gravity
	p.VelY -= Gravity * deltaTime

	// Simple air resistance (makes it feel more paper-like)
	dragForce := AirDensity * deltaTime
	p.VelX *= (1 - dragForce)
	p.VelY *= (1 - dragForce*0.5) // Less drag on Y for better arc
	p.VelZ *= (1 - dragForce)

	// Update position
	p.X += p.VelX * deltaTime
	p.Y += p.VelY * deltaTime
	p.Z += p.VelZ * deltaTime
}
