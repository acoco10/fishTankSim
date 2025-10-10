package entities

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math/rand"
)

func (wb *WhiteBoardSprite) AddCrossOutGraphicEntity(yCoord, xCoord float32, graphicBounds *image.Rectangle, shaderParams map[string]any, drawOptParams map[string]float64) {

	wb.UiSprite.EventHub.Publish(events.WritingToWhiteBoard{
		Msg: CrossOut,
	})

	ySelectOp := rand.Intn(len(wb.YSelectOptionsForCrossOut) - 1)
	ySelect := wb.YSelectOptionsForCrossOut[ySelectOp]
	lowerY := ySelect * 16
	upperY := (ySelect + 1) * 16
	SubRect := image.Rect(0, lowerY, wb.crossOutTexture.Bounds().Dx(), upperY)
	if graphicBounds != nil {
		SubRect = *graphicBounds
	}
	singleCross := wb.crossOutTexture.SubImage(SubRect).(*ebiten.Image)
	//localCoords!!
	crossOutSprite := &sprite.Sprite{Img: singleCross, X: xCoord, Y: yCoord}
	crossOutSprite.Shader = registry.ShaderMap["HandWriting"]
	crossOutSprite.BufferDst = wb.DstImg.Img
	crossOutSprite.ShaderParams = shaderParams
	crossOutSprite.DOptsUpdaterParams = drawOptParams
	crossOutSprite.UpdateShaderParams = shaders.UpdateCounterOneShot

	id := RegisterEntity(&Entity{Sprite: crossOutSprite})
	wb.crossOutGraphicEntityID[id] = struct{}{}

	wb.YSelectOptionsForCrossOut = append(wb.YSelectOptionsForCrossOut[:ySelectOp], wb.YSelectOptionsForCrossOut[ySelectOp+1:]...)
	if len(wb.YSelectOptionsForCrossOut) == 1 {
		wb.YSelectOptionsForCrossOut = []int{0, 1, 2, 3}
	}
}

func (wb *WhiteBoardSprite) processQueuedWriteEvent() {
	if len(wb.taskCreatedEventQueue) >= 1 {
		wb.WriteTextToOpenSlot(wb.taskCreatedEventQueue[0].Task.Text)
		wb.taskCreatedEventQueue = wb.taskCreatedEventQueue[1:]
	} else {
		wb.taskCreatedEventQueue = wb.taskCreatedEventQueue[:0]
	}
}

func (wb *WhiteBoardSprite) processQueuedWriteRequest() {
	if len(wb.writeToWhiteBoardQueue) >= 1 {
		firstInQueue := wb.writeToWhiteBoardQueue[0]
		wb.appendTextToBestSpot(wb.writeToWhiteBoardQueue[0])
		if firstInQueue.EventToPublish != nil {
			for _, ev := range wb.writeToWhiteBoardQueue[0].EventToPublish {
				wb.UiSprite.EventHub.Publish(ev)
			}
		}
		if firstInQueue.EraseAfterFlag {
			wb.eraseRequest = &EraseRequest{time: DayOver}
		}

		wb.writeToWhiteBoardQueue = wb.writeToWhiteBoardQueue[1:]
	} else {
		wb.writeToWhiteBoardQueue = wb.writeToWhiteBoardQueue[:0]
	}
}

func defaultCrossOutShaderParams() map[string]any {
	ShaderParams := make(map[string]any)
	ShaderParams["Speed"] = 15 + rand.Float64()*10
	ShaderParams["MaxCounter"] = 40
	ShaderParams["Counter"] = 0
	maxOp := rand.Float64()*0.1 + 0.9
	ShaderParams["MaxOpacity"] = maxOp
	return ShaderParams
}

func defaultCrossOutDrawOptParams() map[string]float64 {
	DOptsUpdaterParams := make(map[string]float64)
	DOptsUpdaterParams["degree"] = rand.Float64() * -0.01
	return DOptsUpdaterParams
}

func CheckDoneDrawing(entID uint32) bool {
	tEnt, exists := GetEntity(entID)
	if !exists {
		return true
	}

	if tEnt.ShaderTextGraphic != nil && !tEnt.ShaderTextGraphic.FullyDrawn {
		return false
	}

	if tEnt.Sprite != nil {
		return sprite.CheckSpriteWithShaderCounterFinished(tEnt.Sprite)
	}

	return true
}

func (wb *WhiteBoardSprite) underLineHeader() {
	rect := image.Rect(60, 32, 200, 48)
	sParams := defaultCrossOutShaderParams()
	sParams["Speed"] = 5.0
	dParams := defaultCrossOutDrawOptParams()
	wb.AddCrossOutGraphicEntity(20, 74, &rect, sParams, dParams)
}
