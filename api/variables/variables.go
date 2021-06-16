package variables

import "net/http"

type config struct {
  ApiKey            string `env:"API_KEY,notEmpty"`
  ApiSecret         string `env:"API_SECRET,notEmpty"`
  SmsGlobalHost     string `env:"SMS_GLOBAL_HOST" envDefault:"https://api.smsglobal.com"`
  SmsGlobalEndpoint string `env:"SMS_GLOBAL_ENDPOINT" envDefault:"/v2/sms/"`
}

var Cfg config
var Client = &http.Client{}
