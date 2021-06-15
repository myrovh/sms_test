package main

import (
  "bytes"
  "crypto/hmac"
  "crypto/sha256"
  "encoding/base64"
  "flag"
  "fmt"
  "github.com/caarlos0/env/v6"
  "github.com/justinas/alice"
  "github.com/rs/zerolog"
  "github.com/rs/zerolog/hlog"
  "github.com/rs/zerolog/log"
  "math/rand"
  "net/http"
  "net/http/httputil"
  "os"
  "time"
)

type config struct {
  ApiKey    string `env:"API_KEY,notEmpty"`
  ApiSecret string `env:"API_SECRET,notEmpty"`
}

var cfg config

func main() {
  // parse flags
  debug := flag.Bool("debug", false, "sets log level to debug")
  port := flag.Int("port", 8080, "sets port to listen on")
  human := flag.Bool("human", false, "if true logger outputs in human readable format")
  flag.Parse()

  if err := env.Parse(&cfg); err != nil {
    log.Fatal().Err(err).Msg("couldn't parse envs")
  }

  // configure zerolog
  zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
  zerolog.SetGlobalLevel(zerolog.InfoLevel)
  if *debug {
    zerolog.SetGlobalLevel(zerolog.DebugLevel)
  }
  if *human {
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
  }

  // configure middleware
  zerologMiddleware := hlog.NewHandler(log.Logger)
  accessLogMiddleware := hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
    hlog.FromRequest(r).Info().
      Str("method", r.Method).
      Stringer("url", r.URL).
      Int("status", status).
      Dur("duration", duration).
      Msg("http access")
  })
  remoteAddrHandlerMiddleware := hlog.RemoteAddrHandler("ip")
  userAgentMiddleware := hlog.UserAgentHandler("user_agent")
  refererMiddleware := hlog.RefererHandler("referer")
  requestIDMiddleware := hlog.RequestIDHandler("req_id", "Request-Id")

  chain := alice.New(zerologMiddleware, accessLogMiddleware, remoteAddrHandlerMiddleware, userAgentMiddleware, refererMiddleware, requestIDMiddleware)

  // start http server
  log.Info().Int("port", *port).Msg("listening")

  http.Handle("/check", chain.Then(http.HandlerFunc(healthCheck)))
  http.Handle("/api/message", chain.Then(http.HandlerFunc(doSms)))

  if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
    log.Fatal().Err(err).Msg("error on http listener")
  }
}

// httpAuth adds a Authorization header to the provided http request using the oauth http mac specification
// this should be working but get failed to auth responses from sms global api
// possibly:
// - hmac hash + conversion to base64 is wrong
func httpAuth(r *http.Request) (err error) {
  timestamp := time.Now().Unix()
  nonce := rand.Int31()
  macString := fmt.Sprintf("%d\n%d\n%s\n%s\n%s\n433\n\n", timestamp, nonce, r.Method, r.URL.Path, r.URL.Host)
  fmt.Printf("%s", macString)

  h := hmac.New(sha256.New, []byte(cfg.ApiSecret))
  h.Write([]byte(macString))
  mac := base64.StdEncoding.EncodeToString(h.Sum(nil))

  r.Header.Add("Authorization", fmt.Sprintf("MAC id=\"%s\", ts=\"%d\", nonce=\"%d\", mac=\"%s\"", cfg.ApiKey, timestamp, nonce, mac))
  return nil
}

// doSms TODO add helper util to close out http response with error message if function fails before expected return
// test function to try and get communication between program and sms global
func doSms(w http.ResponseWriter, r *http.Request) {
  client := &http.Client{}

  req, err := http.NewRequest("GET", "https://api.smsglobal.com/v2/sms/", bytes.NewReader(make([]byte, 100)))
  if err != nil {
    log.Error().Err(err).Msg("failure to generate request")
    return
  }

  if err := httpAuth(req); err != nil {
    log.Error().Err(err).Msg("failure to generate oauth header")
    return
  }

  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("Accept", "application/json")

  body, err := httputil.DumpRequest(req, true)
  if err != nil {
    log.Error().Err(err).Msg("http request dump")
  }

  fmt.Printf("%s", body)

  resp, err := client.Do(req)
  if err != nil {
    log.Error().Err(err).Msg("http get")
  }
  defer resp.Body.Close()

  body, err = httputil.DumpResponse(resp, true)
  if err != nil {
    log.Error().Err(err).Msg("http response dump")
  }

  fmt.Printf("%s", body)

}

// healthCheck writes a plain text response. Simply used to determine if the
// server is online
func healthCheck(w http.ResponseWriter, r *http.Request) {
  if _, err := fmt.Fprintln(w, "server is online"); err != nil {
    log.Error().Err(err).Msg("attempting to write 'healthCheck' response")
  }
}
