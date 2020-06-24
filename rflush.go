package main

const (
	maxEntries = 40
	minEntries = maxEntries * 40 / 100
)

type RTree struct {
	height   int
	root     BBox
	count    int
	reinsert []BBox
}

// node represents a tree node of an Rtree.
type node struct {
	count    int
	children [maxEntries + 1]BBox
}
//bbox represents a boundary box
type BBox struct {
	Min, Max  [2]float64
	Reference string
	Data      interface{}
}


func (r *RTree) Insert(min, max [2]float64, reference string, value interface{}) {
	var item BBox
	set(min,max, reference,value, &item)
	r.insert(&item)
}

func (r *RTree) insert(item *BBox) {
	if r.root.Data == nil{
		set(item.Min,item.Max, item.Reference, new(node), &r.root)
	}
	grown := r.root.insert(item, r.height)
	if grown {
		r.root.extend(item)
	}
	if r.root.Data.(*node).count == maxEntries+1 {
		newRoot := new(node)
		r.root.split(&newRoot.children[1])
		newRoot.children[0] = r.root
		newRoot.count = 2
		r.root.Data = newRoot
		r.root.adjustParentBBoxes()
		r.height++
	}
	r.count++

}

func (b *BBox) insert (item *BBox, height int)(grown bool){
	n := b.Data.(*node)
	if height == 0 {
		n.children[n.count] = *item
		n.count++
		grown = !b.contains(item)
		return grown
	}

	// ind the best node index for accommodating the item
	index := b.chooseSubtree(item)
	child := &n.children[index]
	grown = child.insert(item, height-1)
	if grown {
		child.extend(item)
		grown = !b.contains(item)
	}
	// split on node overflow; propagate upwards if necessary
	if child.Data.(*node).count == maxEntries+1 {
		child.split(&n.children[n.count])
		n.count++
	}
	return grown
}

func (b *BBox) chooseSubtree(b2 *BBox) int {
	k := -1
	minEnlargement := 0.0
	minArea := 0.0
	n := b.Data.(*node)
	for i := 0; i < n.count; i++ {
		child := n.children[i]
		area := child.bboxArea()
		enlargement := child.enlargedArea(b2) - area
		// choose entry with the least area enlargement
		if k == -1 || enlargement < minEnlargement {
			k = i
			minEnlargement = enlargement
			minArea = area
		} else if enlargement == minEnlargement {
			// otherwise choose one with the smallest area
			if area < minArea {
				k = i
				minEnlargement = enlargement
				minArea = area
			}
		}
	}
	return k
}

func (b *BBox) bboxArea() float64 {
	return (b.Max[0] - b.Min[0]) * (b.Max[1] - b.Min[1])
}

func (b *BBox) enlargedArea(b2 *BBox) float64 {
	area := 1.0
	if b2.Max[0] > b.Max[0] {
		if b2.Min[0] < b.Min[0] {
			area *= b2.Max[0] - b2.Min[0]
		} else {
			area *= b2.Max[0] - b.Min[0]
		}
	} else {
		if b2.Min[0] < b.Min[0] {
			area *= b.Max[0] - b2.Min[0]
		} else {
			area *= b.Max[0] - b.Min[0]
		}
	}
	if b2.Max[1] > b.Max[1] {
		if b2.Min[1] < b.Min[1] {
			area *= b2.Max[1] - b2.Min[1]
		} else {
			area *= b2.Max[1] - b.Min[1]
		}
	} else {
		if b2.Min[1] < b.Min[1] {
			area *= b.Max[1] - b2.Min[1]
		} else {
			area *= b.Max[1] - b.Min[1]
		}
	}
	return area
}

func set(min,max [2]float64, reference string, value interface{}, bbox *BBox) {
	bbox.Min = min
	bbox.Max = max
	bbox.Reference = reference
	bbox.Data = value
}

// contains return struct when b is fully contained inside of n
func (b *BBox) contains(b2 *BBox) bool{
	if b2.Min[0] < b.Min[0] || b2.Max[0] > b.Max[0] {
		return false
	}
	if b2.Min[1] < b.Min[1] || b2.Max[1] > b.Max[1] {
		return false
	}
	return true
}


func (b *BBox) extend(b2 *BBox) {
	if b2.Min[0] < b.Min[0] {
		b.Min[0] = b2.Min[0]
	}
	if b2.Max[0] > b.Max[0] {
		b.Max[0] = b2.Max[0]
	}
	if b2.Min[1] < b.Min[1] {
		b.Min[1] = b2.Min[1]
	}
	if b2.Max[1] > b.Max[1] {
		b.Max[1] = b2.Max[1]
	}
}

