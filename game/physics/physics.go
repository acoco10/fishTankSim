package physics

import (
	"fmt"
	"image"
	"math"
)

// Vector2D represents a 2D vector for position, velocity, acceleration
type Vector2D struct {
	X, Y float64
}

type NetShape struct {
	Center     Vector2D
	Radius     float64
	StartAngle float64  // In radians
	EndAngle   float64  // In radians
	Facing     Vector2D // Direction the net is facing
}

type Shape interface {
	Contains(point Vector2D) bool
	GetBounds() image.Rectangle
	GetCenter() Vector2D
	DistanceToEdge(point Vector2D) float64
}

// CircleShape for round objects (fish, bubbles, etc.)
type CircleShape struct {
	Center Vector2D
	Radius float64
}

func (c CircleShape) Contains(point Vector2D) bool {
	return c.Center.Distance(point) <= c.Radius
}

func (c CircleShape) GetBounds() image.Rectangle {
	return image.Rect(
		int(c.Center.X-c.Radius),
		int(c.Center.Y-c.Radius),
		int(c.Center.X+c.Radius),
		int(c.Center.Y+c.Radius),
	)
}

func (c CircleShape) GetCenter() Vector2D {
	return c.Center
}

func (c CircleShape) DistanceToEdge(point Vector2D) float64 {
	return math.Max(0, c.Center.Distance(point)-c.Radius)
}

// NetShape for scooping - semicircle or arc shape

func (n NetShape) Contains(point Vector2D) bool {
	// Check if point is within radius
	distance := n.Center.Distance(point)
	if distance > n.Radius {
		return false
	}

	// Check if point is within the arc angle
	direction := point.Sub(n.Center)
	angle := math.Atan2(direction.Y, direction.X)

	// Normalize angle to 0-2π
	if angle < 0 {
		angle += 2 * math.Pi
	}

	return angle >= n.StartAngle && angle <= n.EndAngle
}

func (n NetShape) GetBounds() image.Rectangle {
	return image.Rect(
		int(n.Center.X-n.Radius),
		int(n.Center.Y-n.Radius),
		int(n.Center.X+n.Radius),
		int(n.Center.Y+n.Radius),
	)
}

func (n NetShape) GetCenter() Vector2D {
	return n.Center
}

func (n NetShape) DistanceToEdge(point Vector2D) float64 {
	if n.Contains(point) {
		return 0
	}
	return n.Center.Distance(point) - n.Radius
}

// PolygonShape for custom shapes (triangular fins, complex objects)
type PolygonShape struct {
	Vertices []Vector2D
	Center   Vector2D
}

func (p PolygonShape) Contains(point Vector2D) bool {
	// Ray casting algorithm for point-in-polygon
	x, y := point.X, point.Y
	inside := false

	j := len(p.Vertices) - 1
	for i := 0; i < len(p.Vertices); i++ {
		xi, yi := p.Vertices[i].X, p.Vertices[i].Y
		xj, yj := p.Vertices[j].X, p.Vertices[j].Y

		if ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}

	return inside
}

func (p PolygonShape) GetBounds() image.Rectangle {
	if len(p.Vertices) == 0 {
		return image.Rectangle{}
	}

	minX, maxX := p.Vertices[0].X, p.Vertices[0].X
	minY, maxY := p.Vertices[0].Y, p.Vertices[0].Y

	for _, v := range p.Vertices {
		minX = math.Min(minX, v.X)
		maxX = math.Max(maxX, v.X)
		minY = math.Min(minY, v.Y)
		maxY = math.Max(maxY, v.Y)
	}

	return image.Rect(int(minX), int(minY), int(maxX), int(maxY))
}

func (p PolygonShape) GetCenter() Vector2D {
	return p.Center
}

func (v Vector2D) Add(other Vector2D) Vector2D {
	return Vector2D{v.X + other.X, v.Y + other.Y}
}

func (v Vector2D) Sub(other Vector2D) Vector2D {
	return Vector2D{v.X - other.X, v.Y - other.Y}
}

func (v Vector2D) Mul(scalar float64) Vector2D {
	return Vector2D{v.X * scalar, v.Y * scalar}
}

func (v Vector2D) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector2D) Normalize() Vector2D {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector2D{0, 0}
	}
	return Vector2D{v.X / mag, v.Y / mag}
}

