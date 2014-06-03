package near2d

import "math"

const (
	horizontal = iota
	vertical
)

type Point struct {
	X, Y float64
}

type rectangle struct {
	Min, Max Point
}

type split struct {
	kind byte
	at   float64
}

type Tree struct {
	rect        rectangle
	point       *Point
	fence       *split
	firstChild  *Tree
	secondChild *Tree
}

func sq(x float64) float64 {
	return x * x
}

func rect(x0, y0, x1, y1 float64) rectangle {
	return rectangle{Point{x0, y0}, Point{x1, y1}}
}

func Pt(x, y float64) Point {
	return Point{x, y}
}

// Simplified distance
func (p0 Point) dist(p1 Point) float64 {
	return sq(p0.X-p1.X) + sq(p0.Y-p1.Y)
}

func (p Point) dist2(r rectangle) float64 {
	cx := math.Max(r.Min.X, math.Min(p.X, r.Max.X))
	cy := math.Max(r.Min.Y, math.Min(p.Y, r.Max.Y))
	return p.dist(Pt(cx, cy))
}

func NewTree(x1, y1, x2, y2 float64) *Tree {
	return &Tree{rect: rect(x1, y1, x2, y2)}
}

func (t *Tree) Add(p Point) {
	if t.fence == nil && t.point == nil {
		t.point = &p
	} else if t.fence == nil {
		var delta float64
		dX := p.X - t.point.X
		dY := p.Y - t.point.Y
		if math.Abs(dY) > math.Abs(dX) {
			t.fence = &split{horizontal, (t.point.Y + p.Y) / 2}
			t.firstChild = NewTree(t.rect.Min.X, t.rect.Min.Y, t.rect.Max.X, t.fence.at)
			t.secondChild = NewTree(t.rect.Min.X, t.fence.at, t.rect.Max.X, t.rect.Max.Y)
			delta = dY
		} else {
			t.fence = &split{vertical, (t.point.X + p.X) / 2}
			t.firstChild = NewTree(t.rect.Min.X, t.rect.Min.Y, t.fence.at, t.rect.Max.Y)
			t.secondChild = NewTree(t.fence.at, t.rect.Min.Y, t.rect.Max.X, t.rect.Max.Y)
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
	switch t.fence.kind {
	case horizontal:
		at = p.Y
	case vertical:
		at = p.X
	}
	if at < t.fence.at {
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
