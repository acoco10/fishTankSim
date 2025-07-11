package entity

import (
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image/color"
)

type Entity struct {
	Sprite *sprite.AnimatedSprite
	Active bool
}

type PlantNode struct {
	Next     *PlantNode
	Previous *PlantNode
	X, Y     float32
	X1, Y1   float32
	baseX    float32
	Pushed   bool
	Count    int
}

type Plant struct {
	Next          *PlantNode
	Height, Width float32
}

type Force struct {
	Magnitude float32
	Direction string
}

func UpdateSkeletalNode(p *PlantNode) {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		_, y := ebiten.CursorPosition()
		//xf := float32(x)
		yf := float32(y)

		if p.Y1 < yf && p.Y > yf {
			p.Pushed = true
			p.Count = 1
		}
	}

	if p.Pushed && p.Count == 1 {
		p.Pushed = false
		p.X += 10
		p.X1 += 10
	}

	if p.Count >= 1 {
		p.Count++
		lerpFactor := float32(0.25)
		x := Lerp(p.baseX, p.X, lerpFactor)

		p.X += x
		p.X1 += x
	}

	if p.Count == 10 {
		p.Count = 0
	}

	if p.Previous != nil {
		p.Previous.X1 = p.X1
		p.Previous.X = p.Previous.X - (p.Previous.X-p.X)/2
	}

	if p.Next != nil {
		p.Next.X = p.X
		p.Next.X1 = p.Next.X - (p.Next.X-p.X)/2
	}

}

func Lerp(A float32, B float32, t float32) float32 {
	return A + (B-A)*t
}

func MakePlant(size int) *Plant {
	// Root node (base of plant)
	root := &PlantNode{
		X: 600, Y: 600,
		X1: 600, Y1: 550,
		baseX: 600, // Initialize baseX
	}

	plant := &Plant{Width: 10, Height: 50, Next: root}

	current := root
	for i := 0; i < size; i++ {
		newNode := &PlantNode{
			Next:     nil,
			Previous: current,
			X:        current.X1,
			Y:        current.Y1,
			X1:       current.X1,
			Y1:       current.Y1 - plant.Height, // Going upward (or + for downward)
			baseX:    current.X1,                // Set baseX for lerping back
		}
		current.Next = newNode
		current = newNode
	}

	return plant
}

func (p *Plant) Update() {
	head := p.Next.Next

	for head.Next != nil {
		UpdateSkeletalNode(head)
		head = head.Next
	}

}

func (p *Plant) Draw(screen *ebiten.Image) {

	head := p.Next

	colorArr := []color.RGBA{
		colornames.Green,
		colornames.Darkolivegreen,
		colornames.Darkseagreen,
		colornames.Red,
	}

	colorI := 0

	for head.Next != nil {

		if head.Pushed {
			vector.StrokeLine(screen, head.X, head.Y, head.X1, head.Y1, 20, colorArr[3], false)
		} else {
			vector.StrokeLine(screen, head.X, head.Y, head.X1, head.Y1, 20, colorArr[colorI], false)
		}

		head = head.Next
		colorI++

		if colorI == 3 {
			colorI = 0
		}
	}

}

func (e *Entity) Update() {

}
