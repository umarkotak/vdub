package model

import "fmt"

type (
	TaskData struct {
		Name               string `json:"name"`
		Status             string `json:"status"`
		CurrentStatusHuman string `json:"current_status_human"`
		IsRunning          bool   `json:"is_running"`
		ProgressSummary    string `json:"progress_summary"`
	}

	TaskState struct {
		Status        string       `json:"status"`         // Enum: initialized
		Progress      string       `json:"progress"`       //
		Transcripts   []Transcript `json:"transcripts"`    //
		RawTranscript string       `json:"raw_transcript"` //

		YoutubeUrl string `json:"youtube_url" validate:"required"` //
		VoiceName  string `json:"voice_name" validate:"required"`  // eg: id-ID-ArdiNeural
		VoiceRate  string `json:"voice_rate" validate:"required"`  // eg: [-/+]10%
		VoicePitch string `json:"voice_pitch" validate:"required"` // eg: [-/+]10Hz
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

	TaskStateProgress struct {
		Name     string `json:"name"`
		Progress string `json:"progress"` // Enum: not_done, running, done
	}

	GetTaskStateData struct {
		Status             string              `json:"status"`
		CurrentStatusHuman string              `json:"current_status_human"`
		IsRunning          bool                `json:"is_running"`
		ProgressSummary    string              `json:"progress_summary"`
		Progresses         []TaskStateProgress `json:"progresses"`
	}
)

func (ts *TaskState) GetTaskStateData(isRunning bool) GetTaskStateData {
	progresses := []TaskStateProgress{}

	currStateIdx := STATE_IDX_MAP[ts.Status].Idx

	for idx, stateName := range STATE_IDX_ARR {
		progress := "not_done"
		if isRunning && idx == currStateIdx+1 {
			progress = "running"
		} else if idx <= currStateIdx {
			progress = "done"
		}

		progresses = append(progresses, TaskStateProgress{
			Name:     stateName,
			Progress: progress,
		})
	}

	return GetTaskStateData{
		Status:             ts.Status,
		CurrentStatusHuman: STATE_IDX_MAP[ts.Status].StatusHuman,
		IsRunning:          isRunning,
		ProgressSummary:    fmt.Sprintf("%v/%v", currStateIdx, 10),
		Progresses:         progresses,
	}
}
