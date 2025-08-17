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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"log"
	"math"
	"strconv"
	time2 "time"
)

const (
	mapWidth      = 20
	mapHeight     = 11
	time          = 30
	treePositionX = 0
	treePositionY = 0
)

type mowState uint32

const (
	loaded mowState = iota
	started
	finished
)

var SoundPlayer *soundFX.SoundPlayer

type MowingScene struct {
	images             map[string]*ebiten.Image
	state              mowState
	gameState          *SceneStateMachine
	smallerResolution  *ebiten.Image
	brush              *ebiten.Image
	gameLog            *sceneManagement.GameLog
	isLoaded           bool
	gameMap            [mapHeight][mapWidth]int
	locationMap        [mapHeight][mapWidth]int
	colliders          []image.Rectangle
	character          *entities.Entity
	score              int
	time               float64
	timeString         string
	mowedOverLay       *ebiten.Image
	scoreString        string
	timers             map[string]*util.Timer
	allowanceTime      bool
	allowanceString    string
	allowance          float64
	collisionsOccurred []entities.Collision
	direction          string
	frameCount         int
	debug              bool
	ui                 *ebitenui.UI
	sprites            []*sprite.Sprite
	returnScene        sceneManagement.SceneId
	removeWindowFunc   widget.RemoveWindowFunc
}

func NewMowingScene(gameLog *sceneManagement.GameLog) *MowingScene {

	s := &MowingScene{}
	s.images = LoadMowImages()
	s.gameLog = gameLog
	s.gameMap = LoadMap()
	s.sprites = LoadMowSprites(s.images)
	s.colliders = loadMapCollisions(s.gameMap)
	s.mowedOverLay = ebiten.NewImage(ScreenWidth, ScreenHeight)
	s.brush = s.images["grassBrush"]
	s.timers = make(map[string]*util.Timer)
	s.timers["calcAllowance"] = util.NewTimer(0.3)

	s.smallerResolution = ebiten.NewImage(ScreenWidth, ScreenHeight)

	SoundPlayer = s.gameLog.SoundPlayer
	LoadChar(s)

	s.time = time
	s.returnScene = sceneManagement.MowingMiniGameScene

	gameStates := make(map[int]*MowSceneStateHandler)
	gameStates[1] = &MowSceneStateHandler{Updater: GameJustStartedUpdate, TransitionTo: 2, TransitionFunc: mowTransition}
	gameStates[2] = &MowSceneStateHandler{Updater: GameMowingState, TransitionTo: 3, TransitionFunc: overSceneTransition}
	gameStates[3] = &MowSceneStateHandler{Updater: CalcScoreState, TransitionTo: 0}
	s.gameState = &SceneStateMachine{States: gameStates, CurrentState: 1}
	return s
}

type SceneStateMachine struct {
	States       map[int]*MowSceneStateHandler
	CurrentState int
}

type MowSceneStateHandler struct {
	Updater        func(s *MowingScene)
	TransitionFunc func(s *MowingScene)
	TransitionTo   int
}

func (s *SceneStateMachine) Transition(m *MowingScene) {
	if s.States[s.CurrentState].TransitionFunc != nil {
		s.States[s.CurrentState].TransitionFunc(m)
	}
	s.CurrentState = s.States[s.CurrentState].TransitionTo
}

func mowTransition(s *MowingScene) {
	graphics.NewUpdateAbleTextGraphic(&s.scoreString, 10, 10)
	graphics.NewUpdateAbleTextGraphic(&s.timeString, 140, 10)
}

func UpdateMowerCharacter(ent *entities.Entity) {
	if ent.StateMachine != nil {
		ent.StateMachine.States[ent.StateMachine.CurrentState].Updater(ent)
	}
}

func GameJustStartedUpdate(s *MowingScene) {
	s.ui.Update()
}

func GameMowingState(s *MowingScene) {
	s.frameCount++
	UpdateMowerCharacter(s.character)
	updateTimeAndScore(s)
	if s.time <= 0 {
		s.gameState.Transition(s)
	}
}

