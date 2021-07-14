package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	W = 5000
	H = 5000
)

type Ctx struct {
	*image.NRGBA
	*sync.WaitGroup
}

func NewCtx(w, h int) *Ctx {
	return &Ctx{image.NewNRGBA(image.Rect(0, 0, w, h)), &sync.WaitGroup{}}
}
func (c *Ctx) SetPix(x, y int, clr *color.NRGBA) {
	if x >= c.Bounds().Dx()-1 || x <= 0 || y >= c.Bounds().Dy()-1 || y <= 0 {
		return
	}
	c.Pix[y*c.Stride+4*x], c.Pix[y*c.Stride+4*x+1], c.Pix[y*c.Stride+4*x+2], c.Pix[y*c.Stride+4*x+3] = clr.R, clr.G, clr.B, clr.A
}

// (h, k) is the center of the circle and r is the radius, clr is the color, h is the x coordinate, k is the y coordinate
func (c *Ctx) DrawCircle(h, k int, r float64, clr *color.NRGBA) {
	var y int
	for x := h; x-h < int(0.75*r); x++ {
		y = YPosCircle(h, k, x, r)
		c.SetPix(x, y, clr)
		c.SetPix(y, x, clr)
		c.SetPix(x-2*(x-h), y, clr)
		c.SetPix(y, x-2*(x-h), clr)
		c.SetPix(x, y-2*(y-k), clr)
		c.SetPix(y-2*(y-k), x, clr)
		c.SetPix(x-2*(x-h), y-2*(y-k), clr)
		c.SetPix(y-2*(y-k), x-2*(x-h), clr)
		//	fmt.Println(x,y)

	}
}

// Draws a disc or a filled circle
func (c *Ctx) DrawDisc(h0, k0 int, r0 float64, clr0 *color.NRGBA) {
	var y0 int

	for x0 := h0; x0-h0 < int(0.75*r0); x0++ {
		c.Add(1)
		go func(x, y, h, k int, r float64, clr *color.NRGBA) {
			defer c.Done()
			y = YPosCircle(h, k, x, r)
			c.DrawVerticalLine(x, y, y-2*(y-k), clr)
			c.DrawVerticalLine(y, x-2*(x-h), x, clr)
			c.DrawVerticalLine(x-2*(x-h), y, y-2*(y-k), clr)
			c.DrawVerticalLine(y-2*(y-k), x-2*(x-h), x, clr)
		}(x0, y0, h0, k0, r0, clr0)
		//	fmt.Println(x,y)

	}
}

//Deaw the circle with the goven thiccness
func (c *Ctx) DrawThiccCircle(h, k int, r, thiccness float64, clr *color.NRGBA) {
	var y, y0 int
	for x := h - int(r+thiccness); x < h+int(r+thiccness); x++ {
		y = YPosCircle(h, k, x, r-thiccness)
		y0 = YPosCircle(h, k, x, r+thiccness)
		c.DrawVerticalLine(x, y0, y, clr)
		c.DrawVerticalLine(x, y-2*(y-k), y0-2*(y0-k), clr)
	}
}

// Draw an outlined circle withn (h,k) as the center, r radius, thiccness being the thiccness of the outline, total radius is r+thiccness, clr0 is the color of fill color, clr1 is the color of outline
func (c *Ctx) DrawOutlinedDisc(h, k int, r, thiccness float64, clr0, clr1 *color.NRGBA) {
	c.DrawDisc(h, k, r, clr0)
	c.DrawThiccCircle(h, k, r+thiccness/2, thiccness, clr1)
}

func (c *Ctx) DrawEllipse(h, k int, a, b float64, clr *color.NRGBA) {
	y := 0
	x := 0
	//	fmt.Println(a, b)
	for x = h; x < h+int(a*0.75); x += 1 {

		y = YPosEllipse(h, k, x, a, b)
		c.SetPix(x, y, clr)
		//	c.SetPix(y, x, clr)
		//	fmt.Println(y)
		c.SetPix(x-2*(x-h), y, clr)
		//	c.SetPix(y, x-2*(x-h), cl)
		c.SetPix(int(x), y-2*(y-k), clr)
		//	c.SetPix(y-2*(y-k), x, clr)
		c.SetPix(int(x-2*(x-h)), y-2*(y-k), clr)
		//	c.SetPix(y-2*(y-k), x-2*(x-h), clr)
		//	fmt.Println(x,y)

	}

	for y = k; y < k+int(b*075); y += 1 {
		x = XPosEllipse(h, k, y, b, a)
		c.SetPix(x, y, clr)
		//	c.SetPix(y, x, clr)
		//	fmt.Println(y)
		c.SetPix(x-2*(x-h), y, clr)
		//	c.SetPix(y, x-2*(x-h), cl)
		c.SetPix(int(x), y-2*(y-k), clr)
		//	c.SetPix(y-2*(y-k), x, clr)
		c.SetPix(int(x-2*(x-h)), y-2*(y-k), clr)
		//	c.SetPix(y-2*(y-k), x-2*(x-h), clr)
		//	fmt.Println(x,y)

	}
}

