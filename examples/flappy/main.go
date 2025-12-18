package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 400
	screenHeight = 600
	gravity      = 0.5
	jumpForce    = -8
	pipeWidth    = 60
	pipeGap      = 150
	pipeSpeed    = 3
	birdSize     = 30
)

// Bird represents the player.
type Bird struct {
	X, Y      float64
	VelocityY float64
	Rotation  float64
}

// Pipe represents an obstacle.
type Pipe struct {
	X      float64
	GapY   float64
	Passed bool
}

// Game represents the flappy bird game.
type Game struct {
	bird      *Bird
	pipes     []*Pipe
	score     int
	highscore int
	gameOver  bool
	started   bool
	pipeTimer float64
}

// NewGame creates a new game.
func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	return &Game{
		bird: &Bird{
			X: 100,
			Y: float64(screenHeight) / 2,
		},
		pipes: make([]*Pipe, 0),
	}
}

func (g *Game) reset() {
	g.bird = &Bird{
		X: 100,
		Y: float64(screenHeight) / 2,
	}
	g.pipes = make([]*Pipe, 0)
	g.score = 0
	g.gameOver = false
	g.started = false
	g.pipeTimer = 0
}

func (g *Game) spawnPipe() {
	// Random gap position
	minGap := float64(pipeGap/2 + 50)
	maxGap := float64(screenHeight - pipeGap/2 - 50)
	gapY := minGap + rand.Float64()*(maxGap-minGap)

	g.pipes = append(g.pipes, &Pipe{
		X:    float64(screenWidth),
		GapY: gapY,
	})
}

func (g *Game) Update() error {
	// Start game / jump / restart
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.gameOver {
			if g.score > g.highscore {
				g.highscore = g.score
			}
			g.reset()
		} else if !g.started {
			g.started = true
			g.bird.VelocityY = jumpForce
		} else {
			g.bird.VelocityY = jumpForce
		}
	}

	if !g.started || g.gameOver {
		return nil
	}

	// Update bird physics
	g.bird.VelocityY += gravity
	g.bird.Y += g.bird.VelocityY

	// Rotation based on velocity
	g.bird.Rotation = g.bird.VelocityY * 3
	if g.bird.Rotation > 90 {
		g.bird.Rotation = 90
	}
	if g.bird.Rotation < -30 {
		g.bird.Rotation = -30
	}

	// Check floor/ceiling collision
	if g.bird.Y < 0 || g.bird.Y > float64(screenHeight-birdSize) {
		g.gameOver = true
		return nil
	}

	// Spawn pipes
	g.pipeTimer += 1.0 / 60.0
	if g.pipeTimer >= 1.5 {
		g.spawnPipe()
		g.pipeTimer = 0
	}

	// Update pipes
	for i := len(g.pipes) - 1; i >= 0; i-- {
		pipe := g.pipes[i]
		pipe.X -= pipeSpeed

		// Remove off-screen pipes
		if pipe.X < -pipeWidth {
			g.pipes = append(g.pipes[:i], g.pipes[i+1:]...)
			continue
		}

		// Score when passing pipe
		if !pipe.Passed && pipe.X+pipeWidth < g.bird.X {
			pipe.Passed = true
			g.score++
		}

		// Collision detection
		if g.checkCollision(pipe) {
			g.gameOver = true
		}
	}

	return nil
}

