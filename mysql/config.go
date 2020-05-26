package mysql

type Config struct {
	Name       string
	Port       int
	Host       string
	Username   string
	Password   string
	Charset    string
	Debug      bool
	TZLocation string
	DataBaseTZ string
	PoolSize   int
}
