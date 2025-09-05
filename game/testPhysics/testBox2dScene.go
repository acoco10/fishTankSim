package main

import (
	"github.com/ByteArena/box2d"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"log"
)

const (
	screenWidth    = 960
	screenHeight   = 540
	pixelsPerMeter = 50 // scale factor: 1m = 50px
	zoomFactor     = 2  // How much we're scaling up the final image
)

type physicsSprite struct {
	img     *ebiten.Image
	img2    *ebiten.Image
	ball    []*box2d.B2Body
	physics *box2d.B2Body
}

func (ps *physicsSprite) draw(screen *ebiten.Image) {
	// Debug: Draw a small circle where Box2D thinks the center is
	pos := ps.physics.GetPosition()
	centerX := float32(pos.X * pixelsPerMeter)
	centerY := float32(pos.Y * pixelsPerMeter)
	vector.DrawFilledCircle(screen, centerX, centerY, 3, color.RGBA{255, 0, 0, 255}, true)

	// Your existing drawing code...
	dopts := &ebiten.DrawImageOptions{}
	theta := ps.physics.GetAngle()
	imgWidth := float64(ps.img.Bounds().Dx())
	imgHeight := float64(ps.img.Bounds().Dy())

	dopts.GeoM.Translate(-imgWidth/2, -imgHeight/2)
	dopts.GeoM.Rotate(theta)
	dopts.GeoM.Translate(pos.X*pixelsPerMeter, pos.Y*pixelsPerMeter)
	screen.DrawImage(ps.img2, dopts)
	for _, ball := range ps.ball {
		drawFoodFlake(ball, screen)
	}
	screen.DrawImage(ps.img, dopts)
}

type Game struct {
	fishfoodContainer *physicsSprite
	world             *box2d.B2World
	cup               *box2d.B2Body
	tankRect          image.Rectangle
	angle             float64
	lines             [][2]box2d.B2Vec2
	tankLines         [][2]box2d.B2Vec2
	tableLine         [2]box2d.B2Vec2
	tank              *box2d.B2Body
	selected          bool
	scoop             *box2d.B2Body
	scoopLines        [][2]box2d.B2Vec2
	table             *box2d.B2Body
	smallRes          *ebiten.Image
	water             *WaterSimulation
}

func MakeBall(position box2d.B2Vec2, world *box2d.B2World) *box2d.B2Body {

	ballDef := box2d.MakeB2BodyDef()
	ballDef.Type = box2d.B2BodyType.B2_dynamicBody
	ballDef.Position.Set(position.X, position.Y)

	ball := world.CreateBody(&ballDef)

	ballShape := box2d.MakeB2CircleShape()
	ballShape.M_radius = 0.05

	fixtureDef := box2d.MakeB2FixtureDef()
	fixtureDef.Shape = &ballShape
	fixtureDef.Density = 0.25
	fixtureDef.Friction = 1.0
	fixtureDef.Restitution = 0.0
	ball.CreateFixtureFromDef(&fixtureDef)
	return ball
}

