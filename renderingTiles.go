package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	WindowWidth  = 1000
	WindowHeight = 600
	TileSize     = 32
	NumTilesX    = WindowWidth / TileSize
	NumTilesY    = WindowHeight / TileSize
)

var (
	camPos       = pixel.V(100, 290)
	camSpeed     = 500.0
	camZoom      = 1.0
	camZoomSpeed = 1.2
)

type Player struct {
	sprite      *pixel.Sprite
	pos         pixel.Vec
	velocity    pixel.Vec // velocity represents the player's speed and direction of movement
	groundLevel float64   // groundLevel represents the player's y-coordinate when on the ground
	jumping     bool
	frames      []*pixel.Sprite // jumping indicates whether the player is currently jumping
	frame       int32
}

var (
	frames = 0
	second = time.Tick(time.Second)
)

func NewPlayer(sprite *pixel.Sprite, pos pixel.Vec) *Player {
	return &Player{
		sprite:      sprite,
		pos:         pos,
		velocity:    pixel.ZV,
		groundLevel: pos.Y,
		jumping:     false,
		frames:      []*pixel.Sprite{},
		frame:       0,
	}
}

func checkCollision(rect1 pixel.Rect, rect2 pixel.Rect) bool {
	return rect1.Intersects(rect2)
}

func (p *Player) Update(dt float64, win *pixelgl.Window, cam pixel.Vec) {

	newPos := p.pos.Add(p.velocity.Scaled(dt)) // Update player position based on velocity
	newcamPosX := cam.X
	//	newcamPosY := cam.Y
	// Check if player has reached or gone below the ground level
	if !p.jumping {
		p.velocity.Y -= 1000 * dt // Adjust gravity as needed
	}

	if newPos.Y <= p.groundLevel {
		p.jumping = false
		newPos.Y = p.groundLevel
		p.velocity.Y = 0 // Stop vertical velocity when player lands on ground
	}

	if win.Pressed(pixelgl.KeyA) {
		p.frame++ // Increment frame index
		if int(p.frame) >= len(p.frames) {
			p.frame = 0 // Reset frame index if it exceeds the length of frames
		}
		println(p.frame)
		p.frames[int(p.frame)].Draw(win, pixel.IM.Moved(p.pos))
		newcamPosX -= camSpeed * dt
		newPos.X -= camSpeed * dt
	}

	if win.Pressed(pixelgl.KeyD) {
		println(p.frame)
		p.frame++ // Increment frame index
		if int(p.frame) >= len(p.frames) {
			p.frame = 0 // Reset frame index if it exceeds the length of frames
		}
		p.frames[int(p.frame)].Draw(win, pixel.IM.Moved(p.pos))
		newcamPosX += camSpeed * dt
		newPos.X += camSpeed * dt
	}
	if win.Pressed(pixelgl.KeySpace) {
		newPos.Y += camSpeed * dt
		//camPos.Y -= camSpeed * dt
	}
	if win.Pressed(pixelgl.KeyS) {
		//newPos.Y -= camSpeed * dt
		//camPos.Y += camSpeed * dt
	}
	// if win.Pressed(pixelgl.KeySpace) && !p.jumping { // Allow jumping only when not already jumping
	// 	p.jumping = true
	// 	p.velocity.Y = 500 // Set initial vertical velocity for jump
	// }

	// Create a bounding box for the player

	// No collision detected, update player position
	playerBounds := pixel.R(newPos.X-p.sprite.Frame().W()/2, newPos.Y-p.sprite.Frame().H()/2, newPos.X+p.sprite.Frame().W()/2, newPos.Y+p.sprite.Frame().H()/2)

	for x := 0; x < len(MapCordinate); x++ {
		for y := 0; y < len(MapCordinate[0]); y++ {
			if MapCordinate[x][y] == 2 || MapCordinate[x][y] == 4 {
				tileBounds := pixel.R(float64(x*TileSize), float64(y*TileSize), float64((x)*TileSize), float64((y)*TileSize))
				if checkCollision(playerBounds, tileBounds) {
					return
				}
			}
		}
	}
	camPos.X = newcamPosX
	p.pos = newPos
}

func (p *Player) Draw(win *pixelgl.Window) {
	tileSprite, err := LoadPicture("tiles/1 Tiles/Tile_02.png")
	tileDirt, err := LoadPicture("tiles/1 Tiles/Tile_12.png")
	tileDemon, err := LoadPicture("assests/PNG/Objects_separetely/Bones_shadow1_1.png")

	if err != nil {
		log.Fatal(err)
	}

	for x := 0; x < len(MapCordinate); x++ {
		for y := 0; y < len(MapCordinate[0]); y++ {
			if MapCordinate[x][y] == 2 {

				tilePos := pixel.V(float64(x*TileSize), float64(y*TileSize))
				tile := pixel.NewSprite(tileDirt, tileDirt.Bounds())
				tile.Draw(win, pixel.IM.Moved(tilePos))
			} else if MapCordinate[x][y] == 4 {
				tilePos := pixel.V(float64(x*TileSize), float64(y*TileSize))
				tile := pixel.NewSprite(tileSprite, tileSprite.Bounds())
				tile.Draw(win, pixel.IM.Moved(tilePos))
			} else if MapCordinate[x][y] == 6 {
				tilePos := pixel.V(float64(x*TileSize), float64(y*TileSize))
				tile := pixel.NewSprite(tileDemon, tileDemon.Bounds())
				tile.Draw(win, pixel.IM.Moved(tilePos))
			}
		}
	}
	p.frames[p.frame].Draw(win, pixel.IM.Moved(p.pos))
}

func RenderTiles() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1000, 600),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Load tile image

	playerSprite, err := LoadPicture("assests/PNG/Objects_separetely/Bones_shadow1_4.png") // Load the player sprite

	if err != nil {
		log.Fatal(err)
	}

	player := NewPlayer(pixel.NewSprite(playerSprite, playerSprite.Bounds()), pixel.V(100, 100))

	// Draw tiles

	spritesheet, err := LoadPicture("Samurai/Run.png")
	//player.frames = append(player.frames, pixel.NewSprite(tileSprite, tileSprite.Bounds()))
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 125 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 140 {

			player.frames = append(player.frames, pixel.NewSprite(spritesheet, pixel.R(x, y, x+120, y+140)))
		}
	}

	last := time.Now()
	for !win.Closed() {

		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)
		player.Update(dt, win, camPos)
		// Update player
		win.Clear(colornames.White)

		// Draw player
		player.Draw(win)

		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}

		win.Update()

	}
}
