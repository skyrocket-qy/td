package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 600
	playerSpeed  = 3
)

// WeaponType represents weapon types.
type WeaponType int

const (
	WeaponOrb WeaponType = iota
	WeaponWhip
	WeaponAura
	WeaponProjectile
)

// Weapon represents an auto-attacking weapon.
type Weapon struct {
	Type     WeaponType
	Level    int
	Damage   int
	Cooldown float64
	Timer    float64
	Range    float64
	Count    int // Number of projectiles
}

// Projectile represents a weapon projectile.
type Projectile struct {
	X, Y     float64
	VX, VY   float64
	Damage   int
	Lifetime float64
	Radius   float64
	Piercing int
	HitList  map[*Enemy]bool
}

// Enemy represents an enemy.
type Enemy struct {
	X, Y     float64
	HP       int
	MaxHP    int
	Speed    float64
	Damage   int
	XP       int
	Dead     bool
	HitFlash float64
}

// XPGem represents dropped experience.
type XPGem struct {
	X, Y   float64
	Value  int
	Magnet bool
}

// DamageNumber represents floating damage text.
type DamageNumber struct {
	X, Y  float64
	Value int
	Timer float64
}

// Player represents the player.
type Player struct {
	X, Y    float64
	HP      int
	MaxHP   int
	XP      int
	Level   int
	Weapons []*Weapon
}

// Game represents the survivor game.
type Game struct {
	player        *Player
	enemies       []*Enemy
	projectiles   []*Projectile
	xpGems        []*XPGem
	damageNumbers []*DamageNumber

	gameTime       float64
	spawnTimer     float64
	spawnRate      float64
	killCount      int
	gameOver       bool
	paused         bool
	levelingUp     bool
	upgradeOptions []UpgradeOption

	cameraX, cameraY float64
}

// UpgradeOption represents a level-up choice.
type UpgradeOption struct {
	Name        string
	Description string
	Apply       func(*Game)
}

// NewGame creates a new game.
func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	g := &Game{
		player: &Player{
			X:  float64(screenWidth) / 2,
			Y:  float64(screenHeight) / 2,
			HP: 100, MaxHP: 100,
			Level: 1,
			Weapons: []*Weapon{
				{Type: WeaponOrb, Level: 1, Damage: 10, Cooldown: 1.0, Range: 100, Count: 1},
			},
		},
		enemies:     make([]*Enemy, 0),
		projectiles: make([]*Projectile, 0),
		xpGems:      make([]*XPGem, 0),
		spawnRate:   2.0,
	}
	return g
}

func (g *Game) reset() {
	g.player = &Player{
		X: 0, Y: 0,
		HP: 100, MaxHP: 100,
		Level: 1,
		Weapons: []*Weapon{
			{Type: WeaponOrb, Level: 1, Damage: 10, Cooldown: 1.0, Range: 100, Count: 1},
		},
	}
	g.enemies = make([]*Enemy, 0)
	g.projectiles = make([]*Projectile, 0)
	g.xpGems = make([]*XPGem, 0)
	g.damageNumbers = make([]*DamageNumber, 0)
	g.gameTime = 0
	g.spawnTimer = 0
	g.spawnRate = 2.0
	g.killCount = 0
	g.gameOver = false
	g.levelingUp = false
}

