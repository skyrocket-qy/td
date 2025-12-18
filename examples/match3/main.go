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
	screenWidth  = 450
	screenHeight = 550
	gridCols     = 8
	gridRows     = 8
	cellSize     = 50
	gridOffsetX  = 25
	gridOffsetY  = 80
)

// GemType represents a gem color.
type GemType int

const (
	GemRed GemType = iota
	GemGreen
	GemBlue
	GemYellow
	GemPurple
	GemCount
)

// GemColors maps gem types to colors.
var GemColors = []color.RGBA{
	{R: 255, G: 60, B: 60, A: 255},  // Red
	{R: 60, G: 200, B: 60, A: 255},  // Green
	{R: 60, G: 100, B: 255, A: 255}, // Blue
	{R: 255, G: 220, B: 60, A: 255}, // Yellow
	{R: 180, G: 60, B: 200, A: 255}, // Purple
}

// Gem represents a single gem.
type Gem struct {
	Type    GemType
	X, Y    float64 // Animation position
	TargetY float64
	Falling bool
	Matched bool
}

// Game represents the match-3 game.
type Game struct {
	grid           [gridRows][gridCols]*Gem
	selectedX      int
	selectedY      int
	selected       bool
	score          int
	combo          int
	animating      bool
	swapping       bool
	swapX1, swapY1 int
	swapX2, swapY2 int
	swapProgress   float64
}

// NewGame creates a new game.
func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	g := &Game{
		selectedX: -1,
		selectedY: -1,
	}
	g.initGrid()
	return g
}

func (g *Game) initGrid() {
	for y := 0; y < gridRows; y++ {
		for x := 0; x < gridCols; x++ {
			g.grid[y][x] = &Gem{
				Type:    GemType(rand.Intn(int(GemCount))),
				X:       float64(x),
				Y:       float64(y),
				TargetY: float64(y),
			}
		}
	}
	// Remove initial matches
	for g.checkAndMarkMatches() {
		g.removeMatches()
		g.fillGrid()
	}
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0

	// Handle gem falling animation
	g.animating = false
	for y := 0; y < gridRows; y++ {
		for x := 0; x < gridCols; x++ {
			gem := g.grid[y][x]
			if gem == nil {
				continue
			}
			if gem.Falling {
				g.animating = true
				gem.Y += 8 * dt
				if gem.Y >= gem.TargetY {
					gem.Y = gem.TargetY
					gem.Falling = false
				}
			}
		}
	}

	// Handle swapping animation
	if g.swapping {
		g.animating = true
		g.swapProgress += dt * 4
		if g.swapProgress >= 1.0 {
			g.swapping = false
			g.swapProgress = 0

			// Check for matches after swap
			if !g.checkAndMarkMatches() {
				// Swap back if no match
				g.grid[g.swapY1][g.swapX1], g.grid[g.swapY2][g.swapX2] =
					g.grid[g.swapY2][g.swapX2], g.grid[g.swapY1][g.swapX1]
			} else {
				g.combo = 1
				g.processMatches()
			}
		}
	}

	if g.animating {
		return nil
	}

	// Click handling
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		gridX := (mx - gridOffsetX) / cellSize
		gridY := (my - gridOffsetY) / cellSize

		if gridX >= 0 && gridX < gridCols && gridY >= 0 && gridY < gridRows {
			if !g.selected {
				g.selected = true
				g.selectedX = gridX
				g.selectedY = gridY
			} else {
				// Check if adjacent
				dx := abs(gridX - g.selectedX)
				dy := abs(gridY - g.selectedY)
				if (dx == 1 && dy == 0) || (dx == 0 && dy == 1) {
					g.startSwap(g.selectedX, g.selectedY, gridX, gridY)
				}
				g.selected = false
			}
		} else {
			g.selected = false
		}
	}

	return nil
}

func (g *Game) startSwap(x1, y1, x2, y2 int) {
	g.swapping = true
	g.swapX1, g.swapY1 = x1, y1
	g.swapX2, g.swapY2 = x2, y2
	g.swapProgress = 0

	// Actually swap in grid
	g.grid[y1][x1], g.grid[y2][x2] = g.grid[y2][x2], g.grid[y1][x1]
}

func (g *Game) checkAndMarkMatches() bool {
	hasMatch := false

	// Check horizontal matches
	for y := 0; y < gridRows; y++ {
		for x := 0; x < gridCols-2; x++ {
			if g.grid[y][x] == nil {
				continue
			}
			t := g.grid[y][x].Type
			count := 1
			for nx := x + 1; nx < gridCols && g.grid[y][nx] != nil && g.grid[y][nx].Type == t; nx++ {
				count++
			}
			if count >= 3 {
				hasMatch = true
				for i := 0; i < count; i++ {
					g.grid[y][x+i].Matched = true
				}
			}
		}
	}

	// Check vertical matches
	for x := 0; x < gridCols; x++ {
		for y := 0; y < gridRows-2; y++ {
			if g.grid[y][x] == nil {
				continue
			}
			t := g.grid[y][x].Type
			count := 1
			for ny := y + 1; ny < gridRows && g.grid[ny][x] != nil && g.grid[ny][x].Type == t; ny++ {
				count++
			}
			if count >= 3 {
				hasMatch = true
				for i := 0; i < count; i++ {
					g.grid[y+i][x].Matched = true
				}
			}
		}
	}

	return hasMatch
}

