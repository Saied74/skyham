package dataops

import (
	"fmt"
	"os"
	"skyham/pkg/skylog"
	"strconv"
)

//gConvert converts a slice of strings to a slice of float64 with some
//only used in JD Calc.
func gConvert(g []string) (gt []float64, err error) {
	gt = make([]float64, 6)
	if len(g) != 6 {
		return []float64{}, fmt.Errorf("the length of georgian time string not 6")
	}
	for i := 0; i < 5; i++ {
		gt[i], err = strconv.ParseFloat(g[i], 64)
		if err != nil {
			return []float64{}, fmt.Errorf("element %v did not convert %v", g[i], err)
		}
	}
	return gt, nil
}

//gCheck checks to a degree the validity of dates; only used in JD calc.
// TODO: improve the bound checking based on month
// TODO: move some of this to data entry stage
func gCheck(gt []float64) error {
	lowerLimit := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	upperLimit := []float64{100000.0, 12, 31.0, 24.0, 60.0, 60.0}
	for i := 0; i < 5; i++ {
		if gt[i] < lowerLimit[i] || gt[i] > upperLimit[i] {
			return fmt.Errorf("invalid field %v", gt[i])
		}
	}
	return nil
}

// TODO: write an auditData function

//todo this check is poor, re-write
func check(msg string, err error) {
	if err != nil {
		skylog.ErrorLog.Println(msg, err)
		os.Exit(2)
	}
}

func assign(items []string, key string) BaseItem {
	var err error
	var oneItem BaseItem
	oneItem.Name = key
	if len(items) < 3 {
		fmt.Println("bad length", items)
		os.Exit(2)
	}
	oneItem.Value, err = strconv.ParseFloat(items[1], 64)
	check("Item"+" key "+"did not parse float: ", err)
	oneItem.Numonic = items[2]
	oneItem.Description = items[3]
	return oneItem
}
