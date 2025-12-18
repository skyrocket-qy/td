package main

import (
	"fmt"
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
	gridSize     = 4
	tileSize     = 90
	tilePadding  = 10
	gridOffsetX  = 25
	gridOffsetY  = 120
)

// TileColors maps values to colors.
var TileColors = map[int]color.RGBA{
	0:    {R: 205, G: 193, B: 180, A: 255},
	2:    {R: 238, G: 228, B: 218, A: 255},
	4:    {R: 237, G: 224, B: 200, A: 255},
	8:    {R: 242, G: 177, B: 121, A: 255},
	16:   {R: 245, G: 149, B: 99, A: 255},
	32:   {R: 246, G: 124, B: 95, A: 255},
	64:   {R: 246, G: 94, B: 59, A: 255},
	128:  {R: 237, G: 207, B: 114, A: 255},
	256:  {R: 237, G: 204, B: 97, A: 255},
	512:  {R: 237, G: 200, B: 80, A: 255},
	1024: {R: 237, G: 197, B: 63, A: 255},
	2048: {R: 237, G: 194, B: 46, A: 255},
}

// Game represents the 2048 game.
type Game struct {
	grid      [gridSize][gridSize]int
	score     int
	highscore int
	gameOver  bool
	won       bool
	moved     bool
}

// NewGame creates a new game.
func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	g := &Game{}
	g.spawnTile()
	g.spawnTile()
	return g
}

func (g *Game) spawnTile() {
	empty := make([][2]int, 0)
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			if g.grid[i][j] == 0 {
				empty = append(empty, [2]int{i, j})
			}
		}
	}

	if len(empty) == 0 {
		return
	}

	pos := empty[rand.Intn(len(empty))]
	value := 2
	if rand.Float64() < 0.1 {
		value = 4
	}
	g.grid[pos[0]][pos[1]] = value
}

func (g *Game) reset() {
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			g.grid[i][j] = 0
		}
	}
	g.score = 0
	g.gameOver = false
	g.won = false
	g.spawnTile()
	g.spawnTile()
}

func (g *Game) Update() error {
	if g.gameOver || g.won {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			inpututil.IsKeyJustPressed(ebiten.KeyR) {
			if g.score > g.highscore {
				g.highscore = g.score
			}
			g.reset()
		}
		return nil
	}

	g.moved = false

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.moveLeft()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.moveRight()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.moveUp()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.moveDown()
	}

	if g.moved {
		g.spawnTile()
		g.checkGameOver()
	}

	return nil
}

func (g *Game) moveLeft() {
	for i := 0; i < gridSize; i++ {
		row := g.grid[i]
		newRow := g.slideAndMerge(row[:])
		for j := 0; j < gridSize; j++ {
			if g.grid[i][j] != newRow[j] {
				g.moved = true
			}
			g.grid[i][j] = newRow[j]
		}
	}
}

func (g *Game) moveRight() {
	for i := 0; i < gridSize; i++ {
		row := make([]int, gridSize)
		for j := 0; j < gridSize; j++ {
			row[j] = g.grid[i][gridSize-1-j]
		}
		newRow := g.slideAndMerge(row)
		for j := 0; j < gridSize; j++ {
			if g.grid[i][gridSize-1-j] != newRow[j] {
				g.moved = true
			}
			g.grid[i][gridSize-1-j] = newRow[j]
		}
	}
}

func (g *Game) moveUp() {
	for j := 0; j < gridSize; j++ {
		col := make([]int, gridSize)
		for i := 0; i < gridSize; i++ {
			col[i] = g.grid[i][j]
		}
		newCol := g.slideAndMerge(col)
		for i := 0; i < gridSize; i++ {
			if g.grid[i][j] != newCol[i] {
				g.moved = true
			}
			g.grid[i][j] = newCol[i]
		}
	}
}

func (g *Game) moveDown() {
	for j := 0; j < gridSize; j++ {
		col := make([]int, gridSize)
		for i := 0; i < gridSize; i++ {
			col[i] = g.grid[gridSize-1-i][j]
		}
		newCol := g.slideAndMerge(col)
		for i := 0; i < gridSize; i++ {
			if g.grid[gridSize-1-i][j] != newCol[i] {
				g.moved = true
			}
			g.grid[gridSize-1-i][j] = newCol[i]
		}
	}
}