func (g *Game) init() {
	smallRes := ebiten.NewImage(screenWidth, screenHeight)
	// Box2D world with gravity
	gravity := box2d.MakeB2Vec2(0, 9.8)
	world := box2d.MakeB2World(gravity)
	cupDef := box2d.MakeB2BodyDef()
	cupDef.Type = box2d.B2BodyType.B2_kinematicBody
	cupDef.Position.Set(5, 2) // center of cup
	cup := world.CreateBody(&cupDef)

	// Define 3 edges for a "U" shape (local coordinates relative to cup center)
	lines := [][2]box2d.B2Vec2{
		// Main bottom section
		{box2d.MakeB2Vec2(-0.1, 0.59), box2d.MakeB2Vec2(0.1, 0.59)}, // center bottom

		// Left bottom bevel (two segments for smoother curve)
		{box2d.MakeB2Vec2(-0.25, 0.45), box2d.MakeB2Vec2(-0.18, 0.54)}, // left bevel segment 1
		{box2d.MakeB2Vec2(-0.18, 0.54), box2d.MakeB2Vec2(-0.1, 0.59)},  // left bevel segment 2

		// Right bottom bevel (two segments for smoother curve)
		{box2d.MakeB2Vec2(0.1, 0.59), box2d.MakeB2Vec2(0.18, 0.54)},  // right bevel segment 1
		{box2d.MakeB2Vec2(0.18, 0.54), box2d.MakeB2Vec2(0.25, 0.45)}, // right bevel segment 2

		// Left wall (connects to bevel)
		{box2d.MakeB2Vec2(-0.25, -0.5), box2d.MakeB2Vec2(-0.25, 0.45)}, // left wall

		// Right wall (connects to bevel)
		{box2d.MakeB2Vec2(0.25, -0.5), box2d.MakeB2Vec2(0.25, 0.45)}, // right wall

		// Bottom opening edges
		{box2d.MakeB2Vec2(0.25, -0.5), box2d.MakeB2Vec2(0.14, -0.5)}, // right opening edge
		{box2d.MakeB2Vec2(-0.3, -0.5), box2d.MakeB2Vec2(0.0, -0.5)},  // left opening edge
	}

	scooplines := [][2]box2d.B2Vec2{
		// Main bottom section
		{box2d.MakeB2Vec2(0.5, 0.0), box2d.MakeB2Vec2(0.5, 0.5)},   // center bottom
		{box2d.MakeB2Vec2(0.5, 0.0), box2d.MakeB2Vec2(-0.5, -0.5)}, // left bevel segment 1
		{box2d.MakeB2Vec2(0.5, 0.5), box2d.MakeB2Vec2(-0.5, 0.45)}, // left bevel segment 1

	}

	for _, seg := range lines {
		edge := box2d.MakeB2EdgeShape()
		edge.Set(seg[0], seg[1])
		fixtureDef := box2d.MakeB2FixtureDef()
		fixtureDef.Shape = &edge
		fixtureDef.Density = 1.0     // Add density
		fixtureDef.Friction = 0.6    // Optional: add friction
		fixtureDef.Restitution = 0.0 // Optional: add bounce

		cup.CreateFixtureFromDef(&fixtureDef)
	}

	tankRect := image.Rect(200, 200, 550, 300)

	tank, waterSim, tanklines := SetupTankFromRect(&world, tankRect, pixelsPerMeter)

	img, err := util.LoadImageAssetAsEbitenImage("uiSprites/fishFoodOpen")

	if err != nil {
		log.Fatal(err)
	}

	img2, err := util.LoadImageAssetAsEbitenImage("uiSprites/fishFoodBottom")

	tableDef := box2d.MakeB2BodyDef()
	tableDef.Type = box2d.B2BodyType.B2_staticBody
	tableLine := [2]box2d.B2Vec2{box2d.MakeB2Vec2(-100, 0), box2d.MakeB2Vec2(100, 0)} // bottom
	tableDef.Position.Set(5, 10)
	table := world.CreateBody(&tableDef)
	tedge := box2d.MakeB2EdgeShape()
	tedge.Set(tableLine[0], tableLine[1])
	table.CreateFixture(&tedge, 0.0)

	// Attach fixtures
	for _, seg := range lines {
		edge := box2d.MakeB2EdgeShape()
		edge.Set(seg[0], seg[1])
		fixtureDef := box2d.MakeB2FixtureDef()
		fixtureDef.Shape = &edge
		fixtureDef.Density = 1.0     // Add density
		fixtureDef.Friction = 0.6    // Optional: add friction
		fixtureDef.Restitution = 0.0 // Optional: add bounce

		cup.CreateFixtureFromDef(&fixtureDef)
	}

	cupPos := cup.GetPosition()
	var balls []*box2d.B2Body

	ballRadius := 0.05
	layerHeight := ballRadius * 1.1 // Slightly compressed layers

	for i := 0; i < 75; i++ {
		layer := i / 10     // 5 balls per layer
		posInLayer := i % 5 // Position within layer

		// Alternate layer patterns for better stacking
		var x float64
		if layer%2 == 0 {
			// Even layers: spread across cup width
			x = cupPos.X + 0.05 + (float64(posInLayer)-2)*0.08 // -2 to center around 0
		} else {
			// Odd layers: offset pattern
			x = cupPos.X + 0.05 + (float64(posInLayer)-2)*0.08 + 0.04
		}

		y := cupPos.Y - 0.2 + float64(layer)*layerHeight

		posVec := box2d.B2Vec2{x, y}
		ball := MakeBall(posVec, &world)
		balls = append(balls, ball)
	}

	fishFood := &physicsSprite{img: img, img2: img2, ball: balls, physics: cup}

	g.world = &world
	g.cup = cup
	g.fishfoodContainer = fishFood
	g.lines = lines
	g.table = table
	g.tableLine = tableLine
	g.smallRes = smallRes
	g.water = waterSim
	g.tankRect = tankRect
	g.tankLines = tanklines
	g.scoopLines = scooplines
	g.tank = tank

}

