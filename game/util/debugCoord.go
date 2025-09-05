package util

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image"
	"log"
	"math"
	"os"
)

type DebugCoord struct {
	pts      []image.Point
	savedPts []image.Point
	RectState
	Tag                string
	cursorMarkerRadius float32
	relatedObjectPoint image.Point
	coordStateMachine
	eventHub    *tasks.EventHub
	drawMeDaddy bool
}

type coordStateMachine struct {
	States       map[RectState]coordStateHandler
	CurrentState RectState
}

func (s *coordStateMachine) Transition(coord *DebugCoord) {
	if s.States[s.CurrentState].TransitionFunc != nil {
		s.States[s.CurrentState].TransitionFunc(coord)
	}
	s.CurrentState = s.States[s.CurrentState].TransitionTo
}

type coordStateHandler struct {
	Updater        func(entity *DebugCoord)
	TransitionTo   RectState
	TransitionFunc func(rect *DebugCoord)
}

func (c *DebugCoord) GivePoint(pt image.Point) {
	c.relatedObjectPoint = pt
}

func (c *DebugCoord) Init(Tag string, hub *tasks.EventHub) {

	state1 := coordStateHandler{Updater: updateCoordInitState, TransitionTo: Initiated, TransitionFunc: coordTransitionFromIntToDraw}
	state2 := coordStateHandler{Updater: updateCoordInitiated, TransitionTo: Drawn}
	state3 := coordStateHandler{Updater: updateCoordDrawn, TransitionTo: On}

	states := map[RectState]coordStateHandler{
		On:        state1,
		Initiated: state2,
		Drawn:     state3,
	}

	c.Tag = Tag
	sm := coordStateMachine{States: states}
	c.eventHub = hub
	c.coordStateMachine = sm
	c.CurrentState = On

}

func updateCoordInitState(c *DebugCoord) {
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		//point doesnt need a ui for entering custom title as we only edit one point set for each prop
		//an prop title is already defined
		c.drawMeDaddy = true
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && c.drawMeDaddy {
		c.Transition(c)
	}
}

func coordTransitionFromIntToDraw(c *DebugCoord) {
	x, y := GetScaledCursorPosition()
	Pt := image.Point{X: x, Y: y}
	c.pts = append(c.pts, Pt)

}

func updateCoordInitiated(c *DebugCoord) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		c.Transition(c)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		c.Transition(c)
	}
}

func updateCoordDrawn(c *DebugCoord) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		c.Transition(c)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := GetScaledCursorPosition()
		c.AddPoint(image.Point{x, y})
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		c.pts = c.pts[:len(c.pts)-1]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		err := c.Save()
		if err != nil {
			println("cannot save coords")
			log.Fatal(err)
		}
		c.Transition(c)
	}
}

func (c *DebugCoord) Draw(screen *ebiten.Image) {

	if c.CurrentState == Initiated || c.CurrentState == Drawn {
		for _, pt := range c.pts {
			DrawCircleFromPoint(pt, screen)
		}
	}

}

func DrawCircleFromPoint(pt image.Point, screen *ebiten.Image) {
	vector.StrokeCircle(screen, float32(pt.X-2), float32(pt.Y-2), 2, 1, colornames.Darkmagenta, false)
}

func (c *DebugCoord) AddPoint(point image.Point) {
	c.pts = append(c.pts, point)
}

func (c *DebugCoord) Update() error {
	c.coordStateMachine.States[c.CurrentState].Updater(c)
	return nil
}

func (c *DebugCoord) Save() error {

	minX := c.relatedObjectPoint.X
	minY := c.relatedObjectPoint.Y

	for i, _ := range c.pts {
		c.pts[i].X -= minX
		c.pts[i].Y -= minY
	}

	SortRectanglePoints(c.pts)

	existingPos, err := LoadCoords()
	if err != nil {
		return err
	}

	datMap := existingPos // Start with existing data
	datMap[c.Tag] = c.pts

	outputSave, err := json.Marshal(datMap)
	if err != nil {
		return err
	}

	println("saving point data with following Coordinates:", string(outputSave))

	err = os.WriteFile("assets/data/structureCoords.json", outputSave, 999)
	if err != nil {
		return err
	}
	return nil
}

func LoadCoords() (map[string][]image.Point, error) {
	colDat, err := assets.DataDir.ReadFile("data/structureCoords.json")
	if err != nil {
		log.Println("Error opening structure coord file:", err, "Creating a new one")
		file, err2 := os.Create("assets/data/structureCoords.json")
		if err2 != nil {
			fmt.Println("Error creating file:", err2)
			return nil, err2
		}
		defer file.Close()
		return map[string][]image.Point{}, nil
	}

	datMap := make(map[string][]image.Point)

	err = json.Unmarshal(colDat, &datMap)

	if err != nil {
		return map[string][]image.Point{}, err
	}

	return datMap, nil
}

func SortRectanglePoints(pts []image.Point) {
	if len(pts) != 4 {
		return
	}

	// Find the actual min/max coordinates
	minX, minY := pts[0].X, pts[0].Y
	maxX, maxY := pts[0].X, pts[0].Y

	for _, pt := range pts {
		if pt.X < minX {
			minX = pt.X
		}
		if pt.X > maxX {
			maxX = pt.X
		}
		if pt.Y < minY {
			minY = pt.Y
		}
		if pt.Y > maxY {
			maxY = pt.Y
		}
	}

	fmt.Printf("Bounds: minX=%d, maxX=%d, minY=%d, maxY=%d\n", minX, maxX, minY, maxY)

	// Find closest point to each corner
	sorted := make([]image.Point, 4)
	used := make([]bool, 4)

	corners := []image.Point{
		{minX, minY}, // Top-left target
		{maxX, minY}, // Top-right target
		{maxX, maxY}, // Bottom-right target
		{minX, maxY}, // Bottom-left target
	}

	for cornerIdx, target := range corners {
		bestDist := math.MaxFloat64
		bestIdx := -1

		for i, pt := range pts {
			if used[i] {
				continue
			}

			// Calculate distance to target corner
			dx := float64(pt.X - target.X)
			dy := float64(pt.Y - target.Y)
			dist := dx*dx + dy*dy

			if dist < bestDist {
				bestDist = dist
				bestIdx = i
			}
		}

		if bestIdx >= 0 {
			sorted[cornerIdx] = pts[bestIdx]
			used[bestIdx] = true
			fmt.Printf("Corner %d: found %v (dist=%.2f)\n", cornerIdx, pts[bestIdx], math.Sqrt(bestDist))
		}
	}

	copy(pts, sorted)
}
