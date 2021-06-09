package weather

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/buger/jsonparser"
	"github.com/flosch/pongo2"
	"github.com/pexarkh/weathassembly/rest"
)

type PlaceWeatherForecasts struct {
	Place
	Forecasts []WeatherForecast
}

type Place struct {
	Ip                string
	Latitude          string
	Longitude         string
	PlaceName         string
	State             string
	StateAbbreviation string
	Zip               string
	ImageUrl          string
}

type WeatherForecast struct {
	Number           int64
	Name             string
	StartTime        string
	EndTime          string
	IsDaytime        bool
	Temperature      int64
	TemperatureUnit  string
	WindSpeed        string
	WindDirection    string
	ShortForecast    string
	DetailedForecast string
	SvgSymbolId      string
}

func RetrieveForecastForPlace(zip, ip string) (*PlaceWeatherForecasts, error) {
	if zip != "" && !validZip.MatchString(zip) {
		return nil, errors.New(fmt.Sprintf("bad zip format: `%s`", zip))
	}
	if ip != "" && !validIp.MatchString(ip) {
		return nil, errors.New(fmt.Sprintf("bad IP format: `%s`", ip))
	}
	var result = new(PlaceWeatherForecasts)
	var _zip = zip
	var _ip = ip
	result.Zip = _zip
	result.Ip = _ip
	var err error
	if _zip == "" {
		if _ip == "" {
			_ip, err = retrieveExternallyVisibleIp()
			if err != nil {
				return result, err
			}
			log.Printf("api.ipify.org: %s\n", _ip)
		}
		resp, err := retrieveLocationInfoByIP()
		if err != nil {
			// TODO: return result, err
			_zip = "36013" // fall back to Montgomery, Alabama
		} else {
			log.Printf("gd.geobytes.com: %s\n", resp)
			_zip, _ = jsonparser.GetString([]byte(resp), "geobytescityid")
		}
	}
	resp, err := retrieveZippopotam(_zip)
	if err != nil {
		log.Printf("api.zippopotam.us failed: %v\n", err)
		resp, err = retrieveZippopotam("10001") // well-known zip
		if err != nil {
			log.Printf("api.zippopotam.us failed for 10001: %v\n", err)
			return result, err
		}
	}
	log.Printf("api.zippopotam.us: %s\n", resp)
	place, _ := jsonparser.GetString([]byte(resp), "places", "[0]", "place name")
	longitude, _ := jsonparser.GetString([]byte(resp), "places", "[0]", "longitude")
	latitude, _ := jsonparser.GetString([]byte(resp), "places", "[0]", "latitude")
	state, _ := jsonparser.GetString([]byte(resp), "places", "[0]", "state")
	stateAbbreviation, _ := jsonparser.GetString([]byte(resp), "places", "[0]", "state abbreviation")
	result.Ip = _ip
	result.Latitude = latitude
	result.Longitude = longitude
	result.PlaceName = place
	result.State = state
	result.StateAbbreviation = stateAbbreviation
	result.Zip = _zip
	//images
	resp, err = retrieveTeleport(latitude, longitude)
	log.Printf("api.teleport.org: %s\n", resp)
	if err == nil {
		result.ImageUrl, err = jsonparser.GetString([]byte(resp), "_embedded", "location:nearest-urban-areas", "[0]", "_embedded", "location:nearest-urban-area", "_embedded", "ua:images", "photos", "[0]", "image", "web")
	}
	// forecast
	resp, err = retrieveWeatherGovPoints(latitude, longitude)
	if err != nil {
		return result, err
	}
	log.Printf("api.weather.gov: %s\n", resp)
	forecastUrl, _ := jsonparser.GetString([]byte(resp), "properties", "forecast")
	forecast, err := rest.GetUrl(forecastUrl)
	if err != nil {
		return result, err
	}
	log.Printf("%s: %s\n", forecastUrl, forecast)
	_, err = jsonparser.ArrayEach([]byte(forecast), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var fc = new(WeatherForecast)
		fc.Number, _ = jsonparser.GetInt(value, "number")
		fc.Name, _ = jsonparser.GetString(value, "name")
		fc.StartTime, _ = jsonparser.GetString(value, "startTime")
		fc.EndTime, _ = jsonparser.GetString(value, "endTime")
		fc.IsDaytime, _ = jsonparser.GetBoolean(value, "isDaytime")
		fc.Temperature, _ = jsonparser.GetInt(value, "temperature")
		fc.TemperatureUnit, _ = jsonparser.GetString(value, "temperatureUnit")
		fc.WindSpeed, _ = jsonparser.GetString(value, "windSpeed")
		fc.WindDirection, _ = jsonparser.GetString(value, "windDirection")
		fc.ShortForecast, _ = jsonparser.GetString(value, "shortForecast")
		fc.DetailedForecast, _ = jsonparser.GetString(value, "detailedForecast")
		icon, _ := jsonparser.GetString(value, "icon")
		fc.SvgSymbolId = iconToSvgSymbolId(icon)
		result.Forecasts = append(result.Forecasts, *fc)
	}, "properties", "periods")
	return result, nil
}

