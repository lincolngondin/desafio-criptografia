package config

type configs struct {
	DBDataSourceName string
	DBDriverName     string
}

func New() *configs {
	return &configs{
		DBDriverName:     "sqlite",
		DBDataSourceName: "data/database.sqlite",
	}
}
