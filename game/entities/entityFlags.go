package entities

type entFlags uint64

const (
	overUi entFlags = 1 << iota
	overZoom
	reset
	noOffset
	updater
	outline
	autoTransition1
	swirl
	altImage
	lowLight
	opacityFade
	revert
	likeWindow
	clickTransition
	draw
	freeze
	used
	unDraw
	clickForTime
	oneOff
	particlesGenerated
	keepShader
	fishCaught
	textOnSprite
	addClickEffect
	addIdleClickEffect
	noActivationRect
	dontDrawFirstDay
	UsingCounter
	DeInitAfterUsed
)

// Existing methods
func (ent *Entity) SetOverUIEffect() {
	ent.flags |= overUi
}
func (ent *Entity) SetOverZoom() {
	ent.flags |= overZoom
}
func (ent *Entity) IsOverUI() bool {
	return ent.flags&overUi != 0
}

func (ent *Entity) SetReset() {
	ent.flags |= reset
}

func (ent *Entity) IsReset() bool {
	return ent.flags&reset != 0
}

// Setter methods for UI flags
func (ent *Entity) SetNoOffset() {
	ent.flags |= noOffset
}

func (ent *Entity) SetUpdater() {
	ent.flags |= updater
}

func (ent *Entity) SetOutline() {
	ent.flags |= outline
}

func (ent *Entity) SetAutoTransition1() {
	ent.flags |= autoTransition1
}

func (ent *Entity) SetSwirl() {
	ent.flags |= swirl
}

func (ent *Entity) SetAltImage() {
	ent.flags |= altImage
}

func (ent *Entity) SetLowLight() {
	ent.flags |= lowLight
}

func (ent *Entity) SetOpacityFade() {
	ent.flags |= opacityFade
}

func (ent *Entity) SetRevert() {
	ent.flags |= revert
}

func (ent *Entity) SetLikeWindow() {
	ent.flags |= likeWindow
}

func (ent *Entity) SetClickTransition() {
	ent.flags |= clickTransition
}

func (ent *Entity) SetDraw() {
	ent.flags |= draw
}

func (ent *Entity) SetFreeze() {
	ent.flags |= freeze
}

func (ent *Entity) SetUsed() {
	ent.flags |= used
}

func (ent *Entity) SetOneOff() {
	ent.flags |= oneOff
}

func (ent *Entity) SetUnDraw() {
	ent.flags |= unDraw
}

func (ent *Entity) SetClickForTime() {
	ent.flags |= clickForTime
}

func (ent *Entity) SetParticlesGenerated() {
	ent.flags |= particlesGenerated
}

func (ent *Entity) SetKeepShader() {
	ent.flags |= keepShader
}

func (ent *Entity) SetFishCaught() {
	ent.flags |= fishCaught
}

func (ent *Entity) SetText() {
	ent.flags |= textOnSprite
}

func (ent *Entity) SetAddClickEffect() {
	ent.flags |= addClickEffect
}

func (ent *Entity) SetAddIdleClickEffect() {
	ent.flags |= addIdleClickEffect
}

func (ent *Entity) SetNoActivationRect() {
	ent.flags |= noActivationRect
}

func (ent *Entity) SetDontDrawFirstDay() {
	ent.flags |= dontDrawFirstDay
}

func (ent *Entity) SetCounter() {
	ent.flags |= UsingCounter
}
func (ent *Entity) SetDeInitAfterUse() {
	ent.flags |= DeInitAfterUsed
}

// Getter methods
func (ent *Entity) HasNoOffset() bool {
	return ent.flags&noOffset != 0
}

func (ent *Entity) HasUpdater() bool {
	return ent.flags&updater != 0
}

func (ent *Entity) HasOutline() bool {
	return ent.flags&outline != 0
}

func (ent *Entity) HasAutoTransition1() bool {
	return ent.flags&autoTransition1 != 0
}

func (ent *Entity) HasSwirl() bool {
	return ent.flags&swirl != 0
}

func (ent *Entity) HasAltImage() bool {
	return ent.flags&altImage != 0
}

func (ent *Entity) HasLowLight() bool {
	return ent.flags&lowLight != 0
}

func (ent *Entity) HasOpacityFade() bool {
	return ent.flags&opacityFade != 0
}

func (ent *Entity) HasRevert() bool {
	return ent.flags&revert != 0
}

