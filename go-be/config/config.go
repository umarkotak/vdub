package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		BaseDir               string // base dir will contain all tasks folder
		VocalRemoverPy        string
		VocalRemoverModelPath string
		GoogleAiStudioKey     string //
	}
)

var (
	config = Config{}
)

func InitConfig() {
	godotenv.Load()

	config = Config{
		BaseDir:               GetStringWithDefault("BASE_DIR", "/root/shared"),
		VocalRemoverPy:        GetStringWithDefault("VOCAL_REMOVER_PY", "/root/vocal-remover/inference.py"),
		VocalRemoverModelPath: GetStringWithDefault("VOCAL_REMOVER_MODEL_PATH", "/root/vocal-remover/baseline.pth"),
		GoogleAiStudioKey:     os.Getenv("GOOGLE_AI_STUDIO_KEY"),
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
