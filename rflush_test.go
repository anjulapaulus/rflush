package rflush

import (
	"fmt"
	"github.com/anjulapaulus/rflush/test_data"
	"github.com/tidwall/lotsa"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestRTree_Insert(t *testing.T) {
	var tr RTree

	test_data.LoadCSV()
	for i := 0; i < len(test_data.LocationsData.Locations); i++ {
		var point [2]float64
		latitude, _ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Latitude, 64)
		longitude, _ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Longitude, 64)
		point[0] = latitude
		point[1] = longitude
		tr.Insert(point, point, test_data.LocationsData.Locations[i].PlaceName)
	}

	if tr.Len() != len(test_data.LocationsData.Locations) {
		t.Error("Insert Function: Error", len(test_data.LocationsData.Locations))
	}
}

func TestRTree_Len(t *testing.T) {
	var tr RTree

	test_data.LoadCSV()
	for i := 0; i < len(test_data.LocationsData.Locations); i++ {
		var point [2]float64
		latitude, _ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Latitude, 64)
		longitude, _ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Longitude, 64)
		point[0] = latitude
		point[1] = longitude
		tr.Insert(point, point, test_data.LocationsData.Locations[i].PlaceName)
	}

	if tr.Len() != len(test_data.LocationsData.Locations) {
		t.Error("Insert Function Test Failed")
	}
}

func TestRTree_Bounds(t *testing.T) {
	var tr RTree

	test_data.LoadCSV()
	for i := 0; i < len(test_data.LocationsData.Locations); i++ {
		var point [2]float64
		latitude, _ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Latitude, 64)
		longitude, _ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Longitude, 64)
		point[0] = latitude
		point[1] = longitude
		tr.Insert(point, point, test_data.LocationsData.Locations[i].PlaceName)
	}

	min, max := tr.Bounds()
	if min != [2]float64{5.9, 75} && max != [2]float64{10.08333, 81.91667} {
		t.Error("Bounds Function Test: Failed")
	}
}

func TestRTree_All(t *testing.T) {
	var tr RTree

	test_data.LoadCSV()
	for i := 0; i < len(test_data.LocationsData.Locations); i++ {
		var point [2]float64
		latitude, _ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Latitude, 64)
		longitude, _ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Longitude, 64)
		point[0] = latitude
		point[1] = longitude

		tr.Insert(point, point, test_data.LocationsData.Locations[i].PlaceName)
	}

	var bbox []BBox
	tr.All(func(min, max [2]float64, data interface{}) bool {
		var item BBox
		item.Min = min
		item.Max = max
		item.Data = data
		bbox = append(bbox, item)
		return true
	})

	if len(bbox) != len(test_data.LocationsData.Locations) {
		t.Error("All Function Test: Failed", len(bbox))
	}
}

func TestRTree_Search(t *testing.T) {
	var tr RTree

	test_data.LoadCSV()
	for i := 0; i < len(test_data.LocationsData.Locations); i++ {
		var point [2]float64
		latitude, _ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Latitude, 64)
		longitude, _ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Longitude, 64)
		point[0] = latitude
		point[1] = longitude
		tr.Insert(point, point, test_data.LocationsData.Locations[i].PlaceName)
	}

	bbox := PointToBBox([2]float64{6.969658, 79.873308}, 0.2)
	var locationNames []interface{}
	tr.Search(bbox.Min, bbox.Max,
		func(min [2]float64, max [2]float64, data interface{}) bool {
			locationNames = append(locationNames, data)
			return true
		},
	)
	if len(locationNames) != 2234 {
		t.Error("Search Function: Error", len(locationNames))
	}
}

func TestRTree_Remove(t *testing.T) {
	var tr RTree

	test_data.LoadCSV()
	for i := 0; i < len(test_data.LocationsData.Locations); i++ {
		var point [2]float64
		latitude, _ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Latitude, 64)
		longitude, _ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Longitude, 64)
		point[0] = latitude
		point[1] = longitude
		tr.Insert(point, point, test_data.LocationsData.Locations[i].PlaceName)
	}

	tr.Remove([2]float64{7.68333, 80.1}, [2]float64{7.68333, 80.1}, "Kitagama")
	tr.Remove([2]float64{7.88333, 81.53333}, [2]float64{7.88333, 81.53333}, "Chunkankeni")
	tr.Remove([2]float64{6.06667, 80.46667}, [2]float64{6.06667, 80.46667}, "Zowdegala")
	tr.Remove([2]float64{7.5005, 80.3088}, [2]float64{7.5005, 80.3088}, "Yatawehera")
	tr.Remove([2]float64{7.4603, 80.0885}, [2]float64{7.4603, 80.0885}, "Yakarawatta")

	if tr.Len() != 56743 {
		t.Error("Remove Function Test: Failed", len(test_data.LocationsData.Locations))
	}
}

