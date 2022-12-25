package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sync"

	"github.com/lucasb-eyer/go-colorful"
)

const (
	W = 1000
	H = 1000
)

func main() {
	img := image.NRGBA{Pix: make([]uint8, W*H*4), Stride: W * 4, Rect: image.Rect(0, 0, W, H)}
	disc(&img, W/2, H/2, W/4, colorful.HappyColor())
	ring(&img, W/2, H/2, W/3, W/2, colorful.FastHappyColor())
	circle(&img, W/2, H/2, W/3, colorful.HappyColor())
	for i := 0; i < 360; i++ {
		lineThicc(&img, W/2, H/2, int(W/2+(W/3)*math.Cos(math.Pi*float64(i)/180)), int(H/2+(H/3*math.Sin(math.Pi*float64(i)/180))), W/200, colorful.FastHappyColor())
	}
	file, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}
	png.Encode(file, &img)
}

func circle(img *image.NRGBA, cx, cy int, radius float64, c color.Color) {
	y := int(2 * abs(cx))
	for x := int(cx); x-y <= cx-cy; x++ {
		y = int(math.Sqrt(radius*radius-float64((cx-x)*(cx-x)))) + cy
		img.Set(x, y, c)
		img.Set(mirror45(y, cx, cy), mirror45(x, cy, cx), c)
		img.Set(-x+int(2*cx), y, c)
		img.Set(mirror45(y, cx, cy), mirror45(-x+int(2*cx), cy, cx), c)
		img.Set(x, -y+int(2*cy), c)
		img.Set(mirror45(-y+int(2*cy), cx, cy), mirror45(x, cy, cx), c)
		img.Set(-x+int(2*cx), -y+int(2*cy), c)
		img.Set(mirror45(-y+int(2*cy), cx, cy), mirror45(-x+int(2*cx), cy, cx), c)
	}
}

func disc(img *image.NRGBA, cx, cy int, radius float64, c color.Color) {
	y := int(2 * abs(cx))
	var wg sync.WaitGroup
	for x := int(cx); x-y <= cx-cy; x++ {
		y = int(math.Sqrt(radius*radius-float64((cx-x)*(cx-x)))) + cy
		wg.Add(1)
		go func(x, y int) {
			horizontalLine(img, -x+int(2*cx), x, y, c)
			horizontalLine(img, mirror45(-y+int(2*cy), cx, cy), mirror45(y, cx, cy), mirror45(x, cy, cx), c)
			horizontalLine(img, mirror45(-y+int(2*cy), cx, cy), mirror45(y, cx, cy), mirror45(-x+int(2*cx), cy, cx), c)
			horizontalLine(img, -x+int(2*cx), x, -y+int(2*cy), c)
			wg.Done()
		}(x, y)
	}
	wg.Wait()
}

func ring(img *image.NRGBA, cx, cy int, rInner, rOuter float64, c color.Color) {
	yOuter := int(2 * abs(cx))
	var wg sync.WaitGroup
	for x := int(cx); x <= cx+int(rOuter); x++ {
		yOuter = int(math.Sqrt(rOuter*rOuter-float64((cx-x)*(cx-x)))) + cy
		wg.Add(1)
		go func(x, y int) {
			if abs(x-cx) > int(rInner) {
				verticalLine(img, -y+int(2*cy), y, x, c)
				verticalLine(img, y, -y+int(2*cy), -x+int(2*cx), c)
			} else {
				yInner := int(math.Sqrt(rInner*rInner-float64((cx-x)*(cx-x)))) + cy
				verticalLine(img, y, yInner, x, c)
				verticalLine(img, y, yInner, -x+int(2*cx), c)
				verticalLine(img, -yInner+int(2*cy), -y+int(2*cy), -x+int(2*cx), c)
				verticalLine(img, -yInner+int(2*cy), -y+int(2*cy), x, c)
			}
			wg.Done()
		}(x, yOuter)
	}
	wg.Wait()
}

func line0(img *image.NRGBA, x0, y0, x1, y1 int, c color.Color) {
	slope := float64(y1-y0) / float64(x1-x0)
	if math.Abs(slope) < 1 {
		for x := int(x0); x <= int(x1); x++ {
			y := int(slope*float64(x)) - x0 + y0
			img.Set(x, y, c)
		}
	} else {
		for y := int(y0); y <= int(y1); y++ {
			x := int(float64(y-y0)/slope) + x0
			img.Set(x, y, c)
		}
	}
}

func lineThicc(img *image.NRGBA, x0, y0, x1, y1 int, thicc float64, c color.Color) {
	if y0 == y1 {
		for y := y0 - int(thicc); y < y0+int(thicc); y++ {
			horizontalLine(img, int(x0), int(x1), int(y0)+y, c)
		}
		return
	}
	if x0 == x1 {
		for x := x0 - int(thicc); x < x0+int(thicc); x++ {
			verticalLine(img, int(y0), int(y1), int(x0)+x, c)
		}
		return
	}
	if y0 < y1 && x0 > x1 {
		y0 += y1
		y1 = y0 - y1
		y0 -= y1
		x0 += x1
		x1 = x0 - x1
		x0 -= x1
	}
	thicc *= 0.5
	slope := float64(y1-y0) / float64(x1-x0)
	fmt.Printf("Slope: %f\n", slope)
	perpSlope := -1 / slope
	dxE := thicc / math.Sqrt(1+perpSlope*perpSlope)
	dyE := int(dxE * perpSlope)
	y0N := 0
	for x := x0 - int(dxE); x < x0+int(dxE); x++ {
		y0N = int(perpSlope*(float64(x-x0)+dxE)) + y0 + dyE
		verticalLine(img, int(slope*(float64(x)-x0+dxE)+y0+dyE), y0N, x, c)
	}
	dy := int(slope*2*dxE+y0+dyE) - y0N
	for x := int(x0 + dxE); x < int(x1-dxE); x++ {
		y := int(slope*(float64(x)-x0+dxE) + y0 + dyE)
		verticalLine(img, y, y-dy, x, c)
	}
	for x := int(x1 - dxE); x < int(x1+dxE); x++ {
		y0N = int(perpSlope*(float64(x)-x1+dxE) + y1 + dyE)
		verticalLine(img, y0N, int(slope*(float64(x)-x0+dxE)+y0+dyE)-dy, x, c)
	}
}
func verticalLine(img *image.NRGBA, y0, y1, x int, c color.Color) {
	if y0 > y1 {
		y0 += y1
		y1 = y0 - y1
		y0 -= y1
	}
	for y := y0; y <= y1; y += 1 {
		img.Set(x, y, c)
	}
}

func horizontalLine(img *image.NRGBA, x0, x1, y int, c color.Color) {
	if x0 > x1 {
		x0 += x1
		x1 = x0 - x1
		x0 -= x1
	}
	for x := x0; x <= x1; x++ {
		img.Set(x, y, c)
	}
}

func abs(x int) int {
	return x ^ -1 + 1
}

func mirror45(y int, cx, cy int) int {
	return y - cy + cx
}
