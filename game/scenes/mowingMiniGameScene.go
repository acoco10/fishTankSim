package scenes

import (
	"bytes"
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/movement"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/ui"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"log"
	"strconv"
)

const (
	mapWidth      = 15
	mapHeight     = 12
	time          = 10
	treePositionX = 0
	treePositionY = 0
)

type mowState uint32

const (
	loaded mowState = iota
	started
	finished
)

type MowingScene struct {
	images            map[string]*ebiten.Image
	state             mowState
	smallerResolution *ebiten.Image
	gameLog           *sceneManagement.GameLog
	isLoaded          bool
	gameMap           [mapHeight][mapWidth]int
	locationMap       [mapHeight][mapWidth]int
	colliders         []image.Rectangle
	character         *entities.TankCharacter
	score             int
	time              float64
	timeString        string
	scoreString       string
	timers            map[string]*util.Timer
	allowanceTime     bool
	allowanceString   string
	allowance         float64
	mowerStarted      bool
	collisionsOccured []entities.Collision
	direction         string
	debug             bool
	ui                *ebitenui.UI
	sprites           []*sprite.Sprite
	returnScene       sceneManagement.SceneId
	removeWindowFunc  widget.RemoveWindowFunc
}

func NewMowingScene(gameLog *sceneManagement.GameLog) *MowingScene {

	s := &MowingScene{}
	s.images = LoadMowImages()
	s.gameLog = gameLog
	s.gameMap = LoadMap()
	s.sprites = LoadMowSprites(s.images)
	s.colliders = loadMapCollisions(s.gameMap)

	s.timers = make(map[string]*util.Timer)
	s.timers["calcAllowance"] = util.NewTimer(0.3)

	s.smallerResolution = ebiten.NewImage(ScreenWidth, ScreenHeight)

	LoadChar(s)

	//graphics.NewUpdateAbleTextGraphic(&s.direction, 300, 200)
	s.time = time
	s.character.Update([]entities.Collision{})
	s.returnScene = sceneManagement.MowingMiniGameScene
	return s
}

func (s *MowingScene) Update() (sceneManagement.SceneId, error) {
	if s.state == loaded || s.state == finished {
		s.ui.Update()
	}

	updateMowingTimers(s, s.timers)

	if s.time < 0.1 && s.time > 0.0 {
		s.timers["calcAllowance"].TurnOn()
		returnmsg := fmt.Sprintf("Times up! Square footage mowed: %d", s.score)
		graphics.NewFadeInTextGraphic(returnmsg, 400, 200)
	}

	if s.time <= 0.0 {
		updateScoreAfterTimeLimit(s)
	} else {
		updateTimeAndScore(s)
		if s.character.Moving {

			if !ebiten.IsKeyPressed(ebiten.KeySpace) {
				//keep space bar held to keep mower on
				s.mowerStarted = false
				s.character.Sprite = s.character.AnimationMap["StartUp"]
				s.character.Sprite.X = s.character.AnimationMap["Moving"].X
				s.character.Sprite.Y = s.character.AnimationMap["Moving"].Y
			}
		}

		if s.character.Moving == false {
			s.mowerStarted = false
			s.character.Sprite = s.character.AnimationMap["StartUp"]
			if ebiten.IsKeyPressed(ebiten.KeySpace) && ebiten.IsKeyPressed(ebiten.KeyU) {
				if s.character.Sprite.Frame() == s.character.Sprite.LastF {
					//last frame of startup = start mower
					s.mowerStarted = true
					s.character.Moving = true
					s.character.Sprite.Animation.Reset()
				}
			}
		}

		for _, sp := range s.sprites {
			sp.Update()
		}

		s.collisionsOccured = CheckCollision(s.character.Corners, s.colliders)
		s.character.Update(s.collisionsOccured)
		s.updateScore()

		if s.debug {
			s.debugUpdate()
		}

		graphics.UpdateGraphics()
	}

	return s.returnScene, nil
}

