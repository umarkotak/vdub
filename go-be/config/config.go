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
		WhisperBinary         string
		WhisperModelPath      string
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
		WhisperBinary:         GetStringWithDefault("WHISPER_BINARY", "/root/whisper.cpp/main"),
		WhisperModelPath:      GetStringWithDefault("WHISPER_MODEL_PATH", "/root/whisper.cpp/models/ggml-medium.en-q5_0.bin"),
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
