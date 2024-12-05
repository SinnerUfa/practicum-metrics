package server

type Config struct {
	Adress          string `env:"ADDRESS" flag:"a"`
	StoreInterval   uint   `env:"STORE_INTERVAL" flag:"i"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" flag:"f"`
	Restore         bool   `env:"RESTORE" flag:"r"`
	DatabaseDSN     string `env:"DATABASE_DSN" flag:"d"`
}

const (
	DefaultAdress          = "localhost:8080"
	DefaultStoreInterval   = 300
	DefaultFileStoragePath = "/tmp/metrics-db.json"
	DefaultRestore         = true
	DefaultDatabaseDSN     = ""
)
