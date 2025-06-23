package interactableUIObjects

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PillowUI struct {
	*UiSprite
	Triggered bool
}

func (p *PillowUI) Update() {

	if p.SpriteHovered() {
		//shader := shaders.LoadOutlineShader()
		//p.Shader = shader
		//p.ShaderParams["OutlineColor"] = []float64{1, 1, 0, 1}
	}

	if p.Shader != nil && !p.SpriteHovered() {
		//p.Shader = nil
	}

	if p.Triggered && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && p.SpriteHovered() {
		p.EventHub.Publish(events.UISpriteAction{UiSprite: "pillow", UiSpriteAction: "clicked"})
		turnOffClickMeEffect(p.UiSprite)
	}

	p.Sprite.Update()

}

func (p *PillowUI) Draw(screen *ebiten.Image) {
	if p.Triggered == true {
		p.Sprite.Draw(screen)
	}
}

func (p *PillowUI) Subscribe(hub *tasks.EventHub) {
	hub.Subscribe(tasks.AllTasksCompleted{}, func(e tasks.Event) {
		initClickMeEffect(p.UiSprite)
		p.Triggered = true
	})

	hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		p.Triggered = false
	})
}
