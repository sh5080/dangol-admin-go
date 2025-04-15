package config

// DBConfig는 데이터베이스 연결 설정을 위한 구조체입니다.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// NewDBConfig는 환경 변수에서 데이터베이스 설정을 로드합니다.
func (c *Config) NewDBConfig() *DBConfig {
	return &DBConfig{
		Host:     GetEnvOrDefault("DB_HOST", "localhost"),
		Port:     GetEnvOrDefault("DB_PORT", "5432"),
		User:     GetEnvOrDefault("DB_USER", "postgres"),
		Password: GetEnvOrDefault("DB_PASSWORD", "postgres"),
		Database: GetEnvOrDefault("DB_NAME", "postgresdb"),
		SSLMode:  GetEnvOrDefault("DB_SSL_MODE", "disable"),
	}
} 