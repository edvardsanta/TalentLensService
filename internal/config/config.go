package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading .env file: ", err)
	}
}

func GetDBConnectionString() string {
	connectionString := viper.GetString("DB_CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal("DB_CONNECTION_STRING not set")
	}
	return connectionString
}

func GetDBDriver() string {
	driver := viper.GetString("DB_DRIVER")
	if driver == "" {
		log.Fatal("DB_DRIVER not set")
	}
	return driver
}

func GetJWTSecretKey() string {
	return viper.GetString("JWT_SECRET_KEY")
}

func GetAllowedOrigins() []string {
	allowedOrigins := viper.GetString("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		log.Fatal("ALLOWED_ORIGINS not set")
	}
	fmt.Println(allowedOrigins)
	return strings.Split(allowedOrigins, ",")
}

func GetAppTelemetryInfo() (string, string) {
	appName := viper.GetString("SERVICE_NAME")
	appVersion := viper.GetString("SERVICE_VERSION")
	if appName == "" || appVersion == "" {
		log.Fatal("SERVICE_NAME or SERVICE_VERSION not set")
	}
	return appName, appVersion
}

func IsOpenTelemetryDisabled() bool {
	return viper.GetBool("OTEL_SDK_DISABLED")
}