func CalcScoreState(s *MowingScene) {
	s.frameCount++
	s.allowanceString = fmt.Sprintf("Allowance Earned: $ %0.2f", s.allowance)
	msg := fmt.Sprintf("Allowance Earned: $ %0.2f", s.allowance)
	if s.frameCount == 120 {
		graphics.NewUpdateAbleTextGraphic(&s.allowanceString, ScreenWidth/2, ScreenHeight/2+20)
	}
	if s.frameCount > 120 {
		updateScoreAfterTimeLimit(s)
	}
	if s.frameCount == 600 {
		s.returnScene = sceneManagement.FishTank
		/*	msg := "Press r to try again"
			graphics.NewUpdateAbleTextGraphic(&msg, ScreenWidth/2, ScreenHeight/2+45)*/
	}

	fmt.Printf("Local: '%s' (len=%d)\n", msg, len(msg))
	fmt.Printf("Struct: '%s' (len=%d)\n", s.allowanceString, len(s.allowanceString))
	fmt.Printf("Equal: %t\n", msg == s.allowanceString)

}

func overSceneTransition(s *MowingScene) {
	SoundPlayer.Pause()
	s.frameCount = 0
	msg := "Times Up!"
	graphics.NewUpdateAbleTextGraphic(&msg, ScreenWidth/2, ScreenHeight/2-5)
}

func (s *MowingScene) Update() (sceneManagement.SceneId, error) {
	s.gameLog.SoundPlayer.Update()
	graphics.UpdateGraphics()
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		return sceneManagement.Reset, nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		s.debug = true
	}
	s.gameState.States[s.gameState.CurrentState].Updater(s)
	if s.debug {
		s.debugUpdate()
	}

	return s.returnScene, nil
}

func Crash(character *entities.Entity) {
	SoundPlayer.Pause()
	SoundPlayer.Play(soundFX.Crash)
	time2.AfterFunc(1*time2.Second, func() { SoundPlayer.Pause() })
	id := graphics.NewFadeInTextGraphic("Crash!",
		float64(character.Sprite.X+4)*3,
		float64(character.Sprite.Y+2)*3)

	eff := entImportableLoaders.LoadStaticEffect("Crash", character.Sprite.X-32, character.Sprite.Y-32, "")
	id2 := graphics.AddGraphic(&graphics.SpriteGraphic{Sprite: *eff})
	time2.AfterFunc(1*time2.Second, func() { graphics.DeInitGraphicId(id) })
	time2.AfterFunc(1*time2.Second, func() { graphics.DeInitGraphicId(id2) })

	character.Sprite.PublishedGraphicId = append(character.Sprite.PublishedGraphicId, id, id2)
	character.StateMachine.Transition()
}

func Mowing(character *entities.Entity) {
	character.Update(character.TankMovement.WorldBoundaries)
}

func JustStarted(character *entities.Entity) {

	if character.Sprite != character.AnimationMap["StartUp"] {
		character.Sprite = character.AnimationMap["StartUp"]
		character.Sprite.X = character.AnimationMap["Moving"].X
		character.Sprite.Y = character.AnimationMap["Moving"].Y
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) && ebiten.IsKeyPressed(ebiten.KeyU) {
		SoundPlayer.Play(soundFX.FailedStart)
		character.Sprite.Update()
	} else if character.Sprite.Frame() != 0 {
		SoundPlayer.Pause()
		character.Sprite.Animation.Reset()
	}
	if character.Sprite.Frame() == character.Sprite.LastF {
		SoundPlayer.Pause()
		SoundPlayer.Play(soundFX.MowerRunning)
		character.Sprite.Animation.Reset()
		//transferState
		character.Sprite = character.AnimationMap["Moving"]
		character.StateMachine.Transition()
	}
}

