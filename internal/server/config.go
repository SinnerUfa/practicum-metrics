package server

import "fmt"

type Config struct {
	Adress string `env:"ADDRESS" flag:"a"`
}

func (c Config) String() string {
	return fmt.Sprint("\nAdress=", c.Adress)
}

var DefaultConfig = Config{
	Adress: "localhost:8080",
}