func (g *Game) processMatches() {
	g.removeMatches()
	g.dropGems()
	g.fillGrid()

	// Chain reaction
	if g.checkAndMarkMatches() {
		g.combo++
	}
}

func (g *Game) removeMatches() {
	for y := 0; y < gridRows; y++ {
		for x := 0; x < gridCols; x++ {
			if g.grid[y][x] != nil && g.grid[y][x].Matched {
				g.score += 10 * g.combo
				g.grid[y][x] = nil
			}
		}
	}
}

func (g *Game) dropGems() {
	for x := 0; x < gridCols; x++ {
		// Move gems down
		writePos := gridRows - 1
		for y := gridRows - 1; y >= 0; y-- {
			if g.grid[y][x] != nil {
				if writePos != y {
					g.grid[writePos][x] = g.grid[y][x]
					g.grid[writePos][x].TargetY = float64(writePos)
					g.grid[writePos][x].Falling = true
					g.grid[y][x] = nil
				}
				writePos--
			}
		}
	}
}

func (g *Game) fillGrid() {
	for x := 0; x < gridCols; x++ {
		emptyCount := 0
		for y := 0; y < gridRows; y++ {
			if g.grid[y][x] == nil {
				emptyCount++
			}
		}

		fillY := 0
		for y := 0; y < gridRows; y++ {
			if g.grid[y][x] == nil {
				g.grid[y][x] = &Gem{
					Type:    GemType(rand.Intn(int(GemCount))),
					X:       float64(x),
					Y:       float64(-emptyCount + fillY),
					TargetY: float64(y),
					Falling: true,
				}
				fillY++
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(color.RGBA{R: 40, G: 30, B: 50, A: 255})

	// Header
	vector.DrawFilledRect(screen, 0, 0, screenWidth, 70, color.RGBA{R: 60, G: 50, B: 70, A: 255}, false)
	ebitenutil.DebugPrintAt(screen, "Match 3", 20, 15)
	ebitenutil.DebugPrintAt(screen, "Score: "+formatInt(g.score), 20, 40)
	if g.combo > 1 {
		ebitenutil.DebugPrintAt(screen, "Combo x"+formatInt(g.combo), 150, 40)
	}

	// Grid background
	vector.DrawFilledRect(screen, float32(gridOffsetX), float32(gridOffsetY),
		float32(gridCols*cellSize), float32(gridRows*cellSize),
		color.RGBA{R: 30, G: 25, B: 40, A: 255}, false)

	// Draw gems
	for y := 0; y < gridRows; y++ {
		for x := 0; x < gridCols; x++ {
			gem := g.grid[y][x]
			if gem == nil {
				continue
			}

			drawX := float32(gridOffsetX + x*cellSize + cellSize/2)
			drawY := float32(gridOffsetY + int(gem.Y*float64(cellSize)) + cellSize/2)

			// Selected highlight
			if g.selected && x == g.selectedX && y == g.selectedY {
				vector.DrawFilledRect(screen, drawX-float32(cellSize/2)+2, drawY-float32(cellSize/2)+2,
					float32(cellSize-4), float32(cellSize-4), color.RGBA{R: 255, G: 255, B: 255, A: 100}, false)
			}

			// Gem
			gemColor := GemColors[gem.Type]
			vector.DrawFilledCircle(screen, drawX, drawY, float32(cellSize/2-4), gemColor, false)

			// Shine effect
			shineColor := color.RGBA{R: 255, G: 255, B: 255, A: 80}
			vector.DrawFilledCircle(screen, drawX-5, drawY-5, 6, shineColor, false)
		}
	}

	// Grid lines
	for i := 0; i <= gridCols; i++ {
		x := float32(gridOffsetX + i*cellSize)
		vector.DrawFilledRect(screen, x, float32(gridOffsetY), 1, float32(gridRows*cellSize), color.RGBA{R: 60, G: 50, B: 70, A: 255}, false)
	}
	for i := 0; i <= gridRows; i++ {
		y := float32(gridOffsetY + i*cellSize)
		vector.DrawFilledRect(screen, float32(gridOffsetX), y, float32(gridCols*cellSize), 1, color.RGBA{R: 60, G: 50, B: 70, A: 255}, false)
	}

	// Instructions
	ebitenutil.DebugPrintAt(screen, "Click two adjacent gems to swap", 100, screenHeight-25)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func formatInt(n int) string {
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

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Match 3")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
