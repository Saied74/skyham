package dataops

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	"github.com/Saied74/skyham/pkg/skymath"
)

const (
	t0           = 2451545.0       // (t0)Julian days (JD) - EPOCH J2000
	au           = 1.495978707E+11 // (au) meters Astronomical Unit in meters
	sunMass      = 1.989E+30       // (M) Mass of the Sun in kg
	gravity      = 6.67384E-11     // (G) Gravitational constant in SI units
	sideralJ2000 = 280.46          // (Gamma) Greenwich sideral angle at J2000 epoch in degrees
	//EarthTilt is the tilt of the earth's axis of rotation relative to the equotorial plane
	EarthTilt = 23.439     // (Tau) Earth's tilt in degrees
	earthMass = 5.9736E+24 // (m) Mass of earth in kg
	earthSMA  = 1.00000011 // (a) Earth's semi major axis in AU
	//EarthInc is the inclination of the mean earth orbital plane
	EarthInc = 0.00005    // (i) earth inclination
	earthEcc = 0.01671022 // (e) Eccentricity of the Earth
	//EarthNode is the longitude of the Earth's asencing node
	EarthNode          = -11.26064    // (OMEGA) Earth longitude of the ascending node in degrees
	earthLongPrehelian = 102.94719    // longitude of the perihelion
	earthmeanLong      = 100.46435    //mean longitude
	equotorialA        = 6378137.0    // (a) equotorial semi-major axis of earth
	polarB             = 6356752.3142 // (b) polar semi-major axis of earth
	//PlanetName stores the name of the planet of interest for file lookup
	PlanetName = "planetName"
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
	g := pack[0:6]
	err := bd.JTime(g)
	if err != nil {
		return fmt.Errorf("time format is incorrect %v", err)
	}
	lat := strings.Split(pack[6], " ")
	if len(lat) != 2 {
		return fmt.Errorf("Latitude had %d fields it should have 2", len(lat))
	}
	localLat, err := strconv.ParseFloat(lat[0], 64)
	if err != nil {
		return fmt.Errorf("lattitude %s did not convert to number %v", lat[0], err)
	}
	lat[1] = strings.ToUpper(lat[1])
	switch lat[1] {
	case "S":
		localLat = -localLat
	case "N":
	default:
		return fmt.Errorf("latitude must be N or S (uppeer or lower case) %s", lat[1])
	}
	dd["locallat"] = BaseItem{
		Name:        "locallat",
		Value:       localLat,
		Numonic:     "lambda",
		Description: "Degrees: lattitude at the location of interest",
	}

	long := strings.Split(pack[7], " ")
	if len(long) != 2 {
		return fmt.Errorf("longitude had %d fields it should have 2", len(long))
	}
	localLong, err := strconv.ParseFloat(long[0], 64)
	if err != nil {
		return fmt.Errorf("longitude %s did not convert to number %v", long[0], err)
	}
	long[1] = strings.ToUpper(long[1])
	switch long[1] {
	case "W":
		localLong = -localLong
	case "E":
	default:
		return fmt.Errorf("longitude must be E or W (uppeer or lower case) %s", long[1])
	}

	dd["locallong"] = BaseItem{
		Name:        "locallong",
		Value:       localLong,
		Numonic:     "phi",
		Description: "Degrees: longitude at the location of interest",
	}

	ele, err := strconv.ParseFloat(pack[8], 64)
	if err != nil {
		return fmt.Errorf("Elevation did not convert to a number %v", err)
	}
	dd["localelev"] = BaseItem{
		Name:        "localelev",
		Value:       ((ele * 12.0) / 2.54) / 100.0,
		Numonic:     "h",
		Description: "Meters: elevation at the location of interest",
	}
	dd["planetName"] = BaseItem{
		Name:        "planetName",
		Value:       0.0,
		Numonic:     strings.ToLower(pack[9]),
		Description: "Name of the planet of interest",
	}
	return nil
}

//ReadData reads the base data file and returns a map of the input base data
func (bd *BaseItems) ReadData(filename string) error {
	var err error
	dd := *bd
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading planet file %v", err)
	}
	if len(dat) == 0 {
		return fmt.Errorf("planet %v file was empty", filename)
	}
	lines := strings.Split(string(dat), "\n")
	for i, line := range lines {
		lineItems := strings.Split(line, ",")
		if i == 0 || len(lineItems) < 4 {
			continue
		}
		dd[lineItems[0]] = assign(lineItems, lineItems[0])
	}
	return nil
}