func (v Vector2D) Distance(other Vector2D) float64 {
	return v.Sub(other).Magnitude()
}

// PhysicsBody represents any object that can have physics applied
type PhysicsBody struct {
	Position     Vector2D
	Velocity     Vector2D
	Acceleration Vector2D

	// Physical properties
	Mass       float64
	Drag       float64 // Air/water resistance (0.0 = no drag, 1.0 = full stop)
	Bounciness float64 // For collisions (0.0 = no bounce, 1.0 = perfect bounce)

	// Constraints
	MaxVelocity float64
	MinVelocity float64 // Minimum speed to keep moving (prevents infinite slow-down)

	// Boundaries
	Bounds         *Rect // If set, object stays within these bounds
	BoundsBehavior BoundsBehavior
	CollisionShape Shape
}

type Rect struct {
	X, Y, Width, Height float64
}

type BoundsBehavior int

const (
	BoundsWrap   BoundsBehavior = iota // Wrap around edges (like Pac-Man)
	BoundsBounce                       // Bounce off edges
	BoundsStop                         // Stop at edges
	BoundsIgnore                       // Ignore bounds entirely
)

// NewPhysicsBody creates a body with sensible defaults
func NewPhysicsBody(x, y float64) *PhysicsBody {
	return &PhysicsBody{
		Position:       Vector2D{x, y},
		Velocity:       Vector2D{0, 0},
		Acceleration:   Vector2D{0, 0},
		Mass:           1.0,
		Drag:           0.01, // Light drag by default
		Bounciness:     0.8,
		MaxVelocity:    200.0,
		MinVelocity:    0.1,
		BoundsBehavior: BoundsBounce,
		CollisionShape: CircleShape{Center: Vector2D{x, y}, Radius: 25},
	}
}

// Update applies physics for one frame (call this each frame)
func (pb *PhysicsBody) Update() {
	if pb.Velocity.X == 0 {
		pb.ApplyForce(Vector2D{4.0, 9.8}) //gravity
	} else {
		pb.ApplyForce(Vector2D{0, 9.8})
	}
	// Apply acceleration to velocity
	deltaTime := 0.016
	pb.Velocity = pb.Velocity.Add(pb.Acceleration.Mul(deltaTime))

	// Apply drag
	if pb.Drag > 0 {
		dragForce := pb.Velocity.Mul(-pb.Drag)
		pb.Velocity = pb.Velocity.Add(dragForce.Mul(deltaTime))

		// Stop very slow movement to prevent infinite creep
		if pb.Velocity.Magnitude() < pb.MinVelocity {
			pb.Velocity = Vector2D{0, 0}
		}
	}

	// Clamp velocity to max speed
	if pb.MaxVelocity > 0 {
		speed := pb.Velocity.Magnitude()
		if speed > pb.MaxVelocity {
			pb.Velocity = pb.Velocity.Normalize().Mul(pb.MaxVelocity)
		}
	}

	// Apply velocity to position
	pb.Position = pb.Position.Add(pb.Velocity.Mul(deltaTime))

	// Handle bounds if set
	if pb.Bounds != nil {
		pb.handleBounds()
	}

	// Reset acceleration (forces must be applied each frame)
	pb.Acceleration = Vector2D{0, 0}
}

// ApplyForce adds a force to the object (F = ma, so force/mass = acceleration)
func (pb *PhysicsBody) ApplyForce(force Vector2D) {
	acceleration := force.Mul(1.0 / pb.Mass)
	pb.Acceleration = pb.Acceleration.Add(acceleration)
}

// ApplyImpulse instantly changes velocity (for collisions, clicks, etc.)
func (pb *PhysicsBody) ApplyImpulse(impulse Vector2D) {
	velocityChange := impulse.Mul(1.0 / pb.Mass)
	pb.Velocity = pb.Velocity.Add(velocityChange)
}

// Attract applies attraction force toward a target point
func (pb *PhysicsBody) Attract(target Vector2D, strength float64) {
	direction := target.Sub(pb.Position)
	distance := direction.Magnitude()

	if distance > 0 {
		// Normalize and apply strength
		force := direction.Normalize().Mul(strength)
		pb.ApplyForce(force)
	}
}

