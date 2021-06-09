package wssfts

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pexarkh/weathassembly/consts"
	"github.com/pexarkh/weathassembly/weather"
)

//----- THESE FUNCTIONS TURNED OUT TO BE UNNECESSARY -----

func getOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func PostUrlJson(url, body string) (string, error) {
	var hclient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := hclient.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", errors.New("StatusCode: " + strconv.Itoa(resp.StatusCode))
	}
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func retrieveAndRenderServer(zip, ip string) (string, error) {
	pfc, err := weather.RetrieveForecastForPlace(zip, ip)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return "", err
	}
	log.Printf("%s\n", pfc)
	html, err := weather.Render(consts.ForecastFragment, pfc)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return "", err
	}
	log.Printf("%s\n", html)
	return html, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	zips, zipok := r.URL.Query()["zip"]
	gotzip := zipok && len(zips) >= 1
	ips, ipok := r.URL.Query()["ip"]
	gotip := ipok && len(ips) >= 1
	if gotzip && gotip {
		zip := zips[0]
		ip := ips[0]
		log.Printf("zip (%d): [%s] ip (%d): [%s]", len(zip), zip, len(ip), ip)
		out, err := retrieveAndRenderServer(zip, ip)
		if err != nil {
			log.Printf("%+v", err)
			io.WriteString(w, "")
		} else {
			io.WriteString(w, out)
		}
	} else {
		http.Redirect(w, r, "index.html", 301)
	}
}
