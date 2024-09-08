package agent

type Config struct {
	Adress         string `env:"ADDRESS" flag:"a"`
	ReportInterval uint   `env:"REPORT_INTERVAL" flag:"r"`
	PollInterval   uint   `env:"POLL_INTERVAL" flag:"p"`
	ReportNoBatch  bool   `env:"REPORT_NOBATCH" flag:"nob"`
	Key            string `env:"KEY" flag:"k"`
	RateLimit      uint   `env:"RATE_LIMIT" flag:"l"`
}

const (
	DefaultAdress         = "localhost:8080"
	DefaultReportInterval = 10
	DefaultPollInterval   = 2
	DefaultReportNoBatch  = true
	DefaultRateLimit      = 0 // без ограничений
)
