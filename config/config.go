package config

type Server struct {
	// gorm
	Mysql Mysql `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Email Email `mapstructure:"email" json:"email" yaml:"email"`
}
