package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Saied74/skyham/pkg/dataops"
	"github.com/Saied74/skyham/pkg/skymath"
)

//These constants are for the input and output.  The constants for the planet
//data are in the dataops package.
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
	planet    = "planet"
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
	longitude, elevation, planet}

func main() {
	//build the structure to hold base (including earth and sun) and planet data
	basedata := make(dataops.BaseItems)
	bd := &basedata

	//temparirly locate the CLI here
	printIntro()
	reader := bufio.NewReader(os.Stdin)

	var p = *getProfile() // first get a blank profile.

	for {
		//In the following loop, inputs are gathered and the file is read and
		//processed.  If everything is good, it will break and process the data
		for {
			//list data inside the profiles and edit as needed.
			p = *p.getInput(reader)
			pack := p.packageInput()
			err := bd.ProcInputs(pack)
			if err != nil {
				fmt.Println("bad input, try again: ", err)
			}

			planetName, ok := basedata[dataops.PlanetName]
			if !ok {
				check("bad lookup: planetName", nil)
			}
			fileName := "../data/" + planetName.Numonic + "data.csv"
			//get planet (currently Jupiter) data
			err = bd.ReadData(fileName)
			if err != nil {
				fmt.Printf("error in reading the planet file %v\n", err)
			}
			if err == nil && ok {
				break
			}
		}

		//first, we calc the earth period, the angle, and mean and Eccentrioc anomaly
		bd.CalcPeriod()
		bd.CalcOPangles()
		bd.CalcM()
		bd.CalcE()

		//<================= Calculate the Orbital Plane Vectors ==================>
		erSunOP := bd.EarthOPVec()
		prSunOP := bd.PlanetOPVec()

		//<================= Location to "Earth center" vector ====================>
		locTOearthRecef := bd.CalcLocalVec()

		//Temperary print of the primary vectors
		// fmt.Println("Earth OPV", erSunOP)
		// fmt.Println("Planet OPV", prSunOP)
		// fmt.Println("Location to ECEF", locTOearthRecef)

		//Then calculate the Sideral angle at the requested time as well as the earth
		//spin axis precession
		bd.SidAngle()
		bd.EarthPrecession()

		//At this point, all input and intermediate data is read and calculated
		// bd.PrintBaseItems()

		//<==== Calculate Sun Centred Interal to East North Up Transformation ====>

		//Check for any errors in building the basedata data structure
		earthP := bd.GetItem("earthP")
		sidAngle := bd.GetItem("sidAngle")
		locallong := bd.GetItem("locallong")
		locallat := bd.GetItem("locallat")

		//Then build the Euler transformation matrices
		eP3 := skymath.E3(earthP)
		eTau1 := skymath.E1(dataops.EarthTilt)
		eMGamma3 := skymath.E3(-sidAngle)
		eMPhi3 := skymath.E3(-locallong)
		eLam2 := skymath.E2(locallat)
		eNU := skymath.Euler{
			[3]float64{0.0, 1.0, 0.0},
			[3]float64{0.0, 0.0, 1.0},
			[3]float64{1.0, 0.0, 0.0},
		}

		//And multiply them to get the SCI to ENU transformation matrix
		step1 := skymath.Mply(eTau1, eP3)
		step2 := skymath.Mply(eMGamma3, step1)
		step3 := skymath.Mply(eMPhi3, step2)
		step4 := skymath.Mply(eLam2, step3)
		sciTOenu := skymath.Mply(eNU, step4)

		//<============== Calculate the Planet Centered Intertial =================>
		//<=============== to Sun Centric Inertial transformation ================>

		//Check for any errors in building the basedata data structure
		planetArgPre := bd.GetItem("planetArgPre")
		planetInc := bd.GetItem("planetInc")
		planetNode := bd.GetItem("planetNode")

		//Then build the Euler transformation matrices
		p3LittleOmega := skymath.E3(planetArgPre)
		pI1 := skymath.E1(planetInc)
		p3BigOmega := skymath.E3(planetNode)

		//And multiply them to get the OPI to SCI transformation matrix
		step5 := skymath.Mply(pI1, p3LittleOmega)
		oppTOsci := skymath.Mply(p3BigOmega, step5)
		oppTOenu := skymath.Mply(sciTOenu, oppTOsci)

		//<============== Calculate the Earth Centered Intertial ==================>
		//<=============== to Sun Centred Inertial transformation ================>

		//Check for any errors in building the basedata data structure
		argPre := bd.GetItem("argPre")

		//Then build the Euler transformation matrices
		e3LittleOmega := skymath.E3(argPre)
		eI1 := skymath.E1(dataops.EarthInc)
		e3BigOmega := skymath.E3(dataops.EarthNode)

		//And multiply them to get the OPI to SCI transformation matrix
		step6 := skymath.Mply(eI1, e3LittleOmega)
		opeTOsci := skymath.Mply(e3BigOmega, step6)

		//<============= Calculate the Earth Centered Earth Fixed =================>
		//<================== to East North Up transformation =====================>

		// prSunSCI := skymath.Vply(oppTOsci, prSunOP)
		prSunENU := skymath.Vply(oppTOenu, prSunOP)
		erSunSCI := skymath.Vply(opeTOsci, erSunOP)
		erSunENU := skymath.Vply(sciTOenu, erSunSCI)

		step7 := skymath.Mply(eLam2, eMPhi3)
		ecefTOenu := skymath.Mply(eNU, step7)

		erLocENU := skymath.Vply(ecefTOenu, locTOearthRecef)
		// opeTOenu := skymath.Mply(sciTOenu, opeTOsci)

		// prM(eMGamma3, "earthMinusGAMMA3")
		// prM(sciTOenu, "sciTOenu")
		// prM(oppTOsci, "oppTOsci")
		// prM(opeTOsci, "opeTOsci")
		// prM(opeTOenu, "opeTOenu")
		// prM(ecefTOenu, "ecefTOenu")

		//< Add the vectors up to find the location of the planet in the local sky
		second := skymath.Vadd(erSunENU, erLocENU)
		prENU := skymath.Vsub(prSunENU, second)
		fmt.Printf("\n")
		fmt.Printf("Local to Planet vector in ENU: %e   %e   %e\n", prENU[0], prENU[1], prENU[2])
		fmt.Printf("\n")

		beta, epsilon := skymath.CalcBetaEpsilon(prENU)
		fmt.Printf("Beta: %f, Epsilon: %f\n", beta, epsilon)
	}
}
