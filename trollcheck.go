//TODO:
// 1. NMS for multiple box DONE
// 2. Multiple QAs for classification and sentiment , also troll check in classification DONE
// 3. IoU calculator DONE
// 4. DB insert transaction function. DONE

package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"image"
	"strconv"
	"strings"
)

var (
	trollDB postgredb
)

func main() {

	trollDB = postgredb{}
	err := trollDB.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		internalerr := trollDB.DB.Close()
		if internalerr != nil {
			log.Panic(internalerr)

		}

	}()

	var pick DataToVerify
	var dlist []Datas
	dberr := trollDB.DB.Find(&dlist)
	if dberr.Error != nil {
		Custom_panic(dberr.Error)
	}
	for _, v := range dlist {
		if (v.RequiredNumAnswer == 0) && (v.IsFake == true) {
			if v.AnswerType == "1" {
				fmt.Println(v)

				dberr = trollDB.DB.Where("data_id = ?", v.ID).Find(&pick.didx)
				if dberr.Error != nil {
					Custom_panic(dberr.Error)
				}

				d := dataToRect(pick)
				//boxes :=

				selected, _ := IoUcheck(d, int(dberr.RowsAffected), int(v.ID))
				/// Delete and increase ban point to users.
				for i, v := range pick.didx {
					if slice_contains(selected, i) {
						trollDB.DB.Delete(v)
						// For banpoints for wrong data
						trollDB.DB.Table("users").Where("id = ?", v.UserId).UpdateColumn("ban_point", gorm.Expr("ban_point + ?", 1))

					}else{
						trollDB.DB.Delete(v)

					}

				}

			} else {

				dberr = trollDB.DB.Where("data_id = ?", v.ID).Find(&pick.didx)
				if dberr.Error != nil {
					Custom_panic(dberr.Error)
				}
				ans := qsdataToSlice(pick)
				ansarr, ansn := multipleQs_trollcheck(ans, len(ans))
				fmt.Println(ansarr, ansn)

				for i, v := range pick.didx {
					if i != ansn {
						trollDB.DB.Delete(v)
						//For banpoints for wrong data
						trollDB.DB.Table("users").Where("id = ?", v.UserId).UpdateColumn("ban_point", gorm.Expr("ban_point + ?", 1))
					} else{
						trollDB.DB.Delete(v)

					}
				}

			}
		}
	}
}

func multipleQs_trollcheck(candidateQs []int, n int) ([]int, int) {
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

func IoUcheck(boxes []image.Rectangle, n int, dataId int) ([]int, error) {
	//TODO:
	var ans AllAnswers
	var rowtmp [4]int
	var banlist []int
	var iourate float64
	var is image.Rectangle
	dberr := trollDB.DB.Where("data_id = ? AND is_bait = ?", dataId, true).First(&ans)
	if dberr.Error != nil {
		Custom_panic(dberr.Error)
	}

	stringSlice := strings.Split(ans.AnswerData, ",")
	for i, v := range stringSlice {
		rowtmp[i], _ = strconv.Atoi(v)
	}
	ansrect := image.Rect(rowtmp[0], rowtmp[1], rowtmp[2], rowtmp[3])

	for i, v := range boxes {
		is = ansrect.Intersect(v)
		iourate = float64(is.Size().X * is.Size().Y) / float64((v.Size().X * v.Size().Y) + (ansrect.Size().X * ansrect.Size().Y) - (is.Size().X * is.Size().Y))
		fmt.Println(iourate)
		if iourate <= IoUThreshold {
			fmt.Println(v.String(),ansrect.String(), is.Size().String(), iourate)
			fmt.Println((v.Size().X * v.Size().Y) , (ansrect.Size().X * ansrect.Size().Y) , (is.Size().X * is.Size().Y))
			banlist = append(banlist, i)
		}
	}

	return banlist, nil

}
