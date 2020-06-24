package main

import "fmt"

func main(){
	var tr RTree

	tr.Insert([2]float64{6.969658, 79.873308}, [2]float64{6.969658, 79.873308}, "Cargo Logistics",nil) //0.22km
	tr.Insert([2]float64{6.971322, 79.874468}, [2]float64{6.971322, 79.874468}, "Cargills",nil)
	tr.Insert([2]float64{6.970960, 79.873652}, [2]float64{6.970960, 79.873652}, "Vystwyke",nil) //0.1km
	tr.Insert([2]float64{6.973969, 79.876183}, [2]float64{6.973969, 79.876183}, "Xarena",nil) //cannot 0.35km
	tr.Insert([2]float64{6.973597, 79.869424}, [2]float64{6.973597, 79.869424}, "Sea",nil) //Cannot //0.61 km
	tr.Insert([2]float64{6.970263, 79.874112}, [2]float64{6.970263, 79.874112}, "Buhari",nil) //0.12km
	tr.Insert([2]float64{6.967660, 79.872217}, [2]float64{6.967660, 79.872217}, "Laughs",nil) //0.48km

	//var point [2]float64
	//var point1 [2]float64
	//point[0] = 6.973321
	//point[1] = 79.875518
	//point1[0] = 6.966735
	//point1[1] = 79.871175
	//var locationNames []interface{}
	//tr.Search([2]float64{6.96, 79.86}, [2]float64{6.98,79.89}, func(Min, Max [2]float64,Reference string) bool{
	//	fmt.Println(Reference)
	//	return true
	//},)
	bbox1 := tr.All()

	for _,value := range bbox1{
		fmt.Println(value.Min)
	}
}