func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) Update() error {
	// Step physics
	timeStep := 1.0 / 60.0
	g.world.Step(timeStep, 8, 3)

	// Rotate cup on mouse click
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		pos := g.cup.GetPosition()
		mx, my := ebiten.CursorPosition()
		// Convert from scaled screen coordinates to world coordinates
		mouseWorldX := float64(mx) / (pixelsPerMeter * zoomFactor)
		mouseWorldY := float64(my) / (pixelsPerMeter * zoomFactor)

		mouseVec := box2d.B2Vec2{mouseWorldX, mouseWorldY}
		distance := box2d.B2Vec2Distance(mouseVec, pos)
		if distance < 0.5 {
			g.selected = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.init()
		g.selected = false
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.selected {
		g.angle -= 0.05
		cupPos := g.cup.GetPosition()
		g.cup.SetTransform(cupPos, g.angle) // Only for rotation
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && g.selected {
		g.angle += 0.05
		cupPos := g.cup.GetPosition()
		g.cup.SetTransform(cupPos, g.angle) // Only for rotation
	}

	if g.selected {
		mx, my := ebiten.CursorPosition()
		// Convert from scaled screen coordinates to world coordinates
		mouseWorldX := float64(mx) / (pixelsPerMeter * zoomFactor)
		mouseWorldY := float64(my) / (pixelsPerMeter * zoomFactor)

		// Get current cup position
		cupPos := g.cup.GetPosition()

		// Calculate velocity needed to reach mouse position
		speed := 5.0 // Adjust this for responsiveness
		dx := mouseWorldX - cupPos.X
		dy := mouseWorldY - cupPos.Y

		// Set linear velocity instead of teleporting
		g.cup.SetLinearVelocity(box2d.B2Vec2{dx * speed, dy * speed})

	}

	for _, ball := range g.fishfoodContainer.ball {
		g.water.ApplyWaterForces(ball)
	}

	return nil
}

func TransformPoint(xf box2d.B2Transform, v box2d.B2Vec2) box2d.B2Vec2 {
	// rotated point
	vx := xf.Q.C*v.X - xf.Q.S*v.Y
	vy := xf.Q.S*v.X + xf.Q.C*v.Y
	// translated
	return box2d.B2Vec2{X: vx + xf.P.X, Y: vy + xf.P.Y}
}

func drawFoodFlake(flake *box2d.B2Body, screen *ebiten.Image) {
	pos := flake.GetPosition()
	angle := flake.GetAngle()

	x := float32(pos.X * pixelsPerMeter)
	y := float32(pos.Y * pixelsPerMeter)

	// Draw a 1x3 pixel rectangle
	width := float32(3.0)
	height := float32(2.0)

	// Create transformation matrix for rotation
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(-width)/2, float64(-height)/2) // Center
	opts.GeoM.Rotate(angle)                                    // Rotate
	opts.GeoM.Translate(float64(x), float64(y))                // Position

	// Create a small rectangle image or use vector drawing
	vector.DrawFilledRect(screen, x-width/2, y-height/2, width, height,
		color.RGBA{255, 255, 19, 255}, true)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.smallRes.Fill(color.RGBA{30, 30, 30, 255})

	// Draw to the small resolution image (600x600)
	// All coordinates are already correct for this size

	//xf2 := g.table.GetTransform()
	/*DrawLine(xf2, g.tableLine, g.smallRes)*/

	/*xf := g.cup.GetTransform()
	for _, seg := range g.lines {
		DrawLine(xf, seg, g.smallRes)
	}*/

	DrawWaterSurface(g.smallRes, g.tankRect, g.water.waterLevel, pixelsPerMeter)

	xf := g.tank.GetTransform()
	for _, seg := range g.tankLines {
		DrawLine(xf, seg, g.smallRes)
	}

	/*	x1 := float32(g.tank.GetPosition().X) * pixelsPerMeter
		y1 := float32(g.water.waterLevel) * pixelsPerMeter

		x2 := x1 + 1000
		y2 := float32(g.water.waterLevel) * pixelsPerMeter*/

	//vector.StrokeLine(screen, x1, y1, x2, y2, 2, color.RGBA{50, 200, 50, 255}, true)

	g.fishfoodContainer.draw(g.smallRes)

	// Scale up the small image to fill the larger window
	dopts := &ebiten.DrawImageOptions{}
	dopts.GeoM.Scale(zoomFactor, zoomFactor)
	screen.DrawImage(g.smallRes, dopts)
}

func DrawLine(transformationPoint box2d.B2Transform, seg [2]box2d.B2Vec2, screen *ebiten.Image) {
	p1 := TransformPoint(transformationPoint, seg[0])
	p2 := TransformPoint(transformationPoint, seg[1])

	x1 := float32(p1.X) * pixelsPerMeter
	y1 := float32(p1.Y) * pixelsPerMeter
	x2 := float32(p2.X) * pixelsPerMeter
	y2 := float32(p2.Y) * pixelsPerMeter

	vector.StrokeLine(screen, x1, y1, x2, y2, 2, color.RGBA{50, 200, 50, 255}, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth * zoomFactor, screenHeight * zoomFactor
}

func main() {
	ebiten.SetWindowSize(screenWidth*zoomFactor*1.5, screenHeight*zoomFactor)
	ebiten.SetWindowTitle("Rotating Cup with Ball - Box2D + Ebiten")
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func SetupTankFromRect(world *box2d.B2World, tankRect image.Rectangle, pixelsPerMeter float32) (*box2d.B2Body, *WaterSimulation, [][2]box2d.B2Vec2) {
	// Convert pixel coordinates to world coordinates
	leftWorld := float64(tankRect.Min.X) / float64(pixelsPerMeter)
	rightWorld := float64(tankRect.Max.X) / float64(pixelsPerMeter)
	topWorld := float64(tankRect.Min.Y) / float64(pixelsPerMeter)
	bottomWorld := float64(tankRect.Max.Y) / float64(pixelsPerMeter)

	// Create tank body
	tankDef := box2d.MakeB2BodyDef()
	tankDef.Type = box2d.B2BodyType.B2_staticBody
	tankDef.Position.Set((leftWorld+rightWorld)/2, (topWorld+bottomWorld)/2) // Center of tank
	tank := world.CreateBody(&tankDef)

	// Define tank walls relative to tank center
	tankWidth := rightWorld - leftWorld
	tankHeight := bottomWorld - topWorld
	halfWidth := tankWidth / 2
	halfHeight := tankHeight / 2

	tankLines := [][2]box2d.B2Vec2{
		// Bottom wall
		{box2d.MakeB2Vec2(-halfWidth, halfHeight), box2d.MakeB2Vec2(halfWidth, halfHeight)},
		// Left wall
		{box2d.MakeB2Vec2(-halfWidth, -halfHeight), box2d.MakeB2Vec2(-halfWidth, halfHeight)},
		// Right wall
		{box2d.MakeB2Vec2(halfWidth, -halfHeight), box2d.MakeB2Vec2(halfWidth, halfHeight)},
	}

	// Create tank fixtures
	for _, seg := range tankLines {
		edge := box2d.MakeB2EdgeShape()
		edge.Set(seg[0], seg[1])
		fixtureDef := box2d.MakeB2FixtureDef()
		fixtureDef.Shape = &edge
		fixtureDef.Friction = 0.6
		fixtureDef.Restitution = 0.0
		tank.CreateFixtureFromDef(&fixtureDef)
	}

	// Set water level to be at the top of the tank (where water surface would be)
	waterLevel := topWorld

	// Create water simulation
	water := NewWaterSimulation(
		waterLevel, // Water starts at tank top
		0.042,      // Buoyancy
		0.003,      // Drag
		tankRect,
		pixelsPerMeter,
	)

	return tank, water, tankLines
}

// Helper function to draw the water surface line
func DrawWaterSurface(screen *ebiten.Image, tankRect image.Rectangle, waterLevel float64, pixelsPerMeter float32) {
	waterY := float32(waterLevel * float64(pixelsPerMeter))

	x1 := float32(tankRect.Min.X)
	x2 := float32(tankRect.Max.X)

	// Draw water surface
	vector.StrokeLine(screen, x1, waterY, x2, waterY, 2, color.RGBA{100, 150, 255, 255}, true)

	// Optional: Fill water area with semi-transparent blue
	waterHeight := float32(tankRect.Max.Y) - waterY
	if waterHeight > 0 {
		vector.DrawFilledRect(screen, x1, waterY, x2-x1, waterHeight, color.RGBA{100, 150, 255, 100}, false)
	}
}
