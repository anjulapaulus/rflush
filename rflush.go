package rflush

import (
	"fmt"
	"math"
	"sort"
)

type RTree struct {
	MinChildren float64
	MaxChildren float64
	root        node
	size        int
	height      int
}

// node represents a tree node of an Rtree.
type node struct {
	leaf      bool
	children  []node
	count     int
	item	  Item
}
type Item struct {
	MinX      float64
	MinY      float64
	MaxX      float64
	MaxY      float64
	Reference string
}

// NewTree returns an Rtree.
func NewRTree(maxEntries int) *RTree {
	if maxEntries != 0 {
		maxEntriesNode := math.Max(4, float64(maxEntries))
		minEntriesNode := math.Max(2, math.Ceil(maxEntriesNode*0.4))
		return &RTree{
			MinChildren: minEntriesNode,
			MaxChildren: maxEntriesNode,
		}
	}
	maxEntriesNode := math.Max(4, 9)
	minEntriesNode := math.Max(2, math.Ceil(maxEntriesNode*0.4))
	return &RTree{
		MinChildren: minEntriesNode,
		MaxChildren: maxEntriesNode,
	}
}


func (r *RTree) Insert(minX, minY, maxX, maxY float64, reference string) {
	var item node
	set(minX,minY,maxX,maxY,reference,&item)
	r.insert(item, r.height - 1, false)
}

func (r *RTree) insert(item node, level int, isNode bool) {
	var bbox node
	var insertPath []node
	if isNode{
		bbox = item
	}else {
		bbox = item.toBBox()
	}

	// find the best node for accommodating the item, saving all nodes along the path too
	chosenNode := r.root.chooseSubtree(bbox,level,insertPath)

	// put the item into the node
	chosenNode.children = append(chosenNode.children, item)
	chosenNode.extend(&bbox)

	// split on node overflow; propagate upwards if necessary
	for level >=0 {
		if len(insertPath[level].children) > int(r.MaxChildren) {
			r.split(insertPath,level)
			level--
		}else {
			break
		}
	}
	// adjust bboxes along the insertion path
	r.adjustParentBBoxes(&bbox,insertPath,level)
}


func (n *node) chooseSubtree(bbox node, level int,path []node) (chosenNode *node){
	for {
		path = append(path, bbox)
		if n.leaf || len(path) - 1 == level{
			break
		}
		minEnlargement := 0.0
		minArea := 0.0
		var targetNode *node

		for i := 0; i < len(n.children); i++ {
			child := n.children[i]
			area := child.bboxArea()
			enlargement := bbox.enlargedArea(&child) - area

			// choose entry with the least area enlargement
			if enlargement < minEnlargement {
				minEnlargement = enlargement
				if area < minArea {
					minArea = area
				}
				targetNode = &child
			} else if enlargement == minEnlargement {
				// otherwise choose one with the smallest area
				if area < minArea {
					minArea = area
					targetNode = &child
				}
			}
		}
		if targetNode != nil {
			chosenNode = targetNode
		}else {
			chosenNode = &n.children[0]
		}
	}
	return chosenNode
}
func (r *RTree ) adjustParentBBoxes(bbox *node, path []node, count int) {
	// adjust bboxes along the given tree path
	for i := count; i >= 0; i-- {
		path[i].extend(bbox)
	}
}

func set(minX, minY, maxX, maxY float64, reference string, bbox *node){
	bbox.item.MinX = minX
	bbox.item.MinY = minY
	bbox.item.MaxX = maxX
	bbox.item.MaxY = maxY
	bbox.item.Reference = reference
}

func (n node) toBBox() node{
	return n
}
func (n *node) bboxArea() float64 {
	return (n.item.MaxX - n.item.MinX) * (n.item.MaxY - n.item.MinY)
}

func (n *node) enlargedArea(t *node) float64 {
	return (math.Max(t.item.MaxX, n.item.MaxX) - math.Min(t.item.MinX, n.item.MinX)) *
		(math.Max(t.item.MaxY, n.item.MaxY) - math.Min(t.item.MinY, n.item.MinY))
}

func (n *node) extend(b *node) {
	n.item.MinX = math.Min(n.item.MinX, b.item.MinX)
	n.item.MinY = math.Min(n.item.MinY, b.item.MinY)
	n.item.MaxX = math.Max(n.item.MaxX, b.item.MaxX)
	n.item.MaxY = math.Max(n.item.MaxY, b.item.MaxY)
}

