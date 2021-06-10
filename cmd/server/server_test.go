package main

import (
  "net/http"
  "net/http/httptest"
  "testing"
)

func mkGetZipIpApiRequest(zip, ip string) *http.Request {
  req, _ := http.NewRequest("GET", "https://wss.com/api", nil)
  q := req.URL.Query()
  if zip != "<NOTSET>" {
    q.Add("zip", zip)
  }
  if ip != "<NOTSET>" {
    q.Add("ip", ip)
  }
  req.URL.RawQuery = q.Encode()
  return req
}

func Test_wasmHandler(t *testing.T) {
  tests := []struct {
    name          string
    req           *http.Request
    expHttpStatus int
  }{
    {name: "zip: 02135; ip: ''", req: mkGetZipIpApiRequest("02135", ""), expHttpStatus: http.StatusOK},
    {name: "zip: 02135; ip: '<NOTSET>'", req: mkGetZipIpApiRequest("02135", "<NOTSET>"), expHttpStatus: http.StatusOK},
    {name: "zip: 02135; ip: '127.0.0.1'", req: mkGetZipIpApiRequest("02135", "127.0.0.1"), expHttpStatus: http.StatusOK},
    {name: "zip: <NOTSET>; ip: ''", req: mkGetZipIpApiRequest("<NOTSET>", ""), expHttpStatus: http.StatusBadRequest},
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      rec := httptest.NewRecorder()
      handler := http.HandlerFunc(wasmHandler)
      handler.ServeHTTP(rec, tt.req)
      if status := rec.Code; status != tt.expHttpStatus {
        t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expHttpStatus)
      }
    })
  }
}