func (g *Game) Update() error {
	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.reset()
		}
		return nil
	}

	// Level up menu
	if g.levelingUp {
		for i := 0; i < len(g.upgradeOptions) && i < 3; i++ {
			if inpututil.IsKeyJustPressed(ebiten.Key(int(ebiten.Key1) + i)) {
				g.upgradeOptions[i].Apply(g)
				g.levelingUp = false
				break
			}
		}
		return nil
	}

	dt := 1.0 / 60.0
	g.gameTime += dt

	// Player movement
	dx, dy := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx = 1
	}

	// Normalize diagonal movement
	if dx != 0 && dy != 0 {
		dx *= 0.707
		dy *= 0.707
	}
	g.player.X += dx * playerSpeed
	g.player.Y += dy * playerSpeed

	// Update camera
	g.cameraX = g.player.X - float64(screenWidth)/2
	g.cameraY = g.player.Y - float64(screenHeight)/2

	// Spawn enemies
	g.spawnTimer += dt
	if g.spawnTimer >= g.spawnRate {
		g.spawnEnemy()
		g.spawnTimer = 0
		// Speed up spawning over time
		g.spawnRate = 2.0 - g.gameTime*0.01
		if g.spawnRate < 0.3 {
			g.spawnRate = 0.3
		}
	}

	// Update weapons
	for _, w := range g.player.Weapons {
		w.Timer += dt
		if w.Timer >= w.Cooldown {
			g.fireWeapon(w)
			w.Timer = 0
		}
	}

	// Update projectiles
	for i := len(g.projectiles) - 1; i >= 0; i-- {
		p := g.projectiles[i]
		p.X += p.VX
		p.Y += p.VY
		p.Lifetime -= dt

		// Check enemy collisions
		for _, e := range g.enemies {
			if e.Dead {
				continue
			}
			if p.HitList[e] {
				continue
			}
			dist := math.Sqrt(math.Pow(p.X-e.X, 2) + math.Pow(p.Y-e.Y, 2))
			if dist < p.Radius+20 {
				e.HP -= p.Damage
				e.HitFlash = 0.1
				p.HitList[e] = true
				p.Piercing--

				g.damageNumbers = append(g.damageNumbers, &DamageNumber{
					X: e.X, Y: e.Y - 30, Value: p.Damage, Timer: 0.8,
				})

				if e.HP <= 0 {
					e.Dead = true
					g.killCount++
					g.xpGems = append(g.xpGems, &XPGem{X: e.X, Y: e.Y, Value: e.XP})
				}

				if p.Piercing <= 0 {
					p.Lifetime = 0
					break
				}
			}
		}

		if p.Lifetime <= 0 {
			g.projectiles = append(g.projectiles[:i], g.projectiles[i+1:]...)
		}
	}

	// Update enemies
	for _, e := range g.enemies {
		if e.Dead {
			continue
		}
		e.HitFlash -= dt

		// Move toward player
		dx := g.player.X - e.X
		dy := g.player.Y - e.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > 0 {
			e.X += (dx / dist) * e.Speed
			e.Y += (dy / dist) * e.Speed
		}

		// Damage player
		if dist < 25 {
			g.player.HP -= e.Damage
			if g.player.HP <= 0 {
				g.gameOver = true
			}
		}
	}

	// Remove dead enemies
	for i := len(g.enemies) - 1; i >= 0; i-- {
		if g.enemies[i].Dead {
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
		}
	}

	// Collect XP gems
	magnetRange := 100.0
	for i := len(g.xpGems) - 1; i >= 0; i-- {
		gem := g.xpGems[i]
		dx := g.player.X - gem.X
		dy := g.player.Y - gem.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		// Magnet effect
		if dist < magnetRange || gem.Magnet {
			gem.Magnet = true
			speed := 8.0
			gem.X += (dx / dist) * speed
			gem.Y += (dy / dist) * speed
		}

		// Collect
		if dist < 20 {
			g.player.XP += gem.Value
			g.xpGems = append(g.xpGems[:i], g.xpGems[i+1:]...)

			// Check level up
			xpNeeded := g.player.Level * 20
			if g.player.XP >= xpNeeded {
				g.player.XP -= xpNeeded
				g.player.Level++
				g.showLevelUpMenu()
			}
		}
	}

	// Update damage numbers
	for i := len(g.damageNumbers) - 1; i >= 0; i-- {
		d := g.damageNumbers[i]
		d.Y -= 30 * dt
		d.Timer -= dt
		if d.Timer <= 0 {
			g.damageNumbers = append(g.damageNumbers[:i], g.damageNumbers[i+1:]...)
		}
	}

	return nil
}

func (g *Game) spawnEnemy() {
	// Spawn at edge of screen
	angle := rand.Float64() * math.Pi * 2
	dist := float64(screenWidth)/2 + 50

	x := g.player.X + math.Cos(angle)*dist
	y := g.player.Y + math.Sin(angle)*dist

	hp := 20 + int(g.gameTime*0.5)
	speed := 1.0 + g.gameTime*0.01
	if speed > 3 {
		speed = 3
	}

	g.enemies = append(g.enemies, &Enemy{
		X: x, Y: y,
		HP: hp, MaxHP: hp,
		Speed:  speed,
		Damage: 1,
		XP:     5 + int(g.gameTime*0.1),
	})
}

