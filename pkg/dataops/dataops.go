package dataops

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	"github.com/Saied74/skyham/pkg/skymath"
)

//BaseItem is the basic data used for the parameters of the calculations
type BaseItem struct {
	Name        string
	Value       float64
	Numonic     string
	Description string
}

//BaseItems is the main data structure that holds most if not all of the data
type BaseItems map[string]BaseItem

//ProcInputs takes a pointer to the input structure and loads BaseItems
func (bd *BaseItems) ProcInputs(pack []string) error {
	dd := *bd
	g := pack[0:5]
	err := bd.JTime(g)
	if err != nil {
		return fmt.Errorf("time format is incorrect %v", err)
	}
	lat := strings.Split(g[6], " ")
	if len(lat) != 3 {
		return fmt.Errorf("Latitude had %d fields it should have 3", len(lat))
	}
	d, err := strconv.ParseFloat(lat[0], 64)
	if err != nil {
		return fmt.Errorf("lattitude degrees did not convert to number %v", err)
	}
	m, err := strconv.ParseFloat(lat[1], 64)
	if err != nil {
		return fmt.Errorf("lattitude minutes did not convert to number %v", err)
	}
	s, err := strconv.ParseFloat(lat[2], 64)
	if err != nil {
		return fmt.Errorf("lattitude seconds did not convert to number %v", err)
	}
	dd["loclat"] = BaseItem{
		Name:        "locallat",
		Value:       d + m/60 + s/3600.0,
		Numonic:     "lambda",
		Description: "lattitude at the location of interest",
	}

	long := strings.Split(g[7], " ")
	if len(long) != 3 {
		return fmt.Errorf("longitude had %d fields it should have 3", len(long))
	}
	d, err = strconv.ParseFloat(long[0], 64)
	if err != nil {
		return fmt.Errorf("longitude degrees did not convert to number %v", err)
	}
	m, err = strconv.ParseFloat(long[1], 64)
	if err != nil {
		return fmt.Errorf("longitude minutes did not convert to number %v", err)
	}
	s, err = strconv.ParseFloat(long[2], 64)
	if err != nil {
		return fmt.Errorf("longitude seconds did not convert to number %v", err)
	}
	dd["locllong"] = BaseItem{
		Name:        "locallong",
		Value:       d + m/60 + s/3600.0,
		Numonic:     "phi",
		Description: "longitude at the location of interest",
	}

	ele, err := strconv.ParseFloat(g[8], 64)
	if err != nil {
		return fmt.Errorf("Elevation did not convert to a number %v", err)
	}
	dd["localelev"] = BaseItem{
		Name:        "localelev",
		Value:       ((ele * 12.0) / 2.54) / 100.0,
		Numonic:     "h",
		Description: "elevation at the location of interest",
	}
	return nil
}

//ReadData reads the base data file and returns a map of the input base data
func (bd *BaseItems) ReadData(filename string) {
	var err error
	dd := *bd
	dat, err := ioutil.ReadFile(filename)
	check("Base data read error: ", err)
	if len(dat) == 0 {
		check("base data file was empty: ", nil)
	}
	lines := strings.Split(string(dat), "\n")
	for i, line := range lines {
		lineItems := strings.Split(line, ",")
		if i == 0 || len(lineItems) < 4 {
			continue
		}
		dd[lineItems[0]] = assign(lineItems, lineItems[0])
	}
}

//CalcPeriod calculates the period of earth and the designated planet around the sun
func (bd *BaseItems) CalcPeriod() {
	dd := *bd
	numerator := (math.Pi * math.Pi) * 4.0
	denominator := dd["gravity"].Value * (dd["sunMass"].Value + dd["earthMass"].Value)
	x := dd["earthSMA"].Value * dd["au"].Value
	tSquared := (numerator / denominator) * x * x * x
	item := BaseItem{
		Name:        "earthPeriod",
		Value:       math.Sqrt(tSquared) / (24.0 * 3600.0),
		Numonic:     "T",
		Description: "period of earth around the sun",
	}
	dd["earthPeriod"] = item

	denominator = dd["gravity"].Value * (dd["sunMass"].Value + dd["planetMass"].Value)
	x = dd["planetSMA"].Value * dd["au"].Value
	tSquared = (numerator / denominator) * x * x * x
	item = BaseItem{
		Name:        "planetPeriod",
		Value:       math.Sqrt(tSquared) / (24.0 * 3600.0),
		Numonic:     "T",
		Description: "period of planet around the sun",
	}
	dd["planetPeriod"] = item
}

//CalcOPangles calculates the orbital plane angles, mean anomaly and argument of the Perrifocus
func (bd *BaseItems) CalcOPangles() {
	dd := *bd
	m := dd["meanLong"].Value - dd["longPre"].Value
	dd["meanAno"] = BaseItem{
		Name:        "meanAno",
		Value:       m,
		Numonic:     "M0",
		Description: "Mean Anomaly",
	}
	omega := dd["longPre"].Value - dd["earthNode"].Value
	dd["argPre"] = BaseItem{
		Name:        "argPre",
		Value:       omega,
		Numonic:     "omega",
		Description: "argument of the perrifocus",
	}
	mp := dd["planetMeanLong"].Value - dd["planetLongPre"].Value
	dd["planetMeanAno"] = BaseItem{
		Name:        "meanAno",
		Value:       mp,
		Numonic:     "M0",
		Description: "Mean Anomaly",
	}
	omegap := dd["planetLongPre"].Value - dd["planetNode"].Value
	dd["planetArgPre"] = BaseItem{
		Name:        "PlanetArgPre",
		Value:       omegap,
		Numonic:     "omega",
		Description: "argument of the perrifocus",
	}
}