func Render(template string, pfc *PlaceWeatherForecasts) (string, error) {
	tpl, err := pongo2.FromString(template)
	if err != nil {
		return "", err
	}
	out, err := tpl.Execute(pongo2.Context{"pfc": pfc})
	if err != nil {
		return "", err
	}
	return out, nil
}

var (
	validZip = regexp.MustCompile(`^[0-9][0-9][0-9][0-9][0-9]$`)
	validIp  = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)

	icon       = regexp.MustCompile(`icons/land/(?:day|night)/(\w+)`)
	svgSymbols = map[string]string{
		"bkn":             "cloudy",
		"blizzard":        "snowy",
		"cold":            "snowy",
		"dust":            "tornado",
		"few":             "cloudy",
		"fog":             "sunny",
		"fzra":            "snowy",
		"haze":            "sunny",
		"hot":             "sunny",
		"hurricane":       "tornado",
		"ovc":             "cloudy",
		"rain":            "rainy",
		"rain_fzra":       "snowy",
		"rain_showers":    "rainy",
		"rain_showers_hi": "rainy",
		"rain_snow":       "snowy",
		"rain_sleet":      "snowy",
		"sct":             "cloudy",
		"skc":             "sunny",
		"sleet":           "snowy",
		"smoke":           "tornado",
		"snow":            "snowy",
		"snow_fzra":       "snowy",
		"snow_sleet":      "snowy",
		"tornado":         "tornado",
		"tropical_storm":  "stormy",
		"tsra":            "stormy",
		"tsra_sct":        "stormy",
		"tsra_hi":         "stormy",
		"wind_bkn":        "cloudy",
		"wind_ovc":        "cloudy",
		"wind_sct":        "cloudy",
		"wind_skc":        "sunny",
		"wind_few":        "cloudy",
	}
)

func IsValidZip(zip string) bool {
	return validZip.MatchString(zip)
}

func iconToSvgSymbolId(iconUrl string) string {
	res := icon.FindAllStringSubmatch(iconUrl, -1)
	if res == nil {
		return "sunny"
	}
	k, ok := svgSymbols[res[0][1]]
	if ok {
		return k
	} else {
		return "sunny"
	}
}

func retrieveExternallyVisibleIp() (string, error) {
	const url = "http://ipv4bot.whatismyipaddress.com" // "https://api.ipify.org"
	return rest.GetUrl(url)
}

// geo loc from zip
func retrieveZippopotam(zip string) (string, error) {
	const zippopotamUrl = "http://api.zippopotam.us/us/"
	return rest.GetUrl(zippopotamUrl + zip)
}

// geo loc from ip
func retrieveLocationInfoByIP() (string, error) {
	const url = "http://gd.geobytes.com/GetCityDetails"
	return rest.GetUrl(url)
}

// images
func retrieveTeleport(latitude, longitude string) (string, error) {
	const teleportUrl = "https://api.teleport.org/api/locations/"
	return rest.GetUrl(teleportUrl + latitude + "," + longitude + "/?embed=location:nearest-urban-areas/location:nearest-urban-area/ua:images")
}

func retrieveWeatherGovPoints(latitude, longitude string) (string, error) {
	const weatherGovPointsUrl = "https://api.weather.gov/points/"
	return rest.GetUrl(weatherGovPointsUrl + latitude + "," + longitude)
}