func (g *Game) checkCollision(pipe *Pipe) bool {
	birdLeft := g.bird.X
	birdRight := g.bird.X + birdSize
	birdTop := g.bird.Y
	birdBottom := g.bird.Y + birdSize

	pipeLeft := pipe.X
	pipeRight := pipe.X + pipeWidth
	gapTop := pipe.GapY - pipeGap/2
	gapBottom := pipe.GapY + pipeGap/2

	// Check if bird is within pipe X range
	if birdRight > pipeLeft && birdLeft < pipeRight {
		// Check if bird hits top or bottom pipe
		if birdTop < gapTop || birdBottom > gapBottom {
			return true
		}
	}

	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Sky gradient
	for y := 0; y < screenHeight; y++ {
		t := float64(y) / float64(screenHeight)
		r := uint8(100 - t*50)
		gCol := uint8(180 - t*80)
		b := uint8(255 - t*55)
		vector.DrawFilledRect(screen, 0, float32(y), float32(screenWidth), 1, color.RGBA{R: r, G: gCol, B: b, A: 255}, false)
	}

	// Ground
	vector.DrawFilledRect(screen, 0, float32(screenHeight-50), float32(screenWidth), 50, color.RGBA{R: 139, G: 119, B: 101, A: 255}, false)
	vector.DrawFilledRect(screen, 0, float32(screenHeight-50), float32(screenWidth), 5, color.RGBA{R: 34, G: 139, B: 34, A: 255}, false)

	// Pipes
	for _, pipe := range g.pipes {
		g.drawPipe(screen, pipe)
	}

	// Bird
	g.drawBird(screen)

	// Score
	ebitenutil.DebugPrintAt(screen, "Score: "+string(rune('0'+g.score%10)), screenWidth/2-30, 20)
	if g.score >= 10 {
		ebitenutil.DebugPrintAt(screen, string(rune('0'+g.score/10)), screenWidth/2-40, 20)
	}

	// Game state messages
	if !g.started {
		g.drawMessage(screen, "TAP TO START", color.RGBA{R: 255, G: 255, B: 255, A: 255})
	} else if g.gameOver {
		g.drawMessage(screen, "GAME OVER", color.RGBA{R: 255, G: 50, B: 50, A: 255})
		ebitenutil.DebugPrintAt(screen, "Tap to restart", screenWidth/2-50, screenHeight/2+40)
	}

	// Score display (centered at top)
	scoreText := formatScore(g.score)
	ebitenutil.DebugPrintAt(screen, scoreText, screenWidth/2-len(scoreText)*3, 30)
}

func formatScore(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}

func (g *Game) drawBird(screen *ebiten.Image) {
	// Bird body (yellow)
	vector.DrawFilledCircle(screen, float32(g.bird.X+birdSize/2), float32(g.bird.Y+birdSize/2), birdSize/2, color.RGBA{R: 255, G: 220, B: 50, A: 255}, false)

	// Bird eye
	vector.DrawFilledCircle(screen, float32(g.bird.X+birdSize*0.7), float32(g.bird.Y+birdSize*0.3), 5, color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
	vector.DrawFilledCircle(screen, float32(g.bird.X+birdSize*0.75), float32(g.bird.Y+birdSize*0.35), 2, color.RGBA{R: 0, G: 0, B: 0, A: 255}, false)

	// Bird beak
	vector.DrawFilledRect(screen, float32(g.bird.X+birdSize*0.8), float32(g.bird.Y+birdSize*0.45), 10, 6, color.RGBA{R: 255, G: 150, B: 0, A: 255}, false)

	// Wing
	wingY := g.bird.Y + birdSize*0.5
	if g.bird.VelocityY < 0 {
		wingY -= 5 // Wing up when jumping
	}
	vector.DrawFilledCircle(screen, float32(g.bird.X+birdSize*0.3), float32(wingY), 8, color.RGBA{R: 255, G: 180, B: 50, A: 255}, false)
}

func (g *Game) drawPipe(screen *ebiten.Image, pipe *Pipe) {
	pipeColor := color.RGBA{R: 50, G: 180, B: 50, A: 255}
	pipeEdge := color.RGBA{R: 30, G: 140, B: 30, A: 255}

	gapTop := pipe.GapY - pipeGap/2
	gapBottom := pipe.GapY + pipeGap/2

	// Top pipe
	vector.DrawFilledRect(screen, float32(pipe.X), 0, float32(pipeWidth), float32(gapTop), pipeColor, false)
	vector.DrawFilledRect(screen, float32(pipe.X-5), float32(gapTop-30), float32(pipeWidth+10), 30, pipeEdge, false)

	// Bottom pipe
	vector.DrawFilledRect(screen, float32(pipe.X), float32(gapBottom), float32(pipeWidth), float32(screenHeight-50)-float32(gapBottom), pipeColor, false)
	vector.DrawFilledRect(screen, float32(pipe.X-5), float32(gapBottom), float32(pipeWidth+10), 30, pipeEdge, false)
}

func (g *Game) drawMessage(screen *ebiten.Image, text string, clr color.RGBA) {
	// Semi-transparent background
	boxW := float32(200)
	boxH := float32(60)
	boxX := float32(screenWidth-200) / 2
	boxY := float32(screenHeight-60) / 2

	vector.DrawFilledRect(screen, boxX, boxY, boxW, boxH, color.RGBA{R: 0, G: 0, B: 0, A: 180}, false)
	vector.StrokeRect(screen, boxX, boxY, boxW, boxH, 2, clr, false)

	ebitenutil.DebugPrintAt(screen, text, int(boxX)+40, int(boxY)+23)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Flappy Bird Clone")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
