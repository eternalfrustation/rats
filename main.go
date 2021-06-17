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
	W = 1000
	H = 1000
)

type Ctx struct {
	*image.NRGBA
	Wg *sync.WaitGroup
}

func NewCtx(w, h int) *Ctx {
	return &Ctx{image.NewNRGBA(image.Rect(0, 0, W, H)), &sync.WaitGroup{}}
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
		c.Wg.Add(1)
		go func(x, y, h, k int, r float64, clr *color.NRGBA) {
			defer c.Wg.Done()
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
		c.Wg.Add(1)
		go func(x, h, k int, a, b float64, clr *color.NRGBA) {
			defer c.Wg.Done()
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
func (c *Ctx) DrawLine(x0, y0, x1, y1 int, clr *color.NRGBA) {
	x2, y2, x3, y3 := x0, y0, x1, y1
	x0, x1 = Min(x2, x3), Max(x2, x3)
	y0, y1 = Min(y2, y3), Min(y2, y3)

	m := float64(y1-y0) / float64(x1-x0)
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

func (c *Ctx) DrawThiccLine(x0, y0, x1, y1 int, thiccness float64, clr *color.NRGBA) {

	x2, y2, x3, y3 := x0, y0, x1, y1
	x0, x1 = Min(x2, x3), Max(x2, x3)
	y0, y1 = Min(y2, y3), Max(y2, y3)
	if x0 == x1 {
		for x := x0 - int(thiccness); x < x0 + int(thiccness); x ++ {
			c.DrawVerticalLine(x, y0, y1, clr)
		}
		return
	}
	if y0 == y1 {
		for y := y0 - int(thiccness); y < y0 + int(thiccness); y++ {
			c.DrawHorizontalLine(y, x0, x1, clr)
		}
		return
	}
	m := float64(y1-y0) / float64(x1-x0)
	if m < 1 {

		dy := 2 * thiccness / math.Sqrt(m*m+1)
		offY := int(dy / 2)
		offX := int(m * float64(offY))
		for i, y := range c.GetLineY(x0-offX, y0-offY, x1+offX, y1+offY) {
			c.DrawVerticalLine(x0+i, y, y+int(dy), clr)

		}
	} else {
		m0 := -1/m
		dy := 2 * thiccness * math.Sqrt(1+1/m0*m0)
		offY := int(dy / 2)
		offX := int(m * float64(offY))
		for i, y := range c.GetLineY(x0-offX, y0-offY, x1+offX, y1+offY) {
			c.DrawVerticalLine(x0+i, y, y+int(dy), clr)
		}
	}
}

// Draws a line from x0, y0 to y1, y1
// NOTE: y0 < y1 and x0 < x1
func (c *Ctx) GetLineY(x0, y0, x1, y1 int) []int {
	m := float64(y1-y0) / float64(x1-x0)
	ys := make([]int, x1-x0)
	for x := x0; x < x1; x++ {
		ys[x-x0] = int(m*float64(x-x0) + float64(y0))
	}
	return ys
}

// clears the image to alpha 50 black image
func (c *Ctx) Clear() {
	for i := 0; i < len(c.Pix); i++ {
		c.Pix[i] = 50
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

func (c *Ctx) DrawVerticalLine(x, y0, y1 int, clr *color.NRGBA) {
	c.Wg.Add(1)
	go func() {
		defer c.Wg.Done()
	for i := y0; i <= y1; i++ {
		c.SetPix(x, i, clr)
	}}()
}

func (c *Ctx) DrawHorizontalLine(y, x0, x1 int, clr *color.NRGBA) {
	c.Wg.Add(1)
	go func() {
		defer c.Wg.Done()
	for i := x0; i <= x1; i++ {
		c.SetPix(i, y, clr)
	}}()
}
func main() {
	rand.Seed(time.Now().Unix())
	img := NewCtx(W, H)
	img.Clear()
	start := time.Now()
	//	img.DrawFilledEllipse(W/2, H/2, W/4, H/3, &color.NRGBA{62, 255, 255, 255})
//	img.DrawThiccCircle(W/2, H/2, H/4, H/10, &color.NRGBA{76, 100, 220, 255})
//	img.Wg.Wait()
//	fmt.Println("Circle:", time.Now().Sub(start))
//	img.Save("circle.png")
//	img.Clear()
//	start = time.Now()
	//	img.DrawFilledEllipse(W/2, H/2, W/4, H/3, &color.NRGBA{62, 255, 255, 255})
//	img.DrawThiccEllipse(W/2, H/2, H/4, H/6, H/10, &color.NRGBA{76, 200, 120, 255})
//
//	img.Wg.Wait()
//	fmt.Println("Ellipse:", time.Now().Sub(start))
//	img.Save("Ellipse.png")
//	img.Clear()
	start = time.Now()
	img.DrawThiccLine(W/2, 0, W/2, H, 10, &color.NRGBA{27, 219, 150, 255})
	img.Wg.Wait()
	fmt.Println("Line:", time.Now().Sub(start))

	img.Save("Line.png")
}
func Min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}
func Max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}
