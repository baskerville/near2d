package near2d

import "testing"

func TestNearestNeighbor(t *testing.T) {
	tree := NewTree(0, 0, 10, 10)
	tree.Add(Pt(2, 1))
	tree.Add(Pt(1, 5))
	tree.Add(Pt(4, 8))
	tree.Add(Pt(5, 4))
	tree.Add(Pt(7, 3))
	tree.Add(Pt(7, 5))
	tree.Add(Pt(8, 0))
	tree.Add(Pt(9, 9))
	n, _ := tree.NearestNeighbor(Pt(7, 2))
	if n != Pt(7, 3) {
		t.Errorf("Got %v, was expecting %v.\n", n, Pt(7, 3))
	}
}