func (g *Game) fireWeapon(w *Weapon) {
	switch w.Type {
	case WeaponOrb:
		// Orbiting projectiles
		for i := 0; i < w.Count; i++ {
			angle := g.gameTime*2 + float64(i)*(2*math.Pi/float64(w.Count))
			g.projectiles = append(g.projectiles, &Projectile{
				X:  g.player.X + math.Cos(angle)*w.Range,
				Y:  g.player.Y + math.Sin(angle)*w.Range,
				VX: 0, VY: 0,
				Damage:   w.Damage,
				Lifetime: 0.3,
				Radius:   15,
				Piercing: 3,
				HitList:  make(map[*Enemy]bool),
			})
		}

	case WeaponProjectile:
		// Fire at nearest enemy
		var nearest *Enemy
		minDist := 500.0
		for _, e := range g.enemies {
			if e.Dead {
				continue
			}
			dist := math.Sqrt(math.Pow(g.player.X-e.X, 2) + math.Pow(g.player.Y-e.Y, 2))
			if dist < minDist {
				minDist = dist
				nearest = e
			}
		}
		if nearest != nil {
			dx := nearest.X - g.player.X
			dy := nearest.Y - g.player.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			speed := 8.0
			for i := 0; i < w.Count; i++ {
				spread := float64(i-w.Count/2) * 0.2
				g.projectiles = append(g.projectiles, &Projectile{
					X:        g.player.X,
					Y:        g.player.Y,
					VX:       (dx/dist)*speed + spread,
					VY:       (dy/dist)*speed + spread,
					Damage:   w.Damage,
					Lifetime: 2.0,
					Radius:   8,
					Piercing: 1,
					HitList:  make(map[*Enemy]bool),
				})
			}
		}

	case WeaponAura:
		// Damage all enemies in range
		for _, e := range g.enemies {
			if e.Dead {
				continue
			}
			dist := math.Sqrt(math.Pow(g.player.X-e.X, 2) + math.Pow(g.player.Y-e.Y, 2))
			if dist < w.Range {
				e.HP -= w.Damage
				e.HitFlash = 0.1
				g.damageNumbers = append(g.damageNumbers, &DamageNumber{
					X: e.X, Y: e.Y - 30, Value: w.Damage, Timer: 0.5,
				})
				if e.HP <= 0 {
					e.Dead = true
					g.killCount++
					g.xpGems = append(g.xpGems, &XPGem{X: e.X, Y: e.Y, Value: e.XP})
				}
			}
		}
	}
}

