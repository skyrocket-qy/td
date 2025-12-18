# AI-Generation Game Framework

A Go-based ECS (Entity-Component-System) game framework designed for AI-assisted game development.

## Goals

1. **Code-based** — AI can 100% participate, no GUI tools required
2. **Multi-platform** — Windows, Linux, macOS, Android, iOS, Web
3. **Scalable performance** — From mobile phones to high-end PCs
4. **ECS Architecture** — Clean separation of data and logic

---

## Quick Start

### Prerequisites

- Go 1.21+
- For mobile builds: `gomobile` (`make mobile-init`)
- For Windows cross-compile from macOS/Linux: MinGW (`x86_64-w64-mingw32-gcc`)

### Run an Example

```bash
# Run the snake game
make run-snake

# Run other examples
make run-pong
make run-breakout
make run-slots
```

---

## Project Structure

```
├── cmd/game/           # Main game entry point
├── examples/           # Example games (snake, pong, breakout, slots, etc.)
├── internal/
│   ├── components/     # ECS Components (Position, Velocity, Sprite, etc.)
│   ├── systems/        # ECS Systems (Render, Movement, Collision, etc.)
│   ├── engine/         # Core engine functionality
│   ├── game/           # Game-specific logic (cards, monsters, maps)
│   └── assets/         # Asset loading (images, shaders, fonts)
├── assets/             # Game assets (images, sounds, fonts)
├── scripts/            # Build scripts (WASM, etc.)
└── Makefile            # Build commands
```

---

## Developing a New Game

### 1. Create Your Game Directory

```bash
mkdir -p examples/mygame
touch examples/mygame/main.go
```

### 2. Basic Game Template

```go
package main

import (
    "log"
    "github.com/hajimehoshi/ebiten/v2"
)

const (
    screenWidth  = 640
    screenHeight = 480
)

type Game struct {
    // Your game state here
}

func NewGame() *Game {
    return &Game{}
}

func (g *Game) Update() error {
    // Handle input and update game logic
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Render your game
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return screenWidth, screenHeight
}

func main() {
    ebiten.SetWindowSize(screenWidth, screenHeight)
    ebiten.SetWindowTitle("My Game")

    if err := ebiten.RunGame(NewGame()); err != nil {
        log.Fatal(err)
    }
}
```

### 3. Using ECS Components

```go
import (
    "td/internal/components"
    "github.com/mlange-42/ark/ecs"
)

// Create a world
world := ecs.NewWorld()

// Create entities with components
entity := world.NewEntity()
world.Add(entity, 
    components.Position{X: 100, Y: 100},
    components.Velocity{X: 1, Y: 0},
    components.NewSprite(myImage),
)
```

### 4. Creating Systems

```go
import "td/internal/systems"

// Use built-in render system
renderSystem := systems.NewRenderSystem(world)

// In your Draw method:
func (g *Game) Draw(screen *ebiten.Image) {
    renderSystem.Draw(world, screen)
}
```

---

## Available Components

| Component    | Description                          |
|-------------|--------------------------------------|
| `Position`  | 2D position (X, Y float64)           |
| `Velocity`  | 2D velocity (X, Y float64)           |
| `Sprite`    | Renderable image with scale/offset   |
| `Collider`  | Bounding box with collision layers   |
| `Health`    | Current/Max health                   |
| `Tag`       | Entity categorization marker         |

---

## Example Games

| Game          | Description                           | Run Command           |
|---------------|---------------------------------------|----------------------|
| Snake         | Classic snake game                    | `make run-snake`     |
| Pong          | Two-player pong                       | `make run-pong`      |
| Breakout      | Brick breaker                         | `make run-breakout`  |
| Slots         | Slot machine game                     | `make run-slots`     |
| Minesweeper   | Classic minesweeper                   | `go run ./examples/minesweeper` |
| 2048          | Puzzle 2048                           | `go run ./examples/puzzle_2048` |
| Cookie Clicker| Incremental clicker                   | `go run ./examples/cookie_clicker` |
| Flappy        | Flappy bird clone                     | `go run ./examples/flappy` |

---

## Building & Deployment

### Desktop

```bash
make build              # Current platform
make build-darwin       # macOS
make build-windows      # Windows (requires MinGW)
make build-linux        # Linux
make build-all-desktop  # All desktop platforms
```

### WebAssembly

```bash
make build-wasm         # Standard Go WASM
make build-wasm-tiny    # TinyGo (smaller binary)
make serve-wasm         # Serve locally at http://localhost:8080
```

### Mobile

```bash
make mobile-init        # Initialize gomobile (run once)
make build-android      # Android AAR library
make build-android-apk  # Android APK
make build-ios          # iOS framework
```

### Release Package

```bash
make dist VERSION=1.0.0  # Creates all platform packages
```

---

## Development Commands

```bash
make run      # Run main game
make test     # Run tests
make lint     # Run linter
make clean    # Clean build artifacts
make help     # Show all commands
```

---

## Tips for AI-Assisted Development

1. **Keep games in `examples/`** — Each game is self-contained in its own directory
2. **Use the ECS pattern** — Separate data (components) from logic (systems)
3. **Leverage ebiten** — The framework uses [Ebitengine](https://ebitengine.org/) for rendering
4. **Check existing examples** — Snake and Pong are good starting points
5. **Run frequently** — Use `go run ./examples/yourname` to test changes