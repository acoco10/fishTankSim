package assets

import "embed"

//go:embed soundFx/oggs
var OGG embed.FS

//go:embed images
var ImagesDir embed.FS

//go:embed data/animationData
var AnimationDataDir embed.FS

//go:embed fonts
var FontsDir embed.FS

//go:embed data
var DataDir embed.FS
