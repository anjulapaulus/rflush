package rflush

import (
	"rflush/test_data"
	"strconv"
	"testing"
)

func TestRTree_Insert(t *testing.T) {
	var tr RTree

	test_data.LoadCSV()
	for i:=0; i<len(test_data.LocationsData.Locations); i++{
		var point [2]float64
		latitude,_ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Latitude,64)
		longitude,_ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Longitude,64)
		point[0] = latitude
		point[1] = longitude
		tr.Insert(point, point, test_data.LocationsData.Locations[i].PlaceName, nil)
	}

	if tr.Len() != len(test_data.LocationsData.Locations) {
		t.Error("Insert Function: Error",len(test_data.LocationsData.Locations))
	}
}

func TestRTree_Len(t *testing.T) {
	var tr RTree

	test_data.LoadCSV()
	for i:=0; i<len(test_data.LocationsData.Locations); i++{
		var point [2]float64
		latitude,_ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Latitude,64)
		longitude,_ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Longitude,64)
		point[0] = latitude
		point[1] = longitude
		tr.Insert(point, point, test_data.LocationsData.Locations[i].PlaceName, nil)
	}

	if tr.Len() != len(test_data.LocationsData.Locations) {
		t.Error("Insert Function Test Failed")
	}
}

func TestRTree_Bounds(t *testing.T) {
	var tr RTree

	test_data.LoadCSV()
	for i:=0; i<len(test_data.LocationsData.Locations); i++{
		var point [2]float64
		latitude,_ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Latitude,64)
		longitude,_ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Longitude,64)
		point[0] = latitude
		point[1] = longitude
		tr.Insert(point, point, test_data.LocationsData.Locations[i].PlaceName, nil)
	}

	min, max := tr.Bounds()
	if min != [2]float64{5.9, 75} && max != [2]float64{10.08333,81.91667} {
		t.Error("Bounds Function Test: Failed")
	}
}


func TestRTree_All(t *testing.T) {
	var tr RTree

	test_data.LoadCSV()
	for i:=0; i<len(test_data.LocationsData.Locations); i++{
		var point [2]float64
		latitude,_ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Latitude,64)
		longitude,_ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Longitude,64)
		point[0] = latitude
		point[1] = longitude

		tr.Insert(point, point, test_data.LocationsData.Locations[i].PlaceName, nil)
	}

	var bbox []BBox
	tr.All(func(min, max [2]float64, reference string, data interface{}) bool  {
		var item BBox
		item.Min = min
		item.Max = max
		item.Reference = reference
		bbox = append(bbox, item)
		return true
	})

	if len(bbox) != len(test_data.LocationsData.Locations) {
		t.Error("All Function Test: Failed",len(bbox))
	}
}

func TestRTree_Search(t *testing.T) {
	var tr RTree

	test_data.LoadCSV()
	for i:=0; i<len(test_data.LocationsData.Locations); i++{
		var point [2]float64
		latitude,_ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Latitude,64)
		longitude,_ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Longitude,64)
		point[0] = latitude
		point[1] = longitude
		tr.Insert(point, point, test_data.LocationsData.Locations[i].PlaceName, nil)
	}

	bbox := PointToBBox([2]float64{6.969658, 79.873308}, 0.2)
	var locationNames []interface{}
	tr.Search(bbox.Min, bbox.Max,
		func(min [2]float64, max [2]float64, reference string) bool {
			locationNames = append(locationNames, reference)
			return true
		},
	)
	if len(locationNames) != 2234{
		t.Error("Search Function: Error",len(locationNames))
	}
}

func TestRTree_Remove(t *testing.T) {
	var tr RTree

	test_data.LoadCSV()
	for i:=0; i<len(test_data.LocationsData.Locations); i++{
		var point [2]float64
		latitude,_ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Latitude,64)
		longitude,_ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Longitude,64)
		point[0] = latitude
		point[1] = longitude
		tr.Insert(point, point, test_data.LocationsData.Locations[i].PlaceName, nil)
	}

	tr.Remove([2]float64{7.68333,80.1}, [2]float64{7.68333,80.1}, "Kitagama", nil)
	tr.Remove([2]float64{7.88333,81.53333}, [2]float64{7.88333,81.53333}, "Chunkankeni", nil)
	tr.Remove([2]float64{6.06667,80.46667}, [2]float64{6.06667,80.46667}, "Zowdegala", nil)
	tr.Remove([2]float64{7.5005,80.3088}, [2]float64{7.5005,80.3088}, "Yatawehera", nil)
	tr.Remove([2]float64{7.4603,80.0885}, [2]float64{7.4603,80.0885}, "Yakarawatta", nil)

	if tr.Len() != 56743 {
		t.Error("Remove Function Test: Failed",len(test_data.LocationsData.Locations))
	}
}

