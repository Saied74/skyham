package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Saied74/skyham/pkg/cli"
	"github.com/Saied74/skyham/pkg/dataops"
	"github.com/Saied74/skyham/pkg/skymath"
)

func main() {
	//build the structure to hold base (including earth and sun) and planet data
	basedata := make(dataops.BaseItems)
	bd := &basedata

	//temparirly locate the CLI here
	printIntro()
	reader := bufio.NewReader(os.Stdin)
	// first get a blank profile.  False is for satellite, true for planets
	var p = *cli.GetProfile(false)

	for {
		//In the following loop, inputs are gathered and the file is read and
		//processed.  If everything is good, it will break and process the data
		for {
			var line1, line2 string
			//list data inside the profiles and edit as needed.
			p = *p.GetInput(reader)
			pack := p.PackageInput()
			err := bd.ProcSatInputs(pack)
			if err != nil {
				fmt.Println("bad input, try again: ", err)
			}

			satName, ok := basedata[dataops.PlanetName]
			if !ok {
				fmt.Println("bad lookup: planetName")
			}
			fileName := "../data/tle.txt"
			//get planet (currently Jupiter) data
			tleSat, err := ioutil.ReadFile(fileName)
			if err != nil {
				fmt.Printf("error in reading the planet file %v\n", err)
				os.Exit(2)
			}
			if len(tleSat) == 0 {
				fmt.Printf("file %s was empty\n", fileName)
			}
			lines := strings.Split(string(tleSat), "\n")

			for i, line := range lines {
				if i%3 == 0 {
					if strings.ToLower(line) == satName.Numonic {
						line1 = lines[i+1]
						line2 = lines[i+2]
						break
					}
				}
			}
			err = bd.ExtractSatData(line1, line2)
			if err != nil {
				fmt.Printf("error in processing the satellite file: %v\n", err)
				os.Exit(3)
			}
			if err == nil && ok {
				break
			}
		}

		//first, we calc the mean and Eccentrioc anomaly
		bd.CalcSatM()
		bd.CalcE("satellite")

		//Then calculate the Sideral angle at the requested time
		bd.SidAngle()
		bd.CalcA()

		//At this point, all input and intermediate data is read and calculated
		bd.PrintBaseItems()

		//<================= Calculate the Orbital Plane Vectors ==================>

		rStoEinOPS := bd.PlanetOPVec() //srEOPS

		//<================= Location to "Earth center" vector ====================>
		rLtoEinECEF := bd.CalcLocalVec() //locTOearthRecef

		//<==== Calculate Sun Centred Interal to East North Up Transformation ====>

		//Check for any errors in building the basedata data structure
		sidAngle := bd.GetItem("sidAngle")
		locallong := bd.GetItem("locallong")
		locallat := bd.GetItem("locallat")

		//Then build the Euler transformation matrices
		eMGamma3 := skymath.E3(-sidAngle)
		eMPhi3 := skymath.E3(-locallong)
		eLam2 := skymath.E2(locallat)
		eNU := skymath.Euler{
			[3]float64{0.0, 1.0, 0.0},
			[3]float64{0.0, 0.0, 1.0},
			[3]float64{1.0, 0.0, 0.0},
		}

		//And multiply them to get the SCI to ENU transformation matrix
		step3 := skymath.Mply(eMPhi3, eMGamma3)
		step4 := skymath.Mply(eLam2, step3)
		muECItoENU := skymath.Mply(eNU, step4) //eciTOenu

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
		muOPStoECI := skymath.Mply(p3BigOmega, step5)      //opsTOeci
		muOPStoENU := skymath.Mply(muECItoENU, muOPStoECI) //opsTOenu

		rStoEinENU := skymath.Vply(muOPStoENU, rStoEinOPS) //srEENU

		//<============= Calculate the Earth Centered Earth Fixed =================>
		//<================== to East North Up transformation =====================>

		step7 := skymath.Mply(eLam2, eMPhi3)
		muECEFtoENU := skymath.Mply(eNU, step7) //ecefTOenu

		rLtoEinENU := skymath.Vply(muECEFtoENU, rLtoEinECEF)

		muOPStoECEF := skymath.Mply(eMGamma3, muOPStoECI)
		rStoEinECEF := skymath.Vply(muOPStoECEF, rStoEinOPS) //rSEecef

		fmt.Printf("\nearth to satellite vector in ecef %v\n", rStoEinECEF)
		var phi float64
		r, t, phi := skymath.CalcSpherical(rStoEinECEF[0], rStoEinECEF[1], rStoEinECEF[2])
		fmt.Printf("ECEF Frame r: %f ECEF Frame theta: %f ECEF Frame elevation: %f\n ", r, t, phi)

		// Add the vectors up to find the location of the planet in the local sky
		rStoLinENU := skymath.Vsub(rStoEinENU, rLtoEinENU) //prENU
		fmt.Printf("\n")
		fmt.Printf("Local to Satellite vector in ENU: %e   %e   %e\n", rStoLinENU[0], rStoLinENU[1], rStoLinENU[2])
		fmt.Printf("\n")

		_, beta, epsilon := skymath.CalcSpherical(rStoLinENU[0], rStoLinENU[1], rStoLinENU[2])
		fmt.Printf("Beta: %f, Epsilon: %f\n", beta, epsilon)
	}
}

func printIntro() {
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("+   This program calculates the azimuth and elevation of     +")
	fmt.Println("+  a specified satellite at a specified location and time    +")
	fmt.Println("+                                                            +")
	fmt.Println("+     you can get latitude, longitude and elevation from     +")
	fmt.Println("+                       the web                              +")
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
}
