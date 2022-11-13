package env

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Host:     "mouse.db.elephantsql.com",
		Port:     5432,
		User:     "tfalxqpw",
		Password: "P4Ewg8JM_QgiH6pmWhvQpIevj_XmsHvf",
		DBName:   "tfalxqpw",
	}
}
