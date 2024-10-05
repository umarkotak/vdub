package model

type (
	Transcript struct {
		Idx             int64  `json:"idx"`
		TsStart         string `json:"ts_start"`
		TsStop          string `json:"ts_stop"`
		RawText         string `json:"raw_text"`
		TranslatedText  string `json:"translated_text"`
		RawAudio        string `json:"raw_audio"`
		TranslatedAudio string `json:"translated_audio"`
	}

	TranscriptUpdateParams struct {
		TaskName       string           `json:"-"`
		TranscriptData []TranscriptData `json:"transcript_data"`
	}

	TranscriptData struct {
		StartAt string `json:"start_at"`
		EndAt   string `json:"end_at"`
		Value   string `json:"value"`
	}
)
