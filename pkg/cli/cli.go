package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
	planet    = "planet"
	satellite = "satellite"
)

//Profile and profiles are for collecting the user input.  In an excel program
//they would be contents of cells
type profile struct {
	firstMsg  string
	secondMsg string
	input     string
	opSys     string
}

//Profiles holds the data for command line interface data
type Profiles map[string]profile

var profileList = []string{year, month, day, hour, minute, second, latitude,
	longitude, elevation, planet}

func printaid(basedata map[string]dataops.BaseItem) {
	for _, x := range basedata {
		fmt.Println("Name: ", x.Name)
		fmt.Println("Value: ", x.Value)
		fmt.Println("Numonic: ", x.Numonic)
		fmt.Println("Description: ", x.Description)
		fmt.Println()
	}
}

//GetProfile returns a pre populated profile for the command line interface
func GetProfile(p bool) *Profiles {
	var x string
	// var y string
	if p {
		x = planet
		// y = "Jupiter"
	} else {
		x = satellite
		// y = ""
	}
	return &Profiles{
		year: profile{
			firstMsg:  "Enter the year of interest: ",
			secondMsg: "The year of interest is: ",
			input:     "2021",
		},
		month: profile{
			firstMsg:  "Enter the number of the month of interest: ",
			secondMsg: "The month is: ",
			input:     "5",
		},
		day: profile{
			firstMsg:  "Enter the number of the day of interest: ",
			secondMsg: "The day is: ",
			input:     "24",
		},
		hour: {
			firstMsg:  "Enter the hour of interest (24 hour format): ",
			secondMsg: "The hour of interest is: ",
			input:     "4",
		},
		minute: {
			firstMsg:  "Enter the minute of interest: ",
			secondMsg: "The minute of interest is: ",
			input:     "36",
		},
		second: {
			firstMsg:  "Don't enter the second of interest: ",
			secondMsg: "The second of interest is: ",
			input:     "0",
		},
		latitude: {
			firstMsg:  "Enter latitude (N-S) of interest in decimal format: ",
			secondMsg: "Latitude of interest is: ",
			input:     "40.3026 N",
		},
		longitude: {
			firstMsg:  "Enter longitude (E-W) of interest in decimal format: ",
			secondMsg: "Longitude of interest is: ",
			input:     "71.5112 W",
		},
		elevation: {
			firstMsg:  "Enter the elevation of interst in feet: ",
			secondMsg: "Elevation is: ",
			input:     "130",
		},
		planet: {
			firstMsg:  "Enter the " + x + " of interst: ",
			secondMsg: "The " + x + " is: ",
			input:     "AO-109",
		},
	}
}

func (p *Profiles) listItems(reader *bufio.Reader) {
	pp := *p
	if p != nil {
		fmt.Println("the current data is:")
		for i, item := range profileList {
			fmt.Printf("%d. %s: %s\n", i+1, item, pp[item].input)
		}
	}
}

//GetInput for updating the profiles
func (p Profiles) GetInput(reader *bufio.Reader) *Profiles {
	var c = false
	eol := "\n"
	for {
		p.listItems(reader)
		fmt.Println("Enter the number of item you want to change")
		fmt.Println("Enter c to continue, q to quit")
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
		case "10":
			p[planet] = getNumItem(p[planet], reader)
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

func (p Profiles) PackageInput() []string {
	pack := []string{
		p[year].input,
		p[month].input,
		p[day].input,
		p[hour].input,
		p[minute].input,
		p[second].input,
		p[latitude].input,
		p[longitude].input,
		p[elevation].input,
		p[planet].input,
	}
	return pack
}

func check(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		// os.Exit(2)
	}
}

func prM(m skymath.Euler, h string) {
	fmt.Printf("\n%s:\n", h)
	for _, item := range m {
		fmt.Printf("%f    %f    %f\n", item[0], item[1], item[2])
	}
}