//CalcPeriod calculates the period of earth and the designated planet around the sun
func (bd *BaseItems) CalcPeriod() {
	dd := *bd
	numerator := (math.Pi * math.Pi) * 4.0
	denominator := gravity * (sunMass + earthMass)
	x := earthSMA * au
	tSquared := (numerator / denominator) * x * x * x
	item := BaseItem{
		Name:        "earthPeriod",
		Value:       math.Sqrt(tSquared) / (24.0 * 3600.0),
		Numonic:     "T",
		Description: "period of earth around the sun",
	}
	dd["earthPeriod"] = item

	denominator = gravity * (sunMass + dd["planetMass"].Value)
	x = dd["planetSMA"].Value * au
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
	m := earthmeanLong - earthLongPrehelian
	dd["meanAno"] = BaseItem{
		Name:        "meanAno",
		Value:       m,
		Numonic:     "M0",
		Description: "Degrees: Mean Anomaly",
	}
	omega := earthLongPrehelian - EarthNode
	dd["argPre"] = BaseItem{
		Name:        "argPre",
		Value:       omega,
		Numonic:     "omega",
		Description: "Degrees: argument of the perrifocus",
	}
	planetMeanLong, ok := dd["planetMeanLong"]
	if !ok {
		check("bad lookup: planetMeanLong", nil)
	}
	planetLongPre, ok := dd["planetLongPre"]
	if !ok {
		check("bad lookup: planetLongPre", nil)
	}
	mp := planetMeanLong.Value - planetLongPre.Value
	dd["planetMeanAno"] = BaseItem{
		Name:        "meanAno",
		Value:       mp,
		Numonic:     "M0",
		Description: "Degrees: Mean Anomaly",
	}
	omegap := dd["planetLongPre"].Value - dd["planetNode"].Value
	dd["planetArgPre"] = BaseItem{
		Name:        "PlanetArgPre",
		Value:       omegap,
		Numonic:     "omega",
		Description: "Degrees: argument of the perrifocus",
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
func (bd *BaseItems) CalcM(sat string) {
	dd := *bd
	var epoch float64
	if sat == "satellite" || sat == "sat" {
		t, ok := dd["satEpoch"]
		if !ok {
			check("bad lookup: satEpoch", nil)
		}
		epoch = t.Value
	} else {
		epoch = t0
	}

	meanAno, ok := dd["meanAno"]
	if !ok {
		check("bad lookup: meanAno", nil)
	}
	m0 := skymath.ToRadians(meanAno.Value)
	now, ok := dd["now"]
	if !ok {
		check("bad lookup: now", nil)
	}
	earthPeriod, ok := dd["earthPeriod"]
	if !ok {
		check("bad lookup: earthPeriod", nil)
	}
	delT := now.Value - epoch
	// if delT < 0.0 {
	// 	delT = 0.0
	// }
	m := m0 + (2.0*math.Pi*delT)/earthPeriod.Value
	m = math.Mod(m, 2.0*math.Pi)
	dd["earthM"] = BaseItem{
		Name:        "earthM",
		Value:       skymath.ToDegrees(m),
		Numonic:     "M",
		Description: "Degrees: Earth mean anomaly at desired epoch",
	}

	m0 = (dd["planetMeanAno"].Value / 180.0) * math.Pi
	delT = dd["now"].Value - epoch
	if delT < 0.0 {
		delT = 0.0
	}
	m = m0 + (2.0*math.Pi*delT)/dd["planetPeriod"].Value

	m = math.Mod(m, 2.0*math.Pi)

	dd["planetM"] = BaseItem{
		Name:        "planetM",
		Value:       skymath.ToDegrees(m),
		Numonic:     "M",
		Description: "Degrees: Planet mean anomaly at desired epoch",
	}
}

//CalcE calculates earth and planet Eccentricity at specified epoch
func (bd *BaseItems) CalcE(sat string) {
	var bigM, bigE, oldE float64
	dd := *bd
	if sat != "satellite" {
		earthM, ok := dd["earthM"]
		if !ok {
			check("bad lookup: earthM", nil)
		}
		bigM = earthM.Value
		bigM = skymath.ToRadians(bigM)
		bigE = 0.0
		// oldE = 0.0
		for {
			oldE = bigE
			bigE = earthEcc*math.Sin(bigE) + bigM
			test := 100.0 * math.Abs((bigE-oldE)/oldE)
			if test < 0.01 { //0.01% error
				break
			}
		}
		bigE = math.Mod(bigE, 2.0*math.Pi)
		dd["earthE"] = BaseItem{
			Name:        "earthE",
			Value:       skymath.ToDegrees(bigE),
			Numonic:     "E",
			Description: "Degrees: Earth Eccentrioc Anomaly",
		}
	}

	planetM, ok := dd["planetM"]
	if !ok {
		check("bad lookup: planetM", nil)
	}
	bigM = planetM.Value
	bigM = skymath.ToRadians(bigM)
	planetEcc, ok := dd["planetEcc"]
	if !ok {
		check("bad lookup: planetEcc", nil)
	}
	e := planetEcc.Value
	bigE = 0.0

	for {
		oldE = bigE
		bigE = e*math.Sin(bigE) + bigM
		test := 100.0 * math.Abs((bigE-oldE)/oldE)
		if test < 0.01 { //0.01% error
			break
		}
	}

	bigE = math.Mod(bigE, 2.0*math.Pi)
	dd["planetE"] = BaseItem{
		Name:        "planetE",
		Value:       skymath.ToDegrees(bigE),
		Numonic:     "E",
		Description: "Degrees: Planet Eccentrioc Anomaly",
	}
}

//PlanetOPVec returns the planet vector in its own orbital plane inertial frame
func (bd *BaseItems) PlanetOPVec() skymath.Vec {
	dd := *bd
	v := skymath.Vec{}
	planetEcc, ok := dd["planetEcc"]
	if !ok {
		check("bad lookup: planetEcc", nil)
	}
	e := planetEcc.Value
	planetE, ok := dd["planetE"]
	if !ok {
		check("bad lookup: planetE", nil)
	}
	bigE := skymath.ToRadians(planetE.Value)
	planetSMA, ok := dd["planetSMA"]
	if !ok {
		check("Bad lookup: planetSMA", nil)
	}
	a := planetSMA.Value * au
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
	earthE, ok := dd["earthE"]
	if !ok {
		check("bad lookup: earthE", nil)
	}
	bigE := earthE.Value
	bigE = skymath.ToRadians(bigE)
	a := earthSMA * au
	b := a * math.Sqrt(1.0-earthEcc*earthEcc)
	v[0] = a * (math.Cos(bigE) - earthEcc)
	v[1] = b * math.Sin(bigE)
	v[2] = 0.0
	return v
}

//SidAngle updates the sideral angle at the specifed time since J2000
func (bd *BaseItems) SidAngle() {
	dd := *bd
	sidDay := 23.0 + (56.0 / 60.0) + (4.09890 / 3600.0)
	now, ok := dd["now"]
	if !ok {
		check("bad lookup: now", nil)
	}
	deltaT := now.Value - t0
	sidAngle := (360.0 / sidDay) * deltaT * 24.0
	sidAngle += sideralJ2000

	dd["sidAngle"] = BaseItem{
		Name:        "sidAngle",
		Value:       math.Mod(sidAngle, 360),
		Numonic:     "gamma",
		Description: "Degrees - Greenwich Sideral Angle at the specified time",
	}
}

//EarthPrecession updates the presession of the earth since J2000
func (bd *BaseItems) EarthPrecession() {
	dd := *bd
	deltaT := dd["now"].Value - t0
	p := (360.0 * deltaT) / (25770.0 * 365.25)
	dd["earthP"] = BaseItem{
		Name:        "earthP",
		Value:       p,
		Numonic:     "p",
		Description: "Degrees: Earth precession since the J2000 epoch",
	}
}

//CalcLocalVec calculates the coordinates of the local coordinate in earth fixed fra,e
func (bd *BaseItems) CalcLocalVec() skymath.Vec {
	dd := *bd
	a := float64(equotorialA)
	b := float64(polarB)
	locallat, ok := dd["locallat"]
	if !ok {
		check("bad lookup: locallat", nil)
	}
	lambda := skymath.ToRadians(locallat.Value)

	cLambda := math.Cos(lambda) //Cos(lambda)
	sLambda := math.Sin(lambda) //Sin(lambda)
	aSq := a * a
	bSq := b * b
	kc := aSq * cLambda * cLambda
	ks := bSq * sLambda * sLambda
	k := math.Sqrt(kc + ks)

	// fmt.Printf("Lambda: %e, Phi: %e\n", dd["locallat"].Value, dd["locallong"].Value)
	// fmt.Printf("Elevation: %e\n", dd["localelev"].Value)
	// fmt.Printf("a: %f, b:= %f, k:= %f\n", a, b, k)
	locallong, ok := dd["locallong"]
	if !ok {
		check("bad lookup: locallong", nil)
	}
	phi := skymath.ToRadians(locallong.Value)
	cPhi := math.Cos(phi)
	sPhi := math.Sin(phi)
	locallevel, ok := dd["locallev"]
	if !ok {
		check("bad lookup: locallev", nil)
	}
	aCo := (aSq / k) + locallevel.Value
	bCo := (bSq / k) + locallevel.Value

	return skymath.Vec{
		aCo * cLambda * cPhi,
		aCo * cLambda * sPhi,
		bCo * sLambda,
	}
}
