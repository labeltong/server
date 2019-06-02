package main

import (
	"fmt"
	"github.com/mpraski/clusters"
)

func main(){
	//var data [][]float64
	var observation []float64

	nmsa :=[][]float64{
		{12, 84, 140, 212},
		{24, 84, 152, 212},
		{36, 84, 164, 212},
		{12, 96, 140, 224},
		{24, 96, 152, 224},
		{24, 108, 152, 236},
		{32, 84, 120, 202},
		{24, 74, 152, 222},
		{16, 84, 134, 212},
		{12, 96, 140, 214},
		{24, 76, 152, 224},
		{34, 118, 142, 246},
	}

	//nmsb := [][]float64{
	//	{114, 60, 178, 124},
	//	{120, 60, 184, 124},
	//	{114, 66, 178, 130},
	//}
	//nmsc := [][]float64{
	//	{12, 30, 76, 94},
	//	{12, 36, 76, 100},
	//	{72, 36, 200, 164},
	//	{84, 48, 212, 176},
	//}
	//
	//nmsd := [][]float64{
	//	{220,0,550,300},
	//	{0,0,400,400},
	//	{280,50,450,240},
	//	{0,0,500,500},
	//	{200,40,400,300},
	//
	//}




	// Create a new KMeans++ clusterer with 1000 iterations,
	// 8 clusters and a distance measurement function of type func([]float64, []float64) float64).
	// Pass nil to use clusters.EuclideanDistance
	c, e := clusters.KMeans(10, 3, ioudistance)
	if e != nil {
		panic(e)
	}

	// Use the data to train the clusterer
	if e = c.Learn(nmsa); e != nil {
		panic(e)
	}

	fmt.Printf("Clustered data set into %d\n", c.Sizes())

	fmt.Printf("Assigned observation %v to cluster %d\n", observation, c.Predict(observation))

	for index, number := range c.Guesses() {
		fmt.Printf("Assigned data point %v to cluster %d\n", nmsa[index], number)
	}
}

func ioudistance(a []float64, b []float64) float64{
	b1 := []float64{(a[1]-a[0]),(a[3]-a[2])}
	b2 := []float64{(b[1]-b[0]),(b[3]-b[2])}
	iou := (minfloat(b1[0],b2[0])*minfloat(b1[1],b2[1])) / (b1[0]*b2[0]+b1[1]*b2[1]-minfloat(b1[0],b2[0])*minfloat(b1[1],b2[1]))
	fmt.Println(b1,b2,iou)

	return iou

}

func minfloat(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}