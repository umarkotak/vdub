package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		BaseDir           string // base dir will contain all tasks folder
		GoogleAiStudioKey string //
	}
)

var (
	config = Config{}
)

func InitConfig() {
	godotenv.Load()

	config = Config{
		BaseDir:           GetStringWithDefault("BASE_DIR", "/root/shared"),
		GoogleAiStudioKey: os.Getenv("GOOGLE_AI_STUDIO_KEY"),
	}
}

func Get() Config {
	return config
}

func GetStringWithDefault(envKey, defVal string) string {
	if os.Getenv(envKey) != "" {
		return os.Getenv(envKey)
	}
	return defVal
}
