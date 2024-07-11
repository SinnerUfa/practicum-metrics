package server

type Config struct {
	Adress          string `env:"ADDRESS" flag:"a"`
	StoreInterval   uint   `env:"STORE_INTERVAL" flag:"i"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" flag:"f"`
	Restore         bool   `env:"RESTORE" flag:"r"`
	DatabaseDSN     string `env:"DATABASE_DSN" flag:"d"`
}

var DefaultConfig = Config{
	Adress:          "localhost:8080",
	StoreInterval:   2,
	FileStoragePath: "/tmp/metrics-db.json",
	Restore:         true,
	DatabaseDSN:     "",
}
