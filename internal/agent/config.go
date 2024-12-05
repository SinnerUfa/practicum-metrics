package agent

import "fmt"

type Config struct {
	Adress         string `env:"ADRESS" flag:"a"`
	ReportInterval uint   `env:"REPORT_INTERVAL" flag:"r"`
	PollInterval   uint   `env:"POLL_INTERVAL" flag:"p"`
}

func (c Config) String() string {
	return fmt.Sprint("\nAdress=", c.Adress, "\nReportInterval=", c.ReportInterval, "\nPollInterval=", c.PollInterval)
}

var DefaultConfig = Config{
	Adress:         "localhost:8080",
	ReportInterval: 10,
	PollInterval:   2,
}
