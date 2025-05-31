package shaders

import _ "embed"

// OutlineShader solid colored outline shader
//
//go:embed OutlineEffects/outlineVer2.kage
var OutlineShader string

//go:embed OutlineEffects/pulseOutline.kage
var PulseOutline string

//go:embed OutlineEffects/rotatingHighLightOutline.kage
var RotatingHighlightOutline string

//go:embed solidColor.kage
var SolidColor string

//lighting shaders

// OnePointLighting shader intended to provide blue global lighting effect emanating from one point
//
//go:embed lighting/onePointLighting.kage
var OnePointLighting string

//go:embed lighting/spriteLightingEffects.kage
var SpriteLightingEffect string
