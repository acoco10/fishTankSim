package debug

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"image"
	"log"
)

type DebugOption uint8

const (
	Normal DebugOption = iota
	Position
	Print
	ShaderTest
)

type GameMode uint8

const (
	User GameMode = iota
	Debug
)

type DebugData struct {
	DebugText        string
	DebugRect        *util.Rect
	DebugParameter   map[DebugOption]bool
	GameMode         GameMode
	DebugPoint       *util.DebugCoord
	EventHub         *tasks.EventHub
	drawingSomething bool
	showZbounds      bool
	GameState        entities.GameState
	PropData         entities.PropData
	InitiatedProp    *entities.StructureProp
	DebugOption
}

func (d *DebugData) Update() {

	if d.DebugRect != nil {
		err := d.DebugRect.Update()
		if err != nil {
			log.Fatal("Couldnt update debug rect:", err)
		}
	}

	if d.DebugPoint != nil {
		err := d.DebugPoint.Update()
		if err != nil {
			log.Fatal("Couldnt update debug point:", err)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		d.DebugRect = nil
		d.DebugPoint = nil
		d.drawingSomething = false

	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		d.showZbounds = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		d.GameMode = User
		registry.Config.Set(registry.Debug, false)

	}

	if !d.drawingSomething {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			d.DebugText += "Mode: Drawing Rectangle\n"
			d.DebugRect = &util.Rect{}
			d.DebugRect.Init("new", d.EventHub)
			if d.InitiatedProp != nil {
				d.DebugRect.Tag = d.InitiatedProp.Tag
			}
			d.drawingSomething = true
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyP) && !ebiten.IsKeyPressed(ebiten.KeyShift) {
			d.DebugText += "Mode: Drawing Point\n"
			d.DebugPoint = &util.DebugCoord{}
			d.DebugPoint.Init("new", d.EventHub)
			if d.InitiatedProp != nil {
				d.DebugPoint.Tag = d.InitiatedProp.Tag
			}
			d.drawingSomething = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyShift) && inpututil.IsKeyJustPressed(ebiten.KeyP) {
			d.DebugText += "Mode: Editing UI Sprite position\n"
			for _, ent := range entities.LiveList {
				if ent.UiData != nil {
					ent.StateMachine = entities.InitStateMachine(entities.PositionUpdate, entities.AddUiSpriteXYUpdater, entities.UiSpriteTurnOffEverything)
					ent.Sprite.Unfocusable = false
					ent.SetUIState(entities.Idle)
				}
			}
			d.DebugOption = Position
			d.drawingSomething = true
		}

	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		d.DebugText += "Current Initiated Prop: Castle\n"
		id := entities.LoadProp("Castle", d.PropData, d.EventHub, entities.PlacementPicked{}, d.GameState.Zbounds)
		prop, _ := entities.GetEntity(id)
		d.InitiatedProp = prop.PropData

	}

	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		d.DebugText += "Current Initiated Prop: Log\n"
		id := entities.LoadProp("Log", d.PropData, d.EventHub, entities.PlacementPicked{}, d.GameState.Zbounds)
		prop, _ := entities.GetEntity(id)
		d.InitiatedProp = prop.PropData
	}

	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		d.DebugText += "Current Initiated Buddha: Log\n"
		id := entities.LoadProp("Buddha", d.PropData, d.EventHub, entities.PlacementPicked{}, d.GameState.Zbounds)
		prop, _ := entities.GetEntity(id)
		d.InitiatedProp = prop.PropData
	}

	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		d.DebugText += "Current Initiated Bridge: Log\n"
		id := entities.LoadProp("zenBridge", d.PropData, d.EventHub, entities.PlacementPicked{}, d.GameState.Zbounds)
		prop, _ := entities.GetEntity(id)
		d.InitiatedProp = prop.PropData
	}

	if d.InitiatedProp != nil && d.InitiatedProp.State() == entities.SetInPlace {
		if d.DebugRect != nil {
			d.DebugRect.GivePoint(image.Point{X: int(d.InitiatedProp.X), Y: int(d.InitiatedProp.Y)})
		}
		if d.DebugPoint != nil {
			d.DebugPoint.GivePoint(image.Point{X: int(d.InitiatedProp.X), Y: int(d.InitiatedProp.Y)})
		}
	}

}

func (d *DebugData) Draw(screen *ebiten.Image) {

	if d.DebugRect != nil {
		d.DebugRect.Draw(screen)
	}
	if d.DebugPoint != nil {
		d.DebugPoint.Draw(screen)
	}

	if d.showZbounds {
		color := colornames.Orangered
		for _, rect := range d.GameState.Zbounds {
			util.StrokeRectFromImageRect(rect, screen, color, false)
		}

	}

	ebitenutil.DebugPrintAt(screen, d.DebugText, registry.Config.ScreenWidth/2, 100)

}
