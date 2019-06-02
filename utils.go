package main

import (
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"image"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func MakeReverseProxy(rplist []string, rpnum *int) (n int, t []*middleware.ProxyTarget) {
	var urltmp *url.URL
	tm := []*middleware.ProxyTarget{}
	var err error
	n = 0
	for _, s := range rplist {
		urltmp, err = url.Parse(s)
		if err != nil {
			log.Panic(err)
		}
		tm = append(tm, &middleware.ProxyTarget{
			Name: urltmp.String(),
			URL:  urltmp,
		})
		n++
	}
	return n, tm
}

func Custom_fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Custom_panic(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type Slice2 struct {
	sort.Interface
	idx []int
}

func (s Slice2) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}
func NewSlice(n sort.Interface) *Slice2 {
	s := &Slice2{Interface: n, idx: make([]int, n.Len())}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}
func NewIntSlice(n []int) *Slice2 { return NewSlice(sort.IntSlice(n)) }

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

type DataToVerify struct {
	didx []Answers
	dn   int
}

func makecoordslice(sl [][4]int, loc int) []int {
	tmp := make([]int, len(sl))

	for i, v := range sl {
		tmp[i] = v[loc]

	}

	return tmp
}
func rectmakecoordslice(sl []image.Rectangle, loc int) []int {
	tmp := make([]int, len(sl))

	for i, _:= range sl {
		if loc == 0 {
			tmp[i] = sl[i].Min.X
		} else if loc == 1 {
			tmp[i] = sl[i].Min.Y
		} else if loc == 2 {
			tmp[i] = sl[i].Max.X
		} else {
			tmp[i] = sl[i].Max.Y
		}

	}

	return tmp
}

func dataToSlice(d DataToVerify) [][4]int {
	var tmp [][4]int
	var rowtmp [4]int

	for _, v := range d.didx {
		stringSlice := strings.Split(v.AnswerData, ",")
		for i, v := range stringSlice {
			rowtmp[i], _ = strconv.Atoi(v)
		}
		tmp = append(tmp, rowtmp)
	}
	return tmp

}

func dataToRect(d DataToVerify) []image.Rectangle {
	var tmp []image.Rectangle
	var rowtmp [4]int

	for _, v := range d.didx {
		stringSlice := strings.Split(v.AnswerData, ",")
		for i, v := range stringSlice {
			rowtmp[i], _ = strconv.Atoi(v)
		}
		tmp = append(tmp, image.Rect(rowtmp[0], rowtmp[1], rowtmp[2], rowtmp[3]))
	}
	return tmp

}

func qsdataToSlice(d DataToVerify) []int {
	tmp := make([]int, len(d.didx))

	for i, v := range d.didx {
		tmp[i], _ = strconv.Atoi(v.AnswerData)
	}
	return tmp

}

func slice_contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func difference(slice1 []int, slice2 []int) []int {
	var diff []int

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

func deletebyidx(slice1 []int, slice2 []int) []int {
	var lenslice int
	ret := make([]int,len(slice1))
	copy(ret, slice1)
	for _,i := range slice2{
		lenslice = len(ret)
		if i < lenslice-1{
			copy(ret[i:], ret[i+1:])
			ret[len(ret)-1] = 0 // or the zero value of T
			ret = ret[:len(ret)-1]
		}else{
			ret = ret[:len(ret)-1]
		}


	}
	return ret
}