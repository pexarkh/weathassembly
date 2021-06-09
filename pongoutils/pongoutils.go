package pongoutils

import (
  "fmt"
  "time"

  "github.com/flosch/pongo2"
  "github.com/pexarkh/weathassembly/weather"
)

func TopForecast(in *pongo2.Value, params *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
  pfc := in.Interface().(*weather.PlaceWeatherForecasts)
  if len(pfc.Forecasts) > 0 {
    return pongo2.AsSafeValue(pfc.Forecasts[0].DetailedForecast), nil
  } else {
    return pongo2.AsSafeValue(""), nil
  }
}

func TopTemperature(in *pongo2.Value, params *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
  pfc := in.Interface().(*weather.PlaceWeatherForecasts)
  if len(pfc.Forecasts) > 0 {
    return pongo2.AsSafeValue(pfc.Forecasts[0].Temperature), nil
  } else {
    return pongo2.AsSafeValue(""), nil
  }
}

func TopStartTime(in *pongo2.Value, params *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
  t := time.Now()
  return pongo2.AsSafeValue(fmt.Sprintf("%s, %s", t.Weekday().String(), t.Format(time.Kitchen))), nil
}
