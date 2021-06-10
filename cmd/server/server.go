package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/pexarkh/weathassembly/consts"
	"github.com/pexarkh/weathassembly/pongoutils"
	"github.com/pexarkh/weathassembly/weather"
)

func main() {
	port := flag.String("p", "3131", "listen port")
	flag.Parse()
	listen := "localhost:" + *port
	log.Printf("http://%s", listen)

	_ = pongo2.RegisterFilter("TopForecast", pongoutils.TopForecast)
	_ = pongo2.RegisterFilter("TopTemperature", pongoutils.TopTemperature)
	_ = pongo2.RegisterFilter("TopStartTime", pongoutils.TopStartTime)

	http.HandleFunc("/", wasmHandler)

	log.Fatal(http.ListenAndServe(listen, nil))
}
func wasmHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if strings.HasPrefix(r.URL.Path, "/api") {
		zips, zipok := r.URL.Query()["zip"]
		zipok = zipok && len(strings.TrimSpace(zips[0])) >= 1
		ips, ipok := r.URL.Query()["ip"]
		if !zipok {
			http.Error(w,"bad or missing zip code",http.StatusBadRequest)
		} else {
			zip := strings.TrimSpace(zips[0])
			var ip string
			if ipok {
				ip = strings.TrimSpace(ips[0])
			}
			io.WriteString(w, fmt.Sprintf("[%s] [%s]", zip, ip))

			out, err := retrieveAndRenderServer(zip, ip)
			if err != nil {
				http.Error(w,fmt.Sprintf("%v", err),http.StatusInternalServerError)
			} else {
				io.WriteString(w, out)
			}
		}
	} else {
		if strings.HasSuffix(r.URL.Path, ".wasm") {
			w.Header().Set("content-type", "application/wasm")
		}
		http.FileServer(http.Dir(".")).ServeHTTP(w, r)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func retrieveAndRenderServer(zip, ip string) (string, error) {
	pfc, err := weather.RetrieveForecastForPlace(zip, ip)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return "", err
	}
	log.Printf("%+v\n", pfc)
	html, err := weather.Render(consts.ForecastFragment, pfc)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return "", err
	}
	log.Printf("%s\n", html)
	return html, nil
}