func TestRTree_Replace(t *testing.T) {
	var tr RTree

	test_data.LoadCSV()
	for i:=0; i<len(test_data.LocationsData.Locations); i++{
		var point [2]float64
		latitude,_ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Latitude,64)
		longitude,_ := strconv.ParseFloat(test_data.LocationsData.Locations[i].Longitude,64)
		point[0] = latitude
		point[1] = longitude
		tr.Insert(point, point, test_data.LocationsData.Locations[i].PlaceName, nil)
	}

	var newPoint [2]float64
	newPoint[0] = 6.973597
	newPoint[1] = 79.869424
	tr.Replace([2]float64{7.32616,80.62377}, [2]float64{7.32616,80.62377}, "Katugasthota", nil, newPoint, newPoint, "Katugasthota", nil)

	var bbox []BBox
	tr.All(func(min, max [2]float64, reference string, data interface{}) bool  {
		var item BBox
		item.Min = min
		item.Max = max
		item.Reference = reference
		bbox = append(bbox, item)
		return true
	})
	var checkItem BBox
	for i:=0; i< len(bbox); i++{
		if bbox[i].Reference == "Katugasthota"{
			checkItem.Reference = bbox[i].Reference
			checkItem.Min = bbox[i].Min
			checkItem.Max = bbox[i].Max
		}
	}

	if checkItem.Reference == "Katugasthota" && checkItem.Min != newPoint && checkItem.Max != newPoint{
		t.Error("Replace Function Failed",checkItem.Reference,checkItem.Min,checkItem.Max)
	}
}


func TestNewBBox(t *testing.T) {
	p := [2]float64{6.969658, 79.873308}
	lengths := [2]float64{2.5, 8.0}

	bbox := NewBBox(p, lengths)
	if bbox.Min == [2]float64{6.969658, 79.873308} && bbox.Max != [2]float64{9.469657999999999,87.873308} {
		t.Error("NewBBox function Test: failed")
	}
}

func TestPointToBBox(t *testing.T) {
	p := [2]float64{6.969658, 79.873308}

	bbox := PointToBBox(p,1)
	if bbox.Min != [2]float64{5.969658,78.873308}  && bbox.Max != [2]float64{7.969658,80.873308}{
		t.Error("PointToBBox function Test: failed")
	}
}
//
//type tBox struct {
//	min [2]float64
//	max [2]float64
//}
//
//func randBoxes(N int) []tBox {
//	boxes := make([]tBox, N)
//	for i := 0; i < N; i++ {
//		boxes[i].min[0] = rand.Float64()*360 - 180
//		boxes[i].min[1] = rand.Float64()*180 - 90
//		boxes[i].max[0] = boxes[i].min[0] + rand.Float64()
//		boxes[i].max[1] = boxes[i].min[1] + rand.Float64()
//		if boxes[i].max[0] > 180 || boxes[i].max[1] > 90 {
//			i--
//		}
//	}
//	return boxes
//}
//
//func generate(N int, size int) [][]tBox {
//	var data [][]tBox
//	for i := 0; i < N; i++ {
//		data = append(data, randBoxes(size))
//	}
//	return data
//}
//
//func BenchmarkRTree_Insert(b *testing.B) {
//	var tr RTree
//	boxes := randBoxes(b.N)
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		tr.Insert(boxes[i].min, boxes[i].max, string(i),nil)
//	}
//}
//
//
//func BenchmarkRTree_Search(b *testing.B) {
//	var tr RTree
//	boxes := randBoxes(b.N)
//	for i := 0; i < b.N; i++ {
//		tr.Insert(boxes[i].min, boxes[i].max, string(i),nil)
//	}
//	b.ResetTimer()
//	for i:= 0; i< b.N; i++{
//			tr.Search(boxes[i].min, boxes[i].max, func(min [2]float64, max [2]float64, reference string) bool {
//				return true
//			}, )
//
//	}
//}
