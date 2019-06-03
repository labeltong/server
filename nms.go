package main

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"image"
	"sort"
)

var (
	validDB postgredb
)

func main() {

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
	var dlist []Datas
	dberr := validDB.DB.Find(&dlist)
	if dberr.Error != nil {
		Custom_panic(dberr.Error)
	}

	for _, v := range dlist {
		if (v.RequiredNumAnswer == 0) && (v.IsFake != true) {
			if v.AnswerType == "1" {

				dberr = validDB.DB.Where("data_id = ?", v.ID).Find(&pick.didx)
				if dberr.Error != nil {
					Custom_panic(dberr.Error)
				}

				d := dataToRect(pick)

				selected, _ := NMS2(d, int(dberr.RowsAffected), int(v.ID))
				/// Delete and increase ban point to users.
				for i, v := range pick.didx {
					if i == selected {
						validDB.DB.Model(v).Update("is_valid", true)
					} else {
						validDB.DB.Delete(v)
						// For banpoints for wrong data
						//validDB.DB.Table("public.users").Where("id = ?", v.UserId).UpdateColumn("ban_point", gorm.Expr("ban_point + ?", 1))

					}
				}

			} else {

				dberr = validDB.DB.Where("data_id = ?", v.ID).Find(&pick.didx)
				if dberr.Error != nil {
					Custom_panic(dberr.Error)
				}
				ans := qsdataToSlice(pick)
				ansarr, ansn := multipleQs(ans, len(ans))
				fmt.Println(ansarr, ansn)

				for i, v := range pick.didx {
					if i == ansn {
						validDB.DB.Model(v).Update("is_valid", true)
					} else {
						validDB.DB.Delete(v)
						// For banpoints for wrong data
						//validDB.DB.Table("public.users").Where("id = ?", v.UserId).UpdateColumn("ban_point", gorm.Expr("ban_point + ?", 1))

					}
				}

			}
		}
	}
}

func multipleQs(candidateQs []int, n int) ([]int, int) {
	//TODO: MAKE multiple questions verifications.
	tmp := make([]int, n)
	for _, v := range candidateQs {
		tmp[v-1]++
	}
	maxtmp := 0
	maxidx := 0
	for i, v := range tmp {
		if v >= maxtmp {
			maxtmp = v
			maxidx = i
		}
	}
	return tmp, maxidx
}

func NMS(candidateQs []image.Rectangle, n int, dataId int) (int, error) {
	//TODO: MAKE multiple questions verifications.
	if n == 0 {
		return 0, fmt.Errorf("no box to calculate")
	}

	area := make([]int, n)
	var pick []int
	overlapThresh := 0.3
	x1 := rectmakecoordslice(candidateQs, 0)
	x2 := rectmakecoordslice(candidateQs, 1)
	y1 := rectmakecoordslice(candidateQs, 2)
	y2 := rectmakecoordslice(candidateQs, 3)

	for i := 0; i < n; i++ {
		area[i] = (x2[i] - x1[i] + 1) * (y2[i] - y1[i] + 1)
	}
	ds := NewIntSlice(y2)

	sort.Sort(ds)
	idxs := ds.idx
	for len(idxs) > 0 {
		last := len(idxs) - 1
		i := idxs[last]
		pick = append(pick, i)
		suppress := []int{last}
		for pos := 0; pos < last; pos++ {
			j := idxs[pos]
			xx1 := max(x1[i], x1[j])
			yy1 := max(y1[i], y1[j])
			xx2 := min(x2[i], x2[j])
			yy2 := min(y2[i], y2[j])

			w := max(0, xx2-xx1+1)
			h := max(0, yy2-yy1+1)
			overlap := float64(w*h) / float64(area[j])
			if overlap > overlapThresh {
				suppress = append(suppress, pos)

			}
		}
		for k, _ := range suppress {
			if k < len(idxs)-1 {
				idxs = append(idxs[:k], idxs[k+1:]...)
			} else {
				idxs = idxs[:k]
			}
		}

	}
	//fmt.Println(candidateQs[pick[len(pick)-1]])

	return pick[len(pick)-1], nil
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
				fmt.Println("asdf",overlap)

				suppress = append(suppress, pos)

			}
		}
		fmt.Println("before ",idxs, suppress)

		idxs = deletebyidx(idxs, suppress)

		fmt.Println("after ",idxs, suppress)

	}

	fmt.Println(pick)
	return pick[0], nil
}
