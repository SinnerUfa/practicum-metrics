package server

type Config struct {
	Adress string `env:"ADDRESS" flag:"a"`
}

var DefaultConfig = Config{
	Adress: "localhost:8080",
}
