package main

import (
  "flag"
  "fmt"
  "github.com/justinas/alice"
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

  if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
    log.Fatal().Err(err).Msg("error on http listener")
  }
}

// healthCheck writes a plain text response. Simply used to determine if the
// server is online
func healthCheck(w http.ResponseWriter, r *http.Request) {
  if _, err := fmt.Fprintln(w, "server is online"); err != nil {
    log.Error().Err(err).Msg("attempting to write 'healthCheck' response")
  }
}