func (c *Ctx) DrawThiccEllipse(h, k int, a, b, thiccness float64, clr *color.NRGBA) {
	var y, y0 int
	for x := h - int(a+thiccness); x < h+int(a+thiccness); x++ {
		y0 = YPosEllipse(h, k, x, a-thiccness, b-thiccness)
		y = YPosEllipse(h, k, x, a+thiccness, b+thiccness)
		c.DrawVerticalLine(x, y0, y, clr)
		c.DrawVerticalLine(x, y-2*(y-k), y0-2*(y0-k), clr)
	}
}

func (c *Ctx) DrawFilledEllipse(h0, k0 int, a0, b0 float64, clr0 *color.NRGBA) {
	for x0 := h0; x0-h0 < int(a0); x0++ {
		c.Add(1)
		go func(x, h, k int, a, b float64, clr *color.NRGBA) {
			defer c.Done()
			y := YPosEllipse(h, k, x, a, b)
			//	fmt.Printf("h: %d, k: %d, x: %d, a: %f, b: %f\n", h,k,x,a,b)
			c.DrawVerticalLine(x, y-2*(y-k), y, clr)
			c.DrawVerticalLine(x-2*(x-h), y-2*(y-k), y, clr)

			//	c.DrawDisc(x,y,15,clr)
		}(x0, h0, k0, a0, b0, clr0)
		//	fmt.Println(x,y)

	}

}

// Draws a line from x0, y0 to y1, y1
// NOTE: y0 < y1 and x0 < x1

func (c *Ctx) DrawLineLow(x0, y0, x1, y1 int, clr *color.NRGBA) {
	dx := x1 - x0
	dy := y1 - y0
	yi := 1
	if dy < 0 {
		yi = -1
		dy = -dy
	}
	D := (2 * dy) - dx
	y := y0

	for x := x0; x < x1; x++ {
		c.SetPix(x, y, clr)
		if D > 0 {
			y = y + yi
			D = D + (2 * (dy - dx))
		} else {
			D = D + 2*dy
		}
	}
}
func (c *Ctx) DrawLineHigh(x0, y0, x1, y1 int, clr *color.NRGBA) {
	dx := x1 - x0
	dy := y1 - y0
	xi := 1
	if dx < 0 {
		xi = -1
		dx = -dx
	}
	D := (2 * dx) - dy
	x := x0

	for y := y0; y <= y1; y++ {
		c.SetPix(x, y, clr)
		if D > 0 {
			x = x + xi
			D = D + (2 * (dx - dy))
		} else {
			D = D + 2*dx
		}
	}
}

func (c *Ctx) DrawLine(x0, y0, x1, y1 int, clr *color.NRGBA) {
	if Abs(y1-y0) < Abs(x1-x0) {
		if x0 > x1 {
			c.DrawLineLow(x1, y1, x0, y0, clr)
		} else {
			c.DrawLineLow(x0, y0, x1, y1, clr)
		}
	} else {
		if y0 > y1 {
			c.DrawLineHigh(x1, y1, x0, y0, clr)
		} else {
			c.DrawLineHigh(x0, y0, x1, y1, clr)
		}
	}
}
func (c *Ctx) DrawThiccLine(x0, y0, x1, y1 int, wd float64, clr *color.NRGBA) {
	dx := Abs(x1 - x0)
	sx := Sign(x1 - x0)
	dy := Abs(y1 - y0)
	sy := Sign(y1 - y0)
	err := dx - dy
	var e2, x2, y2 int /* error value e_xy */
	var ed float64
	if dx+dy == 0 {
		ed = 1
	} else {
		ed = math.Sqrt(float64(dx*dx + dy + dy))
	}
	r, g, b, a := clr.RGBA()
	for wd = (wd + 1) / 2; true; func() {}() { /* pixel loop */
		c.SetPix(x0, y0, NRGBA(r, g, b, uint32(math.Max(float64(a), 255*(float64(Abs(err-dx+dy))/ed-wd+1)))))
		e2 = err
		x2 = x0
		if 2*e2 >= -dx { /* x step */
			y2 = y0
			for e2 += dy; float64(e2) < ed*wd && (y1 != y2 || dx > dy); e2 += dx {
				y2 += sy
				c.SetPix(x0, y2, NRGBA(r, g, b, uint32(math.Max(float64(a), 255*(float64(Abs(e2))/ed-wd+1)))))
			}
			if x0 == x1 {
				break
			}
			e2 = err
			err -= dy
			x0 += sx
		}
		if 2*e2 <= dy { /* y step */
			for e2 = dx - e2; float64(e2) < ed*wd && (x1 != x2 || dx < dy); e2 += dy {
				x2 += sx
				c.SetPix(x2, y0, NRGBA(r, g, b, uint32(math.Max(float64(a), 255*(float64(Abs(e2))/ed-wd+1)))))
			}
			if y0 == y1 {
				break
			}
			err += dx
			y0 += sy
		}
	}
}

