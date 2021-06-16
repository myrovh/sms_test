package endpoints

import (
  "crypto/hmac"
  "crypto/rand"
  "crypto/sha256"
  "encoding/base64"
  "fmt"
  "github.com/myrovh/sms-sender/api/variables"
  "github.com/rs/zerolog/log"
  "math/big"
  "net/http"
  "time"
)

// httpAuth adds a Authorization header to the provided http request using the oauth http mac specification
func httpAuth(r *http.Request) (err error) {
  timestamp := time.Now().Unix()
  nonce, err := rand.Int(rand.Reader, big.NewInt(9999999999999999))
  if err != nil {
    return err
  }
  macString := fmt.Sprintf(`%d
%d
%s
%s?%s
%s
443

`, timestamp, nonce, r.Method, r.URL.Path, r.URL.RawQuery, r.URL.Host)

  h := hmac.New(sha256.New, []byte(variables.Cfg.ApiSecret))
  h.Write([]byte(macString))
  mac := base64.StdEncoding.EncodeToString(h.Sum(nil))

  r.Header.Add("Authorization", fmt.Sprintf("MAC id=\"%s\", ts=\"%d\", nonce=\"%d\", mac=\"%s\"", variables.Cfg.ApiKey, timestamp, nonce, mac))
  return nil
}

// genericErrorWriter tiny wrapper to quickly write a given string as the response body and set the status code
func genericErrorWriter(w http.ResponseWriter, errorMessage string, statusCode int) {
  w.WriteHeader(statusCode)
  if _, err := fmt.Fprintln(w, errorMessage); err != nil {
    log.Error().Err(err).Str("error_message", errorMessage).Msg("attempting to write response")
  }
}