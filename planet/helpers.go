package main

// func printaid(basedata map[string]dataops.BaseItem) {
// 	for _, x := range basedata {
// 		fmt.Println("Name: ", x.Name)
// 		fmt.Println("Value: ", x.Value)
// 		fmt.Println("Numonic: ", x.Numonic)
// 		fmt.Println("Description: ", x.Description)
// 		fmt.Println()
// 	}
// }

//
// func getProfile() *profiles {
// 	return &profiles{
// 		year: profile{
// 			firstMsg:  "Enter the year of interest: ",
// 			secondMsg: "The year of interest is: ",
// 			input:     "2014",
// 		},
// 		month: profile{
// 			firstMsg:  "Enter the number of the month of interest: ",
// 			secondMsg: "The month is: ",
// 			input:     "3",
// 		},
// 		day: profile{
// 			firstMsg:  "Enter the number of the day of interest: ",
// 			secondMsg: "The day is: ",
// 			input:     "22",
// 		},
// 		hour: {
// 			firstMsg:  "Enter the hour of interest (24 hour format): ",
// 			secondMsg: "The hour of interest is: ",
// 			input:     "10",
// 		},
// 		minute: {
// 			firstMsg:  "Enter the minute of interest: ",
// 			secondMsg: "The minute of interest is: ",
// 			input:     "30",
// 		},
// 		second: {
// 			firstMsg:  "Don't enter the second of interest: ",
// 			secondMsg: "The second of interest is: ",
// 			input:     "0",
// 		},
// 		latitude: {
// 			firstMsg:  "Enter latitude (N-S) of interest in decimal format: ",
// 			secondMsg: "Latitude of interest is: ",
// 			input:     "34.9285 S",
// 		},
// 		longitude: {
// 			firstMsg:  "Enter longitude (E-W) of interest in decimal format: ",
// 			secondMsg: "Longitude of interest is: ",
// 			input:     "138.6007 E",
// 		},
// 		elevation: {
// 			firstMsg:  "Enter the elevation of interst in feet: ",
// 			secondMsg: "Elevation is: ",
// 			input:     "130",
// 		},
// 		planet: {
// 			firstMsg:  "Enter the planet of interst: ",
// 			secondMsg: "The planet is: ",
// 			input:     "Jupiter",
// 		},
// 	}
// }
//
// func (p *profiles) listItems(reader *bufio.Reader) {
// 	pp := *p
// 	if p != nil {
// 		fmt.Println("the current data is:")
// 		for i, item := range profileList {
// 			fmt.Printf("%d. %s: %s\n", i+1, item, pp[item].input)
// 		}
// 	}
// }
//
// //getInput for updating the profiles
// func (p profiles) getInput(reader *bufio.Reader) *profiles {
// 	var c = false
// 	eol := "\n"
// 	for {
// 		p.listItems(reader)
// 		fmt.Println("Enter the number of item you want to change")
// 		fmt.Println("Enter c to continue, q to quit")
// 		input, _ := reader.ReadString('\n')
// 		input = strings.TrimSuffix(input, eol)
// 		switch input {
// 		case "1":
// 			p[year] = getNumItem(p[year], reader)
// 		case "2":
// 			p[month] = getNumItem(p[month], reader)
// 		case "3":
// 			p[day] = getNumItem(p[day], reader)
// 		case "4":
// 			p[hour] = getNumItem(p[hour], reader)
// 		case "5":
// 			p[minute] = getNumItem(p[minute], reader)
// 		case "6":
// 			p[second] = getNumItem(p[minute], reader)
// 		case "7":
// 			p[latitude] = getNumItem(p[latitude], reader)
// 		case "8":
// 			p[longitude] = getNumItem(p[longitude], reader)
// 		case "9":
// 			p[elevation] = getNumItem(p[elevation], reader)
// 		case "10":
// 			p[planet] = getNumItem(p[planet], reader)
// 		case "c":
// 			c = true
// 		case "C":
// 			c = true
// 		case "q":
// 			os.Exit(1)
// 		default:
// 			fmt.Println("You made a mistake, try again")
// 		}
// 		if c {
// 			break
// 		}
// 	}
// 	return &p
// }
//
// //getNumItem for updating individual items
// func getNumItem(p profile, reader *bufio.Reader) profile {
// 	fmt.Println(p.firstMsg)
// 	input, _ := reader.ReadString('\n')
// 	input = strings.TrimSuffix(input, "\n")
// 	p.input = input
// 	return p
// }
//
// func (p profiles) packageInput() []string {
// 	pack := []string{
// 		p[year].input,
// 		p[month].input,
// 		p[day].input,
// 		p[hour].input,
// 		p[minute].input,
// 		p[second].input,
// 		p[latitude].input,
// 		p[longitude].input,
// 		p[elevation].input,
// 		p[planet].input,
// 	}
// 	return pack
// }
//
// func check(msg string, err error) {
// 	if err != nil {
// 		fmt.Println(msg, err)
// 		// os.Exit(2)
// 	}
// }

// func prM(m skymath.Euler, h string) {
// 	fmt.Printf("\n%s:\n", h)
// 	for _, item := range m {
// 		fmt.Printf("%f    %f    %f\n", item[0], item[1], item[2])
// 	}
// }
