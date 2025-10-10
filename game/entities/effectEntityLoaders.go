package entities

import (
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"log"
	"strconv"
)

type EffectParams struct {
	Cycles int
	Scale  float64
	Speed  float32
	Zoom   bool
}

type TurnOffEffect func()

const (
	Smoke = iota
	Spotlight
	CircularProgress
	CursorPower
	ClickHere
	Space
	ArrowKeys
	EnterKey
)

var effects = [10]Entity{}

func LoadEffects(zBounds [13]image.Rectangle) {
	spotLight, err := util.LoadImageAssetAsEbitenImage("staticEffects/spotLight")
	if err != nil {
		log.Fatal("cant load spotlight effect", err)
	}

	spotLightSprite := &sprite.Sprite{Img: spotLight, Y: float32(zBounds[0].Min.Y - 12)}

	ent := Entity{Sprite: spotLightSprite}
	ent.Z = 10
	effects[Spotlight] = ent

	smokeEffect := entImportableLoaders.LoadEffect("smokeEffect")
	sEnt := Entity{Sprite: smokeEffect}
	sEnt.Z = 10
	effects[Smoke] = sEnt

	circularProgressEffect := entImportableLoaders.LoadEffect("circularProgressEffect")
	cEnt := Entity{Sprite: circularProgressEffect}
	sEnt.Z = 15
	effects[CircularProgress] = cEnt

	img := loadCursorImages()[activePressed]
	sp := &sprite.Sprite{Img: img}
	sp.DOptsUpdaterParams = make(map[string]float64)
	sp.Scale = 0.25
	sp.DOptsUpdaterParams["opacity"] = 0.2

	cuEnt := Entity{Sprite: sp}
	cuEnt.Z = 10
	effects[CursorPower] = cuEnt

	clickEffect := entImportableLoaders.LoadEffect("mouseClickEffect")
	clickEnt := Entity{Sprite: clickEffect}
	clickEnt.Z = 15
	effects[ClickHere] = clickEnt

	spaceEffect := entImportableLoaders.LoadEffect("spaceEffect")
	spaceEnt := Entity{Sprite: spaceEffect}
	clickEnt.Z = 15
	effects[Space] = spaceEnt

	arrowKeys, err := util.LoadImageAssetAsEbitenImage("staticEffects/arrowKeys")
	if err != nil {
		log.Fatal("cant load spotlight effect", err)
	}

	akSprite := &sprite.Sprite{Img: arrowKeys}

	akent := Entity{Sprite: akSprite}
	ent.Z = 10
	effects[ArrowKeys] = akent

	enterKey, err := util.LoadImageAssetAsEbitenImage("staticEffects/enterKey")
	if err != nil {
		log.Fatal("cant load spotlight effect", err)
	}

	ekSprite := &sprite.Sprite{Img: enterKey}

	ekEnt := Entity{Sprite: ekSprite}
	ent.Z = 10
	effects[EnterKey] = ekEnt

}

func DrawSpotLight(x float32, params EffectParams) DeInitFunc {
	eff := effects[Spotlight]
	midPoint := eff.Sprite.GetSpriteRect().Dx() / 2
	eff.Sprite.X = x - float32(midPoint)
	RegisterEntity(&eff)
	return func() { RemoveEntity(eff.Id) }
}

func DrawSmoke(x, y float32, params EffectParams) DeInitFunc {
	eff := effects[Smoke]
	eff.Draw = true
	midX := float32(eff.Sprite.MidX())
	eff.Sprite.X = x - midX
	eff.Sprite.Y = y - float32(eff.Sprite.SpriteHeight())
	eff.DeposeAfterNAnimationCycles = params.Cycles
	RegisterEntity(&eff)
	return func() { RemoveEntity(eff.Id) }
}

func DrawCircularProgressEffect(x, y float32, params EffectParams) DeInitFunc {
	eff := effects[CircularProgress]
	eff.Draw = true
	quarterX := float32(eff.Sprite.MidX()) / 2
	eff.Sprite.X = x - quarterX - float32(eff.Sprite.GetSpriteRect().Dx())
	eff.Sprite.Y = y - float32(eff.Sprite.MidY())
	eff.Z = 15
	eff.DeposeAfterNAnimationCycles = params.Cycles
	eff.Sprite.ChangeAnimationSpeed(params.Speed)
	RegisterEntity(&eff)
	return func() { RemoveEntity(eff.Id) }
}

func DrawCursorPowerEffect(params EffectParams) (DeInitFunc, func(speed float64)) {
	eff := effects[CursorPower]
	eff.Draw = true
	x, y := util.GetScaledCursorPosition()
	width := float32(eff.Sprite.GetSpriteRect().Dx()) * float32(eff.Sprite.Scale)
	height := float32(eff.Sprite.GetSpriteRect().Dy()) * float32(eff.Sprite.Scale)

	eff.Sprite.X = float32(x) - width/2
	eff.Sprite.Y = float32(y) - height/2

	eff.Z = 15
	sprite.SetScaleUpdater(eff.Sprite, 2, 0.01, 0.25)
	RegisterEntity(&eff)
	speedUpdater := func(speed float64) {
		if eff.Sprite.DOptsUpdaterParams["Speed"] <= 0.5 {
			eff.Sprite.DOptsUpdaterParams["Speed"] += speed
		}
	}

	return func() { RemoveEntity(eff.Id) }, speedUpdater
}

