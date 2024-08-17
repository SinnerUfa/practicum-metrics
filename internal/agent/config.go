package agent

type Config struct {
	Adress         string `env:"ADDRESS" flag:"a"`
	ReportInterval uint   `env:"REPORT_INTERVAL" flag:"r"`
	PollInterval   uint   `env:"POLL_INTERVAL" flag:"p"`
	ReportNoBatch  bool   `env:"REPORT_NOBATCH" flag:"nob"`
}

var DefaultConfig = Config{
	Adress:         "localhost:8080",
	ReportInterval: 10,
	PollInterval:   2,
	ReportNoBatch:  true,
}
