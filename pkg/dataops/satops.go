package dataops

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Saied74/skyham/pkg/skymath"
)

//ProcSatInputs takes a pointer to the input structure and loads BaseItems
func (bd *BaseItems) ProcSatInputs(pack []string) error {
	//monthDays is the cumulative number of days to that month.  For example,
	//the second element contains the cumulative days of January and February.
	monthDays := []float64{
		31.0,
		31.0 + 28.0,
		31.0 + 28.0 + 31.0,
		31.0 + 28.0 + 31.0 + 30,
		31.0 + 28.0 + 31.0 + 30 + 31,
		31.0 + 28.0 + 31.0 + 30 + 31 + 30,
		31.0 + 28.0 + 31.0 + 30 + 31 + 30 + 31,
		31.0 + 28.0 + 31.0 + 30 + 31 + 30 + 31 + 31,
		31.0 + 28.0 + 31.0 + 30 + 31 + 30 + 31 + 31 + 30,
		31.0 + 28.0 + 31.0 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
		31.0 + 28.0 + 31.0 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
		31.0 + 28.0 + 31.0 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
	}
	dd := *bd
	// g := pack[0:6]
	gTime, err := gConvert(pack[0:6])
	if err != nil {
		return fmt.Errorf("Time strings did not convert to numbers %v", err)
	}
	//Then check numbers for validity
	err = gCheck(gTime)
	if err != nil {
		return fmt.Errorf("Time strings were not valid %v", err)
	}
	// TODO: Leave out the year for now and come back to it later
	// year := gTime[0]
	// month := gTime[1]
	day := gTime[2]
	hour := gTime[3]
	minute := gTime[4]
	second := gTime[5]
	intMonth, _ := strconv.ParseInt(pack[1], 10, 32)
	// TODO: This is a problem for January - need to fix
	days := monthDays[intMonth-2] + day + (hour+(minute+second/60)/60)/24

	dd["tempNow"] = BaseItem{
		Name:        "tempNow",
		Value:       days,
		Numonic:     "t",
		Description: "desired time temporary",
	}

	g := pack[0:6]
	err = bd.JTime(g)
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

//ExtractSatData extracts satellite orbital data from the two lines
//and populates basedata
func (bd *BaseItems) ExtractSatData(line1, line2 string) error {
	dd := *bd
	if len(line1) != 69 {
		return fmt.Errorf("Length of line 1 is not 69, it is %d", len(line1))
	}
	if len(line2) != 69 {
		return fmt.Errorf("Length of line 2 is not 69, it is %d", len(line2))
	}
	lineNo := string(line1[0])
	if lineNo != "1" {
		return fmt.Errorf("number of first line was not 1, it was %s", lineNo)
	}
	satNo1 := line1[2:7]

	// TODO: I am going to make simple assumptin that the epoch year is the same
	//as current year and will fix this later after the program works.
	// epochYear, err := tleFloat(line1[18:20])
	// if err != nil {
	// 	return fmt.Errorf("Epoch year %s %v", line1[18:20], err)
	// }
	// if epochYear >= 57 && epochYear <= 99 {
	// 	epochYear += 1900
	// }
	// if epochYear >= 0 && epochYear < 57 {
	// 	epochYear += 2000
	// }
	// if epochYear < 0 {
	// 	return fmt.Errorf("error epoch year %s was negative", line1[18:20])
	// }
	// yearData := []float64{epochYear, 1, 1, 0, 0, 0}

	epochDay, err := tleFloat(line1[20:32])
	if err != nil {
		return fmt.Errorf("Epoch day %s %v", line1[20:32], err)
	}
	// fmt.Println("calculating epoch", yearData, epochDay)
	// epoch, err := JSatTime(yearData, epochDay)
	// if err != nil {
	// 	return fmt.Errorf("Satellite epoch %s did not convert to JD time", line1[18:32])
	// }
	dd["satEpoch"] = BaseItem{
		Name:        "satEpoch",
		Value:       epochDay,
		Numonic:     "satT0",
		Description: "used with the satellite mean anomaly",
	}
	lineNo = string(line2[0])
	if lineNo != "2" {
		return fmt.Errorf("number of second line was not 2, it was %s", lineNo)
	}
	satNo2 := line2[2:7]
	if satNo1 != satNo2 {
		return fmt.Errorf("satelline numbers did not match, %s and %s", satNo1, satNo2)
	}
	inc, err := tleFloat(line2[8:16])
	if err != nil {
		return fmt.Errorf("inclination %s %v", line2[8:16], err)
	}
	dd["planetInc"] = BaseItem{
		Name:        "planetInc",
		Value:       inc,
		Numonic:     "i",
		Description: "Degrees: inclination of the satellite orbital plane",
	}
	node, err := tleFloat(line2[17:25])
	if err != nil {
		return fmt.Errorf("satellite asending node %s %v", line2[17:25], err)
	}
	dd["planetNode"] = BaseItem{
		Name:        "planetNode",
		Value:       node,
		Numonic:     "OMEGA",
		Description: "Degrees: longitude of the satellite asending node",
	}
	e, err := tleFloat("0." + line2[26:33])
	if err != nil {
		return fmt.Errorf("satellite eccentricity %s %v", line2[26:33], err)
	}
	dd["planetEcc"] = BaseItem{
		Name:        "planetEcc",
		Value:       e, // * (10e-8),
		Numonic:     "e",
		Description: "Eccentricity of the satellite orbit",
	}
	littleOmega, err := tleFloat(line2[34:42])
	if err != nil {
		return fmt.Errorf("satellite argument of prefocus %s %v", line2[34:42], err)
	}
	dd["planetArgPre"] = BaseItem{
		Name:        "planetArgPre",
		Value:       littleOmega,
		Numonic:     "omega",
		Description: "Degrees: satellite argument of prefocus",
	}
	meanAno, err := tleFloat(line2[43:51])
	if err != nil {
		return fmt.Errorf("mean anomaly %s %v", line2[43:51], err)
	}
	dd["planetMeanAno"] = BaseItem{
		Name:        "planetMeanAno",
		Value:       meanAno,
		Numonic:     "M0",
		Description: "satellite mean anomaly is directly give, not calculated",
	}
	meanMotion, err := tleFloat(line2[52:63])
	if err != nil {
		return fmt.Errorf("mean motion %s %v", line2[52:63], err)
	}
	dd["planetPeriod"] = BaseItem{
		Name:        "planetPeriod",
		Value:       1 / meanMotion,
		Numonic:     "T",
		Description: "Mean motion in revolutions per day, no need to calc T",
	}
	return nil
}

//CalcA calculates the major axis or the Satellite
func (bd *BaseItems) CalcA() {
	t := bd.GetItem("planetPeriod")
	//ignoring the satellite mass compared to the earth mass
	t = t * 24.0 * 3600.0
	aCubed := (earthMass * gravity * t * t) / (4.0 * math.Pi * math.Pi)
	p := 1.0 / 3.0
	a := math.Pow(aCubed, p)
	dd := *bd
	dd["planetSMA"] = BaseItem{
		Name:        "planetSMA",
		Value:       a / au,
		Numonic:     "a",
		Description: "Satellite Semi Major Axis (SMA)",
	}
}

//CalcSatM calculates the earth and planet mean anomaly at the specified epoch
func (bd *BaseItems) CalcSatM() {
	dd := *bd
	epoch := bd.GetItem("satEpoch")
	m0 := (dd["planetMeanAno"].Value / 180.0) * math.Pi
	delT := dd["tempNow"].Value - epoch // TODO: need to add error checking
	fmt.Println("Delta T", delT, epoch, dd["tempNow"].Value)
	// if delT < 0.0 {
	// 	delT = 0.0
	// }
	m := m0 + (2.0*math.Pi*delT)/dd["planetPeriod"].Value

	m = math.Mod(m, 2.0*math.Pi)

	dd["planetM"] = BaseItem{
		Name:        "planetM",
		Value:       skymath.ToDegrees(m),
		Numonic:     "M",
		Description: "Degrees: Planet mean anomaly at desired epoch",
	}
}