func (pb *PhysicsBody) handleBounds() {
	bounds := pb.Bounds

	switch pb.BoundsBehavior {
	case BoundsWrap:
		if pb.Position.X < bounds.X {
			pb.Position.X = bounds.X + bounds.Width
		} else if pb.Position.X > bounds.X+bounds.Width {
			pb.Position.X = bounds.X
		}

		if pb.Position.Y < bounds.Y {
			pb.Position.Y = bounds.Y + bounds.Height
		} else if pb.Position.Y > bounds.Y+bounds.Height {
			pb.Position.Y = bounds.Y
		}

	case BoundsBounce:
		if pb.Position.X < bounds.X || pb.Position.X > bounds.X+bounds.Width {
			pb.Velocity.X *= -pb.Bounciness
			pb.Position.X = math.Max(bounds.X, math.Min(bounds.X+bounds.Width, pb.Position.X))
		}

		if pb.Position.Y < bounds.Y || pb.Position.Y > bounds.Y+bounds.Height {
			pb.Velocity.Y *= -pb.Bounciness
			pb.Position.Y = math.Max(bounds.Y, math.Min(bounds.Y+bounds.Height, pb.Position.Y))
		}

	case BoundsStop:
		pb.Position.X = math.Max(bounds.X, math.Min(bounds.X+bounds.Width, pb.Position.X))
		pb.Position.Y = math.Max(bounds.Y, math.Min(bounds.Y+bounds.Height, pb.Position.Y))

		if pb.Position.X <= bounds.X || pb.Position.X >= bounds.X+bounds.Width {
			pb.Velocity.X = 0
		}
		if pb.Position.Y <= bounds.Y || pb.Position.Y >= bounds.Y+bounds.Height {
			pb.Velocity.Y = 0
		}
	}
}

// Utility functions for common behaviors
func (pb *PhysicsBody) Stop() {
	pb.Velocity = Vector2D{0, 0}
	pb.Acceleration = Vector2D{0, 0}
}

func (pb *PhysicsBody) SetPosition(x, y float64) {
	pb.Position = Vector2D{x, y}
}

func (pb *PhysicsBody) SetVelocity(x, y float64) {
	pb.Velocity = Vector2D{x, y}
}

func (pb *PhysicsBody) GetIntPosition() (int, int) {
	return int(pb.Position.X), int(pb.Position.Y)
}

func (pb *PhysicsBody) CollidesWith(other *PhysicsBody) bool {
	if pb.CollisionShape == nil || other.CollisionShape == nil {
		return false
	}

	// Check if the other body's center is within this body's shape
	return pb.CollisionShape.Contains(other.Position) ||
		other.CollisionShape.Contains(pb.Position)
}

// CollisionInfo contains details about a collision
type CollisionInfo struct {
	Normal       Vector2D // Direction to separate objects
	Penetration  float64  // How deep the collision is
	ContactPoint Vector2D // Where the collision occurred
}

// GetCollisionInfo returns detailed collision information
func (pb *PhysicsBody) GetCollisionInfo(other *PhysicsBody) *CollisionInfo {
	if !pb.CollidesWith(other) {
		return nil
	}

	// For now, implement circle-circle collision (most common case)
	shape1, ok1 := pb.CollisionShape.(CircleShape)
	shape2, ok2 := other.CollisionShape.(CircleShape)

	if ok1 && ok2 {
		println("detecting circle circle collision")
		return getCircleCircleCollision(shape1, shape2)
	}

	// Fallback: simple separation based on centers
	direction := pb.Position.Sub(other.Position)
	distance := direction.Magnitude()

	if distance == 0 {
		direction = Vector2D{1, 0} // Arbitrary direction
		distance = 1
	}

	return &CollisionInfo{
		Normal:       direction.Normalize(),
		Penetration:  10.0, // Arbitrary penetration
		ContactPoint: pb.Position.Add(other.Position).Mul(0.5),
	}
}

func NewNetBody(x, y, radius float64, startAngle, endAngle float64) *PhysicsBody {
	body := NewPhysicsBody(x, y)
	body.CollisionShape = NetShape{
		Center:     Vector2D{x, y},
		Radius:     radius,
		StartAngle: startAngle,
		EndAngle:   endAngle,
		Facing:     Vector2D{-1, 0},
	}
	body.Mass = 0.5 // Nets are lighter
	return body
}

