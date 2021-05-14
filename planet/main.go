package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Saied74/skyham/pkg/dataops"
)

const (
	year      = "year"
	month     = "month"
	day       = "day"
	hour      = "hour"
	minute    = "minute"
	second    = "second"
	latitude  = "latitude"
	longitude = "longitude"
	elevation = "elevation"
)

//profile and profiles are for collecting the user input.  In an excel program
//they would be contents of cells
type profile struct {
	firstMsg  string
	secondMsg string
	input     string
	opSys     string
}

type profiles map[string]profile

//When indexing through map keys, the order in indeterminant.  This list is to
//make sure that it is not indeterminant and the map is traversed in this order.
var profileList = []string{year, month, day, hour, minute, latitude,
	longitude, elevation}

func main() {
	//build the structure to hold base (including earth and sun) and planet data
	basedata := make(dataops.BaseItems)
	// planetdata := make(fileops.BaseItems)
	bd := &basedata
	// pd := &planetdata

	//get base (including Earth) and planet (currently Jupiter) data
	bd.ReadData("../data/basedata.csv")
	bd.ReadData("../data/jupiterdata.csv")

	//temparirly locate the CLI here
	printIntro()
	reader := bufio.NewReader(os.Stdin)

	var p = *getProfile() // first get a blank profile.

	for {
		//list data inside the profiles and edit as needed.
		p.listItems(reader)
		p = *p.getInput(reader)

		pack := p.packageInput()
		// gt := p.makeGt()
		err := basedata.ProcInputs(pack)
		// err := basedata.JTime(gt)
		if err != nil {
			fmt.Println("bad input: ", err)
		}

		bd.CalcPeriod()
		bd.CalcOPangles()
		bd.CalcM()
		bd.CalcE()
		bd.PrintBaseItems()

		eOPV := bd.EarthOPVec()
		fmt.Println("Earth OPV", eOPV)
		pOPV := bd.PlanetOPVec()
		fmt.Println("Planet OPV", pOPV)

		bd.SidAngle()
		bd.EarthPrecession()

		// eP3 := skymath.E3(basedata["p"].Value)
		// eTau1 := skymath.E1(basedata["earthTilt"].Value)
		// eMGamma3 := skymath.E3(-basedata["sidAngle"].Value)
		// eMPhi3 := skymath.E3(basedata[])
	}

}
