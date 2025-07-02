package interactableUIObjects

import (
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PiggyBankSprite struct {
	*UiSprite
	amountAvailable    float32
	AnimationMap       map[string]drawables.Drawable
	triggeredAnimation string
}

func (pb *PiggyBankSprite) Update() {
	aniSprite, ok := pb.AnimationMap[pb.triggeredAnimation].(*sprite.AnimatedSprite)
	if ok {
		if aniSprite.LastF == aniSprite.Frame() {
			aniSprite.Reset()
			pb.triggeredAnimation = ""
		}
	}

	if pb.triggeredAnimation != "" {
		pb.AnimationMap[pb.triggeredAnimation].Update()
	}

	if pb.state == Clickable {
		if pb.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			turnOffClickMeEffect(pb.UiSprite)
			ev := events.MoneyAdded{Amount: pb.amountAvailable}
			pb.EventHub.Publish(ev)
			pb.triggeredAnimation = "allowance"
			pb.Shader = nil
			pb.state = Idle
		}
	}

	pb.Sprite.Update()
	pb.stateWas = pb.state
}

func (pb *PiggyBankSprite) Draw(screen *ebiten.Image) {
	if pb.triggeredAnimation != "" {
		pb.AnimationMap[pb.triggeredAnimation].Draw(screen)
		return
	}

	pb.Sprite.Draw(screen)
}

func (pb *PiggyBankSprite) Subscribe(hub *tasks.EventHub) {

	hub.Subscribe(events.MoneyAvailable{}, func(e tasks.Event) {
		ev := e.(events.MoneyAvailable)
		initClickMeEffect(pb.UiSprite)
		pb.state = Clickable
		pb.amountAvailable = ev.Amount
	})
}