// Circle-circle collision detection and info
func getCircleCircleCollision(c1, c2 CircleShape) *CollisionInfo {
	direction := c1.Center.Sub(c2.Center)
	distance := direction.Magnitude()
	totalRadius := c1.Radius + c2.Radius

	if distance >= totalRadius {
		return nil // No collision
	}

	if distance == 0 {
		println("setting direction to 1")
		// Circles are exactly on top of each other
		direction = Vector2D{1, 0}
		distance = 1
	}

	normal := direction.Normalize()
	penetration := totalRadius - distance
	contactPoint := c1.Center.Add(normal.Mul(-c1.Radius))

	fmt.Printf("collision detected with normal %f: Penetration %f, contact Point X: %f, Y:%f", normal, penetration, contactPoint.X, contactPoint.Y)

	return &CollisionInfo{
		Normal:       normal,
		Penetration:  penetration,
		ContactPoint: contactPoint,
	}
}

// ResolveCollision handles the physics response between two bodies
func (pb *PhysicsBody) ResolveCollision(other *PhysicsBody, info *CollisionInfo) {
	if info == nil {
		return
	}

	// Separate the objects first
	pb.SeparateFrom(other, info)

	// Calculate relative velocity
	relativeVelocity := pb.Velocity.Sub(other.Velocity)
	velocityAlongNormal := relativeVelocity.X*info.Normal.X + relativeVelocity.Y*info.Normal.Y

	// Don't resolve if velocities are separating
	if velocityAlongNormal > 0 {
		return
	}

	// Calculate restitution (bounciness)
	restitution := math.Min(pb.Bounciness, other.Bounciness)

	// Calculate impulse magnitude
	impulseScalar := -(1 + restitution) * velocityAlongNormal
	impulseScalar /= (1/pb.Mass + 1/other.Mass)

	// Apply impulse
	impulse := info.Normal.Mul(impulseScalar)
	pb.ApplyImpulse(impulse.Mul(1.0 / pb.Mass))
	other.ApplyImpulse(impulse.Mul(-1.0 / other.Mass))
}

// SeparateFrom pushes two overlapping objects apart
func (pb *PhysicsBody) SeparateFrom(other *PhysicsBody, info *CollisionInfo) {
	// Calculate how much to move each object based on their masses
	totalMass := pb.Mass + other.Mass
	myPercent := other.Mass / totalMass // Heavier objects move less
	otherPercent := pb.Mass / totalMass

	// Move objects apart
	separation := info.Normal.Mul(info.Penetration)
	pb.Position = pb.Position.Add(separation.Mul(myPercent))
	other.Position = other.Position.Add(separation.Mul(-otherPercent))
}

// ResolveCollisionWith is a convenience method that does everything
func (pb *PhysicsBody) ResolveCollisionWith(other *PhysicsBody) bool {
	info := pb.GetCollisionInfo(other)
	if info != nil {
		pb.ResolveCollision(other, info)
		return true
	}
	return false
}

func (pb *PhysicsBody) ResolveStaticCollision(staticPosition Vector2D, staticRadius float64) {
	direction := pb.Position.Sub(staticPosition)
	distance := direction.Magnitude()

	// Get radius of this object
	radius := 10.0 // Default
	if shape, ok := pb.CollisionShape.(CircleShape); ok {
		radius = shape.Radius
	}

	minDistance := radius + staticRadius

	if distance < minDistance {
		if distance == 0 {
			direction = Vector2D{1, 0}
			distance = 1
		}

		normal := direction.Normalize()
		penetration := minDistance - distance

		// Push object away from static object
		pb.Position = pb.Position.Add(normal.Mul(penetration))

		// Bounce velocity off the static object
		velocityAlongNormal := pb.Velocity.X*normal.X + pb.Velocity.Y*normal.Y
		if velocityAlongNormal < 0 {
			pb.Velocity = pb.Velocity.Add(normal.Mul(-velocityAlongNormal * (1 + pb.Bounciness)))
		}
	}
}

// Check if a point is inside this body's collision shape
func (pb *PhysicsBody) ContainsPoint(point Vector2D) bool {
	if pb.CollisionShape == nil {
		return false
	}
	return pb.CollisionShape.Contains(point)
}

// Get all objects within this body's collision area (useful for net scooping)
func (pb *PhysicsBody) GetObjectsWithin(others []*PhysicsBody) []*PhysicsBody {
	var contained []*PhysicsBody
	for _, other := range others {
		if pb.ContainsPoint(other.Position) {
			contained = append(contained, other)
		}
	}
	return contained
}
