package env

import (
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"regexp"
	"sso/internal/config"
	"strconv"
	"strings"
	"time"
)

var lettersL = []rune("abcdefghijklmnopqrstuvwxyz")
var lettersU = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var digit = []rune("1234567890")
var special = []rune("!@#$%^&*_+-")

func genPassword(n int) string {
	b := make([]rune, n)
	for i := 0; i < n; i += 4 {
		b[i] = lettersL[rand.Intn(len(lettersL))]
		if n >= i+1 {
			b[i+1] = lettersU[rand.Intn(len(lettersU))]
		}
		if n >= i+2 {
			b[i+2] = digit[rand.Intn(len(digit))]
		}
		if n >= i+3 {
			b[i+3] = special[rand.Intn(len(special))]
		}
	}
	return string(b)
}

func New() *config.Config {
	if err := godotenv.Load(); err != nil {
		log.Print("Load config: No .env file found")
	}
	return &config.Config{
		DebugLevel:   getEnv("DEBUG_LEVEL", "local"),
		RootPassword: getEnv("INIT_ROOT_PASSWORD", genPassword(16)),
		Server: config.ServerConfig{
			Address:      getEnv("SERVER_ADDRESS", "0.0.0.0:8085"),
			ReadTimeout:  getEnvDuration("SERVER_TIMEOUT_READ", 5*time.Second),
			WriteTimeout: getEnvDuration("SERVER_TIMEOUT_WRITE", 5*time.Second),
		},
		Db: config.DbConfig{
			Server:   getEnv("DB_SERVER", "localhost:27017"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_DATABASE", "SSO"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")
	if valStr == "" {
		return defaultVal
	}
	val := strings.Split(valStr, sep)
	return val
}

func getEnvDuration(name string, defaultVal time.Duration) time.Duration {
	valStr := getEnv(name, "")
	if valStr != "" {
		rResult := regexp.MustCompile(`^([0-9]*)(m|s|ms|us|ns)?$`).FindStringSubmatch(valStr)
		if len(rResult) != 0 {
			t, _ := strconv.Atoi(rResult[1]) //regex гарантирует что будет число при совпадении и будет минимум 3 подстроки
			switch rResult[2] {
			case "m":
				t *= 1000000000 * 60
			case "s":
				t *= 1000000000
			case "ms":
				t *= 1000000
			case "us":
				t *= 1000
			}
			return time.Duration(t)
		}
	}
	return defaultVal
}
