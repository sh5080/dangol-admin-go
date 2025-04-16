package config

// DBConfig는 데이터베이스 연결 설정을 위한 구조체입니다.
type DBConfig struct {
	DatabaseURL string
}

// NewDBConfig는 환경 변수에서 데이터베이스 설정을 로드합니다.
func (c *Config) NewDBConfig() *DBConfig {
	databaseURL := GetEnvOrDefault("DATABASE_URL", "postgresql://neondb_owner:npg_et5nvG0lMogX@ep-nameless-snowflake-a15r97wt-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=verify-full")
	if databaseURL != "" {
		return &DBConfig{
			DatabaseURL: databaseURL,
		}
	}
	return nil
}
