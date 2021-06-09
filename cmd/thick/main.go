package main

import (
	"log"
	"syscall/js"

	"github.com/flosch/pongo2"
	"github.com/pexarkh/weathassembly/consts"
	"github.com/pexarkh/weathassembly/pongoutils"
	"github.com/pexarkh/weathassembly/weather"
)

func main() {
	js.Global().Set("WssNavToggleOnClick", js.FuncOf(WssNavToggleOnClick))
	//cb0 := js.NewEventCallback(js.PreventDefault, WssSearchOnSubmit)
	cb0 := js.FuncOf(WssSearchOnSubmit)
	wssSearch := js.Global().Get("document").Call("getElementById", "WssSearch")
	wssSearch.Set("onsubmit", cb0)
	wssSearch.Set("onclick", cb0)
	//cb1 := js.NewEventCallback(js.PreventDefault, WssCurrentLocation)
	cb1 := js.FuncOf(WssCurrentLocation)
	wssLoc := js.Global().Get("document").Call("getElementById", "WssCurrentLocation")
	wssLoc.Set("onsubmit", cb1)
	wssLoc.Set("onclick", cb1)

	_ = pongo2.RegisterFilter("TopForecast", pongoutils.TopForecast)
	_ = pongo2.RegisterFilter("TopTemperature", pongoutils.TopTemperature)
	_ = pongo2.RegisterFilter("TopStartTime", pongoutils.TopStartTime)

	forever := make(chan bool)
	retrieveAndRenderWeb("10001", "165.225.39.62")
	<-forever
}

func retrieveAndRenderWeb(zip, ip string) {
	pfc, err := weather.RetrieveForecastForPlace(zip, ip)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}
	htmlView, err := weather.Render(consts.ForecastFragment, pfc)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}
	js.Global().Get("document").Call("getElementById", "root").Set("innerHTML", htmlView)
}

func WssCurrentLocation(this js.Value, args []js.Value) interface{} {
	log.Printf("WssCurrentLocation: this: %v", this)
	log.Printf("WssCurrentLocation: args: %v", args)
	go func() {
		retrieveAndRenderWeb("02135", "165.225.39.62")
	}()
	f := js.Global().Get("document").Call("querySelector", "#WssNav")
	f.Get("classList").Call("toggle", "is-open")
	return nil
}

func WssSearchOnSubmit(this js.Value, args []js.Value) interface{} {
	log.Printf("WssSearchOnSubmit: this: %v", this)
	log.Printf("WssSearchOnSubmit: args: %v", args)
	wssZip := js.Global().Get("document").Call("getElementById", "WssZip")
	if zip := wssZip.Get("value").String(); weather.IsValidZip(zip) {
		go func() {
			retrieveAndRenderWeb(zip, "165.225.39.62")
		}()
		f := js.Global().Get("document").Call("querySelector", "#WssNav")
		f.Get("classList").Call("toggle", "is-open")
	} else {
		js.Global().Call("alert", "ERROR: invalid zip: "+zip)
	}
	return nil
}

func WssNavToggleOnClick(this js.Value, args []js.Value) interface{} {
	f := js.Global().Get("document").Call("querySelector", "#WssNav")
	f.Get("classList").Call("toggle", "is-open")
	return nil
}
