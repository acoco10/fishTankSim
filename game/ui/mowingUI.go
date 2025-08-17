package ui

import (
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

func LoadMowingUI(eHub *tasks.EventHub) *ebitenui.UI {
	rootContainer := widget.NewContainer(
		//widget.ContainerOpts.BackgroundImage(nineSliceImage),
		// the container will use a plain color as its background
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	return &ebitenui.UI{Container: rootContainer}
}

func MowingUiSubs(hub *tasks.EventHub) {

}
