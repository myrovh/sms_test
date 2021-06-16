package endpoints

import (
  "fmt"
  "github.com/rs/zerolog/log"
  "net/http"
)

// HealthCheck writes a plain text response. Simply used to determine if the
// server is online
func HealthCheck(w http.ResponseWriter, r *http.Request) {
  if _, err := fmt.Fprintln(w, "server is online"); err != nil {
    log.Error().Err(err).Msg("attempting to write 'healthCheck' response")
  }
}
