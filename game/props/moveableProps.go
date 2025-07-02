package props

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
)

type State uint8

const (
	Moveable State = iota
	SetInPlace
)

type Prop interface {
	Draw(screen *ebiten.Image)
	Update()
}

type StructureProp struct {
	state    State
	stateWas State
	*sprite.Sprite
	shadowPoint  image.Point
	boundaries   image.Rectangle
	StaticShadow bool
	baseY        float32
}

func NewStructureProp(x float32, y float32, img *ebiten.Image, normal *ebiten.Image, hub *tasks.EventHub) *StructureProp {

	p := StructureProp{}

	sp := &sprite.Sprite{Img: img, NormalMap: normal, X: x, Y: y}

	if normal != nil {
		normalMapShader := registry.ShaderMap["NormalMap"]
		sp.Shader = normalMapShader
		sp.ShaderParams = make(map[string]any)
		sp.ShaderParams["Cursor"] = []float64{440, 600}
	}

	subscribe(&p, hub)
	p.Sprite = sp
	p.state = Moveable
	sprite.LoadPulseOutlineShader(p.Sprite)
	p.shadowPoint = image.Point{X: int(x), Y: int(y)}

	return &p
}

func (p *StructureProp) Draw(screen *ebiten.Image) {
	baseOffset := float32(10.0)
	if p.state == SetInPlace && p.Y >= float32(p.boundaries.Max.Y-p.Img.Bounds().Dy())-30 && p.StaticShadow {
		//static shadow
		vector.StrokeRect(screen, p.X+20, p.Y+float32(p.Img.Bounds().Dy())-6, float32(p.Img.Bounds().Dx())-60, 2, 4, color.RGBA{0, 0, 0, 100}, false)
	}

	p.Sprite.Draw(screen)
	if p.state == Moveable {
		x := p.X + baseOffset
		y := float32(p.boundaries.Max.Y - 35)
		height := float32(2)
		width := float32(p.Img.Bounds().Dx()) - 2*baseOffset
		vector.StrokeRect(screen, x, y, width, height, 4, color.RGBA{0, 0, 0, 100}, false)
	}

	if p.state == SetInPlace && p.Y < float32(p.boundaries.Max.Y-p.Img.Bounds().Dy())-35 {

		dist := float32(p.boundaries.Max.Y-p.Img.Bounds().Dy()) - p.Y
		//fishtank base - image height = base comparable to image y
		//base-y = distance between height and current y
		//positive number between 50 and 30
		dist = 50 - dist

		x := p.X + dist + baseOffset
		//increase offset from origin
		y := float32(p.boundaries.Max.Y - 35)
		height := float32(2)

		if dist < 35 {
			y += 1
			height = 1.0
		}
		width := float32(p.Img.Bounds().Dx()) - 2*dist - 2*baseOffset
		vector.StrokeRect(screen, x, y, width, height, 4, color.RGBA{0, 0, 0, 100}, false)
	}

}

func (p *StructureProp) Update() {
	p.Sprite.Update()
	if p.state == Moveable {
		p.Y = float32(p.boundaries.Max.Y-p.Img.Bounds().Dy()) - 50
		//y1 = max of fishtank y - height of the image - 50
		//fishtank y - 50 = "ceiling"
		//fishtank y - 30 = "floor"
		x, _ := ebiten.CursorPosition()
		if p.X > float32(p.boundaries.Min.X-p.Img.Bounds().Max.X)-10 {
			p.X = min(float32(x), float32(p.boundaries.Max.X-p.Img.Bounds().Max.X)-10)
		}
		if p.X < float32(p.boundaries.Min.X)+20 {
			p.X = max(float32(x), float32(p.boundaries.Min.X)+20)
		}

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			p.state = SetInPlace
		}
	}
	if p.state == SetInPlace && p.Y < float32(p.boundaries.Max.Y-p.Img.Bounds().Dy())-27 {
		p.Sprite.UnLoadShader()
		p.Y++
	}

	if p.state == SetInPlace && p.SpriteHovered() && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && p.stateWas == SetInPlace {

	}

	p.stateWas = p.state

}

func subscribe(P *StructureProp, hub *tasks.EventHub) {
	hub.Subscribe(events.FishTankLayout{}, func(e tasks.Event) {

		ev := e.(events.FishTankLayout)
		println("recebed fish tank BOUNDARIES =", ev.Rectangle.Min.X, ev.Rectangle.Max.X)
		P.boundaries = ev.Rectangle
	})
}

func GamePropToSaveProp(prop *StructureProp, name string) entities.TankObject {
	save := entities.TankObject{}
	save.X = prop.X
	save.Y = prop.Y
	save.Name = name
	return save
}
