RFlush
======

RFlush is high performance 2D RTree spatial indexing structure for points and rectangles written in golang.

## Usage

### Install
````
go get github.com/anjulapaulus/rflush
````

### Creating a Tree
````
var tr rflush.RTree
````

### Insert Items

Rectangles are data structures for representing spatial objects, while Points represent spatial locations.
````
//Insert a point
tr.Insert([2]float64{-112.0078, 33.4373}, [2]float64{-112.0078, 33.4373}, "Somewhere")

//Insert a rectangle
tr.Insert([2]float64{10, 10}, [2]float64{20, 20}, "rect")

````
A helper function NewBBox helps you create a rectangle by specifying locations and lengths of sides or you could use the PointToBBox function by providing the length.

````
bbox := rflush.NewBBox([2]float64{-112.0078, 33.4373},[2]float64{1, 2})
fmt.Println(bbox)
//[-112.0078 33.4373] [-111.0078 35.4373]
	

pointToBBox := rflush.PointToBBox([2]float64{-112.0078, 33.4373},0.01)
fmt.Println(pointToBBox)
//[-112.01780000000001 33.4273] [-111.9978 33.4473]

````
### Search Items

By providing an iterator function and providing bounding box intersects you could get the items within the specified bounding box area.
````
tr.Search([2]float64{-112.1, 33.4}, [2]float64{-112.0, 33.5},
		func(min, max [2]float64, value interface{}) bool {
			println(value.(string)) // prints "PHX"
			return true
		},
	)
````

### Remove Item

Remove a previously inserted item.


````
tr.Remove([2]float64{-112.0078, 33.4373}, [2]float64{-112.0078, 33.4373}, "Somewhere")
````

### Replace Item

Replace function removes the old location from the tree and inserts the new location based on reference.

````
tr.Replace([2]float64{-112.0078, 33.4373}, [2]float64{-112.0078, 33.4373}, "Somewhere",
		       [2]float64{-113.0078, 34.4373}, [2]float64{-113.0078, 34.4373}, "Somewhere")
````

### Get all items in tree

This iterates through the tree and passes values to the iterator function provided in no specific order.

````
tr.All(func(min, max [2]float64, data interface{}) bool {
		fmt.Println(data.(string))
		return true
	})
````


#### K-Nearest Neighbors
For "k nearest neighbors around a point" type of queries for RBush. [rflush_knn](github.com/anjulapaulus/rflush_knn)


## Performance

The benchmarks were conducted on a Macbook Pro 2019 RAM 8GB Mac Os Catalina for max entries to be 16.

Test                         | RBush  | [old RTree](https://github.com/imbcmdth/RTree) | RFlush
---------------------------- | ------ | ------ | -------
insert 1M items one by one   | 3.18s  | 7.83s  | 1.321s
1000 searches                | 0.03s  | 0.93s  | 0.001s
remove 1000 items one by one | 0.02s  | 1.18s  | 0.001s


## Algorithms Used
* Single Insertion : The insertion is the same as the algorithm defined in Guttman’s Rtree. From the root to the leaf, the boxes which will incur the least enlargement are chosen. Ties go to boxes with the smallest area.

* Node Splitting : The node splitting is based on R* Tree algorithm Greene's split which acts as the overflow treatment when a boundary box has reached maximum entries. The Split axis is calculated and the box is split into two small boxes. The children are distributed into the relevant small box according to the distance from parent’s minimum(min-dist) and maximum(max-dist) values of its axis to the child’s minimum(min-dist) and maximum(max-dist) values.
If min-dist < max-dist, then the child is placed in the left box.
If min-dist > max-dist, then the child is placed in the right box.
If min-dist = max-dist, then the child is kept until all children are evaluated and placed in the box with less children.

* Single Deletion : The deletion is the same as the algorithm defined by Guttman.The item is directly removed from the tree. When the number of children in a box falls below its minimum entries, it is removed from the tree and it's items are re-inserted.

* Search : The same algorithm suggested by Guttman.


## Papers
* [R-trees: a Dynamic Index Structure For Spatial Searching](http://www-db.deis.unibo.it/courses/SI-LS/papers/Gut84.pdf)
* [The R*-tree: An Efficient and Robust Access Method for Points and Rectangles+](http://dbs.mathematik.uni-marburg.de/publications/myPapers/1990/BKSS90.pdf)
* [R-Trees: Theory and Applications (book)](http://www.apress.com/9781852339777)
* [R-Trees: Theory and Applications (book)](http://www.apress.com/9781852339777)