func (s *MowingScene) debugUpdate() {
	switch s.character.Direction {
	case entities.CharRight:
		s.direction = "right"
	case entities.CharLeft:
		s.direction = "left"
	case entities.Down:
		s.direction = "down"
	case entities.Up:
		s.direction = "up "
	}
}

func (s *MowingScene) updateScore() {

	if s.character.Sprite.X < 0 || s.character.Sprite.Y < 0 {
		println("coordinates are fucked")
	}
	charTileY := int(s.character.Sprite.Y / 16)
	charTileX := int(s.character.Sprite.X / 16)
	if s.gameMap[charTileY][charTileX] == 1 {
		s.gameMap[charTileY][charTileX] = 2
		s.score++
	}

}

func (s *MowingScene) Draw(screen *ebiten.Image) {
	if s.images["map"] == nil {
		log.Print("Map image is nil, cannot draw tilemap")
		return
	}

	drawSimpleTileMap(s.smallerResolution, s.images["map"], s.gameMap, 16)

	if s.character != nil && s.character.Sprite.Img != nil {
		dopts := &ebiten.DrawImageOptions{}
		dopts.GeoM.Translate(float64(-s.character.Sprite.SpriteWidth/2), -float64(s.character.Sprite.SpriteHeight/2))
		dopts.GeoM.Rotate(s.character.MovementSystem.Params.Direction)
		dopts.GeoM.Translate(float64(s.character.Sprite.X), float64(s.character.Sprite.Y))
		s.character.Sprite.UpdateOpts(dopts)
		s.character.Sprite.Draw(s.smallerResolution)
	}

	for _, sp := range s.sprites {
		sp.Draw(s.smallerResolution)
	}

	dOpts := &ebiten.DrawImageOptions{}
	xOffset := float64(registry.Config.ResolutionWidth) - (mapWidth*16)*registry.Config.ResolutionScalingF*2
	xOffset = xOffset / 8
	dOpts.GeoM.Translate(float64(xOffset), registry.Config.YOffsetF)
	dOpts.GeoM.Scale(registry.Config.ResolutionScalingF*2, registry.Config.ResolutionScalingF*2)

	screen.DrawImage(s.smallerResolution, dOpts)
	graphics.DrawUnScaledGraphics(screen)

	if s.debug {
		s.DebugDraw(screen)
	}

	if s.state == loaded || s.state == finished {
		s.ui.Draw(screen)
	}

}

func (s *MowingScene) DebugDraw(screen *ebiten.Image) {
	drawCollisionMap(MakeCollisionMap(s.colliders, s.character.Corners), s.character, s.smallerResolution)
	debugPrintCollisions(s, screen)
}

func debugPrintCollisions(s *MowingScene, screen *ebiten.Image) {
	for _, col := range s.collisionsOccured {
		switch col.Corner {
		case entities.FrontLeft:
			{
				ebitenutil.DebugPrint(screen, "Front Left Collision!")
				return
			}
		case entities.FrontRight:
			{
				ebitenutil.DebugPrint(screen, "Front Right Collision!")
				return
			}
		case entities.RearRight:
			{
				ebitenutil.DebugPrint(screen, "rear right collision!")
				return
			}
		case entities.RearLeft:
			{
				ebitenutil.DebugPrint(screen, "rear left collison!")
				return
			}
		}
	}
}

func (s *MowingScene) FirstLoad() {
	s.ui = ui.LoadMowingUI(s.gameLog.GlobalEventHub)
	s.isLoaded = true
}

func (s *MowingScene) OnEnter() {
	graphics.NewUpdateAbleTextGraphic(&s.scoreString, 10, 10)
	graphics.NewUpdateAbleTextGraphic(&s.timeString, 150, 10)
	log.Printf("Entering Mowing Scene")
	s.gameLog.SongPlayer.Play(soundFX.IndieCafe)

	stringSlice := []string{
		"1. Press Space and U to start your mower",
		" 2. Hold Space to keep it running",
		" 3. Mow as much grass as possible to\n earn a higher allowance"}

	s.removeWindowFunc = ui.TriggerTextWindow(s.gameLog.GlobalEventHub, s.ui, "How To Play", stringSlice)

	s.subs(s.gameLog.GlobalEventHub)
}

