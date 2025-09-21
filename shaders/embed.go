package shaders

import _ "embed"

// OutlineShader solid colored outline shader
//
//go:embed outlineEffects/outlineVer2.kage
var OutlineShader string

//go:embed outlineEffects/pulseOutline.kage
var PulseOutline string

//go:embed outlineEffects/rotatingHighLightOutline.kage
var RotatingHighlightOutline string

//go:embed solidColor.kage
var SolidColor string

//lighting shaders

// OnePointLightingBlue shader intended to provide blue global lighting effect emanating from one point
//
//go:embed lighting/onePointLightingBlue.kage
var OnePointLightingBlue string

//go:embed lighting/onePointLightingBlueSecondLayer.kage
var OnePointLightingBlueSecondLayer string

//go:embed lighting/onePointLightingNeutral.kage
var OnePointLightingNeutral string

//go:embed lighting/spriteLightingEffects.kage
var SpriteLightingEffect string

//go:embed lighting/normalMap.kage
var NormalMap string

//go:embed uiEffects/handwriting.kage
var HandWritingEffect string

//go:embed uiEffects/erase.kage
var EraseEffect string

//go:embed lighting/turnOffTheLights.kage
var TurnOffLight string

//go:embed lighting/dayLight.kage
var DayLight string

//go:embed water/water.kage
var Water string

//go:embed lighting/matteWallDayLight.kage
var Wall string

//go:embed uiEffects/phEffect.kage
var PHEffect string

//go:embed outlineEffects/pixelPerfectText.kage
var textOutline string

//go:embed lighting/figureOutCoordinateSystem.kage
var FigureOutCoords string

//go:embed uiEffects/highLight.kage
var Highlight string

//go:embed uiEffects/lowLight.kage
var LowLight string

//go:embed uiEffects/stomachFill.kage
var StomachFill string

//go:embed outlineEffects/normalMap+Outline.kage
var NormalMapOutlined string

//go:embed uiEffects/pulseHighlight.kage
var PulseHighlight string

//go:embed outlineEffects/normalMap+pulseOutline.kage
var PulseOutlineNormal string
