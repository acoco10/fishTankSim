package main

import (
	"flag"
	"log"
	"os"
)

const gameSkeleton = `package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Vec2 struct {
	X, Y float64
}

type Entity struct {
	Position Vec2
	Velocity Vec2
	Size     float64
	Active   bool
}

func (e *Entity) Update() {
	e.Position.X += e.Velocity.X
	e.Position.Y += e.Velocity.Y

	if e.Position.X < 0 || e.Position.X > 1280 {
		e.Velocity.X = -e.Velocity.X
	}
	if e.Position.Y < 0 || e.Position.Y > 960 {
		e.Velocity.Y = -e.Velocity.Y
	}
}

func (e *Entity) Draw(screen *ebiten.Image) {
	x := int(e.Position.X)
	y := int(e.Position.Y)
	size := int(e.Size)

	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(size), float64(size), colornames.Blue)
}

type Game struct {
	Entities []Entity
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())

	entities := make([]Entity, 100)
	for i := range entities {
		entities[i] = Entity{
			Position: Vec2{X: rand.Float64() * 1280, Y: rand.Float64() * 960},
			Velocity: Vec2{X: (rand.Float64() - 0.5) * 5, Y: (rand.Float64() - 0.5) * 5},
			Size:     10,
			Active:   true,
		}
	}

	return &Game{
		Entities: entities,
	}
}

func (g *Game) Update() error {
	for i := range g.Entities {
		if g.Entities[i].Active {
			g.Entities[i].Update()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		for i := range g.Entities {
			g.Entities[i].Velocity = Vec2{X: (rand.Float64() - 0.5) * 5, Y: (rand.Float64() - 0.5) * 5}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Black)

	for i := range g.Entities {
		if g.Entities[i].Active {
			g.Entities[i].Draw(screen)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1280, 960
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Generated Ebiten Skeleton")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
`

func main() {
	filename := flag.String("file", "game_skeleton.go", "Output filename for generated game code")
	flag.Parse()

	f, err := os.Create(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(gameSkeleton)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully generated %s", *filename)
}
