package endpoints

import (
  "bufio"
  "fmt"
  "github.com/myrovh/sms-sender/api/variables"
  "github.com/rs/zerolog/log"
  "io"
  "net/http"
)

// SmsEndpointHandler handles proxying GET and POSTS requests to the sms global endpoint set by Cfg.SmsGlobalEndpoint
// Adds a hmac auth header to any request
func SmsEndpointHandler(w http.ResponseWriter, r *http.Request) {
  var req *http.Request
  var err error

  switch r.Method {
  case http.MethodGet:
    // forward query string from incoming request to sms global
    req, err = http.NewRequest("GET", fmt.Sprintf("%s%s?%s", variables.Cfg.SmsGlobalHost, variables.Cfg.SmsGlobalEndpoint, r.URL.RawQuery), nil)
    if err != nil {
      log.Error().Err(err).Msg("generating get request")
      genericErrorWriter(w, "internal server error", http.StatusInternalServerError)
      return
    }
  case http.MethodPost:
    req, err = http.NewRequest("POST", fmt.Sprintf("%s%s?", variables.Cfg.SmsGlobalHost, variables.Cfg.SmsGlobalEndpoint), r.Body)
    if err != nil {
      log.Error().Err(err).Msg("generating post request")
      genericErrorWriter(w, "internal server error", http.StatusInternalServerError)
      return
    }
  default:
    genericErrorWriter(w, "request must be GET or POST", http.StatusBadRequest)
    return
  }

  if err := httpAuth(req); err != nil {
    log.Error().Err(err).Msg("generating auth header")
    genericErrorWriter(w, "internal server error", http.StatusInternalServerError)
    return
  }

  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("Accept", "application/json")

  resp, err := variables.Client.Do(req)
  if err != nil {
    log.Error().Err(err).Msg("sending sms global request")
    genericErrorWriter(w, "bad response from sms global", http.StatusInternalServerError)
    return
  }
  defer func(Body io.ReadCloser) {
    err := Body.Close()
    if err != nil {
      log.Error().Err(err).Msg("couldn't close body")
    }
  }(resp.Body)

  w.WriteHeader(resp.StatusCode)

  // write json response from sms global directly into our response
  reader := bufio.NewReader(resp.Body)
  if _, err = reader.WriteTo(w); err != nil {
    log.Error().Err(err).Msg("writing sms global body to response body")
    genericErrorWriter(w, "internal server error", http.StatusInternalServerError)
    return
  }
}