func (s *MowingScene) debugUpdate() {
	switch s.character.TankMovement.Direction {
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

func (s *MowingScene) Draw(screen *ebiten.Image) {
	if s.images["map"] == nil {
		log.Print("Map image is nil, cannot draw tilemap")
		return
	}

	drawSimpleTileMap(s.smallerResolution, s.images["map"], s.gameMap, 16)

	for _, sp := range s.sprites {
		sp.Draw(s.smallerResolution)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Rotate(s.character.MovementSystem.Params.Direction)
	angle := s.character.MovementSystem.Params.Direction

	x := float64(s.character.TankMovement.Corners.FrontRight.X + s.character.TankMovement.Corners.FrontLeft.X - s.character.TankMovement.Corners.FrontRight.X)
	y := float64(s.character.TankMovement.Corners.FrontRight.Y + s.character.TankMovement.Corners.FrontLeft.Y - s.character.TankMovement.Corners.FrontRight.Y)

	frontDistance := float64(0) // how far forward
	sideDistance := float64(0)  // how far to the side (adjust as needed)

	// Calculate forward and perpendicular directions
	forwardX := math.Cos(angle) * frontDistance
	forwardY := math.Sin(angle) * frontDistance
	sideX := math.Cos(angle+math.Pi/2) * sideDistance // perpendicular
	sideY := math.Sin(angle+math.Pi/2) * sideDistance

	op.GeoM.Translate(x+forwardX+sideX, y+forwardY+sideY)
	if s.character.StateMachine.CurrentState == 2 {
		s.mowedOverLay.DrawImage(s.brush, op)
	}
	op.GeoM.Reset()
	// In draw function, draw mowed layer before character

	s.smallerResolution.DrawImage(s.mowedOverLay, op)
	if s.character != nil && s.character.Sprite.Img != nil {
		dopts := &ebiten.DrawImageOptions{}
		dopts.GeoM.Translate(float64(-s.character.Sprite.SpriteWidth/2), -float64(s.character.Sprite.SpriteHeight/2))
		dopts.GeoM.Rotate(s.character.MovementSystem.Params.Direction)
		dopts.GeoM.Translate(float64(s.character.Sprite.X), float64(s.character.Sprite.Y))
		s.character.Sprite.UpdateOpts(dopts)
		s.character.Sprite.Draw(s.smallerResolution)
	}

	graphics.DrawScaledGraphics(s.smallerResolution)

	if s.debug {
		s.DebugDraw(s.smallerResolution)
	}
	dOpts := &ebiten.DrawImageOptions{}

	dOpts.GeoM.Scale(registry.Config.ResolutionScalingF*3, registry.Config.ResolutionScalingF*3)

	screen.DrawImage(s.smallerResolution, dOpts)

	graphics.DrawUnScaledGraphics(screen)

	if s.state == loaded || s.state == finished {
		s.ui.Draw(screen)
	}
}

func (s *MowingScene) DebugDraw(screen *ebiten.Image) {
	drawCollisionMap(MakeDebugCollisionMap(s.colliders, s.character.TankMovement.Corners), s.character.TankMovement, s.smallerResolution)
	debugPrintCollisions(s, screen)
	vector.StrokeRect(
		s.smallerResolution,
		float32(s.character.TankMovement.WorldBoundaries.Min.X),
		float32(s.character.TankMovement.WorldBoundaries.Min.Y),
		float32(s.character.TankMovement.WorldBoundaries.Dx()),
		float32(s.character.TankMovement.WorldBoundaries.Dy()),
		1,
		colornames.Red,
		false)
}

func debugPrintCollisions(s *MowingScene, screen *ebiten.Image) {
	for _, col := range s.collisionsOccurred {
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

	log.Printf("Entering Mowing Scene")
	s.gameLog.SongPlayer.Play(soundFX.IndieCafe)

	stringSlice := []string{
		"1. Press Space and U to start your mower",
		"2. WASD or Arrow Keys to move: forward and backwards for acceleration, Left and Right for Direction",
		"3. Mow as much grass as possible to\n earn a higher allowance",
		"4. If you crash you'll need to start the mower again"}

	s.removeWindowFunc = ui.TriggerTextWindow(s.gameLog.GlobalEventHub, s.ui, "How To Play", stringSlice)

	s.subs(s.gameLog.GlobalEventHub)
}

func (s *MowingScene) OnExit() {
	log.Printf("Leaving Mowing Scene")
	s.gameLog.SongPlayer.Pause()
	graphics.DeInitAllGraphics()
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
			s.gameState.Transition(s)
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

func MakeDebugCollisionMap(cols []image.Rectangle, corners *entities.TankCorners) map[image.Rectangle]bool {
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
	if s.allowance < float64(s.score/3)*0.05 {
		s.allowance += 0.005
		s.gameLog.SoundPlayer.Play(soundFX.Kaching)
	}
	if s.allowance >= float64(s.score/3)*0.05 {
		s.gameLog.SoundPlayer.Pause()
	}

	s.ui.Update()
}

func updateTimeAndScore(s *MowingScene) {
	if s.frameCount%30 == 0 { // every 0.5 seconds at 60fps
		pixels := make([]byte, ScreenWidth*ScreenHeight*4)
		s.mowedOverLay.ReadPixels(pixels)
		mowedPixels := 0
		for i := 0; i < len(pixels); i += 4 {
			if pixels[i+3] > 0 { // alpha > 0 means mowed
				mowedPixels++
			}
		}
		s.score = mowedPixels / 100 // scale it down
	}
	s.scoreString = fmt.Sprintf("Score: %d", s.score)
	s.timeString = fmt.Sprintf("Time: " + strconv.FormatFloat(s.time, 'f', 2, 32))
	s.time -= 0.016
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
	characterOrignX := float32(mapWidth-2) * 16
	characterOriginY := float32(mapHeight-3) * 16

	println("charx at time of load:", characterOrignX, "char y at time of load", characterOriginY)
	if s.images["characterSpriteSheet"] == nil {
		log.Print("ERROR: TankCharacter sprite image not loaded in map")
	}

	charAnimation, charSpriteSheet, err := entImportableLoaders.LoadAnimation("data/animationData/lawnMowingCharacterSprite.json")

	if err != nil {
		log.Fatal(err)
	}

	asp := &sprite.Sprite{Img: s.images["characterSpriteSheet"], X: characterOrignX, Y: characterOriginY, SpriteSheet: charSpriteSheet, Animation: charAnimation}

	startUpAnimation2, startUpSpriteSheet2, err := entImportableLoaders.LoadAnimation("data/animationData/lawnMowerStartAnimation.json")
	startUpAnimation2.SpeedInTPS = 15

	if err != nil {
		log.Fatal(err)
	}

	asp2 := &sprite.Sprite{Img: s.images["lawnMowerStartSpriteSheet"], X: characterOrignX, Y: characterOriginY, SpriteSheet: startUpSpriteSheet2, Animation: startUpAnimation2}

	animationMap := make(map[string]*sprite.Sprite)

	animationMap["StartUp"] = asp2
	animationMap["Moving"] = asp

	movementParams := movement.Params{
		MaxSpeed:     1.9, // Slower for a mowing game
		Acceleration: 0.0, // Moderate acceleration
		Friction:     0.5, // High friction for precise control
	}

	movementS := movement.NewMovementSystem(movementParams, &movement.WASDInputHandler{})

	tankMove := entities.TankCharacter{}
	character := &entities.Entity{MovementSystem: movementS, Sprite: asp2, AnimationMap: animationMap, MovementState: &movement.State{}, TankMovement: &tankMove}
	character.TankMovement.WorldBoundaries = image.Rect(16, -8, (mapWidth+2)*16, (mapHeight+2)*16)
	character.MovementSystem = movementS
	character.TankMovement.Corners = entities.GetCharCorners(character)
	s.character = character
	mowerState1 := entities.StateHandler{Updater: JustStarted, TransitionTo: 2}
	mowerState2 := entities.StateHandler{Updater: Mowing, TransitionTo: 3}
	mowerState3 := entities.StateHandler{Updater: Crash, TransitionTo: 1}

	states := make(map[int]*entities.StateHandler)
	states[1] = &mowerState1
	states[2] = &mowerState2
	states[3] = &mowerState3
	charSM := &entities.StateMachine{States: states, CurrentState: 1}

	s.character.StateMachine = charSM
	s.character.TankMovement.Colliders = s.colliders
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
			if i >= mapHeight-2 {
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
	arr[3][7] = 8
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
	for _ = range 3 {
		treeSprite := &sprite.Sprite{Img: imgMap["treeSprite"], X: treePositionX, Y: treePositionY + yoffset}
		yoffset += float32(treeSprite.Img.Bounds().Dy()/2) + float32(20)
		sprites = append(sprites, treeSprite)
	}

	return sprites
}
