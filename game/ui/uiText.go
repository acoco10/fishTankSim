package ui

import (
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
	"log"
)

func MakeOutlineText(txt string) *widget.Container {
	face, err := util.LoadFont(32, "reglisseOutline") //white center text
	if err != nil {
		log.Fatal(err)
	}
	face2, err := util.LoadFont(32, "reglisseOutlined") //black outline
	if err != nil {
		log.Fatal(err)
	}
	headerLbl := widget.NewText(
		widget.TextOpts.Text(txt, face, colornames.Aliceblue),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		widget.TextOpts.Insets(widget.Insets{Top: 0, Left: 2}))

	headerLblOutline := widget.NewText(
		widget.TextOpts.Text(txt, face2, colornames.Black),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		widget.TextOpts.Insets(widget.Insets{Top: 0, Left: 2}),
	)
	headerContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	headerContainer.AddChild(headerLbl)
	headerContainer.AddChild(headerLblOutline)

	return headerContainer

}
