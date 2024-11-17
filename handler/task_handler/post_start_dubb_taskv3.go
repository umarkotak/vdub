package task_handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/model"
	"github.com/umarkotak/vdub-go/service"
	"github.com/umarkotak/vdub-go/utils"
)

func PostStartTask(w http.ResponseWriter, r *http.Request) {
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
		TaskType:   params.TaskType,
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

	if params.ForceStartFrom != "" {
		state.Status = params.ForceStartFrom

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
	bgCtx := context.WithValue(context.Background(), chiMiddleware.RequestIDKey, reqID)
	go func() {
		logrusProc := logrus.WithContext(bgCtx).WithField("task_dir", params.TaskDir)

		handlerState.RunningTask[params.TaskName] = true
		logrusProc.Infof("DUBBING TASK START: %s", params.TaskName)
		defer func() {
			handlerState.RunningTask[params.TaskName] = false

			if state.Status == model.STATE_DUBBED_VIDEO_GENERATED {
				logrusProc.Infof("DUBBING TASK FINISH: %s", params.TaskName)
			} else {
				logrusProc.Infof("DUBBING TASK STOPEPD: %s", params.TaskName)
			}
		}()

		if state.Status == model.STATE_INITIALIZED {
			logrusProc.Infof("DUBBING TASK RUNNING: %s; (1/11) %s", params.TaskName, "download youtube video")
			err = service.DownloadYoutubeVideo(bgCtx, state.YoutubeUrl, params.RawVideoPath, params.TaskDir)
			if err != nil {
				logrusProc.Error(err)
				return
			}

			err = service.GenerateVideoSnapshot(bgCtx, params.RawVideoPath, params.VideoScreenshotPath)
			if err != nil {
				logrusProc.Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_VIDEO_DOWNLOADED)
			if err != nil {
				logrusProc.Error(err)
				return
			}
		}

		if state.Status == model.STATE_VIDEO_DOWNLOADED {
			logrusProc.Infof("DUBBING TASK RUNNING: %s; (2/11) %s", params.TaskName, "separating video and audio")
			err = service.SeparateAudio(bgCtx, params.RawVideoPath, params.RawVideoAudioPath)
			if err != nil {
				logrusProc.Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_VIDEO_AUDIO_GENERATED)
			if err != nil {
				logrusProc.Error(err)
				return
			}
		}

		if state.Status == model.STATE_VIDEO_AUDIO_GENERATED {
			logrusProc.Infof("DUBBING TASK RUNNING: %s; (3/11) %s", params.TaskName, "separate audio vocal and instrument")
			err = service.SeparateVocal(bgCtx, params.RawVideoAudioPath, params.TaskDir)
			if err != nil {
				logrusProc.Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_VIDEO_AUDIO_SEPARATED)
			if err != nil {
				logrusProc.Error(err)
				return
			}
		}

		if state.Status == model.STATE_VIDEO_AUDIO_SEPARATED {
			logrusProc.Infof("DUBBING TASK RUNNING: %s; (4/11) %s", params.TaskName, "convert vocal audio to 16KHz")
			err = service.Generate16KHzAudio(bgCtx, params.AudioVocalPath, params.Vocal16KHzPath)
			if err != nil {
				logrusProc.Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_AUDIO_16KHZ_GENERATED)
			if err != nil {
				logrusProc.Error(err)
				return
			}
		}

		if state.Status == model.STATE_AUDIO_16KHZ_GENERATED {
			logrusProc.Infof("DUBBING TASK RUNNING: %s; (5/11) %s", params.TaskName, "merging video with instrument")
			err = service.MergeVideoWithAudio(bgCtx, params.RawVideoPath, params.AudioInstrumentPath, params.InstrumentVideoPath)
			if err != nil {
				logrusProc.Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_VIDEO_WITH_INSTRUMENT_GENERATED)
			if err != nil {
				logrusProc.Error(err)
				return
			}
		}

		if state.Status == model.STATE_VIDEO_WITH_INSTRUMENT_GENERATED {
			logrusProc.Infof("DUBBING TASK RUNNING: %s; (6/11) %s", params.TaskName, "diarization started")

			err = service.DiarizeVoice(bgCtx, params.TaskDir)
			if err != nil {
				logrusProc.Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_AUDIO_DIARIZED)
			if err != nil {
				logrusProc.Error(err)
				return
			}
		}

		if state.Status == model.STATE_AUDIO_DIARIZED {
			logrusProc.Infof("DUBBING TASK RUNNING: %s; (7/11) %s", params.TaskName, "transcript audio with diarization")

			err = service.TranscriptAudioWithDiarization(bgCtx, params.TaskDir, params.Vocal16KHzPath, params.TranscriptPath, params.SegmentedSpeechDir)
			if err != nil {
				logrusProc.Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_AUDIO_TRANSCRIPTED)
			if err != nil {
				logrusProc.Error(err)
				return
			}
		}

		if state.Status == model.STATE_AUDIO_TRANSCRIPTED {
			logrusProc.Infof("DUBBING TASK RUNNING: %s; (8/11) %s", params.TaskName, "translating transcript")
			err = service.TranslateTranscript(bgCtx, params.TaskDir, params.TranscriptVttPath, params.TranscriptTranslatedPath)
			if err != nil {
				logrusProc.Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_TRANSCRIPT_TRANSLATED)
			if err != nil {
				logrusProc.Error(err)
				return
			}
		}

		if state.Status == model.STATE_TRANSCRIPT_TRANSLATED {
			logrusProc.Infof("DUBBING TASK RUNNING: %s; (9/11) %s", params.TaskName, "generating translated audio")
			err = service.GenerateVoice(bgCtx, params.TaskDir, params.TranscriptTranslatedPath, params.GeneratedSpeechDir, service.VoiceOpts{
				Name:  state.VoiceName,
				Rate:  state.VoiceRate,
				Pitch: state.VoicePitch,
			})
			if err != nil {
				logrusProc.Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_AUDIO_GENERATED)
			if err != nil {
				logrusProc.Error(err)
				return
			}
		}

		if state.Status == model.STATE_AUDIO_GENERATED {
			logrusProc.Infof("DUBBING TASK RUNNING: %s; (10/11) %s", params.TaskName, "adjusting audio speed")
			err = service.AdjustVoiceSpeed(bgCtx, params.TranscriptTranslatedPath, params.GeneratedSpeechDir, params.SpeechAdjustedDir)
			if err != nil {
				logrusProc.Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_AUDIO_ADJUSTED)
			if err != nil {
				logrusProc.Error(err)
				return
			}
		}

		if state.Status == model.STATE_AUDIO_ADJUSTED {
			logrusProc.Infof("DUBBING TASK RUNNING: %s; (11/11) %s", params.TaskName, "merge video with translatted audio")
			err = service.MergeVideoWithDubb(
				bgCtx,
				params.TranscriptTranslatedPath,
				params.SpeechAdjustedDir,
				params.InstrumentVideoPath,
				params.DubbedVideoPath,
			)
			if err != nil {
				logrusProc.Error(err)
				return
			}

			err = service.SaveStateStatus(bgCtx, params.TaskDir, &state, model.STATE_DUBBED_VIDEO_GENERATED)
			if err != nil {
				logrusProc.Error(err)
				return
			}
		}
	}()

	utils.Render(w, r, 200, state, nil)
}

func PatchUpdateTaskSetting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commonCtx := utils.GetCommonCtx(r)

	params := StartDubbTaskParams{}
	err := utils.BindJson(r, &params)
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 400, err)
		return
	}
	params.TaskName = chi.URLParam(r, "task_name")
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

	state.YoutubeUrl = params.YoutubeUrl
	state.VoiceName = params.VoiceName
	state.VoiceRate = params.VoiceRate
	state.VoicePitch = params.VoicePitch

	err = service.SaveState(ctx, params.TaskDir, state)
	if err != nil {
		logrus.WithContext(r.Context()).Error(err)
		utils.RenderError(w, r, 422, err)
		return
	}

	utils.Render(w, r, 200, nil, nil)
}
