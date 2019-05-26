package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

type Slice struct {
	sort.Interface
	idx []int
}

func (s Slice) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}
func NewSlice(n sort.Interface) *Slice {
	s := &Slice{Interface: n, idx: make([]int, n.Len())}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}
func NewIntSlice(n []int) *Slice         { return NewSlice(sort.IntSlice(n)) }

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}


func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}


var (
	validDB postgredb

)


type DataToVerify struct{
	didx []Answers
	dn int
}

func main(){
	

	validDB = postgredb{}
	err := validDB.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		internalerr := validDB.DB.Close()
		if internalerr != nil {
			log.Panic(internalerr)

		}

	}()

	var pick DataToVerify
	dberr := validDB.DB.Where("data_id = ?",2).Find(&pick.didx)
	Custom_panic(dberr.Error)

	fmt.Println(pick.didx)
	fmt.Println(dataToSlice(pick))
}


func makecoordslice(sl [][]int, loc int) []int {
	tmp := make([]int, len(sl))

	for i,v := range sl{
		tmp[i] = v[loc]

	}

	return tmp
}

func dataToSlice(d DataToVerify) [][4]int{
	var tmp [][4]int
	var rowtmp [4]int

	for _,v := range d.didx{
		stringSlice := strings.Split(v.AnswerData, ",")
		for i,v := range stringSlice{
			rowtmp[i],_ = strconv.Atoi(v)
		}
		tmp = append(tmp, rowtmp)
	}
	return tmp

}







func multipleQs (candidateQs []int, n int) error {
	//TODO: MAKE multiple questions verifications.
	return nil
}


func NMS(candidateQs [][]int, n int, dataId int ) error{
	//TODO: MAKE multiple questions verifications.
	if n==0{
		return fmt.Errorf("no box to calculate")
	}



	area := make([]int, 3)
	var pick []int
	overlapThresh := 0.3
	x1 := makecoordslice(candidateQs,0)
	x2 := makecoordslice(candidateQs,1)
	y1 := makecoordslice(candidateQs,2)
	y2 := makecoordslice(candidateQs,3)

	for i :=0;i<3;i++{
		area[i] = (x2[i] - x1[i] + 1) * (y2[i] - y1[i] + 1)
	}
	ds := NewIntSlice(y2)

	sort.Sort(ds)
	idxs := ds.idx
	for len(idxs) > 0{
		last := len(idxs) -1
		i := idxs[last]
		pick = append(pick,i)
		suppress := []int{last}
		for pos:=0;pos<last;pos++{
			j := idxs[pos]
			xx1 := max(x1[i], x1[j])
			yy1 := max(y1[i], y1[j])
			xx2 := min(x2[i], x2[j])
			yy2 := min(y2[i], y2[j])

			w := max(0, xx2 - xx1 + 1)
			h := max(0, yy2 - yy1 + 1)

			overlap := float64(w * h) / float64(area[j])
			if overlap > overlapThresh{
				suppress = append(suppress,pos)

			}
		}
		for k,_:= range suppress{
			if k < len(idxs) -1 {
				idxs =  append(idxs[:k], idxs[k+1:]...)
			}else {
				idxs = idxs[:k]
			}
		}

	}
	fmt.Println(candidateQs[pick[len(pick)-1]])





	return nil
}