func (s *MowingScene) OnExit() {
	log.Printf("Leaving Mowing Scene")
	s.gameLog.SongPlayer.Pause()
	s.gameLog.GlobalEventHub.Publish(events.MoneyAvailable{Amount: s.allowance})
}

func (s *MowingScene) IsLoaded() bool {
	return s.isLoaded
}

func (s *MowingScene) subs(eventHub *tasks.EventHub) {
	eventHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		println(ev.ButtonText, "button event received")
		switch ev.ButtonText {
		case "Continue":
			println("switching back to fish tank ")
			s.returnScene = sceneManagement.FishTank
		case "Lets Mow!":
			if s.removeWindowFunc != nil {
				s.removeWindowFunc()
			}
			s.state = started
		}

	})
}

func loadMapCollisions(tileMap [mapHeight][mapWidth]int) []image.Rectangle {

	var cols []image.Rectangle

	for i := range tileMap {
		for j := range tileMap[i] {
			if tileMap[i][j] != 1 && tileMap[i][j] != 2 {
				x0 := j * 16
				y0 := i * 16
				colRect := image.Rect(x0, y0, x0+16, y0+16)
				cols = append(cols, colRect)
			}
		}
	}

	return cols
}

func DrawRectangleFromPoints(screen *ebiten.Image, corners *entities.TankCorners, strokeColor color.Color, strokeWidth float32) {
	// DrawRectangleFromPoints draws lines connecting the 4 corner points to form a rectangle
	// Draw lines connecting the points in order:
	// 0 -> 1 -> 2 -> 3 -> 0 (back to start)

	// Top edge: point 0 to point 1
	vector.StrokeLine(screen,
		float32(corners.FrontLeft.X), float32(corners.FrontLeft.Y),
		float32(corners.FrontRight.X), float32(corners.FrontRight.Y),
		strokeWidth, strokeColor, false)

	// Right edge: point 1 to point 2
	vector.StrokeLine(screen,
		float32(corners.FrontRight.X), float32(corners.FrontRight.Y),
		float32(corners.RearRight.X), float32(corners.RearRight.Y),
		strokeWidth, strokeColor, false)

	// Bottom edge: point 2 to point 3
	vector.StrokeLine(screen,
		float32(corners.RearRight.X), float32(corners.RearRight.Y),
		float32(corners.RearLeft.X), float32(corners.RearLeft.Y),
		strokeWidth, strokeColor, false)

	// Left edge: point 3 to point 0
	vector.StrokeLine(screen,
		float32(corners.RearLeft.X), float32(corners.RearLeft.Y),
		float32(corners.FrontLeft.X), float32(corners.FrontLeft.Y),
		strokeWidth, strokeColor, false)
}

func MakeCollisionMap(cols []image.Rectangle, corners *entities.TankCorners) map[image.Rectangle]bool {
	colMap := make(map[image.Rectangle]bool)

	for _, col := range cols {

		colMap[col] = false
		if corners.FrontRight.In(col) {
			colMap[col] = true
		}
		if corners.FrontLeft.In(col) {
			colMap[col] = true
		}
		if corners.RearLeft.In(col) {
			colMap[col] = true
		}
		if corners.RearRight.In(col) {
			colMap[col] = true
		}
	}
	return colMap
}

