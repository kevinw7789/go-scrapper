package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var halt bool = false

type FlightData struct {
	num, date, origin, destination, reg, actype, std, atd, sta, ata string
}

func trimAirport(text string) string {
	stext := strings.Split(text, " ")

	if len(stext) > 1 {
		var ftext string = strings.Replace(stext[1], "(", "", -1)
		var fftext string = strings.Replace(ftext, ")", "", -1)

		return fftext
	}

	return "-----"
}

func main() {
	var Flights []FlightData

	file, err := os.Create("../output/flights.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close()

	flightNum := 1746
	writer := csv.NewWriter(file)

	//retrieve stuffs
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("GET: ", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("ERR: ", err)
		halt = true
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("RES: ", r.StatusCode)
	})
	c.OnHTML("table#tbl-datatable > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			/*
				fmt.Println(el.ChildText("td:nth-child(3)"))  //date
				fmt.Println(el.ChildText("td:nth-child(4)"))  //origin
				fmt.Println(el.ChildText("td:nth-child(5)"))  //destination
				fmt.Println(el.ChildText("td:nth-child(6)"))  //AC
				fmt.Println(el.ChildText("td:nth-child(8)"))  //STD
				fmt.Println(el.ChildText("td:nth-child(9)"))  //ATD
				fmt.Println(el.ChildText("td:nth-child(10)")) //STA
				fmt.Println(el.ChildText("td:nth-child(12)")) //ATA
			*/

			date := strings.Split(el.ChildText("td:nth-child(3)"), " ")
			day, err := strconv.Atoi(date[0])

			ac1 := strings.Replace(el.ChildText("td:nth-child(6)"), "(", "", -1)
			ac2 := strings.Replace(ac1, ")", "", -1)

			arrt := strings.Replace(el.ChildText("td:nth-child(12)"), "Landed ", "", -1)

			if err != nil {
				println("Could not convert to int")
				defer writer.Flush()
			}

			//format to struc
			//set specific data conditions
			if date[1] == "Sep" && (day == 17 || day == 16) {
				newac := strings.Split(ac2, " ")
				data := FlightData{
					num:         strconv.Itoa(flightNum),
					date:        el.ChildText("td:nth-child(3)"),
					reg:         newac[1],
					actype:      newac[0],
					origin:      trimAirport(el.ChildText("td:nth-child(4)")),
					destination: trimAirport(el.ChildText("td:nth-child(5)")),
					std:         el.ChildText("td:nth-child(8)"),
					atd:         el.ChildText("td:nth-child(9)"),
					sta:         el.ChildText("td:nth-child(10)"),
					ata:         arrt,
				}

				Flights = append(Flights, data)
				//fmt.Println(data)
			}

			//fmt.Println(table)
		})
	})

	min := 1
	max := 10

	//number of entries to check
	for i := 0; i < 5; {
		URL := "https://" + strconv.Itoa(flightNum)
		c.Visit(URL)

		i = i + rand.Intn(max-min) + min
		flightNum = flightNum + i

		if halt {
			time.Sleep(10 * time.Second)

			halt = false
		}
		time.Sleep(1 * time.Second)
	}

	for _, data := range Flights {
		record := []string{
			data.num,
			data.reg,
			data.actype,
			data.date,
			data.origin,
			data.destination,
			data.std,
			data.atd,
			data.sta,
			data.ata,
		}
		writer.Write(record)
	}

	defer writer.Flush()
}
