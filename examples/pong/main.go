package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	paddleWidth  = 15
	paddleHeight = 80
	ballSize     = 12
	paddleSpeed  = 6.0
	ballSpeed    = 5.0
)

// Paddle represents a player paddle.
type Paddle struct {
	X, Y   float64
	Width  float64
	Height float64
	Score  int
}

// Ball represents the game ball.
type Ball struct {
	X, Y   float64
	VX, VY float64
	Size   float64
}

// Pong represents the pong game.
type Pong struct {
	player1  *Paddle
	player2  *Paddle
	ball     *Ball
	paused   bool
	gameOver bool
	winScore int
}

// NewPong creates a new pong game.
func NewPong() *Pong {
	p := &Pong{
		player1: &Paddle{
			X:      30,
			Y:      float64(screenHeight-paddleHeight) / 2,
			Width:  paddleWidth,
			Height: paddleHeight,
		},
		player2: &Paddle{
			X:      float64(screenWidth) - 30 - paddleWidth,
			Y:      float64(screenHeight-paddleHeight) / 2,
			Width:  paddleWidth,
			Height: paddleHeight,
		},
		ball: &Ball{
			Size: ballSize,
		},
		winScore: 5,
	}
	p.resetBall(1)
	return p
}

func (p *Pong) resetBall(direction float64) {
	p.ball.X = float64(screenWidth) / 2
	p.ball.Y = float64(screenHeight) / 2
	p.ball.VX = ballSpeed * direction
	p.ball.VY = ballSpeed * 0.5
}

func (p *Pong) Update() error {
	if p.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			*p = *NewPong()
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		p.paused = !p.paused
	}

	if p.paused {
		return nil
	}

	// Player 1 controls (W/S)
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.player1.Y -= paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.player1.Y += paddleSpeed
	}

	// Player 2 controls (Up/Down arrows)
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.player2.Y -= paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.player2.Y += paddleSpeed
	}

	// Clamp paddles to screen
	p.player1.Y = clamp(p.player1.Y, 0, float64(screenHeight)-p.player1.Height)
	p.player2.Y = clamp(p.player2.Y, 0, float64(screenHeight)-p.player2.Height)

	// Update ball position
	p.ball.X += p.ball.VX
	p.ball.Y += p.ball.VY

	// Ball collision with top/bottom walls
	if p.ball.Y <= 0 || p.ball.Y+p.ball.Size >= float64(screenHeight) {
		p.ball.VY = -p.ball.VY
		p.ball.Y = clamp(p.ball.Y, 0, float64(screenHeight)-p.ball.Size)
	}

	// Ball collision with paddles
	if p.ball.X <= p.player1.X+p.player1.Width &&
		p.ball.Y+p.ball.Size >= p.player1.Y &&
		p.ball.Y <= p.player1.Y+p.player1.Height &&
		p.ball.VX < 0 {
		p.ball.VX = -p.ball.VX * 1.05 // Speed up slightly
		// Add angle based on where ball hits paddle
		relativeY := (p.ball.Y + p.ball.Size/2 - p.player1.Y) / p.player1.Height
		p.ball.VY = (relativeY - 0.5) * ballSpeed * 2
	}

	if p.ball.X+p.ball.Size >= p.player2.X &&
		p.ball.Y+p.ball.Size >= p.player2.Y &&
		p.ball.Y <= p.player2.Y+p.player2.Height &&
		p.ball.VX > 0 {
		p.ball.VX = -p.ball.VX * 1.05
		relativeY := (p.ball.Y + p.ball.Size/2 - p.player2.Y) / p.player2.Height
		p.ball.VY = (relativeY - 0.5) * ballSpeed * 2
	}

	// Clamp ball speed
	maxSpeed := 12.0
	p.ball.VX = clamp(p.ball.VX, -maxSpeed, maxSpeed)
	p.ball.VY = clamp(p.ball.VY, -maxSpeed, maxSpeed)

	// Scoring
	if p.ball.X < 0 {
		p.player2.Score++
		p.resetBall(-1)
	}
	if p.ball.X > float64(screenWidth) {
		p.player1.Score++
		p.resetBall(1)
	}

	// Check win condition
	if p.player1.Score >= p.winScore || p.player2.Score >= p.winScore {
		p.gameOver = true
	}

	return nil
}

func (p *Pong) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(color.RGBA{R: 10, G: 10, B: 20, A: 255})

	// Draw center line
	for y := 0; y < screenHeight; y += 20 {
		line := ebiten.NewImage(4, 10)
		line.Fill(color.RGBA{R: 50, G: 50, B: 60, A: 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(screenWidth)/2-2, float64(y))
		screen.DrawImage(line, op)
	}

	// Draw paddles
	paddle1 := ebiten.NewImage(int(p.player1.Width), int(p.player1.Height))
	paddle1.Fill(color.RGBA{R: 100, G: 150, B: 255, A: 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.player1.X, p.player1.Y)
	screen.DrawImage(paddle1, op)

	paddle2 := ebiten.NewImage(int(p.player2.Width), int(p.player2.Height))
	paddle2.Fill(color.RGBA{R: 255, G: 100, B: 100, A: 255})
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.player2.X, p.player2.Y)
	screen.DrawImage(paddle2, op)

	// Draw ball
	ball := ebiten.NewImage(int(p.ball.Size), int(p.ball.Size))
	ball.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.ball.X, p.ball.Y)
	screen.DrawImage(ball, op)

	// Draw scores
	p.drawScore(screen, p.player1.Score, screenWidth/4, 40, color.RGBA{R: 100, G: 150, B: 255, A: 255})
	p.drawScore(screen, p.player2.Score, screenWidth*3/4, 40, color.RGBA{R: 255, G: 100, B: 100, A: 255})

	// Pause overlay
	if p.paused {
		overlay := ebiten.NewImage(screenWidth, screenHeight)
		overlay.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 150})
		screen.DrawImage(overlay, nil)

		pauseBox := ebiten.NewImage(200, 50)
		pauseBox.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(screenWidth-200)/2, float64(screenHeight-50)/2)
		screen.DrawImage(pauseBox, op)
	}

	// Game over overlay
	if p.gameOver {
		overlay := ebiten.NewImage(screenWidth, screenHeight)
		overlay.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 180})
		screen.DrawImage(overlay, nil)

		var winColor color.RGBA
		if p.player1.Score >= p.winScore {
			winColor = color.RGBA{R: 100, G: 150, B: 255, A: 255}
		} else {
			winColor = color.RGBA{R: 255, G: 100, B: 100, A: 255}
		}

		winBox := ebiten.NewImage(200, 80)
		winBox.Fill(winColor)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(screenWidth-200)/2, float64(screenHeight-80)/2)
		screen.DrawImage(winBox, op)
	}
}

func (p *Pong) drawScore(screen *ebiten.Image, score int, x, y int, clr color.RGBA) {
	// Simple score display as rectangles
	size := 30
	scoreImg := ebiten.NewImage(size, size)
	scoreImg.Fill(clr)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x-size/2), float64(y))
	screen.DrawImage(scoreImg, op)

	// Draw score count as small boxes
	for i := 0; i < score; i++ {
		dot := ebiten.NewImage(10, 10)
		dot.Fill(clr)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x-25+i*12), float64(y+40))
		screen.DrawImage(dot, op)
	}
}

func (p *Pong) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func clamp(v, min, max float64) float64 {
	return math.Max(min, math.Min(max, v))
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong - Framework Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(NewPong()); err != nil {
		log.Fatal(err)
	}
}