func DrawControlEffect(x float64, y float64, params EffectParams, Button int, text string, updateFunc EntityUpdater) DeInitFunc {
	eff := effects[Button]
	eff.Draw = true
	eff.UpdateFunc = updateFunc

	eff.Sprite.X = float32(x) - float32(eff.Sprite.SpriteWidth()/2)
	eff.Sprite.Y = float32(y)

	var buffId uint32
	if text != "" {
		Xmid, _ := eff.Sprite.MidPointf()
		textX := Xmid
		textY := eff.Sprite.Y

		if params.Zoom {
			w, h := util.MeasureText(text, 18, "RockSalt_18")

			textX -= float32(w / 8)
			textY -= float32(h / 2)

			textX, textY = util.TranslateOverZoom(textX, textY)
		}

		buffId = MakeTextGraphicEntity(text, textX, textY, params.Zoom)
	}

	eff.Sprite.ChangeAnimationSpeed(15)
	eff.Z = 15
	eff.Parameters.EntIds[LinkedGraphic1] = buffId
	RegisterEntity(&eff)

	eff.labelForDebugging = "control effect graphic" + text + strconv.Itoa(int(eff.Id))

	return func() {
		RemoveEntity(eff.Id)
		RemoveEntity(buffId)
	}
}

func DrawRectEntity(rect image.Rectangle, color color.RGBA, filled bool) DeInitFunc {
	sp := &sprite.Sprite{Img: ebiten.NewImage(rect.Dx(), rect.Dy()), X: float32(rect.Min.X), Y: float32(rect.Min.Y), IsBuffer: true}
	g := &graphics.RectGraphic{
		Rectangle: image.Rect(0, 0, rect.Dx(), rect.Dy()),
		Color:     color,
		Filled:    filled,
		Dst:       sp.Img,
	}

	rectEnt := &Entity{RectGraphic: g, Sprite: sp}

	rectEnt.labelForDebugging = "rectangle entity graphic"
	RegisterEntity(rectEnt)

	return func() {
		RemoveEntity(rectEnt.Id)
	}
}

func MakeTextGraphicEntity(text string, x, y float32, overZoom bool) uint32 {
	//to do make this shit work
	//w, _ := util.MeasureText(text, 12, "RockSalt_12")
	//halfWidth := w / 2
	var ft *graphics.FadeInText

	buffEnt := &Entity{Sprite: &sprite.Sprite{
		Img: ebiten.NewImage(500, 150),
		X:   x,
		Y:   y}}
	buffEnt.Sprite.IsBuffer = true
	if overZoom {
		buffEnt.SetOverZoom()
		ft = graphics.NewFadeInTextReturnText(text, 120, registry.FontMap["RockSalt_18"])
	} else {
		ft = graphics.NewFadeInTextReturnText(text, 120, registry.FontMap["RockSalt_12"])
		buffEnt.Z = 15
	}

	buffEnt.graphicText = ft
	ft.Dst = buffEnt.Sprite.Img
	buffEnt.labelForDebugging = "text graphic for control effect"
	RegisterEntity(buffEnt)

	return buffEnt.Id
}

func ClearEnterGraphicOnClick(ent *Entity, state GameState) {
	ent.Parameters.Ints[IndexCounter]++
	if registry.ClickCheck() && ent.Parameters.Ints[IndexCounter] > 60 {
		RemoveEntity(ent.Id)
		RemoveEntity(ent.Parameters.EntIds[LinkedGraphic1])
	}
}

func StructurePlacementArrowKeysEnterChain(ent *Entity, gs GameState) {
	ent.Parameters.Ints[IndexCounter]++
	if ent.Parameters.Ints[IndexCounter] == 60 && registry.AnyArrowCheck() {
		DrawControlEffect(float64(ent.Sprite.X+100), float64(gs.Zbounds[12].Max.Y), EffectParams{Zoom: true}, EnterKey, "Confirm Position", ClearEnterGraphicOnClick)
	}

	if registry.ClickCheck() {
		RemoveEntity(ent.Id)
		RemoveEntity(ent.Parameters.EntIds[LinkedGraphic1])
	}
}

func ClearSpaceGraphic(ent *Entity, state GameState) {
	ent.Parameters.Ints[IndexCounter]++
	if registry.ZoomCheck() && ent.Parameters.Ints[IndexCounter] > 60 {
		RemoveEntity(ent.Id)
		RemoveEntity(ent.Parameters.EntIds[LinkedGraphic1])
	}
}
