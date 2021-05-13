package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"

	"skyham/pkg/dataops"
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
		gt := p.makeGt()
		fmt.Println("GT: ", gt)
		err := basedata.JTime(gt)
		if err != nil {
			fmt.Println("did not get a JD: ", err)
		}
		fmt.Printf("Julian Days JD: = %10.4f\n", basedata["now"].Value)

		bd.CalcPeriod()
		fmt.Println("Earth period: ", basedata["earthPeriod"])
		fmt.Println("Planet period: ", basedata["planetPeriod"])

		bd.CalcOPangles()
		fmt.Println("Earth M0: ", basedata["meanAno"])
		fmt.Println("Earth Arg of Perrifocus: ", basedata["argPre"])
		fmt.Println("Planet M0: ", basedata["planetMeanAno"])
		fmt.Println("Planet Arg of Perrifocus: ", basedata["planetArgPre"])

		bd.CalcM()
		fmt.Println("Earth M: ", basedata["earthM"])
		fmt.Println("Planet M: ", basedata["planetM"])

		bd.CalcE()
		fmt.Println("Earth E: ", (basedata["earthE"].Value/math.Pi)*180.0)
		fmt.Println("Planet E: ", (basedata["planetE"].Value/math.Pi)*180.0)
		// printaid(basedata)
	}

}

func printaid(basedata map[string]dataops.BaseItem) {
	for _, x := range basedata {
		fmt.Println("Name: ", x.Name)
		fmt.Println("Value: ", x.Value)
		fmt.Println("Numonic: ", x.Numonic)
		fmt.Println("Description: ", x.Description)
		fmt.Println()
	}
}

func printIntro() {
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("+   This program calculates the azimuth and elevation of a   +")
	fmt.Println("+     specified planet at a specified location and time      +")
	fmt.Println("+                                                            +")
	fmt.Println("+     you can get latitude, longitude and elevation from     +")
	fmt.Println("+       your smart phone compass application or the web      +")
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
}

func getProfile() *profiles {
	return &profiles{
		year: profile{
			firstMsg:  "Enter the year of interest: ",
			secondMsg: "The year of interest is: ",
			input:     "2014",
		},
		month: profile{
			firstMsg:  "Enter the number of the month of interest: ",
			secondMsg: "The month is: ",
			input:     "3",
		},
		day: profile{
			firstMsg:  "Enter the number of the day of interest: ",
			secondMsg: "The day is: ",
			input:     "22",
		},
		hour: {
			firstMsg:  "Enter the hour of interest (24 hour format): ",
			secondMsg: "The hour of interest is: ",
			input:     "21",
		},
		minute: {
			firstMsg:  "Enter the minute of interest: ",
			secondMsg: "The minute of interest is: ",
			input:     "0",
		},
		second: {
			firstMsg:  "Don't enter the second of interest: ",
			secondMsg: "The second of interest is: ",
			input:     "0",
		},
		latitude: {
			firstMsg:  "Enter latitude (N-S) of interest (see example): ",
			secondMsg: "Latitude of interest is: ",
			input:     "40 90 17 N",
		},
		longitude: {
			firstMsg:  "Enter longitude (E-W) of interest (see exmaple): ",
			secondMsg: "Longitude of interest is: ",
			input:     "74 30 37 W",
		},
		elevation: {
			firstMsg:  "Enter the elevation of interst in feet: ",
			secondMsg: "Elevation is: ",
			input:     "130",
		},
	}
}

func (p *profiles) listItems(reader *bufio.Reader) {
	pp := *p
	if p != nil {
		fmt.Println("the current data is:")
		for i, item := range profileList {
			fmt.Printf("%d. %s: %s\n", i+1, item, pp[item].input)
		}
	}
}

//getInput for updating the profiles
func (p profiles) getInput(reader *bufio.Reader) *profiles {
	var c = false
	eol := "\n"
	for {
		fmt.Println("Enter the number of item you want to change")
		fmt.Println("Enter c to continue")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, eol)
		switch input {
		case "1":
			p[year] = getNumItem(p[year], reader)
		case "2":
			p[month] = getNumItem(p[month], reader)
		case "3":
			p[day] = getNumItem(p[day], reader)
		case "4":
			p[hour] = getNumItem(p[hour], reader)
		case "5":
			p[minute] = getNumItem(p[minute], reader)
		case "6":
			p[second] = getNumItem(p[minute], reader)
		case "7":
			p[latitude] = getNumItem(p[latitude], reader)
		case "8":
			p[longitude] = getNumItem(p[longitude], reader)
		case "9":
			p[elevation] = getNumItem(p[elevation], reader)
		case "c":
			c = true
		case "C":
			c = true
		case "q":
			os.Exit(1)
		default:
			fmt.Println("You made a mistake, try again")
		}
		if c {
			break
		}
	}
	return &p
}

//getNumItem for updating individual items
func getNumItem(p profile, reader *bufio.Reader) profile {
	fmt.Println(p.firstMsg)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSuffix(input, "\n")
	p.input = input
	return p
}

func (p profiles) makeGt() []string {
	g := make([]string, 6)
	g[0] = p[year].input
	g[1] = p[month].input
	g[2] = p[day].input
	g[3] = p[hour].input
	g[4] = p[minute].input
	g[5] = p[second].input
	return g
}
