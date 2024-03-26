package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func runFont() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	face, err := loadTTF("intuitive.ttf", 80)
	if err != nil {
		panic(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)
	txt := text.New(pixel.V(50, 500), atlas)

	txt.Color = colornames.Lightgrey

	fps := time.Tick(time.Second / 120)

	for !win.Closed() {
		txt.WriteString(win.Typed())
		if win.JustPressed(pixelgl.KeyEnter) || win.Repeated(pixelgl.KeyEnter) {
			txt.WriteRune('\n')
		}

		win.Clear(colornames.Darkcyan)
		txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))
		win.Update()

		<-fps
	}
}

func runShapes() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	imd.Color = colornames.Blueviolet
	imd.EndShape = imdraw.RoundEndShape
	imd.Push(pixel.V(100, 100), pixel.V(700, 100))
	imd.EndShape = imdraw.SharpEndShape
	imd.Push(pixel.V(100, 500), pixel.V(700, 500))
	imd.Line(30)

	imd.Color = colornames.Limegreen
	imd.Push(pixel.V(500, 500))
	imd.Circle(300, 50)
	imd.Color = colornames.Navy
	imd.Push(pixel.V(200, 500), pixel.V(800, 500))
	imd.Ellipse(pixel.V(120, 80), 0)

	imd.Color = colornames.Red
	imd.EndShape = imdraw.RoundEndShape
	imd.Push(pixel.V(500, 350))
	imd.CircleArc(150, -math.Pi, 0, 30)

	for !win.Closed() {
		win.Clear(colornames.Aliceblue)
		imd.Draw(win)
		win.Update()
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	spritesheet, err := LoadPicture("trees.png")
	if err != nil {
		panic(err)
	}

	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)

	var treesFrames []pixel.Rect
	for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 32 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 32 {
			treesFrames = append(treesFrames, pixel.R(x, y, x+32, y+32))
		}
	}

	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
	)

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.Pressed(pixelgl.MouseButtonLeft) {
			tree := pixel.NewSprite(spritesheet, treesFrames[rand.Intn(len(treesFrames))])
			mouse := cam.Unproject(win.MousePosition())
			tree.Draw(batch, pixel.IM.Scaled(pixel.ZV, 4).Moved(mouse))
		}
		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		win.Clear(colornames.Forestgreen)
		batch.Draw(win)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	//pixelgl.Run(run)
	//	pixelgl.Run(runShapes)

	//pixelgl.Run(runFont)

	pixelgl.Run(RenderTiles)
}
