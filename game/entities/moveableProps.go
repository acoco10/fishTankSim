package entities

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"log"
)

type PropState uint8

const (
	Moveable PropState = iota
	SettingInPlace
	SetInPlace
)

type Prop interface {
	Draw(screen *ebiten.Image)
	Update()
}

type PropQueue struct {
	ActiveProp  *StructureProp
	DrawProps   []*StructureProp
	QueuedProps []*StructureProp
}

type StructureProp struct {
	state    PropState
	stateWas PropState
	*sprite.Sprite
	shadowPoint  image.Point
	boundaries   image.Rectangle
	StaticShadow bool
	baseY        float32
}

func (p *StructureProp) State() PropState {
	return p.state
}

func NewStructureProp(x float32, y float32, img *ebiten.Image, normal *ebiten.Image, hub *tasks.EventHub, bounds image.Rectangle) *StructureProp {

	p := StructureProp{}

	sp := &sprite.Sprite{Img: img, NormalMap: normal, X: x, Y: y, Z: 0}

	if normal != nil {
		normalMapShader := registry.ShaderMap["NormalMap"]
		sp.Shader = normalMapShader
		sp.ShaderParams = make(map[string]any)
		sp.ShaderParams["Cursor"] = []float64{440, 600}
	} else {
		p.Sprite.ShaderParams = make(map[string]any)
	}

	subscribe(&p, hub)
	p.Sprite = sp
	p.state = Moveable
	sprite.LoadPulseOutlineShader(p.Sprite)
	p.shadowPoint = image.Point{X: int(x), Y: int(y)}
	p.boundaries = bounds

	return &p
}

func (p *StructureProp) Draw(screen *ebiten.Image) {

	baseOffset := float32(10.0)
	if p.state == SetInPlace {
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

	if p.state == SettingInPlace {

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

	if p.state == Moveable {
		p.Y = float32(p.boundaries.Max.Y-p.Img.Bounds().Dy()) - 50
		//y1 = max of fishtank y - height of the image - 50
		//fishtank y - 50 = "ceiling"
		//fishtank y - 30 = "floor"
		x, _ := util.GetScaledCursorPosition()
		if p.X > float32(p.boundaries.Min.X-p.Img.Bounds().Max.X)-10 {
			p.X = min(float32(x), float32(p.boundaries.Max.X-p.Img.Bounds().Max.X)-10)
		}
		if p.X < float32(p.boundaries.Min.X)+20 {
			p.X = max(float32(x), float32(p.boundaries.Min.X)+20)
		}

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			p.state = SettingInPlace
		}
	}
	if p.state == SettingInPlace {
		p.Sprite.UnLoadShader()
		p.Y++
	}
	if p.Y >= float32(p.boundaries.Max.Y-p.Img.Bounds().Dy())-29 {
		p.state = SetInPlace
	}

	p.Sprite.Update()
	p.stateWas = p.state

}

func subscribe(P *StructureProp, hub *tasks.EventHub) {
	hub.Subscribe(events.FishTankLayout{}, func(e tasks.Event) {
		ev := e.(events.FishTankLayout)
		println("recebed fish tank BOUNDARIES =", ev.Rectangle.Min.X, ev.Rectangle.Max.X)
		P.boundaries = ev.Rectangle
	})
}

func UpdateProps(queue PropQueue) {
	//everything in  draw queue will be drawn
	if queue.ActiveProp == nil {
		return
	}

	queue.ActiveProp.Update()

	if len(queue.DrawProps) == 0 {
		//add first active prop to be draw
		log.Println("appending first active prop to draw prop slice")
		queue.DrawProps = append(queue.DrawProps, queue.ActiveProp)
	}

	if queue.ActiveProp.state == SetInPlace {
		if len(queue.QueuedProps) > 0 {
			queue.ActiveProp = queue.QueuedProps[0]
			queue.QueuedProps = queue.QueuedProps[1:]
			//add secondary props when they reach the front of the queue
			queue.DrawProps = append(queue.DrawProps, queue.ActiveProp)
		}
	}
}

func DrawProps(queue PropQueue, screen *ebiten.Image) {
	for _, prop := range queue.DrawProps {
		prop.Draw(screen)
	}
}

func LoadProp(propName string, tankBoundaries image.Rectangle, eventhub *tasks.EventHub) *StructureProp {
	var prop *StructureProp
	switch propName {
	case "Log":
		logPropImg, err := util.LoadImageAssetAsEbitenImage("tankProps/logProp")
		logNormal, err := util.LoadImageAssetAsEbitenImage("tankProps/logProp_n")
		logProp := NewStructureProp(0, 0, logPropImg, logNormal, eventhub, tankBoundaries)

		if err != nil {
			log.Fatal(err)
		}
		prop = logProp
	case "Castle":

		log.Println("returning castle prop from load prop call")
		castleImg, err := util.LoadImageAssetAsEbitenImage("tankProps/castleProp")
		castleNormal, err := util.LoadImageAssetAsEbitenImage("tankProps/castleProp_n")
		if err != nil {
			log.Fatal(err)
		}
		castleProp := NewStructureProp(0, 0, castleImg, castleNormal, eventhub, tankBoundaries)
		prop = castleProp

	case "Grass":
		grassImg, err := util.LoadImageAssetAsEbitenImage("tankProps/grass")
		if err != nil {
			log.Fatal(err, "tried to load plant img from wrong place")
		}

		grassNormal, err := util.LoadImageAssetAsEbitenImage("tankProps/grass_n")
		if err != nil {
			log.Fatal(err)
		}

		grassProp := NewStructureProp(0, 0, grassImg, grassNormal, eventhub, tankBoundaries)
		prop = grassProp

	case "Rock":
		//rockImg, err := util.LoadImageAssetAsEbitenImage("uiSprites/redRock.png")
	}

	prop.Sprite.AbleToBeUnfocusedAutomatically = true

	return prop
}
