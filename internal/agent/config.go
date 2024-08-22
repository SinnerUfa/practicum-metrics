package agent

type Config struct {
	Adress         string `env:"ADDRESS" flag:"a"`
	ReportInterval uint   `env:"REPORT_INTERVAL" flag:"r"`
	PollInterval   uint   `env:"POLL_INTERVAL" flag:"p"`
	ReportNoBatch  bool   `env:"REPORT_NOBATCH" flag:"nob"`
}

const (
	DefaultAdress         = "localhost:8080"
	DefaultReportInterval = 10
	DefaultPollInterval   = 2
	DefaultReportNoBatch  = true
)
