package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Data struct {
	Status `json:"status"`
}

type StatusBencana struct {
	StatusWater string
	StatusWind  string
}

var status = StatusBencana{}

func updateData() {
	for {
		var data = Data{Status: Status{}}
		min, max := 1, 30

		data.Status.Water = rand.Intn(max-min) + min
		data.Status.Wind = rand.Intn(max-min) + min

		b, err := json.MarshalIndent(&data, "", "")

		if err != nil {
			log.Fatalln("error while marshalling json data =>", err.Error())
		}
		err = os.WriteFile("data.json", b, 0644)

		if err != nil {
			log.Fatalln("error while writing value to data.json file")
		}
		if data.Status.Water <= 5 {
			status.StatusWater = "Aman"

		} else if data.Status.Water >= 6 && data.Status.Water <= 8 {
			status.StatusWater = "Siaga"
		} else {
			status.StatusWater = "Bahaya"
		}

		if data.Status.Wind <= 6 {
			status.StatusWind = "Aman"
		} else if data.Status.Wind >= 7 && data.Status.Wind <= 15 {
			status.StatusWind = "Siaga"
		} else {
			status.StatusWind = "Bahaya"
		}

		fmt.Println("menunggu 15 detik")
		time.Sleep(time.Second * 15)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	go updateData()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tl, _ := template.ParseFiles("index.html")

		var data = Data{Status: Status{}}

		b, err := os.ReadFile("data.json")

		if err != nil {
			fmt.Fprint(w, "Error cannot read file ")
			return
		}

		err = json.Unmarshal(b, &data)

		if err != nil {
			fmt.Fprintf(w, "unmarshal parses error")
			return
		}

		err = tl.ExecuteTemplate(w, "index.html", status)

		if err != nil {
			fmt.Fprint(w, "Cannot Execute to template")
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