func (g *Game) slideAndMerge(line []int) []int {
	// Remove zeros
	nonZero := make([]int, 0)
	for _, v := range line {
		if v != 0 {
			nonZero = append(nonZero, v)
		}
	}

	// Merge adjacent equal tiles
	merged := make([]int, 0)
	skip := false
	for i := 0; i < len(nonZero); i++ {
		if skip {
			skip = false
			continue
		}
		if i+1 < len(nonZero) && nonZero[i] == nonZero[i+1] {
			merged = append(merged, nonZero[i]*2)
			g.score += nonZero[i] * 2
			if nonZero[i]*2 == 2048 {
				g.won = true
			}
			skip = true
		} else {
			merged = append(merged, nonZero[i])
		}
	}

	// Pad with zeros
	result := make([]int, gridSize)
	for i, v := range merged {
		result[i] = v
	}
	return result
}

func (g *Game) checkGameOver() {
	// Check for empty cells
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			if g.grid[i][j] == 0 {
				return
			}
		}
	}

	// Check for possible merges
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			if j+1 < gridSize && g.grid[i][j] == g.grid[i][j+1] {
				return
			}
			if i+1 < gridSize && g.grid[i][j] == g.grid[i+1][j] {
				return
			}
		}
	}

	g.gameOver = true
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(color.RGBA{R: 250, G: 248, B: 239, A: 255})

	// Header
	vector.DrawFilledRect(screen, 0, 0, screenWidth, 100, color.RGBA{R: 187, G: 173, B: 160, A: 255}, false)

	// Title
	ebitenutil.DebugPrintAt(screen, "2048", 20, 30)

	// Score boxes
	g.drawScoreBox(screen, 200, 20, "SCORE", g.score)
	g.drawScoreBox(screen, 320, 20, "BEST", g.highscore)

	// Grid background
	gridW := float32(gridSize*tileSize + (gridSize+1)*tilePadding)
	vector.DrawFilledRect(screen, float32(gridOffsetX), float32(gridOffsetY), gridW, gridW, color.RGBA{R: 187, G: 173, B: 160, A: 255}, false)

	// Draw tiles
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			g.drawTile(screen, i, j)
		}
	}

	// Game over / win overlay
	if g.gameOver {
		g.drawOverlay(screen, "Game Over!", color.RGBA{R: 119, G: 110, B: 101, A: 200})
	} else if g.won {
		g.drawOverlay(screen, "You Win!", color.RGBA{R: 237, G: 194, B: 46, A: 200})
	}

	// Controls hint
	ebitenutil.DebugPrintAt(screen, "Arrow Keys / WASD to move | R to restart", 60, screenHeight-25)
}

func (g *Game) drawScoreBox(screen *ebiten.Image, x, y int, label string, value int) {
	vector.DrawFilledRect(screen, float32(x), float32(y), 100, 55, color.RGBA{R: 187, G: 173, B: 160, A: 255}, false)
	ebitenutil.DebugPrintAt(screen, label, x+35, y+5)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", value), x+35, y+25)
}

func (g *Game) drawTile(screen *ebiten.Image, row, col int) {
	value := g.grid[row][col]
	x := float32(gridOffsetX + tilePadding + col*(tileSize+tilePadding))
	y := float32(gridOffsetY + tilePadding + row*(tileSize+tilePadding))

	// Tile background
	tileColor := TileColors[value]
	if _, ok := TileColors[value]; !ok {
		tileColor = color.RGBA{R: 60, G: 58, B: 50, A: 255} // For values > 2048
	}
	vector.DrawFilledRect(screen, x, y, tileSize, tileSize, tileColor, false)

	// Tile number
	if value > 0 {
		textColor := color.RGBA{R: 119, G: 110, B: 101, A: 255}
		if value >= 8 {
			textColor = color.RGBA{R: 249, G: 246, B: 242, A: 255}
		}
		_ = textColor

		// Center text (simple approach)
		text := fmt.Sprintf("%d", value)
		textX := int(x) + (tileSize-len(text)*8)/2
		textY := int(y) + tileSize/2 - 6
		ebitenutil.DebugPrintAt(screen, text, textX, textY)
	}
}

func (g *Game) drawOverlay(screen *ebiten.Image, message string, bgColor color.RGBA) {
	vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, bgColor, false)

	boxW := float32(250)
	boxH := float32(100)
	boxX := float32(screenWidth-250) / 2
	boxY := float32(screenHeight-100) / 2

	vector.DrawFilledRect(screen, boxX, boxY, boxW, boxH, color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
	ebitenutil.DebugPrintAt(screen, message, int(boxX)+80, int(boxY)+30)
	ebitenutil.DebugPrintAt(screen, "Press SPACE to restart", int(boxX)+40, int(boxY)+60)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("2048")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
