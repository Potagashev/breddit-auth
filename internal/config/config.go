package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	
	DbName string
	DbUser string
	DbPassword string
	DbHost string
	DbPort string
	DbUrl string

	JWTSecret string
	JWTTokenExpMinutes int
}

func LoadConfig() (*Config, error) {
	needToParseEnvFileString := getEnv("NEED_TO_PARSE_ENV_FILE", "false")
	needToParseEnvFileBool, err := strconv.ParseBool(needToParseEnvFileString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Environmental variable 'NEED_TO_PARSE_ENV_FILE' cannot be parsed as boolean")
		os.Exit(1)
	}
	if needToParseEnvFileBool {
		err := godotenv.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occured while loading environmental variables from file")
			os.Exit(1)
		}	
	}

	dbName := mustGetEnv("POSTGRES_DB")
	dbUser := mustGetEnv("POSTGRES_USER")
	dbPassword := mustGetEnv("POSTGRES_PASSWORD")
	dbHost := mustGetEnv("DATABASE_HOST")
	dbPort := mustGetEnv("DATABASE_PORT")
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	appPort := getEnv("APP_PORT", "8080")

	jwtTokenExpMinutes, err := strconv.Atoi(getEnv("JWT_TOKEN_EXP_MINUTES", "60"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "JWT_TOKEN_EXP_MINUTES must be a valid integer value")
		os.Exit(1)
	}
	config := &Config{
		DbName: dbName,
		DbUser: dbUser,
		DbPassword: dbPassword,
		DbHost: dbHost,
		DbPort: dbPort,
		DbUrl: dbUrl,
		
		AppPort: appPort,

		JWTSecret: mustGetEnv("JWT_SECRET"),
		JWTTokenExpMinutes: jwtTokenExpMinutes,
	}
	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func mustGetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	fmt.Fprintf(os.Stderr, "Fatal error: environment variable %s is not set\n", key)
	os.Exit(1)
	return ""
}