func CheckCollision(corners *entities.TankCorners, colliders []image.Rectangle) []entities.Collision {
	var collisions []entities.Collision

	for _, col := range colliders {
		if corners.FrontRight.In(col) {
			collision := entities.Collision{Corner: entities.FrontRight, Object: col}
			collisions = append(collisions, collision)
		}
		if corners.FrontLeft.In(col) {
			collision := entities.Collision{Corner: entities.FrontLeft, Object: col}
			collisions = append(collisions, collision)
		}

		if corners.RearLeft.In(col) {
			collision := entities.Collision{Corner: entities.RearRight, Object: col}
			collisions = append(collisions, collision)
		}
		if corners.RearRight.In(col) {
			collision := entities.Collision{Corner: entities.RearLeft, Object: col}
			collisions = append(collisions, collision)
		}
	}
	return collisions
}

func drawCollisionMap(colMap map[image.Rectangle]bool, character *entities.TankCharacter, screen *ebiten.Image) {
	for col, collided := range colMap {
		switch collided {
		case true:
			vector.StrokeRect(
				screen,
				float32(col.Min.X),
				float32(col.Min.Y),
				16,
				16,
				1,
				colornames.Red,
				false)
		case false:
			vector.StrokeRect(
				screen,
				float32(col.Min.X),
				float32(col.Min.Y),
				16,
				16,
				1,
				colornames.Yellow,
				false)
		}

	}

	DrawRectangleFromPoints(screen, character.Corners, colornames.Blue, float32(1))

}

func updateScoreAfterTimeLimit(s *MowingScene) {

	graphics.UpdateGraphics()

	if s.allowanceTime {
		if float64(s.score)*0.05-s.allowance >= 0.05 {
			s.allowance += 0.05
			s.allowanceString = "Allowance Earned: $" + strconv.FormatFloat(s.allowance, 'f', 2, 32)
		}

	}

	if s.ui != nil {
		s.ui.Update()
	}
}

func updateMowingTimers(scene *MowingScene, timers map[string]*util.Timer) {
	for key, timer := range timers {
		state := timer.Update()
		switch key {
		case "calcAllowance":
			if state == util.Done {
				timer.TurnOff()
				scene.allowanceTime = true
				graphics.NewUpdateAbleTextGraphic(&scene.allowanceString, 400, 100)
			}
		}
	}
}

func updateTimeAndScore(s *MowingScene) {
	if s.state == started {
		if s.time > 0.0 {
			s.time = s.time - 0.016 //0.016 seconds per tick
		}

		if s.time-0.02 <= 0.0 {
			s.state = finished
			s.time = 0.0
		}

		s.scoreString = "Score: " + strconv.Itoa(s.score)

		s.timeString = "Time: " + strconv.FormatFloat(s.time, 'f', 2, 32)
	}
}

func drawSimpleTileMap(screen *ebiten.Image, mapImage *ebiten.Image, arr [mapHeight][mapWidth]int, tileSize int) {

	// DEBUG: Check mapImage bounds
	bounds := mapImage.Bounds()

	for y, row := range arr {
		for x, id := range row {

			if id == 0 {
				continue // Skip empty tiles
			}
			// Calculate source rectangle for the tile
			srcX := (id - 1) * tileSize // Assuming tiles are arranged horizontally
			srcY := 0                   // All tiles in first row

			// DEBUG: Check if source rectangle is within image bounds
			if srcX+tileSize > bounds.Dx() {
				continue
			}

			srcRect := image.Rect(srcX, srcY, srcX+tileSize, srcY+tileSize)

			// Create sub-image for this tile
			tileImg := mapImage.SubImage(srcRect)

			// Calculate destination position
			opts := &ebiten.DrawImageOptions{}

			opts.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
			// Draw the tile
			screen.DrawImage(tileImg.(*ebiten.Image), opts)
		}
	}
}

