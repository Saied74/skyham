package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Saied74/skyham/pkg/dataops"
	"github.com/Saied74/skyham/pkg/skymath"
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
var profileList = []string{year, month, day, hour, minute, second, latitude,
	longitude, elevation}

func main() {
	//build the structure to hold base (including earth and sun) and planet data
	basedata := make(dataops.BaseItems)
	// planetdata := make(fileops.BaseItems)
	bd := &basedata
	// pd := &planetdata

	//get base (including Earth) and planet (currently Jupiter) data
	// bd.ReadData("../data/basedata.csv")
	bd.ReadData("../data/jupiterdata.csv")

	//temparirly locate the CLI here
	printIntro()
	reader := bufio.NewReader(os.Stdin)

	var p = *getProfile() // first get a blank profile.

	for {
		//list data inside the profiles and edit as needed.
		p.listItems(reader)
		p = *p.getInput(reader)
		// p.listItems(reader)

		pack := p.packageInput()
		// gt := p.makeGt()
		err := bd.ProcInputs(pack)
		// err := basedata.JTime(gt)
		if err != nil {
			fmt.Println("bad input: ", err)
			os.Exit(2)
		}

		bd.CalcPeriod()
		bd.CalcOPangles()
		bd.CalcM()

		bd.CalcE()
		erSunOP := bd.EarthOPVec()
		fmt.Println("Earth OPV", erSunOP)
		prSunOP := bd.PlanetOPVec()
		fmt.Println("Planet OPV", prSunOP)

		bd.SidAngle()
		bd.EarthPrecession()

		bd.PrintBaseItems()

		earthP, ok := basedata["earthP"]
		if !ok {
			check("bad lookup: earthP", nil)
		}
		eP3 := skymath.E3(earthP.Value)
		eTau1 := skymath.E1(dataops.EarthTilt)

		sidAngle, ok := basedata["sidAngle"]
		if !ok {
			check("bad lookup: sidAngle", nil)
		}
		eMGamma3 := skymath.E3(-sidAngle.Value)

		locallong, ok := basedata["locallong"]
		if !ok {
			check("bad lookup: locallong", nil)
		}
		eMPhi3 := skymath.E3(-locallong.Value)

		locallat, ok := basedata["locallat"]
		if !ok {
			check("bad lookup: locallat", nil)
		}
		eLam2 := skymath.E2(locallat.Value)
		eNU := skymath.Euler{
			[3]float64{0.0, 1.0, 0.0},
			[3]float64{0.0, 0.0, 1.0},
			[3]float64{1.0, 0.0, 0.0},
		}

		step1 := skymath.Mply(eTau1, eP3)
		step2 := skymath.Mply(eMGamma3, step1)
		step3 := skymath.Mply(eMPhi3, step2)
		step4 := skymath.Mply(eLam2, step3)
		sciTOenu := skymath.Mply(eNU, step4)

		planetArgPre, ok := basedata["planetArgPre"]
		if !ok {
			check("bad lookup: planetArgPre", nil)
		}
		p3LittleOmega := skymath.E3(planetArgPre.Value)

		planetInc, ok := basedata["planetInc"]
		if !ok {
			check("bad lookup: planetInc", nil)
		}
		pI1 := skymath.E1(planetInc.Value)

		planetNode, ok := basedata["planetNode"]
		if !ok {
			check("bad lookup: planetNode", nil)
		}
		p3BigOmega := skymath.E3(planetNode.Value)

		step5 := skymath.Mply(pI1, p3LittleOmega)
		oppTOsci := skymath.Mply(p3BigOmega, step5)

		prSunSCI := skymath.Vply(oppTOsci, prSunOP)
		prSunENU := skymath.Vply(sciTOenu, prSunSCI)
		// fmt.Printf("Planet to sun vector in ENU: %e   %e   %e\n", prSunENU[0], prSunENU[1], prSunENU[2])

		//======================
		argPre, ok := basedata["argPre"]
		if !ok {
			check("bad lookup: argPre", nil)
		}
		e3LittleOmega := skymath.E3(argPre.Value)
		eI1 := skymath.E1(dataops.EarthInc)
		e3BigOmega := skymath.E3(dataops.EarthNode)

		step6 := skymath.Mply(eI1, e3LittleOmega)
		opeTOsci := skymath.Mply(e3BigOmega, step6)

		erSunSCI := skymath.Vply(opeTOsci, erSunOP)

		//==================   Important Vector   =================================
		erSunENU := skymath.Vply(sciTOenu, erSunSCI)

		locTOearthRecef := bd.CalcLocalVec()
		// fmt.Printf("Local vetor in ECEF frame: %e   %e   %e\n", locTOearthRecef[0], locTOearthRecef[1], locTOearthRecef[2])

		step7 := skymath.Mply(eLam2, eMPhi3)
		ecefTOenu := skymath.Mply(eNU, step7)

		erLocENU := skymath.Vply(ecefTOenu, locTOearthRecef)

		second := skymath.Vadd(erSunENU, erLocENU)
		prENU := skymath.Vsub(prSunENU, second)
		fmt.Printf("\n")
		fmt.Printf("Local to Planet vector in ENU: %e   %e   %e\n", prENU[0], prENU[1], prENU[2])
		fmt.Printf("\n")

		beta, epsilon := skymath.CalcBetaEpsilon(prENU)
		fmt.Printf("Beta: %f, Epsilon: %f\n", beta, epsilon)
	}
}
