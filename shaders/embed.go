package shaders

import _ "embed"

//go:embed outlineVer2.kage
var OutlineShader string

//go:embed test.kage
var Test string

//go:embed solidColor.kage
var SolidColor string

//go:embed pulseOutline.kage
var PulseOutline string