func (ent *Entity) HasLikeWindow() bool {
	return ent.flags&likeWindow != 0
}

func (ent *Entity) HasClickTransition() bool {
	return ent.flags&clickTransition != 0
}

func (ent *Entity) HasDraw() bool {
	return ent.flags&draw != 0
}

func (ent *Entity) HasFreeze() bool {
	return ent.flags&freeze != 0
}

func (ent *Entity) HasUsed() bool {
	return ent.flags&used != 0
}

func (ent *Entity) HasOneOff() bool {
	return ent.flags&oneOff != 0
}

func (ent *Entity) HasClickForTime() bool {
	return ent.flags&clickForTime != 0
}

func (ent *Entity) HasUnDraw() bool {
	return ent.flags&unDraw != 0
}

func (ent *Entity) HasParticlesGenerated() bool {
	return ent.flags&particlesGenerated != 0
}

func (ent *Entity) HasKeepShader() bool {
	return ent.flags&keepShader != 0
}

func (ent *Entity) HasFishCaught() bool {
	return ent.flags&fishCaught != 0
}

func (ent *Entity) HasText() bool {
	return ent.flags&textOnSprite != 0
}

func (ent *Entity) HasOverZoom() bool {
	return ent.flags&overZoom != 0
}

func (ent *Entity) HasOverUi() bool {
	return ent.flags&overUi != 0
}

func (ent *Entity) HasAddClickEffect() bool {
	return ent.flags&addClickEffect != 0
}

func (ent *Entity) HasIdleClickEffect() bool {
	return ent.flags&addIdleClickEffect != 0

}

func (ent *Entity) HasNoActivationRect() bool {
	return ent.flags&noActivationRect != 0
}

func (ent *Entity) HasDontDrawFirstDay() bool {
	return ent.flags&dontDrawFirstDay != 0
}

func (ent *Entity) HasCounter() bool {
	return ent.flags&UsingCounter != 0
}

func (ent *Entity) HasDeInitAfterUSed() bool {
	return ent.flags&DeInitAfterUsed != 0
}

// Clear methods
func (ent *Entity) ClearNoOffset() {
	ent.flags &^= noOffset
}

func (ent *Entity) ClearUpdater() {
	ent.flags &^= updater
}

func (ent *Entity) ClearOutline() {
	ent.flags &^= outline
}

func (ent *Entity) ClearAutoTransition1() {
	ent.flags &^= autoTransition1
}

func (ent *Entity) ClearSwirl() {
	ent.flags &^= swirl
}

func (ent *Entity) ClearAltImage() {
	ent.flags &^= altImage
}

func (ent *Entity) ClearLowLight() {
	ent.flags &^= lowLight
}

func (ent *Entity) ClearOpacityFade() {
	ent.flags &^= opacityFade
}

func (ent *Entity) ClearRevert() {
	ent.flags &^= revert
}

func (ent *Entity) ClearLikeWindow() {
	ent.flags &^= likeWindow
}

func (ent *Entity) ClearClickTransition() {
	ent.flags &^= clickTransition
}

func (ent *Entity) ClearDraw() {
	ent.flags &^= draw
}

func (ent *Entity) ClearFreeze() {
	ent.flags &^= freeze
}

func (ent *Entity) ClearUsed() {
	ent.flags &^= used
}

func (ent *Entity) ClearClickForTime() {
	ent.flags &^= clickForTime
}

func (ent *Entity) ClearOneOff() {
	ent.flags &^= oneOff
}

func (ent *Entity) ClearUnDraw() {
	ent.flags &^= unDraw
}

func (ent *Entity) ClearParticlesGenerated() {
	ent.flags &^= particlesGenerated
}

func (ent *Entity) ClearKeepShader() {
	ent.flags &^= keepShader
}

func (ent *Entity) ClearFishCaught() {
	ent.flags &^= fishCaught
}

func (ent *Entity) ClearText() {
	ent.flags &^= textOnSprite
}

func (ent *Entity) ClearOverZoom() {
	ent.flags &^= overZoom
}

func (ent *Entity) ClearAddClickEffect() {
	ent.flags &^= addClickEffect
}

func (ent *Entity) ClearIdleClickEffect() {
	ent.flags &^= addIdleClickEffect
}

func (ent *Entity) ClearNoActivationRect() {
	ent.flags &^= noActivationRect
}
