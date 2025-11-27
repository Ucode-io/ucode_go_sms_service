package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/spf13/cast"
)

const (
	// DebugMode indicates service mode is debug.
	DebugMode = "debug"
	// TestMode indicates service mode is test.
	TestMode = "test"
	// ReleaseMode indicates service mode is release.
	ReleaseMode = "release"
)

// Config ...
type Config struct {
	ServiceName string
	Environment string // develop, staging, production

	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string

	LogLevel    string
	RPCPort     string
	AccountSid  string
	AuthToken   string
	PhoneNumber string
	SMSSender   string

	PlayMobileLogin      string
	PlayMobilePassword   string
	PlayMobileUrl        string
	PlayMobileOriginator string

	PostgresMaxConnections int32

	SmsServiceHost string
	SmsGRPCPort    string
}

// Load loads environment vars and inflates Config
func Load() Config {

	if err := godotenv.Load("/app/.env"); err != nil {
		fmt.Println("No .env file found")
	}

	c := Config{}

	c.ServiceName = cast.ToString(getOrReturnDefault("SERVICE_NAME", "ucode/ucode_go_sms_service"))

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", TestMode))
	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", ""))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "sms_service"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "sms_service"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", ""))
	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.RPCPort = cast.ToString(getOrReturnDefault("RPC_PORT", ":5004"))

	c.AccountSid = cast.ToString(getOrReturnDefault("ACCOUNT_SID", ""))
	c.AuthToken = cast.ToString(getOrReturnDefault("AUTH_TOKEN", ""))
	c.PhoneNumber = cast.ToString(getOrReturnDefault("PHONE_NUMBER", ""))

	c.SmsServiceHost = cast.ToString(getOrReturnDefault("SMS_SERVICE_HOST", "0.0.0.0"))
	c.SmsGRPCPort = cast.ToString(getOrReturnDefault("SMS_GRPC_PORT", ":9105"))

	c.SMSSender = cast.ToString(getOrReturnDefault("SMS_SENDER", "play_mobile"))

	c.PostgresMaxConnections = cast.ToInt32(getOrReturnDefault("POSTGRES_MAX_CONNECTIONS", 30))

	c.PlayMobileLogin = cast.ToString(getOrReturnDefault("PLAY_MOBILE_LOGIN", ""))
	c.PlayMobilePassword = cast.ToString(getOrReturnDefault("PLAY_MOBILE_PASSWORD", ""))
	c.PlayMobileUrl = cast.ToString(getOrReturnDefault("PLAY_MOBILE_URL", "https://send.smsxabar.uz/broker-api/send"))
	c.PlayMobileOriginator = cast.ToString(getOrReturnDefault("PLAY_MOBILE_ORIGINATOR", "3700"))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
