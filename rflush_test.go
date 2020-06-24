package rflush

import (
	"testing"
)

//func TestRTree_Insert(t *testing.T) {
//	var tr RTree
//
//	tr.Insert([2]float64{6.969658, 79.873308}, [2]float64{6.969658, 79.873308}, "Cargo Logistics",nil) //0.22km
//	tr.Insert([2]float64{6.971322, 79.874468}, [2]float64{6.971322, 79.874468}, "Cargills",nil)
//	tr.Insert([2]float64{6.970960, 79.873652}, [2]float64{6.970960, 79.873652}, "Vystwyke",nil) //0.1km
//	tr.Insert([2]float64{6.973969, 79.876183}, [2]float64{6.973969, 79.876183}, "Xarena",nil) //cannot 0.35km
//	tr.Insert([2]float64{6.973597, 79.869424}, [2]float64{6.973597, 79.869424}, "Sea",nil) //Cannot //0.61 km
//	tr.Insert([2]float64{6.970263, 79.874112}, [2]float64{6.970263, 79.874112}, "Buhari",nil) //0.12km
//	tr.Insert([2]float64{6.967660, 79.872217}, [2]float64{6.967660, 79.872217}, "Laughs",nil) //0.48km
//
//}

func TestRTree_Search(t *testing.T) {
		var tr RTree

		tr.Insert([2]float64{6.969658, 79.873308}, [2]float64{6.969658, 79.873308}, "Cargo Logistics",nil) //0.22km
		tr.Insert([2]float64{6.971322, 79.874468}, [2]float64{6.971322, 79.874468}, "Cargills",nil)
		tr.Insert([2]float64{6.970960, 79.873652}, [2]float64{6.970960, 79.873652}, "Vystwyke",nil) //0.1km
		tr.Insert([2]float64{6.973969, 79.876183}, [2]float64{6.973969, 79.876183}, "Xarena",nil) //cannot 0.35km
		tr.Insert([2]float64{6.973597, 79.869424}, [2]float64{6.973597, 79.869424}, "Sea",nil) //Cannot //0.61 km
		tr.Insert([2]float64{6.970263, 79.874112}, [2]float64{6.970263, 79.874112}, "Buhari",nil) //0.12km
		tr.Insert([2]float64{6.967660, 79.872217}, [2]float64{6.967660, 79.872217}, "Laughs",nil) //0.48km

		var point [2]float64
		point[0] = 6.973597
		point[1] = 79.869424
		var locationNames []interface{}
		tr.Search(point,point, func(min, max [2]float64,reference string) bool{
			locationNames = append(locationNames,reference)
			return false
		},)
		if len(locationNames) != 1{
			t.Error("Search Function Test: Failed")
		}
}
