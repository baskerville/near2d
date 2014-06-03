package near2d

import "testing"

func TestNearestNeighbor(t *testing.T) {
	tr := NewTree(0, 0, 10, 10)
	tr.Add(Pt(2, 1))
	tr.Add(Pt(1, 5))
	tr.Add(Pt(4, 8))
	tr.Add(Pt(5, 4))
	tr.Add(Pt(7, 3))
	tr.Add(Pt(7, 5))
	tr.Add(Pt(8, 0))
	tr.Add(Pt(9, 9))
	n, _ := tr.NearestNeighbor(Pt(7, 2))
	if n != Pt(7, 3) {
		t.Errorf("Got %v, was expecting %v.\n", n, Pt(7, 3))
	}
}
