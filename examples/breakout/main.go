package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	paddleWidth  = 100
	paddleHeight = 15
	ballSize     = 10
	brickRows    = 5
	brickCols    = 10
	brickWidth   = 56
	brickHeight  = 20
	brickPadding = 4
	brickOffsetX = 20
	brickOffsetY = 60
)

// Brick represents a breakable brick.
type Brick struct {
	X, Y   float64
	Width  float64
	Height float64
	Color  color.RGBA
	Alive  bool
	Points int
}

// Ball represents the game ball.
type Ball struct {
	X, Y   float64
	VX, VY float64
	Size   float64
}

// Paddle represents the player paddle.
type Paddle struct {
	X, Y   float64
	Width  float64
	Height float64
}

// Breakout represents the breakout game.
type Breakout struct {
	paddle   *Paddle
	ball     *Ball
	bricks   []*Brick
	score    int
	lives    int
	gameOver bool
	victory  bool
	launched bool
}

// NewBreakout creates a new breakout game.
func NewBreakout() *Breakout {
	rand.Seed(time.Now().UnixNano())

	b := &Breakout{
		paddle: &Paddle{
			X:      float64(screenWidth-paddleWidth) / 2,
			Y:      float64(screenHeight) - 40,
			Width:  paddleWidth,
			Height: paddleHeight,
		},
		ball: &Ball{
			Size: ballSize,
		},
		lives:    3,
		launched: false,
	}

	b.resetBall()
	b.createBricks()

	return b
}

func (b *Breakout) resetBall() {
	b.ball.X = float64(screenWidth)/2 - ballSize/2
	b.ball.Y = b.paddle.Y - ballSize - 5
	b.ball.VX = 0
	b.ball.VY = 0
	b.launched = false
}

func (b *Breakout) launchBall() {
	if !b.launched {
		angle := (rand.Float64()*60 - 30) * math.Pi / 180 // -30 to 30 degrees
		speed := 5.0
		b.ball.VX = math.Sin(angle) * speed
		b.ball.VY = -math.Cos(angle) * speed
		b.launched = true
	}
}

func (b *Breakout) createBricks() {
	colors := []color.RGBA{
		{R: 255, G: 50, B: 50, A: 255},  // Red
		{R: 255, G: 150, B: 50, A: 255}, // Orange
		{R: 255, G: 255, B: 50, A: 255}, // Yellow
		{R: 50, G: 255, B: 50, A: 255},  // Green
		{R: 50, G: 150, B: 255, A: 255}, // Blue
	}
	points := []int{50, 40, 30, 20, 10}

	b.bricks = make([]*Brick, 0, brickRows*brickCols)

	for row := 0; row < brickRows; row++ {
		for col := 0; col < brickCols; col++ {
			brick := &Brick{
				X:      float64(brickOffsetX + col*(brickWidth+brickPadding)),
				Y:      float64(brickOffsetY + row*(brickHeight+brickPadding)),
				Width:  brickWidth,
				Height: brickHeight,
				Color:  colors[row],
				Points: points[row],
				Alive:  true,
			}
			b.bricks = append(b.bricks, brick)
		}
	}
}

