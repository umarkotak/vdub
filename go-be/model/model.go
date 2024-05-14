package model

type (
	TaskState struct {
		Status        string       `json:"status"`         // Enum: initialized
		Progress      string       `json:"progress"`       //
		Transcripts   []Transcript `json:"transcripts"`    //
		RawTranscript string       `json:"raw_transcript"` //
	}

	Transcript struct {
		Idx             int64  `json:"idx"`
		TsStart         string `json:"ts_start"`
		TsStop          string `json:"ts_stop"`
		RawText         string `json:"raw_text"`
		TranslatedText  string `json:"translated_text"`
		RawAudio        string `json:"raw_audio"`
		TranslatedAudio string `json:"translated_audio"`
	}
)
