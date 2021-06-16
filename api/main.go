package main

import (
  "flag"
  "fmt"
  "github.com/caarlos0/env/v6"
  "github.com/justinas/alice"
  "github.com/myrovh/sms-sender/api/endpoints"
  "github.com/myrovh/sms-sender/api/variables"
  "github.com/rs/cors"
  "github.com/rs/zerolog"
  "github.com/rs/zerolog/hlog"
  "github.com/rs/zerolog/log"
  "net/http"
  "os"
  "time"
)

func main() {
  // parse flags
  debug := flag.Bool("debug", false, "sets log level to debug")
  port := flag.Int("port", 8080, "sets port to listen on")
  human := flag.Bool("human", false, "if true logger outputs in human readable format")
  flag.Parse()

  if err := env.Parse(&variables.Cfg); err != nil {
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
  corsMiddleware := cors.New(cors.Options{
    AllowedOrigins: variables.Cfg.AllowedOrigins,
  }).Handler

  chain := alice.New(zerologMiddleware, accessLogMiddleware, remoteAddrHandlerMiddleware, userAgentMiddleware, refererMiddleware, requestIDMiddleware, corsMiddleware)

  // start http server
  log.Info().Int("port", *port).Msg("listening")

  http.Handle("/check", chain.Then(http.HandlerFunc(endpoints.HealthCheck)))
  http.Handle("/api/message", chain.Then(http.HandlerFunc(endpoints.SmsEndpointHandler)))

  if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
    log.Fatal().Err(err).Msg("error on http listener")
  }
}
