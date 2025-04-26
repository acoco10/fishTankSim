package assets

import (
	"embed"
)

//go:embed images
var ImagesDir embed.FS

//go:embed data
var DataDir embed.FS

//go:embed fonts
var FontsDir embed.FS

//go:embed soundFx
var SoundDir embed.FS
