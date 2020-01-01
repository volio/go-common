package mysql

type Config struct {
	Name     string
	Port     int
	Host     string
	Username string
	Password string
	Charset  string
	Debug    bool
	Location string
	PoolSize int
}
