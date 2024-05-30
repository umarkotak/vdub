package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/model"
	"github.com/umarkotak/vdub-go/service"
	"github.com/umarkotak/vdub-go/utils"
)

type (
	StartDubbTaskParams struct {
		TaskName       string `json:"task_name"`        // must unique - it will determine the task folder
		YoutubeUrl     string `json:"youtube_url"`      //
		VoiceName      string `json:"voice_name"`       // eg: id-ID-ArdiNeural
		VoiceRate      string `json:"voice_rate"`       // eg: [-/+]10%
		VoicePitch     string `json:"voice_pitch"`      // eg: [-/+]10Hz
		ForceStartFrom string `json:"force_start_from"` // used to run from certain state

		TaskDir                  string
		RawVideoName             string
		RawVideoPath             string
		RawVideoAudioName        string
		RawVideoAudioPath        string
		AudioInstrumentPath      string
		AudioVocalPath           string
		Vocal16KHzName           string
		Vocal16KHzPath           string
		InstrumentVideoPath      string
		TranscriptPath           string
		TranscriptVttPath        string
		TranscriptTranslatedPath string
		GeneratedSpeechDir       string
		SpeechAdjustedDir        string
		DubbedVideoPath          string
	}
)

const (
	TOTAL_STEPS = "10"
)

func (p *StartDubbTaskParams) Gen() {
	p.TaskDir = fmt.Sprintf("%s/%s", config.Get().BaseDir, p.TaskName)
	p.RawVideoName = "raw_video.mp4"
	p.RawVideoPath = fmt.Sprintf("%s/%s", p.TaskDir, p.RawVideoName)
	p.RawVideoAudioName = "raw_video_audio.wav"
	p.RawVideoAudioPath = fmt.Sprintf("%s/%s", p.TaskDir, p.RawVideoAudioName)
	p.AudioInstrumentPath = fmt.Sprintf("%s_Instruments.wav", strings.TrimSuffix(p.RawVideoAudioPath, ".wav"))
	p.AudioVocalPath = fmt.Sprintf("%s_Vocals.wav", strings.TrimSuffix(p.RawVideoAudioPath, ".wav"))
	p.Vocal16KHzName = "raw_video_audio_Vocals_16KHz.wav"
	p.Vocal16KHzPath = fmt.Sprintf("%s/%s", p.TaskDir, p.Vocal16KHzName)
	p.InstrumentVideoPath = fmt.Sprintf("%s/%s", p.TaskDir, "instrument_video.mp4")
	p.TranscriptPath = fmt.Sprintf("%s/%s", p.TaskDir, "transcript")
	p.TranscriptVttPath = fmt.Sprintf("%s/%s", p.TaskDir, "transcript.vtt")
	p.TranscriptTranslatedPath = fmt.Sprintf("%s/%s", p.TaskDir, "transcript_translated.vtt")
	p.GeneratedSpeechDir = fmt.Sprintf("%s/generated_speech", p.TaskDir)
	p.SpeechAdjustedDir = fmt.Sprintf("%s/adjusted_speech", p.TaskDir)
	p.DubbedVideoPath = fmt.Sprintf("%s/%s", p.TaskDir, "dubbed_video.mp4")
}

func PostStartDubbTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := StartDubbTaskParams{}
	err := utils.BindJson(r, &params)
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 400, err)
		return
	}
	params.Gen()

	state, err := service.GetState(ctx, params.TaskDir)
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}

	if handlerState.RunningTask[params.TaskName] {
		err = fmt.Errorf("task is still running")
		utils.RenderError(w, r, 422, err)
		return
	}

	if params.ForceStartFrom != "" {
		state.Status = params.ForceStartFrom
	}

	reqID := chiMiddleware.GetReqID(ctx)
	go func() {
		bgCtx := context.Background()
		bgCtx = context.WithValue(bgCtx, chiMiddleware.RequestIDKey, reqID)

		handlerState.RunningTask[params.TaskName] = true
		logrus.Infof("DUBBING TASK START: %s", params.TaskName)
		defer func() {
			handlerState.RunningTask[params.TaskName] = false
			logrus.Infof("DUBBING TASK FINISH: %s", params.TaskName)
		}()

		if state.Status == model.STATE_INITIALIZED {
			logrus.Infof("DUBBING TASK RUNNING: %s; (1/10) %s", params.TaskName, "download youtube video")
			err = service.DownloadYoutubeVideo(bgCtx, params.YoutubeUrl, params.RawVideoPath)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_VIDEO_DOWNLOADED)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}
		}

		if state.Status == model.STATE_VIDEO_DOWNLOADED {
			logrus.Infof("DUBBING TASK RUNNING: %s; (2/10) %s", params.TaskName, "separating video and audio")
			err = service.SeparateAudio(bgCtx, params.RawVideoPath, params.RawVideoAudioPath)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_VIDEO_AUDIO_GENERATED)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}
		}

		if state.Status == model.STATE_VIDEO_AUDIO_GENERATED {
			logrus.Infof("DUBBING TASK RUNNING: %s; (3/10) %s", params.TaskName, "separate audio vocal and instrument")
			err = service.SeparateVocal(bgCtx, params.RawVideoAudioPath, params.TaskDir)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_VIDEO_AUDIO_SEPARATED)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}
		}

		if state.Status == model.STATE_VIDEO_AUDIO_SEPARATED {
			logrus.Infof("DUBBING TASK RUNNING: %s; (4/10) %s", params.TaskName, "convert vocal audio to 16KHz")
			err = service.Generate16KHzAudio(bgCtx, params.AudioVocalPath, params.Vocal16KHzPath)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_AUDIO_16KHZ_GENERATED)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}
		}

		if state.Status == model.STATE_AUDIO_16KHZ_GENERATED {
			logrus.Infof("DUBBING TASK RUNNING: %s; (5/10) %s", params.TaskName, "merging video with instrument")
			err = service.MergeVideoWithAudio(bgCtx, params.RawVideoPath, params.AudioInstrumentPath, params.InstrumentVideoPath)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_VIDEO_WITH_INSTRUMENT_GENERATED)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}
		}

		if state.Status == model.STATE_VIDEO_WITH_INSTRUMENT_GENERATED {
			logrus.Infof("DUBBING TASK RUNNING: %s; (6/10) %s", params.TaskName, "transcript audio")
			err = service.TranscriptAudio(bgCtx, params.Vocal16KHzPath, params.TranscriptPath)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_AUDIO_TRANSCRIPTED)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}
		}

		if state.Status == model.STATE_AUDIO_TRANSCRIPTED {
			logrus.Infof("DUBBING TASK RUNNING: %s; (7/10) %s", params.TaskName, "translating transcript")
			err = service.TranslateTranscript(bgCtx, params.TranscriptVttPath, params.TranscriptTranslatedPath)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_TRANSCRIPT_TRANSLATED)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}
		}

		if state.Status == model.STATE_TRANSCRIPT_TRANSLATED {
			logrus.Infof("DUBBING TASK RUNNING: %s; (8/10) %s", params.TaskName, "generating translated audio")
			err = service.GenerateVoice(bgCtx, params.TranscriptTranslatedPath, params.GeneratedSpeechDir, service.VoiceOpts{
				Name:  params.VoiceName,
				Rate:  params.VoiceRate,
				Pitch: params.VoicePitch,
			})
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_AUDIO_GENERATED)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}
		}

		if state.Status == model.STATE_AUDIO_GENERATED {
			logrus.Infof("DUBBING TASK RUNNING: %s; (9/10) %s", params.TaskName, "adjusting audio speed")
			err = service.AdjustVoiceSpeed(bgCtx, params.TranscriptTranslatedPath, params.GeneratedSpeechDir, params.SpeechAdjustedDir)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_AUDIO_ADJUSTED)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}
		}

		if state.Status == model.STATE_AUDIO_ADJUSTED {
			logrus.Infof("DUBBING TASK RUNNING: %s; (10/10) %s", params.TaskName, "merge video with translatted audio")
			err = service.MergeVideoWithDubb(
				bgCtx,
				params.TranscriptTranslatedPath,
				params.SpeechAdjustedDir,
				params.InstrumentVideoPath,
				params.DubbedVideoPath,
			)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_DUBBED_VIDEO_GENERATED)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}
		}

	}()

	utils.Render(w, r, 200, state, nil)
}