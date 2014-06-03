package near2d

import "math"

const (
	horizontal = iota
	vertical
)

type Point struct {
	X, Y float64
}

type Rectangle struct {
	Min, Max Point
}

type Split struct {
	kind byte
	at   float64
}

type Tree struct {
	rect        Rectangle
	point       *Point
	split       *Split
	firstChild  *Tree
	secondChild *Tree
}

func sq(x float64) float64 {
	return x * x
}

func Rect(x0, y0, x1, y1 float64) Rectangle {
	return Rectangle{Point{x0, y0}, Point{x1, y1}}
}

func Pt(x, y float64) Point {
	return Point{x, y}
}

func (p0 Point) Add(p1 Point) Point {
	return Point{p0.X + p1.X, p0.Y + p1.Y}
}

// Simplified distance
func (p0 Point) dist(p1 Point) float64 {
	return sq(p0.X-p1.X) + sq(p0.Y-p1.Y)
}

func (p Point) dist2(r Rectangle) float64 {
	cx := math.Max(r.Min.X, math.Min(p.X, r.Max.X))
	cy := math.Max(r.Min.Y, math.Min(p.Y, r.Max.Y))
	return p.dist(Pt(cx, cy))
}

func NewTree(x1, y1, x2, y2 float64) *Tree {
	return &Tree{rect: Rect(x1, y1, x2, y2)}
}

func (t *Tree) Add(p Point) {
	if t.split == nil && t.point == nil {
		t.point = &p
	} else if t.split == nil {
		var delta float64
		dX := p.X - t.point.X
		dY := p.Y - t.point.Y
		if math.Abs(dY) > math.Abs(dX) {
			t.split = &Split{horizontal, (t.point.Y + p.Y) / 2}
			t.firstChild = NewTree(t.rect.Min.X, t.rect.Min.Y, t.rect.Max.X, t.split.at)
			t.secondChild = NewTree(t.rect.Min.X, t.split.at, t.rect.Max.X, t.rect.Max.Y)
			delta = dY
		} else {
			t.split = &Split{vertical, (t.point.X + p.X) / 2}
			t.firstChild = NewTree(t.rect.Min.X, t.rect.Min.Y, t.split.at, t.rect.Max.Y)
			t.secondChild = NewTree(t.split.at, t.rect.Min.Y, t.rect.Max.X, t.rect.Max.Y)
			delta = dX
		}
		if delta > 0 {
			t.firstChild.point = t.point
			t.secondChild.point = &p
		} else {
			t.firstChild.point = &p
			t.secondChild.point = t.point
		}
		t.point = nil
	} else {
		p.nearestChild(t).Add(p)
	}
}

func (p Point) nearestChild(t *Tree) *Tree {
	var at float64
	switch t.split.kind {
	case horizontal:
		at = p.Y
	case vertical:
		at = p.X
	}
	if at < t.split.at {
		return t.firstChild
	} else {
		return t.secondChild
	}
}

func (p Point) nearestPoint(t *Tree, nearest *Point, dmin float64, remains []*Tree) (*Point, float64, []*Tree) {
	if t.point != nil {
		d := p.dist(*t.point)
		if d < dmin {
			return t.point, d, remains
		} else {
			return nearest, dmin, remains
		}
	} else {
		child := p.nearestChild(t)
		if child == t.firstChild {
			remains = append(remains, t.secondChild)
		} else {
			remains = append(remains, t.firstChild)
		}
		return p.nearestPoint(child, nearest, dmin, remains)
	}
}

func (t *Tree) NearestNeighbor(p Point) (Point, float64) {
	var (
		remains []*Tree
		nearest *Point
		dmin    = math.Inf(1)
	)
	nearest, dmin, remains = p.nearestPoint(t, nearest, dmin, remains)
	for len(remains) > 0 {
		last := remains[len(remains)-1]
		remains = remains[:len(remains)-1]
		if p.dist2(last.rect) < dmin {
			nearest, dmin, remains = p.nearestPoint(last, nearest, dmin, remains)
		}
	}
	return *nearest, dmin
}