func LoadChar(s *MowingScene) {
	//start character in bottom leftish corner
	characterOrignX := float32((mapWidth - 3) * 16)
	characterOriginY := float32((mapHeight - 3) * 16)

	if s.images["characterSpriteSheet"] == nil {
		log.Print("ERROR: TankCharacter sprite image not loaded in map")
	}

	charAnimation, charSpriteSheet, err := entImportableLoaders.LoadAnimation("data/animationData/lawnMowingCharacterSprite.json")

	if err != nil {
		log.Fatal(err)
	}

	asp := &sprite.Sprite{Img: s.images["characterSpriteSheet"], X: characterOrignX, Y: characterOriginY, SpriteSheet: charSpriteSheet, Animation: charAnimation}

	startUpAnimation2, startUpSpriteSheet2, err := entImportableLoaders.LoadAnimation("data/animationData/lawnMowerStartAnimation.json")

	if err != nil {
		log.Fatal(err)
	}

	asp2 := &sprite.Sprite{Img: s.images["lawnMowerStartSpriteSheet"], X: characterOrignX, Y: characterOriginY, SpriteSheet: startUpSpriteSheet2, Animation: startUpAnimation2}

	animationMap := make(map[string]*sprite.Sprite)

	animationMap["StartUp"] = asp2
	animationMap["Moving"] = asp

	movementParams := movement.Params{
		MaxSpeed:     1.2, // Slower for a mowing game
		Acceleration: 0.0, // Moderate acceleration
		Friction:     0.5, // High friction for precise control
	}

	movementS := movement.NewMovementSystem(movementParams, &movement.WASDInputHandler{})

	character := entities.NewCharacter(100, 100, asp)

	character.AnimationMap = animationMap

	character.MovementSystem = movementS

	character.Sprite = character.AnimationMap["StartUp"]

	s.character = character
}

func LoadMap() [mapHeight][mapWidth]int {
	var arr [mapHeight][mapWidth]int

	for i := range arr {
		for j := range arr[i] {
			//default grass
			arr[i][j] = 1
			if j == 1 {
				//left side hedge
				arr[i][j] = 5
			}
			if j == mapWidth-1 {
				//right side hedge
				arr[i][j] = 5
			}
			if i == mapHeight-1 {
				//back line roof
				arr[i][j] = 3
			}
			if i == 0 && j > 0 {
				//back row
				arr[i][j] = 4
			}
			if i == 0 && j == mapWidth-1 {
				//right corner
				arr[i][j] = 6
			}
			if i == 0 && j == 1 {
				//left corner
				arr[i][j] = 7
			}
		}
	}
	return arr
}

func LoadMowImages() map[string]*ebiten.Image {

	dir, err := assets.ImagesDir.ReadDir("images/miniGames/mowingGame")

	if err != nil {
		log.Printf("Error reading directory: %v", err)
		return nil
	}

	imageMap := make(map[string]*ebiten.Image)

	for _, file := range dir {
		if file.IsDir() {
			continue // Skip subdirectories
		}

		// Read the file
		data, err := assets.ImagesDir.ReadFile("images/miniGames/mowingGame/" + file.Name())
		if err != nil {
			log.Printf("Error reading file %s: %v", file.Name(), err)
			continue // Skip files that can't be read
		}

		// Decode the image
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			log.Printf("Error decoding image %s: %v", file.Name(), err)
			continue // Skip files that can't be decoded as images
		}

		// Convert to ebiten image
		ebitenImg := ebiten.NewImageFromImage(img)

		// Use filename as key (you might want to remove extension)
		fileName := file.Name()[:len(file.Name())-4]
		log.Printf("Successfully loaded image: %s in mowing game minigame", fileName)
		imageMap[fileName] = ebitenImg
	}

	return imageMap
}

func LoadMowSprites(imgMap map[string]*ebiten.Image) []*sprite.Sprite {
	var sprites []*sprite.Sprite

	yoffset := float32(0.0)
	for _ = range 4 {
		treeSprite := &sprite.Sprite{Img: imgMap["treeSprite"], X: treePositionX, Y: treePositionY + yoffset}
		yoffset += float32(treeSprite.Img.Bounds().Dy()/2) + float32(20)
		sprites = append(sprites, treeSprite)
	}

	return sprites
}
