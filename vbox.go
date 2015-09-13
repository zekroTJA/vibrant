package vibrant

import "math"

const (
	COMPONENT_RED   int = -3
	COMPONENT_GREEN int = -2
	COMPONENT_BLUE      = -1
)

type Vbox struct {
	lowerIndex  int
	upperIndex  int
	minRed      int
	maxRed      int
	minGreen    int
	maxGreen    int
	minBlue     int
	maxBlue     int
	colors      map[int]int
	populations map[int]int
}

func NewVbox(lowerIndex, upperIndex int, colors map[int]int, populations map[int]int) *Vbox {
	v := Vbox{lowerIndex: lowerIndex, upperIndex: upperIndex, colors: colors, populations: populations}
	v.fitBox()
	return &v
}

func (v *Vbox) fitBox() {
	v.minRed = 255
	v.minGreen = 255
	v.minBlue = 255
	v.maxRed = 0
	v.maxGreen = 0
	v.maxBlue = 0

	for i := v.lowerIndex; i <= v.upperIndex; i++ {
		r, g, b := unpackColor(v.colors[i])
		if r > v.maxRed {
			v.maxRed = r
		}
		if r < v.minRed {
			v.minRed = r
		}
		if g > v.maxGreen {
			v.maxGreen = g
		}
		if g < v.minGreen {
			v.minGreen = g
		}
		if b > v.maxBlue {
			v.maxBlue = b
		}
		if b < v.minBlue {
			v.minBlue = b
		}
	}
}

func (v *Vbox) Volume() int {
	return (v.maxRed - v.minRed + 1) * (v.maxGreen - v.minGreen + 1) * (v.maxBlue - v.minBlue + 1)
}

func (v *Vbox) CanSplit() bool {
	return (v.upperIndex - v.lowerIndex + 1) > 1
}

func (v *Vbox) Split() *Vbox {
	if !v.CanSplit() {
		panic("Cannot split a box with only 1 color!")
	}
	lenRed := v.maxRed - v.minRed
	lenGreen := v.maxGreen - v.minGreen
	lenBlue := v.maxBlue - v.minBlue

	var longestDim, midPoint int
	switch {
	case lenRed >= lenGreen && lenRed >= lenBlue:
		longestDim = COMPONENT_RED
		// Already in RGB, no need to do anything
		midPoint = (v.minRed + v.maxRed) / 2
	case lenGreen >= lenRed && lenGreen >= lenBlue:
		longestDim = COMPONENT_GREEN
		// We need to do a RGB to GRB swap, or vice-versa
		for i := v.lowerIndex; i <= v.upperIndex; i++ {
			r, g, b := unpackColor(v.colors[i])
			v.colors[i] = packColor(g, r, b)
		}
		midPoint = (v.minGreen + v.maxGreen) / 2
	default:
		longestDim = COMPONENT_BLUE
		// We need to do a RGB to BGR swap, or vice-versa
		for i := v.lowerIndex; i <= v.upperIndex; i++ {
			r, g, b := unpackColor(v.colors[i])
			v.colors[i] = packColor(b, g, r)
		}
		midPoint = (v.minBlue + v.maxBlue) / 2
	}
	splitPoint := v.lowerIndex
loop:
	for i := v.lowerIndex; i <= v.upperIndex; i++ {
		r, g, b := unpackColor(v.colors[i])
		switch longestDim {
		case COMPONENT_RED:
			if r >= midPoint {
				splitPoint = i
				break loop
			}
		case COMPONENT_GREEN:
			if g >= midPoint {
				splitPoint = i
				break loop
			}
		case COMPONENT_BLUE:
			if b >= midPoint {
				splitPoint = i
				break loop
			}
		}
	}

	vbox := NewVbox(splitPoint+1, v.upperIndex, v.colors, v.populations)
	v.upperIndex = splitPoint
	v.fitBox()
	return vbox
}

func (v *Vbox) AverageColor() *Swatch {
	sumRed := 0
	sumGreen := 0
	sumBlue := 0
	pop := 0
	for i := v.lowerIndex; i <= v.upperIndex; i++ {
		color := v.colors[i]
		r, g, b := unpackColor(color)
		pop += v.populations[color]
		sumRed += r
		sumGreen += g
		sumBlue += b
	}
	avgRed := int(math.Floor(float64(sumRed) / float64(pop)))
	avgGreen := int(math.Floor(float64(sumGreen) / float64(pop)))
	avgBlue := int(math.Floor(float64(sumBlue) / float64(pop)))

	return &Swatch{Color: packColor(avgRed, avgGreen, avgBlue), Population: pop}
}
