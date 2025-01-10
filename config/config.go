package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		BaseDir                     string // base dir will contain all tasks folder
		VocalRemoverPy              string //
		VocalRemoverModelPath       string //
		WhisperBinary               string //
		WhisperModelPath            string //
		GoogleAiStudioKey           string //
		HuggingFaceDiarizationToken string //
		PythonDiarizationPath       string //
	}
)

var (
	config = Config{}
)

func InitConfig() {
	godotenv.Load()

	config = Config{
		BaseDir:                     GetStringWithDefault("BASE_DIR", "/root/shared"),
		VocalRemoverPy:              GetStringWithDefault("VOCAL_REMOVER_PY", "/root/vocal-remover/inference.py"),
		VocalRemoverModelPath:       GetStringWithDefault("VOCAL_REMOVER_MODEL_PATH", "/root/vocal-remover/baseline.pth"),
		WhisperBinary:               GetStringWithDefault("WHISPER_BINARY", "/root/whisper.cpp/build/bin/whisper-cli"),
		WhisperModelPath:            GetStringWithDefault("WHISPER_MODEL_PATH", "/root/whisper.cpp/models/ggml-large-v3-turbo-q5_0.bin"),
		PythonDiarizationPath:       GetStringWithDefault("PYTHON_DIARIZATION_PATH", "/root/vdub/bin/speaker_diarization/main.py"),
		GoogleAiStudioKey:           os.Getenv("GOOGLE_AI_STUDIO_KEY"),
		HuggingFaceDiarizationToken: os.Getenv("HUGGING_FACE_DIARIZATION_TOKEN"),
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