func (g *Game) showLevelUpMenu() {
	g.levelingUp = true
	g.upgradeOptions = []UpgradeOption{
		{
			Name:        "Orb +1",
			Description: "Add another orbiting orb",
			Apply: func(g *Game) {
				for _, w := range g.player.Weapons {
					if w.Type == WeaponOrb {
						w.Count++
						return
					}
				}
			},
		},
		{
			Name:        "Damage +20%",
			Description: "Increase all damage",
			Apply: func(g *Game) {
				for _, w := range g.player.Weapons {
					w.Damage = int(float64(w.Damage) * 1.2)
				}
			},
		},
		{
			Name:        "Max HP +20",
			Description: "Increase max health",
			Apply: func(g *Game) {
				g.player.MaxHP += 20
				g.player.HP += 20
			},
		},
		{
			Name:        "Fire Rate +15%",
			Description: "Weapons attack faster",
			Apply: func(g *Game) {
				for _, w := range g.player.Weapons {
					w.Cooldown *= 0.85
				}
			},
		},
		{
			Name:        "New: Aura",
			Description: "Damage nearby enemies",
			Apply: func(g *Game) {
				g.player.Weapons = append(g.player.Weapons, &Weapon{
					Type: WeaponAura, Level: 1, Damage: 5, Cooldown: 0.5, Range: 80,
				})
			},
		},
		{
			Name:        "New: Projectile",
			Description: "Fire at nearest enemy",
			Apply: func(g *Game) {
				g.player.Weapons = append(g.player.Weapons, &Weapon{
					Type: WeaponProjectile, Level: 1, Damage: 15, Cooldown: 0.8, Count: 1,
				})
			},
		},
	}
	// Shuffle and pick 3
	rand.Shuffle(len(g.upgradeOptions), func(i, j int) {
		g.upgradeOptions[i], g.upgradeOptions[j] = g.upgradeOptions[j], g.upgradeOptions[i]
	})
	if len(g.upgradeOptions) > 3 {
		g.upgradeOptions = g.upgradeOptions[:3]
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background grid
	screen.Fill(color.RGBA{R: 30, G: 35, B: 40, A: 255})
	gridSize := 50.0
	offsetX := math.Mod(g.cameraX, gridSize)
	offsetY := math.Mod(g.cameraY, gridSize)
	for x := -gridSize; x < float64(screenWidth)+gridSize; x += gridSize {
		vector.DrawFilledRect(screen, float32(x-offsetX), 0, 1, screenHeight, color.RGBA{R: 40, G: 45, B: 50, A: 255}, false)
	}
	for y := -gridSize; y < float64(screenHeight)+gridSize; y += gridSize {
		vector.DrawFilledRect(screen, 0, float32(y-offsetY), screenWidth, 1, color.RGBA{R: 40, G: 45, B: 50, A: 255}, false)
	}

	// Draw XP gems
	for _, gem := range g.xpGems {
		screenX := gem.X - g.cameraX
		screenY := gem.Y - g.cameraY
		if screenX >= -20 && screenX <= screenWidth+20 && screenY >= -20 && screenY <= screenHeight+20 {
			vector.DrawFilledRect(screen, float32(screenX)-5, float32(screenY)-5, 10, 10, color.RGBA{R: 100, G: 200, B: 255, A: 255}, false)
		}
	}

	// Draw enemies
	for _, e := range g.enemies {
		if e.Dead {
			continue
		}
		screenX := e.X - g.cameraX
		screenY := e.Y - g.cameraY
		if screenX >= -30 && screenX <= screenWidth+30 && screenY >= -30 && screenY <= screenHeight+30 {
			c := color.RGBA{R: 200, G: 50, B: 50, A: 255}
			if e.HitFlash > 0 {
				c = color.RGBA{R: 255, G: 255, B: 255, A: 255}
			}
			vector.DrawFilledCircle(screen, float32(screenX), float32(screenY), 15, c, false)

			// Health bar
			hpRatio := float32(e.HP) / float32(e.MaxHP)
			vector.DrawFilledRect(screen, float32(screenX)-15, float32(screenY)-25, 30, 4, color.RGBA{R: 50, G: 50, B: 50, A: 255}, false)
			vector.DrawFilledRect(screen, float32(screenX)-15, float32(screenY)-25, 30*hpRatio, 4, color.RGBA{R: 255, G: 50, B: 50, A: 255}, false)
		}
	}

	// Draw projectiles
	for _, p := range g.projectiles {
		screenX := p.X - g.cameraX
		screenY := p.Y - g.cameraY
		vector.DrawFilledCircle(screen, float32(screenX), float32(screenY), float32(p.Radius), color.RGBA{R: 255, G: 200, B: 50, A: 255}, false)
	}

	// Draw aura effect if player has aura weapon
	for _, w := range g.player.Weapons {
		if w.Type == WeaponAura {
			screenX := g.player.X - g.cameraX
			screenY := g.player.Y - g.cameraY
			vector.StrokeCircle(screen, float32(screenX), float32(screenY), float32(w.Range), 2, color.RGBA{R: 100, G: 150, B: 255, A: 100}, false)
		}
	}

	// Draw player
	playerScreenX := g.player.X - g.cameraX
	playerScreenY := g.player.Y - g.cameraY
	vector.DrawFilledCircle(screen, float32(playerScreenX), float32(playerScreenY), 18, color.RGBA{R: 50, G: 150, B: 255, A: 255}, false)
	vector.StrokeCircle(screen, float32(playerScreenX), float32(playerScreenY), 18, 3, color.RGBA{R: 100, G: 200, B: 255, A: 255}, false)

	// Draw damage numbers
	for _, d := range g.damageNumbers {
		screenX := d.X - g.cameraX
		screenY := d.Y - g.cameraY
		alpha := uint8(d.Timer / 0.8 * 255)
		_ = alpha
		ebitenutil.DebugPrintAt(screen, formatInt(d.Value), int(screenX)-10, int(screenY))
	}

	// UI
	g.drawUI(screen)

	// Level up menu
	if g.levelingUp {
		g.drawLevelUpMenu(screen)
	}

	// Game over
	if g.gameOver {
		vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 0, G: 0, B: 0, A: 180}, false)
		ebitenutil.DebugPrintAt(screen, "GAME OVER", screenWidth/2-40, screenHeight/2-30)
		ebitenutil.DebugPrintAt(screen, "Survived: "+formatTime(g.gameTime), screenWidth/2-50, screenHeight/2)
		ebitenutil.DebugPrintAt(screen, "Kills: "+formatInt(g.killCount), screenWidth/2-35, screenHeight/2+20)
		ebitenutil.DebugPrintAt(screen, "Press SPACE to restart", screenWidth/2-80, screenHeight/2+50)
	}
}

