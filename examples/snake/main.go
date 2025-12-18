package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	gridSize     = 20
	gridWidth    = screenWidth / gridSize
	gridHeight   = screenHeight / gridSize
)

// Direction represents movement direction.
type Direction int

const (
	DirUp Direction = iota
	DirDown
	DirLeft
	DirRight
)

// Point represents a grid position.
type Point struct {
	X, Y int
}

// Snake represents the snake game.
type Snake struct {
	body      []Point
	direction Direction
	nextDir   Direction
	food      Point
	score     int
	gameOver  bool
	moveTimer float64
	moveDelay float64
}

// NewSnake creates a new snake game.
func NewSnake() *Snake {
	rand.Seed(time.Now().UnixNano())
	s := &Snake{
		body: []Point{
			{X: gridWidth / 2, Y: gridHeight / 2},
			{X: gridWidth/2 - 1, Y: gridHeight / 2},
			{X: gridWidth/2 - 2, Y: gridHeight / 2},
		},
		direction: DirRight,
		nextDir:   DirRight,
		moveDelay: 0.1, // Move every 0.1 seconds
	}
	s.spawnFood()
	return s
}

func (s *Snake) spawnFood() {
	for {
		s.food = Point{
			X: rand.Intn(gridWidth),
			Y: rand.Intn(gridHeight),
		}
		// Make sure food doesn't spawn on snake
		onSnake := false
		for _, p := range s.body {
			if p.X == s.food.X && p.Y == s.food.Y {
				onSnake = true
				break
			}
		}
		if !onSnake {
			break
		}
	}
}

func (s *Snake) Update() error {
	if s.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			*s = *NewSnake()
		}
		return nil
	}

	// Handle input
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if s.direction != DirDown {
			s.nextDir = DirUp
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if s.direction != DirUp {
			s.nextDir = DirDown
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if s.direction != DirRight {
			s.nextDir = DirLeft
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if s.direction != DirLeft {
			s.nextDir = DirRight
		}
	}

	// Update movement timer
	s.moveTimer += 1.0 / 60.0
	if s.moveTimer >= s.moveDelay {
		s.moveTimer = 0
		s.direction = s.nextDir
		s.move()
	}

	return nil
}

func (s *Snake) move() {
	head := s.body[0]
	var newHead Point

	switch s.direction {
	case DirUp:
		newHead = Point{X: head.X, Y: head.Y - 1}
	case DirDown:
		newHead = Point{X: head.X, Y: head.Y + 1}
	case DirLeft:
		newHead = Point{X: head.X - 1, Y: head.Y}
	case DirRight:
		newHead = Point{X: head.X + 1, Y: head.Y}
	}

	// Wrap around screen edges
	if newHead.X < 0 {
		newHead.X = gridWidth - 1
	}
	if newHead.X >= gridWidth {
		newHead.X = 0
	}
	if newHead.Y < 0 {
		newHead.Y = gridHeight - 1
	}
	if newHead.Y >= gridHeight {
		newHead.Y = 0
	}

	// Check self collision
	for _, p := range s.body {
		if p.X == newHead.X && p.Y == newHead.Y {
			s.gameOver = true
			return
		}
	}

	// Move snake
	s.body = append([]Point{newHead}, s.body...)

	// Check food collision
	if newHead.X == s.food.X && newHead.Y == s.food.Y {
		s.score += 10
		s.spawnFood()
		// Speed up slightly
		if s.moveDelay > 0.05 {
			s.moveDelay -= 0.005
		}
	} else {
		// Remove tail if no food eaten
		s.body = s.body[:len(s.body)-1]
	}
}

func (s *Snake) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(color.RGBA{R: 20, G: 20, B: 30, A: 255})

	// Draw grid lines (subtle)
	gridColor := color.RGBA{R: 30, G: 30, B: 40, A: 255}
	for x := 0; x < gridWidth; x++ {
		line := ebiten.NewImage(1, screenHeight)
		line.Fill(gridColor)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x*gridSize), 0)
		screen.DrawImage(line, op)
	}
	for y := 0; y < gridHeight; y++ {
		line := ebiten.NewImage(screenWidth, 1)
		line.Fill(gridColor)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, float64(y*gridSize))
		screen.DrawImage(line, op)
	}

	// Draw food
	food := ebiten.NewImage(gridSize-2, gridSize-2)
	food.Fill(color.RGBA{R: 255, G: 50, B: 50, A: 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.food.X*gridSize+1), float64(s.food.Y*gridSize+1))
	screen.DrawImage(food, op)

	// Draw snake
	for i, p := range s.body {
		segment := ebiten.NewImage(gridSize-2, gridSize-2)
		if i == 0 {
			// Head is brighter
			segment.Fill(color.RGBA{R: 100, G: 255, B: 100, A: 255})
		} else {
			segment.Fill(color.RGBA{R: 50, G: 200, B: 50, A: 255})
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(p.X*gridSize+1), float64(p.Y*gridSize+1))
		screen.DrawImage(segment, op)
	}

	// Draw score
	scorePanel := ebiten.NewImage(100, 30)
	scorePanel.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 180})
	screen.DrawImage(scorePanel, nil)

	// Game over overlay
	if s.gameOver {
		overlay := ebiten.NewImage(screenWidth, screenHeight)
		overlay.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 180})
		screen.DrawImage(overlay, nil)

		gameOverBox := ebiten.NewImage(200, 80)
		gameOverBox.Fill(color.RGBA{R: 255, G: 50, B: 50, A: 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(screenWidth-200)/2, float64(screenHeight-80)/2)
		screen.DrawImage(gameOverBox, op)
	}
}

func (s *Snake) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake - Framework Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(NewSnake()); err != nil {
		log.Fatal(err)
	}
}