//JTime converts Georgian time to Julian Days (JD)
func (bd *BaseItems) JTime(gTStrings []string) error {
	//gTStrings is a slice of strings composed of year, month, day, minute and second
	//first, we convert it to
	gTime, err := gConvert(gTStrings)
	if err != nil {
		return fmt.Errorf("Time strings did not convert to numbers %v", err)
	}
	//Then check numbers for validity
	err = gCheck(gTime)
	if err != nil {
		return fmt.Errorf("Time strings were not valid %v", err)
	}
	year := gTime[0]
	month := gTime[1]
	day := gTime[2]
	hour := gTime[3]
	minute := gTime[4]
	second := gTime[5]

	//This is the formula from Don Koks' paper
	a := math.Floor((14.0 - month) / 12)
	y := year + 4800.0 - a
	m := month + 12.0*a - 3.0

	jd1 := day + math.Floor((153*m+2.0)/5.0) + 365*y + math.Floor(y/4.0)
	jd2 := math.Floor(y/100.0) - math.Floor(y/400.0) + 32045
	jd3 := (hour - 12.0 + minute/60.0 + second/3600.0) / 24.0
	dd := *bd
	dd["now"] = BaseItem{
		Name:        "now",
		Value:       jd1 - jd2 + jd3,
		Numonic:     "t",
		Description: "desired time",
	}
	return nil
}

//CalcM calculates the earth and planet mean anomaly at the specified epoch
func (bd *BaseItems) CalcM() {
	dd := *bd
	m0 := (dd["meanAno"].Value / 180.0) * math.Pi
	m := m0 + (2.0*math.Pi*(dd["now"].Value-dd["t0"].Value))/dd["earthPeriod"].Value
	dd["earthM"] = BaseItem{
		Name:        "earthM",
		Value:       m,
		Numonic:     "M",
		Description: "Radians: Earth mean anomaly at desired epoch",
	}

	m0 = (dd["planetMeanAno"].Value / 180.0) * math.Pi
	m = m0 + (2.0*math.Pi*(dd["now"].Value-dd["t0"].Value))/dd["planetPeriod"].Value
	dd["planetM"] = BaseItem{
		Name:        "planetM",
		Value:       m,
		Numonic:     "M",
		Description: "Radians: Planet mean anomaly at desired epoch",
	}
}

//CalcE calculates earth and planet Eccentricity at specified epoch
func (bd *BaseItems) CalcE() {
	dd := *bd
	bigM := dd["earthM"].Value
	e := dd["earthEcc"].Value
	bigE := 0.0
	oldE := 0.0
	for {
		oldE = bigE
		bigE = e*math.Sin(bigE) + bigM
		test := 100.0 * math.Abs((bigE-oldE)/oldE)
		if test < 0.0001 { //0.01% error
			break
		}
		dd["earthE"] = BaseItem{
			Name:        "earthE",
			Value:       bigE,
			Numonic:     "E",
			Description: "Earth Eccentrioc Anomaly",
		}
	}

	bigM = dd["planetM"].Value
	e = dd["planetEcc"].Value
	bigE = 0.0
	for {
		oldE = bigE
		bigE = e*math.Sin(bigE) + bigM
		test := 100.0 * math.Abs((bigE-oldE)/oldE)
		if test < 0.01 { //0.01% error
			break
		}
		dd["planetE"] = BaseItem{
			Name:        "planetE",
			Value:       bigE,
			Numonic:     "E",
			Description: "Planet Eccentrioc Anomaly",
		}
	}
}

//PlanetOPVec returns the planet vector in its own orbital plane inertial frame
func (bd *BaseItems) PlanetOPVec() skymath.Vec {
	dd := *bd
	v := skymath.Vec{}
	e := dd["planetEcc"].Value
	bigE := dd["planetE"].Value
	a := dd["planetSMA"].Value * dd["au"].Value
	b := a * math.Sqrt(1-e*e)
	v[0] = a * (math.Cos(bigE) - e)
	v[1] = b * math.Sin(bigE)
	v[2] = 0.0
	return v
}

//EarthOPVec returns the planet vector in its own orbital plane inertial frame
func (bd *BaseItems) EarthOPVec() skymath.Vec {
	dd := *bd
	v := skymath.Vec{}
	e := dd["earthEcc"].Value
	bigE := dd["earthE"].Value
	a := dd["earthSMA"].Value * dd["au"].Value
	b := a * math.Sqrt(1-e*e)
	v[0] = a * (math.Cos(bigE) - e)
	v[1] = b * math.Sin(bigE)
	v[2] = 0.0
	return v
}

//SidAngle updates the sideral angle at the specifed time since J2000
func (bd *BaseItems) SidAngle() {
	dd := *bd
	sidDay := 23.0 + (56.0 / 60.0) + (4.09890 / 3600.0)
	deltaT := dd["now"].Value - dd["t0"].Value
	sidAngle := (360.0 / sidDay) * deltaT * 24.0
	sidAngle += dd["sideralJ2000"].Value
	dd["sidAngle"] = BaseItem{
		Name:        "sidAngle",
		Value:       sidAngle,
		Numonic:     "gamma",
		Description: "Greenwich Sideral Angle at the specified time",
	}
}

//EarthPrecession updates the presession of the earth since J2000
func (bd *BaseItems) EarthPrecession() {
	dd := *bd
	deltaT := dd["now"].Value - dd["t0"].Value
	p := (360.0 * deltaT) / (25770.0 * 365.25)
	dd["earthP"] = BaseItem{
		Name:        "earthP",
		Value:       p,
		Numonic:     "p",
		Description: "Earth precession since the J2000 epoch",
	}
}