func (b *BBox) split(right *BBox) {
	axis, _ := b.chooseSplitAxis()
	left := b
	leftNode := left.Data.(*node)
	rightNode := new(node)
	right.Data = rightNode

	var equals []BBox
	for i := 0; i < leftNode.count; i++ {
		minDist := leftNode.children[i].Min[axis] - left.Min[axis]
		maxDist := left.Max[axis] - leftNode.children[i].Max[axis]
		if minDist < maxDist {
			// stay left
		} else {
			if minDist > maxDist {
				// move to right
				rightNode.children[rightNode.count] = leftNode.children[i]
				rightNode.count++
			} else {
				// move to equals, at the end of the left array
				equals = append(equals, leftNode.children[i])
			}
			leftNode.children[i] = leftNode.children[leftNode.count-1]
			leftNode.children[leftNode.count-1].Data = nil
			leftNode.count--
			i--
		}
	}
	for _, b := range equals {
		if leftNode.count < rightNode.count {
			leftNode.children[leftNode.count] = b
			leftNode.count++
		} else {
			rightNode.children[rightNode.count] = b
			rightNode.count++
		}
	}
	left.adjustParentBBoxes()
	right.adjustParentBBoxes()
}

func (b *BBox) chooseSplitAxis() (axis int, size float64) {
	if b.Max[1]-b.Min[1] > b.Max[0]-b.Min[0] {
		return 1, b.Max[1] - b.Min[1]
	}
	return 0, b.Max[0] - b.Min[0]
}

func (b *BBox) adjustParentBBoxes() {
	n := b.Data.(*node)
	b.Min = n.children[0].Min
	b.Max = n.children[0].Max
	for i := 1; i < n.count; i++ {
		b.extend(&n.children[i])
	}
}

// Search ...
func (r *RTree) Search (
	min, max [2]float64,
	iter func(min, max [2]float64,reference string) bool,
) {
	var target BBox
	set(min, max, "",nil, &target)
	r.search(&target, iter)
}

func (r *RTree) search (
	target *BBox,
	iter func(min, max [2]float64,reference string) bool,
) {
	if r.root.Data == nil {
		return
	}
	if target.intersects(&r.root) {
		r.root.search(target, r.height, iter)
	}
}

func (b *BBox) intersects(a *BBox) bool {
	if a.Min[0] > b.Max[0] || a.Max[0] < b.Min[0] {
		return false
	}
	if a.Min[1] > b.Max[1] || a.Max[1] < b.Min[1] {
		return false
	}
	return true
}


func (b *BBox) search(
	target *BBox, height int,
	iter func(min, max [2]float64,reference string) bool,
) bool {
	n := b.Data.(*node)
	if height == 0 {
		for i := 0; i < n.count; i++ {
			child := n.children[i]
			if target.intersects(&child) {
				if !iter(child.Min, child.Max,child.Reference) {
					return false
				}
			}
		}
	} else if height == 1 {
		for i := 0; i < n.count; i++ {
			child := n.children[i]
			if target.intersects(&child) {
				cn := child.Data.(*node)
				for i := 0; i < cn.count; i++ {
					if target.intersects(&cn.children[i]) {
						if !iter(cn.children[i].Min, cn.children[i].Max,child.Reference) {
							return false
						}
					}
				}
			}
		}
	} else {
		for i := 0; i < n.count; i++ {
			child := n.children[i]
			if target.intersects(&child) {
				if !child.search(target, height-1, iter) {
					return false
				}
			}
		}
	}
	return true
}

// All function returns all the entries in the tree.
func (r *RTree) All() []BBox {
	if r.root.Reference == "" {
		return nil
	}
	all := r.root.all(r.height)
	return all
}

func (b *BBox) all(height int) []BBox {
	var all []BBox
	n := b.Data.(*node)
	if height == 0 {
		for i := 0; i < n.count; i++ {
			all = append(all, n.children[i])
		}
	} else if height == 1 {
		for i := 0; i < n.count; i++ {
			cn := n.children[i].Data.(*node)
			for j := 0; j < cn.count; j++ {
				all = append(all, cn.children[i])
			}
		}
	} else {
		for i := 0; i < n.count; i++ {
			n.children[i].all(height - 1)
		}
	}
	return all
}

// Bounds returns the minimum bounding rect
func (r *RTree) Bounds() (min, max [2]float64) {
	if r.root.Reference == "" {
		return
	}
	return r.root.Min, r.root.Max
}

// Len returns the number of items in the tree

func (r *RTree) Len() int{
	return r.count
}
