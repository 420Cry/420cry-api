package types

type SMTPConfig struct {
	Host string
	Port string
}

type EnvConfig struct {
	APIPort      int
	App          string
	DBHost       string
	DBPort       int
	DBDatabase   string
	DBUserName   string
	DBPassword   string
	SMTPConfig   SMTPConfig
	NoReplyEmail string
}
