package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Saied74/skyham/pkg/dataops"
)

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
			input:     "10",
		},
		minute: {
			firstMsg:  "Enter the minute of interest: ",
			secondMsg: "The minute of interest is: ",
			input:     "30",
		},
		second: {
			firstMsg:  "Don't enter the second of interest: ",
			secondMsg: "The second of interest is: ",
			input:     "0",
		},
		latitude: {
			firstMsg:  "Enter latitude (N-S) of interest (see example): ",
			secondMsg: "Latitude of interest is: ",
			input:     "34 55 42 S",
		},
		longitude: {
			firstMsg:  "Enter longitude (E-W) of interest (see exmaple): ",
			secondMsg: "Longitude of interest is: ",
			input:     "138 36 3 E",
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

func (p profiles) packageInput() []string {
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
	}
	return pack
}

func check(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(2)
	}
}
