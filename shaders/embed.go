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