func (b *Breakout) Update() error {
	if b.gameOver || b.victory {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			*b = *NewBreakout()
		}
		return nil
	}

	// Paddle movement with mouse
	mx, _ := ebiten.CursorPosition()
	b.paddle.X = float64(mx) - b.paddle.Width/2

	// Clamp paddle to screen
	if b.paddle.X < 0 {
		b.paddle.X = 0
	}
	if b.paddle.X > float64(screenWidth)-b.paddle.Width {
		b.paddle.X = float64(screenWidth) - b.paddle.Width
	}

	// Launch ball
	if !b.launched {
		b.ball.X = b.paddle.X + b.paddle.Width/2 - b.ball.Size/2
		b.ball.Y = b.paddle.Y - b.ball.Size - 2

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
			inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			b.launchBall()
		}
		return nil
	}

	// Update ball position
	b.ball.X += b.ball.VX
	b.ball.Y += b.ball.VY

	// Ball collision with walls
	if b.ball.X <= 0 || b.ball.X+b.ball.Size >= float64(screenWidth) {
		b.ball.VX = -b.ball.VX
		b.ball.X = clamp(b.ball.X, 0, float64(screenWidth)-b.ball.Size)
	}
	if b.ball.Y <= 0 {
		b.ball.VY = -b.ball.VY
		b.ball.Y = 0
	}

	// Ball falls below screen
	if b.ball.Y > float64(screenHeight) {
		b.lives--
		if b.lives <= 0 {
			b.gameOver = true
		} else {
			b.resetBall()
		}
		return nil
	}

	// Ball collision with paddle
	if b.ball.Y+b.ball.Size >= b.paddle.Y &&
		b.ball.Y <= b.paddle.Y+b.paddle.Height &&
		b.ball.X+b.ball.Size >= b.paddle.X &&
		b.ball.X <= b.paddle.X+b.paddle.Width &&
		b.ball.VY > 0 {

		// Calculate bounce angle based on hit position
		hitPos := (b.ball.X + b.ball.Size/2 - b.paddle.X) / b.paddle.Width
		angle := (hitPos - 0.5) * math.Pi * 0.6 // -54 to 54 degrees
		speed := math.Sqrt(b.ball.VX*b.ball.VX + b.ball.VY*b.ball.VY)

		b.ball.VX = math.Sin(angle) * speed
		b.ball.VY = -math.Abs(math.Cos(angle) * speed)
		b.ball.Y = b.paddle.Y - b.ball.Size
	}

	// Ball collision with bricks
	for _, brick := range b.bricks {
		if !brick.Alive {
			continue
		}

		if b.ball.X+b.ball.Size >= brick.X &&
			b.ball.X <= brick.X+brick.Width &&
			b.ball.Y+b.ball.Size >= brick.Y &&
			b.ball.Y <= brick.Y+brick.Height {

			brick.Alive = false
			b.score += brick.Points

			// Determine bounce direction
			overlapLeft := b.ball.X + b.ball.Size - brick.X
			overlapRight := brick.X + brick.Width - b.ball.X
			overlapTop := b.ball.Y + b.ball.Size - brick.Y
			overlapBottom := brick.Y + brick.Height - b.ball.Y

			minOverlapX := math.Min(overlapLeft, overlapRight)
			minOverlapY := math.Min(overlapTop, overlapBottom)

			if minOverlapX < minOverlapY {
				b.ball.VX = -b.ball.VX
			} else {
				b.ball.VY = -b.ball.VY
			}

			break // Only hit one brick per frame
		}
	}

	// Check victory
	allDestroyed := true
	for _, brick := range b.bricks {
		if brick.Alive {
			allDestroyed = false
			break
		}
	}
	if allDestroyed {
		b.victory = true
	}

	return nil
}

func (b *Breakout) Draw(screen *ebiten.Image) {
	// Background with gradient effect
	screen.Fill(color.RGBA{R: 15, G: 10, B: 30, A: 255})

	// Draw bricks
	for _, brick := range b.bricks {
		if !brick.Alive {
			continue
		}

		brickImg := ebiten.NewImage(int(brick.Width), int(brick.Height))
		brickImg.Fill(brick.Color)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(brick.X, brick.Y)
		screen.DrawImage(brickImg, op)
	}

	// Draw paddle with gradient
	paddleImg := ebiten.NewImage(int(b.paddle.Width), int(b.paddle.Height))
	paddleImg.Fill(color.RGBA{R: 100, G: 200, B: 255, A: 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.paddle.X, b.paddle.Y)
	screen.DrawImage(paddleImg, op)

	// Draw ball with glow effect
	ballImg := ebiten.NewImage(int(b.ball.Size), int(b.ball.Size))
	ballImg.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.ball.X, b.ball.Y)
	screen.DrawImage(ballImg, op)

	// Draw UI
	uiPanel := ebiten.NewImage(screenWidth, 30)
	uiPanel.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 180})
	screen.DrawImage(uiPanel, nil)

	// Draw lives as hearts
	for i := 0; i < b.lives; i++ {
		heart := ebiten.NewImage(15, 15)
		heart.Fill(color.RGBA{R: 255, G: 50, B: 100, A: 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(10+i*20), 7)
		screen.DrawImage(heart, op)
	}

	// Launch hint
	if !b.launched {
		hint := ebiten.NewImage(200, 30)
		hint.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 100})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(screenWidth-200)/2, float64(screenHeight)/2+50)
		screen.DrawImage(hint, op)
	}

	// Game over overlay
	if b.gameOver {
		overlay := ebiten.NewImage(screenWidth, screenHeight)
		overlay.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 180})
		screen.DrawImage(overlay, nil)

		gameOverBox := ebiten.NewImage(200, 80)
		gameOverBox.Fill(color.RGBA{R: 255, G: 50, B: 50, A: 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(screenWidth-200)/2, float64(screenHeight-80)/2)
		screen.DrawImage(gameOverBox, op)
	}

	// Victory overlay
	if b.victory {
		overlay := ebiten.NewImage(screenWidth, screenHeight)
		overlay.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 180})
		screen.DrawImage(overlay, nil)

		victoryBox := ebiten.NewImage(200, 80)
		victoryBox.Fill(color.RGBA{R: 50, G: 255, B: 100, A: 255})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(screenWidth-200)/2, float64(screenHeight-80)/2)
		screen.DrawImage(victoryBox, op)
	}
}

func (b *Breakout) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func clamp(v, min, max float64) float64 {
	return math.Max(min, math.Min(max, v))
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Breakout - Framework Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(NewBreakout()); err != nil {
		log.Fatal(err)
	}
}