// split overflowed node into two
func (r *RTree) split(insertPath []node, level int) {
	node := insertPath[level]
	M := len(node.children)
	m := int(r.MinChildren)

	r.chooseSplitAxis(&node, m, M)

	//splitIndex := r.chooseSplitIndex(&node, m, M)
	//
	//newNode := createNode(node.children[:splitIndex])
	//newNode.count = node.count
	//newNode.leaf = node.leaf
	//
	//node.calcBBox()
	//newNode.calcBBox()
	//
	//if level != 0 {
	//	insertPath[level-1].children = append(insertPath[level-1].children, newNode)
	//} else {
	//	r.splitRoot(node, newNode)
	}

}

func (r *RTree) splitRoot(oldNode, newNode node) {
	var children []node
	children = append(children, oldNode, newNode)
	r.root = createNode(children)
	r.height = oldNode.count + 1
	r.root.leaf = false
	r.root.calcBBox()
}

func (r *RTree) chooseSplitAxis(checkNode *node, m, M int) {

	xMargin := allDistMargin(checkNode, m, M, 1)
	yMargin := allDistMargin(checkNode, m, M, 2)

	// if total distributions margin value is minimal for x, sort by minX,
	// otherwise it's already sorted by minY
	if xMargin < yMargin {
		sort.SliceStable(checkNode.children, func(i, j int) bool {
			return checkNode.children[i].item.MinX < checkNode.children[j].item.MinX
		})
	}
}


// total margin of all possible split distributions where each node is at least m full
func allDistMargin(checkNode *node, m, M int, compare int) float64 {

	if compare == 1 {
		sort.SliceStable(checkNode.children, func(i, j int) bool {
			return checkNode.children[i].item.MinX < checkNode.children[j].item.MinX
		})
	} else if compare == 2 {
		sort.SliceStable(checkNode.children, func(i, j int) bool {
			return checkNode.children[i].item.MinY < checkNode.children[j].item.MinY
		})
	}

	leftBBox := distBBox(checkNode, 0, m, nil)
	rightBBox := distBBox(checkNode, M-m, M, nil)
	margin := bboxMargin(&leftBBox) + bboxMargin(&rightBBox)

	for i := m; i < M-m; i++ {
		child := checkNode.children[i]
		leftBBox.extend(&child)
		margin += bboxMargin(&leftBBox)
	}

	for i := M - m - 1; i >= m; i-- {
		child := checkNode.children[i]
		rightBBox.extend(&child)
		margin += bboxMargin(&rightBBox)
	}
	return margin
}



func distBBox(checkNode *node, k, p int, destNode *node) node {
	var destiNode node
	if destNode == nil {
		destiNode = createNode(nil)
	}

	for i := k; i < p; i++ {
		child := checkNode.children[i]
		destiNode.extend(&child)
	}
	return destiNode
}

func bboxMargin(a *node) float64 {
	return (a.item.MaxX - a.item.MinX) + (a.item.MaxY - a.item.MinY)
}

func createNode(children []node) node {
	return node{
		leaf:     true,
		children: children,
	}
}

func (r *RTree) chooseSplitIndex(checkNode *node, m, M int) int {
	var index int
	minOverlap := 0.0
	minArea := 0.0
	for i := m; i <= M-m; i++ {
		bbox1 := distBBox(checkNode, 0, i, nil)
		bbox2 := distBBox(checkNode, i, M, nil)
		overlap := bbox1.intersectionArea(&bbox2)
		area := bbox1.bboxArea() + bbox2.bboxArea()

		// choose distribution with minimum overlap
		if overlap < minOverlap {
			minOverlap = overlap
			index = i

			if area < minArea {
				minArea = area
			}

		} else if overlap == minOverlap {
			// otherwise choose distribution with minimum area
			if area < minArea {
				minArea = area
				index = i
			}
		}
	}
	if index != 0 {
		return index
	} else {
		return M - m
	}
}

func (n *node) intersectionArea(b *node) float64 {
	minX := math.Max(n.item.MinX, b.item.MinX)
	minY := math.Max(n.item.MinY, b.item.MinY)
	maxX := math.Min(n.item.MaxX, b.item.MaxX)
	maxY := math.Min(n.item.MaxY, b.item.MaxY)

	return math.Max(0, maxX-minX) *
		math.Max(0, maxY-minY)
}

// calculate node's bbox from bboxes of its children
func (n *node) calcBBox() {
	distBBox(n, 0, len(n.children), n)
}

func (n *node) contains(b *node) bool{
	if b.MinX < n.MinX || b.MaxX > n.MaxX {
		return false
	}
	if b.MinY < n.MinY || b.MaxY > n.MaxY {
		return false
	}
	return true
}
