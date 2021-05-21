package dataops

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"
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

//GetItem checks to make sure the item exists in basedata and returns the item.
//If item does not exist, it writes an error log and exits the program.
func (bd *BaseItems) GetItem(item string) float64 {
	dd := *bd
	d, ok := dd[item]
	if !ok {
		fmt.Printf("item %s did not exist - terminating program\n", item)
		os.Exit(2)
	}
	return d.Value
}

//todo this check is poor, re-write
func check(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
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

//PrintBaseItems prints the contents of the baseitems structured tabbed.
func (bd *BaseItems) PrintBaseItems() {
	dd := *bd
	//we will need this later for generating tabbed output
	w := new(tabwriter.Writer)
	sortKey := []string{}
	for key := range dd {
		sortKey = append(sortKey, key)
	}
	sk := sort.StringSlice(sortKey)
	sk.Sort()
	fmt.Printf("\n")
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	fmt.Fprintln(w, "Key\tName\tValue\tNumonic\tDescription")
	for _, key := range sk {
		s := dd[key]
		buf := []byte(fmt.Sprintf("%s\t%s\t%e\t%s\t%s\n", key, s.Name, s.Value, s.Numonic, s.Description))
		// fmt.Fprintln(w, dp.Name, dp.Value, dp.Numonic, dp.Description)
		w.Write(buf)
	}
	fmt.Printf("\n")
	w.Flush()
}
