package datastore

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"google.golang.org/api/option"
)

type (
	DataStore struct {
		genaiGiminiProVision *genai.GenerativeModel
	}
)

var (
	dataStore = DataStore{}
)

func InitDataStore() {
	genaiClient, err := genai.NewClient(
		context.TODO(),
		option.WithAPIKey(config.Get().GoogleAiStudioKey),
	)
	if err != nil {
		logrus.Error(err)
	}

	// genaiGiminiProVision = genaiClient.GenerativeModel("gemini-pro-vision")
	// genaiGiminiProVision.SafetySettings = []*genai.SafetySetting{
	// 	{
	// 		Category:  genai.HarmCategoryUnspecified,
	// 		Threshold: genai.HarmBlockNone,
	// 	},
	// }

	dataStore = DataStore{
		genaiGiminiProVision: genaiClient.GenerativeModel("gemini-pro"),
	}
}

func Get() DataStore {
	return dataStore
}
