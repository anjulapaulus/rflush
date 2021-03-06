package rflush

const (
	maxEntries = 16
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
	Data      interface{}
}


func (r *RTree) Insert(min, max [2]float64,value interface{}) {
	var item BBox
	set(min,max,value, &item)
	r.insert(&item)
}

func (r *RTree) insert(item *BBox) {
	if r.root.Data == nil{
		set(item.Min,item.Max,new(node), &r.root)
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

func set(min,max [2]float64,value interface{}, bbox *BBox) {
	bbox.Min = min
	bbox.Max = max
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
	iter func(min, max [2]float64,data interface{}) bool,
) {
	var target BBox
	set(min, max,nil, &target)
	r.search(&target, iter)
}

func (r *RTree) search (
	target *BBox,
	iter func(min, max [2]float64,data interface{}) bool,
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


func (b *BBox) search (
	target *BBox, height int,
	iter func(min, max [2]float64,Data interface{}) bool,
) bool {
	n := b.Data.(*node)
	if height == 0 {
		for i := 0; i < n.count; i++ {
			child := n.children[i]
			if target.intersects(&child) {
				if !iter(child.Min, child.Max,child.Data) {
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
						if !iter(cn.children[i].Min, cn.children[i].Max,child.Data) {
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
func (r *RTree) All(iter func(min, max [2]float64,data interface{}) bool){
	if r.root.Data == nil {
		return
	}
	r.root.all(r.height,iter)
}

func (b *BBox) all(
	height int,
	iter func(min, max [2]float64,value interface{}) bool,
	)bool{
	n := b.Data.(*node)
	if height == 0 {
		for i := 0; i < n.count; i++ {
			if !iter(n.children[i].Min, n.children[i].Max,n.children[i].Data) {
				return false
			}
		}
	} else if height == 1 {
		for i := 0; i < n.count; i++ {
			cn := n.children[i].Data.(*node)
			for j := 0; j < cn.count; j++ {
				if !iter(cn.children[i].Min, cn.children[j].Max, cn.children[j].Data) {
					return false
				}
			}
		}
	} else {
		for i := 0; i < n.count; i++ {
			if !n.children[i].all(height-1, iter) {
				return false
			}
		}
	}
	return true
}

// Bounds returns the minimum bounding rect
func (r *RTree) Bounds() (min, max [2]float64) {
	if r.root.Data == nil {
		return
	}
	return r.root.Min, r.root.Max
}

// Len returns the number of items in the tree

func (r *RTree) Len() int{
	return r.count
}

// Remove data from tree
func (r *RTree) Remove(min, max [2]float64, data interface{}) {
	var item BBox
	set(min, max,data, &item)
	if r.root.Data == nil || !r.root.contains(&item) {
		return
	}
	var removed, recalced bool
	removed, recalced, r.reinsert =
		r.root.remove(&item, r.height, r.reinsert[:0])
	if !removed {
		return
	}
	r.count -= len(r.reinsert) + 1
	if r.count == 0 {
		r.root = BBox{}
		recalced = false
	} else {
		for r.height > 0 && r.root.Data.(*node).count == 1 {
			r.root = r.root.Data.(*node).children[0]
			r.height--
			r.root.adjustParentBBoxes()
		}
	}
	if recalced {
		r.root.adjustParentBBoxes()
	}
	for i := range r.reinsert {
		r.insert(&r.reinsert[i])
		r.reinsert[i].Data = nil
	}
}

func (b *BBox) remove(item *BBox, height int, reinsert []BBox) (
	removed, recalced bool, reinsertOut []BBox,
) {
	n := b.Data.(*node)
	if height == 0 {
		for i := 0; i < n.count; i++ {
			if n.children[i].Data == item.Data {
				// found the target item to remove
				recalced = b.onEdge(&n.children[i])
				n.children[i] = n.children[n.count-1]
				n.children[n.count-1].Data = nil
				n.count--
				if recalced {
					b.adjustParentBBoxes()
				}
				return true, recalced, reinsert
			}
		}
	} else {
		for i := 0; i < n.count; i++ {
			if !n.children[i].contains(item) {
				continue
			}
			removed, recalced, reinsert =
				n.children[i].remove(item, height-1, reinsert)
			if !removed {
				continue
			}
			if n.children[i].Data.(*node).count < minEntries {
				// underflow
				if !recalced {
					recalced = b.onEdge(&n.children[i])
				}
				reinsert = n.children[i].flatten(reinsert, height-1)
				n.children[i] = n.children[n.count-1]
				n.children[n.count-1].Data = nil
				n.count--
			}
			if recalced {
				b.adjustParentBBoxes()
			}
			return removed, recalced, reinsert
		}
	}
	return false, false, reinsert
}


// flatten flattens all leaf children into a single list
func (b *BBox) flatten(all []BBox, height int) []BBox {
	n := b.Data.(*node)
	if height == 0 {
		all = append(all, n.children[:n.count]...)
	} else {
		for i := 0; i < n.count; i++ {
			all = n.children[i].flatten(all, height-1)
		}
	}
	return all
}

// onEdge returns true when b2 is on the edge of b
func (b *BBox) onEdge(b2 *BBox) bool {
	if b.Min[0] == b2.Min[0] || b.Max[0] == b2.Max[0] {
		return true
	}
	if b.Min[1] == b2.Min[1] || b.Max[1] == b2.Max[1] {
		return true
	}
	return false
}

// Children is a utility function that returns all children for parent node.
// If parent node is nil then the root nodes should be returned. The min, max,
// data, and items slices all must have the same lengths. And, each element
// from all slices must be associated. Returns true for `items` when the the
// item at the leaf level. The reuse buffers are empty length slices that can
// optionally be used to avoid extra allocations.
func (r *RTree) Children(
	parent interface{},
	reuse []Child,
) []Child {
	children := reuse
	if parent == nil {
		if r.Len() > 0 {
			// fill with the root
			children = append(children, Child{
				Min:  r.root.Min,
				Max:  r.root.Max,
				Data: r.root.Data,
				Item: false,
			})
		}
	} else {
		// fill with child items
		n := parent.(*node)
		item := true
		if n.count > 0 {
			if _, ok := n.children[0].Data.(*node); ok {
				item = false
			}
		}
		for i := 0; i < n.count; i++ {
			children = append(children, Child{
				Min:  n.children[i].Min,
				Max:  n.children[i].Max,
				Data: n.children[i].Data,
				Item: item,
			})
		}
	}
	return children
}


// Replace an item in the structure. This is effectively just a Remove
// followed by an Insert.
func (r *RTree) Replace(
	oldMin, oldMax [2]float64,oldData interface{},
	newMin, newMax [2]float64,newData interface{},
) {
	r.Remove(oldMin, oldMax, oldData)
	r.Insert(newMin, newMax,newData)
}

// NewRect constructs and returns a pointer to a Rect given a corner point and
// the lengths of each dimension.  The point p should be the most-negative point
// on the rectangle (in every dimension) and every length should be positive.
func NewBBox(p [2]float64, lengths [2]float64) (r *BBox) {
	r = new(BBox)
	r.Min = p

	r.Max = [2]float64{}
	for i := range p {
		if lengths[i] <= 0 {
			return
		}

		r.Max[i] = p[i] + lengths[i]
	}
	return
}

//PointToBBox returns BBox from point given
func PointToBBox(p [2]float64, tol float64) *BBox {
	var a, b [2]float64
	for i := range p {
		a[i] = p[i] - tol
		b[i] = p[i] + tol
	}
	return &BBox{a, b,nil}
}
