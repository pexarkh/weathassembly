package main

import (
	"log"
	"regexp"
	"syscall/js"

	"github.com/pexarkh/weathassembly/rest"
)

func main() {
	fofWssNavToggleOnClick := js.FuncOf(WssNavToggleOnClick)
	js.Global().Set("WssNavToggleOnClick", fofWssNavToggleOnClick)

	cb0 := js.FuncOf(WssSearchOnSubmit)
	wssSearch := js.Global().Get("document").Call("getElementById", "WssSearch")
	wssSearch.Set("onsubmit", cb0)
	wssSearch.Set("onclick", cb0)

	cb1 := js.FuncOf(WssCurrentLocation)
	wssLoc := js.Global().Get("document").Call("getElementById", "WssCurrentLocation")
	wssLoc.Set("onsubmit", cb1)
	wssLoc.Set("onclick", cb1)

	forever := make(chan bool)
	retrieveServerRendered("", "")
	<-forever
}

func retrieveServerRendered(zip, ip string) {
	html, err := rest.GetUrl("http://localhost:3131/?zip=" + zip + "&ip=" + ip)
	if err != nil {
		log.Printf("%+v", err)
	} else {
		js.Global().Get("document").Call("getElementById", "root").Set("innerHTML", html)
	}
}

func WssCurrentLocation(this js.Value, args []js.Value) interface{} {
	go func() {
		retrieveServerRendered("", "")
	}()
	f := js.Global().Get("document").Call("querySelector", "#WssNav")
	f.Get("classList").Call("toggle", "is-open")
	return nil
}

func WssSearchOnSubmit(this js.Value, args []js.Value) interface{} {
	wssZip := js.Global().Get("document").Call("getElementById", "WssZip")
	if zip := wssZip.Get("value").String(); isValidZip(zip) {
		go func() {
			retrieveServerRendered(zip, "")
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

func isValidZip(zip string) bool {
	validZip := regexp.MustCompile(`^[0-9][0-9][0-9][0-9][0-9]$`)
	return validZip.MatchString(zip)
}
