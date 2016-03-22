package mscolor

import (
	"fmt"
	"math"
)

// Color is RGB<=>HSV color converter
type Color struct {
	// heap ARGB
	R uint8
	G uint8
	B uint8
	H int
	S float64
	V float64
	A uint8
}

// FromRGB returns new Color by red, green, blue
func FromRGB(r uint8, g uint8, b uint8) *Color {
	return FromARGB(255, r, g, b)
}

// FromARGB returns new Color by alpha, red, green, blue
func FromARGB(a uint8, r uint8, g uint8, b uint8) *Color {
	c := &Color{A: a, R: r, G: g, B: b}
	return c.MakeHSV()
}

// FromHSV returns new Color by hue, saturation, brightness
func FromHSV(h int, s float64, v float64) *Color {
	return FromAHSV(255, h, s, v)
}

// FromAHSV returns new Color by alpha, hue, saturation, brightness
func FromAHSV(a uint8, h int, s float64, v float64) *Color {
	c := &Color{A: a, H: h, S: s, V: v}
	return c.MakeRGB()
}

// MakeRGB makes red, green, blue member
func (c *Color) MakeRGB() *Color {
	// 彩度がないときはすべて同値になる
	if c.S == 0 {
		v := uint8(round(c.V * 255))
		c.R = v
		c.G = v
		c.B = v
		return c
	}

	H := float64(c.H)
	S := c.S

	Hi := math.Floor(H / 60)
	F := (H / 60.0) - Hi
	M := uint8(round((c.V * (1.0 - S)) * 255.0))
	N := uint8(round((c.V * (1.0 - (S * F))) * 255.0))
	K := uint8(round((c.V * (1.0 - (S * (1 - F)))) * 255.0))
	V := uint8(round(c.V * 255.0))

	switch {
	case Hi == 0:
		c.R = V
		c.G = K
		c.B = M
	case Hi == 1:
		c.R = N
		c.G = V
		c.B = M
	case Hi == 2:
		c.R = M
		c.G = V
		c.B = K
	case Hi == 3:
		c.R = M
		c.G = N
		c.B = V
	case Hi == 4:
		c.R = K
		c.G = M
		c.B = V
	case Hi == 5:
		c.R = V
		c.G = M
		c.B = N
	}

	return c
}

// MakeHSV makes hue, saturation, brightness member
func (c *Color) MakeHSV() *Color {
	max := maxVal(c.R, c.G, c.B)
	min := minVal(c.R, c.G, c.B)
	R := float64(c.R)
	G := float64(c.G)
	B := float64(c.B)

	// Vは最大値
	c.V = max / 255

	// maxがゼロであればその色は黒なので S,H共に0
	if c.V == 0 {
		c.S = 0
		c.H = 0
		return c
	}

	// 彩度は最大値の比率と等しい
	c.S = (max - min) / max

	var H float64
	if max == R {
		H = 60.0 * ((G - B) / (max - min))
	} else if max == G {
		H = 60.0 * (2.0 + ((B - R) / (max - min)))
	} else {
		H = 60.0 * (4.0 + ((R - G) / (max - min)))
	}

	c.H = round(regularize(H)) % 360

	return c
}

func maxVal(vals ...uint8) float64 {
	t := uint8(0)
	for _, i := range vals {
		if t < i {
			t = i
		}
	}
	return float64(t)
}

func minVal(vals ...uint8) float64 {
	t := vals[0]
	for _, i := range vals {
		if t > i {
			t = i
		}
	}
	return float64(t)
}

// hueが0以下の場合正の数になるまで回す
func regularize(hue float64) float64 {
	if math.IsNaN(hue) {
		return 0
	}
	H := hue
	for {
		if H >= 0 {
			break
		} else {
			H = H + 360
		}
	}
	return H
}

func round(n float64) int {
	if n < 0 {
		return int(math.Ceil(n - 0.5))
	}
	return int(math.Floor(n + 0.5))
}

// String returns RGB string
func (c *Color) String() string {
	return fmt.Sprintf("#%X%X%X", c.R, c.G, c.B)
}
