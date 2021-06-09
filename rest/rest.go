package rest

import (
  "errors"
  "fmt"
  "io/ioutil"
  "net/http"
  "strconv"
  "time"
)

func GetUrl(url string) (string, error) {
  var hclient = &http.Client{
    Timeout: time.Second * 10,
  }

  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    fmt.Printf("NewRequest: %s\n", err)
    return "", err
  }
  //req.Header.Add("js.fetch:mode", "no-cors")

  resp, err := hclient.Do(req)
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
