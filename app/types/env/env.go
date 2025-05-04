package types

type SMTPConfig struct {
	Host string
	Port string
}

type EnvConfig struct {
	CryAppURL    string
	APIPort      int
	DBHost       string
	DBPort       int
	DBDatabase   string
	DBUserName   string
	DBPassword   string
	SMTPConfig   SMTPConfig
	NoReplyEmail string
}
