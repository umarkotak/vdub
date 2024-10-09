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

	TaskStateProgress struct {
		Name      string `json:"name"`
		NameHuman string `json:"name_human"`
		Progress  string `json:"progress"` // Enum: not_done, running, done
	}

	GetTaskStateData struct {
		Status             string              `json:"status"`
		CurrentStatusHuman string              `json:"current_status_human"`
		IsRunning          bool                `json:"is_running"`
		ProgressSummary    string              `json:"progress_summary"`
		Progresses         []TaskStateProgress `json:"progresses"`
		Completed          bool                `json:"completed"`
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
			Name:      stateName,
			NameHuman: STATE_IDX_MAP[stateName].StatusHuman,
			Progress:  progress,
		})
	}

	return GetTaskStateData{
		Status:             ts.Status,
		CurrentStatusHuman: STATE_IDX_MAP[ts.Status].StatusHuman,
		IsRunning:          isRunning,
		ProgressSummary:    fmt.Sprintf("%v/%v", currStateIdx, 10),
		Progresses:         progresses,
		Completed:          currStateIdx == (len(STATE_IDX_ARR) - 1),
	}
}