// Draws a line from x0, y0 to y1, y1
// with the given slope
// NOTE: y0 < y1 and x0 < x1
func (c *Ctx) DrawLineSlope(x0, y0, x1, y1 int, m float64, clr *color.NRGBA) {
	if m > 1 {
		for x := x0; x < x1; x++ {
			c.SetPix(int(m*float64(x-x0)+float64(y0)), x, clr)
		}
	} else {
		for x := x0; x < x1; x++ {
			c.SetPix(x, int(m*float64(x-x0)+float64(y0)), clr)
		}
	}
}

func NRGBA(r, g, b, a uint32) *color.NRGBA {
	return &color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

// Draws a line from x0, y0 to y1, y1
// NOTE: y0 < y1 and x0 < x1
func GetLineY(x0, y0, x1, y1 int) []int {
	if x1 < x0 {
		x2, y2, x3, y3 := x0, y0, x1, y1
		x0, y0, x1, y1 = x3, y3, x2, y2
	}
	m := float64(y1-y0) / float64(x1-x0)
	ys := make([]int, x1-x0+1)

	for x := x0; x <= x1; x++ {
		ys[x-x0] = int(m*float64(x-x0) + float64(y0))
	}
	return ys
}

// clears the image to alpha 50 black image
func (c *Ctx) Clear(r, g, b, a uint8) {
	for i := 0; i < len(c.Pix)/4; i++ {
		c.Pix[4*i], c.Pix[4*i+1], c.Pix[4*i+2], c.Pix[4*i+3] = r, g, b, a
	}
}

// Saves the image as a png
func (c *Ctx) Save(path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	png.Encode(file, c)
}

// Returns a y coords for the given circle at x0
func YPosCircle(h, k, x0 int, r float64) int {
	return -int(math.Sqrt(r*r-float64((x0-h)*(x0-h)))) + k
}

func YPosEllipse(h, k, x int, a, b float64) int {
	return int(math.Sqrt(1-float64((x-h)*(x-h))/(a*a))*b) + k

}

func XPosEllipse(h, k, y int, a, b float64) int {
	return int(math.Sqrt(1-float64((y-k)*(y-k))/(a*a))*b) + h

}

func (c *Ctx) DrawHorizontalLine(y, x0, x1 int, clr *color.NRGBA) {
	c.Add(1)
	if x1 < x0 {
		x2, x3 := x0, x1
		x0, x1 = x3, x2
	}

	go func() {
		defer c.Done()
		for i := x0; i <= x1; i++ {
			c.SetPix(i, y, clr)
		}
	}()
}
func (c *Ctx) DrawVerticalLine(x, y0, y1 int, clr *color.NRGBA) {
	c.Add(1)

	if y1 < H/2-H/4 {
		fmt.Printf("y1: %d \n", y1)
	}
	if y0 < H/2-H/4 {
		fmt.Printf("y0: %d \n", y0)
	}
	if y1 < y0 {
		y2, y3 := y0, y1
		y0, y1 = y3, y2
	}

	go func() {
		defer c.Done()
		for i := y0; i <= y1; i++ {
			c.SetPix(x, i, clr)
		}
	}()
}

func Sign(i int) int {
	if i >= 0 {
		return 1
	}
	return -1
}

func Abs(i int) int {
	if i >= 0 {
		return i
	}
	return -i
}
func main() {
	rand.Seed(time.Now().Unix())
	img := NewCtx(W, H)
	img.Clear(0, 0, 0, 255)
	start := time.Now()
	//	img.DrawFilledEllipse(W/2, H/2, W/4, H/3, &color.NRGBA{62, 255, 255, 255})
	//img.DrawThiccCircle(W/2, H/2, H/4, H/10, &color.NRGBA{76, 100, 220, 255})
	//	img.Wg.Wait()
	//	fmt.Println("Circle:", time.Now().Sub(start))
	//	img.Save("circle.png")
	//	img.Clear(0, 0, 0, 255)
	//	start = time.Now()
	//	img.DrawFilledEllipse(W/2, H/2, W/4, H/3, &color.NRGBA{62, 255, 255, 255})
	//	img.DrawThiccEllipse(W/2, H/2, H/4, H/6, H/10, &color.NRGBA{76, 200, 120, 255})
	//
	//	img.Wg.Wait()
	//	fmt.Println("Ellipse:", time.Now().Sub(start))
	//	img.Save("Ellipse.png")
	//	img.Clear(0, 0, 0, 255)
	//	start = time.Now()
	for i := 0.0; i < 2*math.Pi; i+= math.Pi/360 {
		x := W/2 + W/4*math.Cos(i)
		y := H/2 + H/4*math.Sin(i)
		img.DrawThiccLine(W/2, H/2, int(x), int(y), 20, &color.NRGBA{byte(rand.Intn(255)), byte(rand.Intn(255)), byte(rand.Intn(255)), 255})
	}
	//img.DrawThiccLine(50, H/2-2, H-50, H/2+3, 50, &color.NRGBA{210, 192, 27, 255})

	//img.DrawThiccLine(W/2, H/2, W - W/8, H/2 - H/8, 10, &color.NRGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255})
	img.Wait()
	img.DrawDisc(W/2, H/2, H/150, &color.NRGBA{76, 100, 220, 255})
	img.Wait()
	fmt.Println("Line:", time.Now().Sub(start))

	img.Save("Line.png")
}
