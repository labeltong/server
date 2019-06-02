package main

import (
	"fmt"
	"image"
	"sort"
)

func main(){
	//a := image.Rect(0,0,200,200)
	//b:= image.Rect(100,100,300,300)
	//c :=image.Rect(300,300,500,500)
	//answer := image.Rect(10,10,200,200)

	nmsa :=[][]int{
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

	nmsb := [][]int{
		{114, 60, 178, 124},
		{120, 60, 184, 124},
		{114, 66, 178, 130},
	}
	nmsc := [][]int{
		{12, 30, 76, 94},
		{12, 36, 76, 100},
		{72, 36, 200, 164},
		{84, 48, 212, 176},
	}
	nmsd := [][]int{
		{220,0,550,300},
		{0,0,400,400},
		{280,50,450,240},
		{0,0,500,500},
		{200,40,400,300},

	}


	reta,_ := NMS2(makerects(nmsa),12,1)
	retb,_ := NMS2(makerects(nmsb),3,2)
	retc,_ := NMS2(makerects(nmsc),4,3)
	retd,_ := NMS2(makerects(nmsd),5,4)

	fmt.Println(reta, retb, retc, retd)



}


func makerects(d [][]int) []image.Rectangle{
	var ret []image.Rectangle
	for _,v :=range d{
		ret = append(ret,image.Rect(v[0],v[1],v[2],v[3]))
	}
	return ret
}


func NMS2(candidateQs []image.Rectangle, n int, dataId int) (int, error) {
	//TODO: MAKE multiple questions verifications.
	if n == 0 {
		return 0, fmt.Errorf("no box to calculate")
	}

	area := make([]int, n)
	score := make([]int, n)
	var pick []int

	x1 := rectmakecoordslice(candidateQs, 0)
	y1 := rectmakecoordslice(candidateQs, 1)
	x2 := rectmakecoordslice(candidateQs, 2)
	y2 := rectmakecoordslice(candidateQs, 3)

	for i := 0; i < n; i++ {
		area[i] = (x2[i] - x1[i] + 1) * (y2[i] - y1[i] + 1)
	}
	////naive approach
	//ds := NewIntSlice(y2)
	//sort.Sort(ds)
	//idxs := ds.idx
	//For score by distance of mean size
	sum := 0
	for _, num := range area {
		sum += num
	}
	sum = sum/len(area)
	for i, num := range area{
		t := num-sum
		if t<0{
			score[i] = -t
		}else{
			score[i] = t
		}

	}
	fmt.Println(sum,area)
	scorelists := NewIntSlice(score)
	fmt.Println(scorelists)
	sort.Sort(scorelists)
	idxs := scorelists.idx
	fmt.Println(scorelists)

	for len(idxs) > 0 {
		last := len(idxs) - 1
		i := idxs[last]
		pick = append(pick, i)
		suppress := []int{last}
		fmt.Println(dataId,pick, suppress)
		for pos := 0; pos < last; pos++ {
			j := idxs[pos]
			xx1 := max(x1[i], x1[j])
			yy1 := max(y1[i], y1[j])
			xx2 := min(x2[i], x2[j])
			yy2 := min(y2[i], y2[j])

			w := max(0, xx2-xx1+1)
			h := max(0, yy2-yy1+1)
			overlap := float64(w*h) / float64(area[j])
			if overlap > OverlapThresh  {
				suppress = append(suppress, pos)

			}
		}
		fmt.Println("before ",idxs, suppress)

		idxs = deletebyidx(idxs, suppress)

		fmt.Println("after ",idxs, suppress)

	}

	fmt.Println(pick)
	return pick[len(pick)-1], nil
}
