package task_handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/model"
	"github.com/umarkotak/vdub-go/service"
	"github.com/umarkotak/vdub-go/utils"
)

func PostStartDubbTaskV2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commonCtx := utils.GetCommonCtx(r)

	params := StartDubbTaskParams{}
	err := utils.BindJson(r, &params)
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 400, err)
		return
	}
	params.Gen(commonCtx.DirectUsername)

	state, err := service.GetState(ctx, params.TaskDir, model.TaskState{
		YoutubeUrl: params.YoutubeUrl,
		VoiceName:  params.VoiceName,
		VoiceRate:  params.VoiceRate,
		VoicePitch: params.VoicePitch,
	})
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}
	if state.YoutubeUrl == "" {
		state.YoutubeUrl = params.YoutubeUrl
	}
	if state.VoiceName == "" {
		state.VoiceName = params.VoiceName
	}
	if state.VoiceRate == "" {
		state.VoiceRate = params.VoiceRate
	}
	if state.VoicePitch == "" {
		state.VoicePitch = params.VoicePitch
	}

	if params.ForceStartFrom != "" {
		state.Status = params.ForceStartFrom

		if params.VoiceName != "" {
			state.VoiceName = params.VoiceName
		}
		if params.VoiceRate != "" {
			state.VoiceRate = params.VoiceRate
		}
		if params.VoicePitch != "" {
			state.VoicePitch = params.VoicePitch
		}

		err = service.SaveStateStatus(ctx, params.TaskDir, &state, params.ForceStartFrom)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return
		}
	}

	if handlerState.RunningTask[params.TaskName] {
		err = fmt.Errorf("task is still running")
		utils.RenderError(w, r, 422, err)
		return
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
			err = service.DownloadYoutubeVideo(bgCtx, state.YoutubeUrl, params.RawVideoPath)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
				return
			}

			err = service.GenerateVideoSnapshot(bgCtx, params.RawVideoPath, params.VideoScreenshotPath)
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
			logrus.Infof("DUBBING TASK RUNNING: %s; (6/10) %s", params.TaskName, "transcript audio with diarization")

			logrus.Infof("DUBBING TASK RUNNING: %s; (6/10) %s", params.TaskName, fmt.Sprintf("diarization started: %s", time.Now().Format(time.RFC3339)))
			err = service.DiarizeVoice(bgCtx, params.TaskDir)
			if err != nil {
				logrus.WithContext(bgCtx).Error(err)
			}
			logrus.Infof("DUBBING TASK RUNNING: %s; (6/10) %s", params.TaskName, fmt.Sprintf("diarization finished: %s", time.Now().Format(time.RFC3339)))

			err = service.TranscriptAudioWithDiarization(bgCtx, params.TaskDir, params.Vocal16KHzPath, params.TranscriptPath, params.SegmentedSpeechDir)
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
				Name:  state.VoiceName,
				Rate:  state.VoiceRate,
				Pitch: state.VoicePitch,
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