func TestRTree_Replace(t *testing.T) {
	var tr RTree

	test_data.LoadCSV()
	for i := 0; i < len(test_data.LocationsData.Locations); i++ {
		var point [2]float64
		latitude, _ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Latitude, 64)
		longitude, _ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Longitude, 64)
		point[0] = latitude
		point[1] = longitude
		tr.Insert(point, point, test_data.LocationsData.Locations[i].PlaceName)
	}

	var newPoint [2]float64
	newPoint[0] = 6.973597
	newPoint[1] = 79.869424
	tr.Replace([2]float64{7.32616, 80.62377}, [2]float64{7.32616, 80.62377}, "Katugasthota", newPoint, newPoint, "Katugasthota")

	var bbox []BBox
	tr.All(func(min, max [2]float64,data interface{}) bool {
		var item BBox
		item.Min = min
		item.Max = max
		item.Data = data
		bbox = append(bbox, item)
		return true
	})
	var checkItem BBox
	for i := 0; i < len(bbox); i++ {
		if bbox[i].Data == "Katugasthota" {
			checkItem.Data = bbox[i].Data
			checkItem.Min = bbox[i].Min
			checkItem.Max = bbox[i].Max
		}
	}

	if checkItem.Data == "Katugasthota" && checkItem.Min != newPoint && checkItem.Max != newPoint {
		t.Error("Replace Function Failed", checkItem.Data, checkItem.Min, checkItem.Max)
	}
}

func TestNewBBox(t *testing.T) {
	p := [2]float64{6.969658, 79.873308}
	lengths := [2]float64{2.5, 8.0}

	bbox := NewBBox(p, lengths)
	if bbox.Min == [2]float64{6.969658, 79.873308} && bbox.Max != [2]float64{9.469657999999999, 87.873308} {
		t.Error("NewBBox function Test: failed")
	}
}

func TestPointToBBox(t *testing.T) {
	p := [2]float64{6.969658, 79.873308}

	bbox := PointToBBox(p, 1)
	if bbox.Min != [2]float64{5.969658, 78.873308} && bbox.Max != [2]float64{7.969658, 80.873308} {
		t.Error("PointToBBox function Test: failed")
	}
}

//
type tBox struct {
	min [2]float64
	max [2]float64
}

func randBoxes(N int) []tBox {
	boxes := make([]tBox, N)
	for i := 0; i < N; i++ {
		boxes[i].min[0] = rand.Float64()*360 - 180
		boxes[i].min[1] = rand.Float64()*180 - 90
		boxes[i].max[0] = boxes[i].min[0] + rand.Float64()
		boxes[i].max[1] = boxes[i].min[1] + rand.Float64()
		if boxes[i].max[0] > 180 || boxes[i].max[1] > 90 {
			i--
		}
	}
	return boxes
}


func benchVarious(t *testing.T, tr RTree, numPoints int) {
	N := numPoints
	rand.Seed(time.Now().UnixNano())
	points := make([][2]float64, N)
	for i := 0; i < N; i++ {
		points[i][0] = rand.Float64()*360 - 180
		points[i][1] = rand.Float64()*180 - 90
	}
	pointsReplace := make([][2]float64, N)
	for i := 0; i < N; i++ {
		pointsReplace[i][0] = points[i][0] + rand.Float64()
		if pointsReplace[i][0] > 180 {
			pointsReplace[i][0] = points[i][0] - rand.Float64()
		}
		pointsReplace[i][1] = points[i][1] + rand.Float64()
		if pointsReplace[i][1] > 90 {
			pointsReplace[i][1] = points[i][1] - rand.Float64()
		}
	}
	lotsa.Output = os.Stdout
	fmt.Printf("insert:  ")
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Insert(points[i], points[i], i)
	})
	fmt.Printf("search:  ")
	var count int
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Search(points[i], points[i],
			func(min, max [2]float64, value interface{}) bool {
				count++
				return true
			},
		)
	})
	if count != N {
		t.Fatalf("expected %d, got %d", N, count)
	}
	fmt.Printf("replace: ")
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Replace(
			points[i], points[i], i,
			pointsReplace[i], pointsReplace[i], i,
		)
	})
	if tr.Len() != N {
		t.Fatalf("expected %d, got %d", N, tr.Len())
	}

	fmt.Printf("delete:  ")
	lotsa.Ops(N, 1, func(i, _ int) {
		tr.Remove(pointsReplace[i], pointsReplace[i], i)
	})
	if tr.Len() != 0 {
		t.Fatalf("expected %d, got %d", 0, tr.Len())
	}
}

var Tests = struct {
	TestBenchVarious func(t *testing.T, tr RTree, numPoints int)
}{
	benchVarious,
}

func init() {
	seed := time.Now().UnixNano()
	//println("seed:", seed)
	rand.Seed(seed)
}

func TestRtree(t *testing.T) {
	t.Run("BenchVarious", func(t *testing.T) {
		Tests.TestBenchVarious(t, RTree{}, 1000)
	})
}




