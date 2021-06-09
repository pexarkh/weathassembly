package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pexarkh/weathassembly/weather"
)

func main() {
	var zip, ip string
	flag.StringVar(&zip, "z", "", "Zip of a place")
	flag.StringVar(&ip, "i", "", "IP address")
	flag.Parse()

	err := retrieveAndRender(zip, ip)
	if err != nil {
		log.Fatal(err)
	}
	zip = "-"
	for zip != "" {
		fmt.Print("\nZip: ")
		reader := bufio.NewReader(os.Stdin)
		zip, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		zip = strings.Replace(zip, "\n", "", -1)
		if zip != "" {
			err := retrieveAndRender(zip, "")
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func retrieveAndRender(zip, ip string) error {
	pfc, err := weather.RetrieveForecastForPlace(zip, ip)
	if err != nil {
		return err
	}
	log.Printf("%+v\n", pfc) // #

	consoleView, err := weather.Render(consolePage, pfc)
	if err != nil {
		return err
	}
	log.Printf("%+v\n", pfc)
	fmt.Print(consoleView)
	return nil
}

var consolePage = onBlue("Place:") + " {{ pfc.PlaceName }}, {{ pfc.StateAbbreviation }} {{ pfc.Zip }} " +
	onBlue("Lat/Long:") + " {{ pfc.Latitude }} / {{ pfc.Longitude }} " +
	"{% if pfc.Ip != \"\" %}" + onBlue("Ip:") + " {{ pfc.Ip }}" + "" + "{% endif %}" + "\n" +
	"{% for fc in pfc.Forecasts %}" + onBlue("{{ fc.Name }}:") + " " + onGreen("Temp:") + " {{ fc.Temperature }}\u00B0{{ fc.TemperatureUnit }} " + onGreen("Wind:") + " {{ fc.WindSpeed }} {{ fc.WindDirection }} " +
	onGreen("Forecast:") + " {{ fc.DetailedForecast }}\n{% endfor %}"

func onBlue(s string) string {
	return "\u001B[48;5;21m" + s + "\u001B[0m"
}

func onGreen(s string) string {
	return "\u001B[48;5;22m" + s + "\u001B[0m"
}