func (g *Game) drawUI(screen *ebiten.Image) {
	// Top bar
	vector.DrawFilledRect(screen, 0, 0, screenWidth, 50, color.RGBA{R: 0, G: 0, B: 0, A: 180}, false)

	// Health bar
	vector.DrawFilledRect(screen, 10, 10, 200, 15, color.RGBA{R: 50, G: 50, B: 50, A: 255}, false)
	hpRatio := float32(g.player.HP) / float32(g.player.MaxHP)
	vector.DrawFilledRect(screen, 10, 10, 200*hpRatio, 15, color.RGBA{R: 255, G: 50, B: 50, A: 255}, false)
	ebitenutil.DebugPrintAt(screen, formatInt(g.player.HP)+"/"+formatInt(g.player.MaxHP), 85, 10)

	// XP bar
	xpNeeded := g.player.Level * 20
	xpRatio := float32(g.player.XP) / float32(xpNeeded)
	vector.DrawFilledRect(screen, 10, 30, 200, 10, color.RGBA{R: 50, G: 50, B: 50, A: 255}, false)
	vector.DrawFilledRect(screen, 10, 30, 200*xpRatio, 10, color.RGBA{R: 100, G: 200, B: 255, A: 255}, false)

	// Stats
	ebitenutil.DebugPrintAt(screen, "Lv "+formatInt(g.player.Level), 220, 15)
	ebitenutil.DebugPrintAt(screen, "Time: "+formatTime(g.gameTime), 300, 15)
	ebitenutil.DebugPrintAt(screen, "Kills: "+formatInt(g.killCount), 450, 15)
	ebitenutil.DebugPrintAt(screen, "Enemies: "+formatInt(len(g.enemies)), 600, 15)
}

func (g *Game) drawLevelUpMenu(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{R: 0, G: 0, B: 0, A: 150}, false)

	boxW := float32(350)
	boxH := float32(200)
	boxX := float32(screenWidth-350) / 2
	boxY := float32(screenHeight-200) / 2

	vector.DrawFilledRect(screen, boxX, boxY, boxW, boxH, color.RGBA{R: 40, G: 45, B: 55, A: 255}, false)
	vector.StrokeRect(screen, boxX, boxY, boxW, boxH, 3, color.RGBA{R: 255, G: 215, B: 0, A: 255}, false)

	ebitenutil.DebugPrintAt(screen, "LEVEL UP! Choose an upgrade:", int(boxX)+60, int(boxY)+15)

	for i, opt := range g.upgradeOptions {
		y := int(boxY) + 50 + i*50
		ebitenutil.DebugPrintAt(screen, "["+formatInt(i+1)+"] "+opt.Name, int(boxX)+20, y)
		ebitenutil.DebugPrintAt(screen, "    "+opt.Description, int(boxX)+20, y+15)
	}
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

func formatTime(t float64) string {
	mins := int(t) / 60
	secs := int(t) % 60
	return formatInt(mins) + ":" + formatIntPad(secs)
}

func formatIntPad(n int) string {
	if n < 10 {
		return "0" + formatInt(n)
	}
	return formatInt(n)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Endless Swarm - Survivor")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
