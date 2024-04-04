package dbs

type DBConfig struct {
	Addr     string
	Port     uint16
	User     string
	Password string
	DB       string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Addr:     "localhost",
		Port:     5432,
		User:     "admin",
		Password: "12345",
		DB:       "postgres",
	}
}
