package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	astisub "github.com/asticode/go-astisub"
	"github.com/bregydoc/gtranslate"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/datastore"
	"github.com/umarkotak/vdub-go/handler"
	"github.com/umarkotak/vdub-go/model"
	"github.com/umarkotak/vdub-go/service"
	"github.com/umarkotak/vdub-go/utils"
)

var (
	whisperModelPath = "/root/whisper.cpp/models/ggml-medium.en-q5_0.bin"

	youtubeVideoURL = "https://www.youtube.com/watch?v=SLsElgfhZtM"
	taskName        = "p-1-avicen-1"

	taskDir   = fmt.Sprintf("%s/%s", config.Get().BaseDir, taskName)
	stateName = "state.json"
	statePath = fmt.Sprintf("%s/%s", taskDir, stateName)
)

func initialize() {
	logrus.SetReportCaller(true)
	config.InitConfig()
	datastore.InitDataStore()
}

func main() {
	var cmd *exec.Cmd

	initialize()

	r := chi.NewRouter()

	r.Use(
		chiMiddleware.RequestID,
		chiMiddleware.RealIP,
		chiMiddleware.Recoverer,
	)

	handler.Initialize()

	r.Get("/", handler.Ping)

	r.Post("/vdub/api/dubb/start", handler.PostStartDubbTask)
	r.Get("/vdub/api/dubb/{task_name}/status", handler.GetTaskStatus)

	port := ":29000"
	logrus.Infof("Listening on port %s", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		logrus.Fatal(err)
	}

	state, err := service.GetState(context.TODO(), taskDir)
	if err != nil {
		logrus.Error(err)
		return
	}

	//MARK: 3. Download youtube video
	var stderr bytes.Buffer
	rawVideoName := "raw_video.mp4"
	rawVideoPath := fmt.Sprintf("%s/%s", taskDir, rawVideoName)
	logrus.Info("3. DOWNLOADING VIDEO")
	if state.Status == "initialized" {
		cmd = exec.Command("yt-dlp", "-S", "ext", "-o", rawVideoPath, youtubeVideoURL)
		cmd.Stderr = &stderr
		_, err = cmd.Output()
		if err != nil {
			logrus.Errorf("%v - %v", err.Error(), stderr.String())
			return
		}

		state.Status = "video_downloaded"
		saveState(state)
	}

	//MARK: 4. Generate video audio in wav format
	rawVideoAudioName := "raw_video_audio.wav"
	rawVideoAudioPath := fmt.Sprintf("%s/%s", taskDir, rawVideoAudioName)
	logrus.Info("4. GENERATING AUDIO")
	if state.Status == "video_downloaded" {
		cmd = exec.Command("ffmpeg", "-i", rawVideoPath, "-vn", "-acodec", "pcm_s16le", "-ar", "44100", "-ac", "2", rawVideoAudioPath)
		cmd.Stderr = &stderr
		_, err = cmd.Output()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"rawVideoPath":      rawVideoPath,
				"rawVideoAudioPath": rawVideoAudioPath,
			}).Errorf("%v - %v", err.Error(), stderr.String())
			return
		}

		state.Status = "video_audio_generated"
		saveState(state)
	}

	//MARK: 5. Separate video vocal and sound
	logrus.Info("5. SEPARATE VIDEO VOCAL AND SOUND")
	if state.Status == "video_audio_generated" {
		cmd = exec.Command(
			"python", config.Get().VocalRemoverPy,
			"--input", rawVideoAudioPath,
			"-P", config.Get().VocalRemoverModelPath,
			"-o", taskDir,
		)

		// stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		err = cmd.Start()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"cmd": cmd.String(),
			}).Error(err)
			return
		}

		utils.StreamStd(stderr)

		err = cmd.Wait()
		fmt.Printf("\n")
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"cmd": cmd.String(),
			}).Error(err)
			return
		}

		state.Status = "video_audio_separated"
		saveState(state)
	}

	//MARK: 6. Convert vocal sound to 16KHz
	vocalPath := fmt.Sprintf("%s_Vocals.wav", strings.TrimSuffix(rawVideoAudioPath, ".wav"))
	vocal16KHzName := "raw_video_audio_Vocals_16KHz.wav"
	vocal16KHzPath := fmt.Sprintf("%s/%s", taskDir, vocal16KHzName)
	logrus.Info("6. CONVERTING VOCAL SOUND TO 16KHz")
	if state.Status == "video_audio_separated" {
		cmd = exec.Command("ffmpeg", "-i", vocalPath, "-acodec", "pcm_s16le", "-ac", "1", "-ar", "16000", vocal16KHzPath)
		cmd.Stderr = &stderr
		_, err = cmd.Output()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"vocalPath": vocalPath,
			}).Errorf("%v - %v", err.Error(), stderr.String())
			return
		}

		state.Status = "audio_16khz_generated"
		saveState(state)
	}

	//MARK: 7. Transcript vocal
	logrus.Info("7. TRANSCRIPTING VOCAL TO TEXT")
	transcriptPath := fmt.Sprintf("%s/%s", taskDir, "transcript")
	whisperExe := "/root/whisper.cpp/main"
	if state.Status == "audio_16khz_generated" {
		// ./main -m /root/shared/models/ggml-base.en.bin -l en -bs 7 -bo 7 -wt 0.04 -ovtt /root/shared/test3-kurzgesagt/transcript-t /root/shared/test3-kurzgesagt/raw_video_audio_Vocals_16KHz.wav
		// cmdTranscript := exec.Command(whisperExe, "-m", whisperModelPath, "-oj", "-otxt", "-ovtt", "-of", transcriptPath, vocal16KHzPath)
		// cmdTranscript := exec.Command(whisperExe, "-m", whisperModelPath, "-ovtt", "-of", transcriptPath, vocal16KHzPath)
		cmdTranscript := exec.Command(
			whisperExe,
			"-m", whisperModelPath,
			"-ovtt",
			"--beam-size", "6",
			"--entropy-thold", "2.8",
			"--max-context", "128",
			"-of", transcriptPath,
			vocal16KHzPath,
		)

		stdout, _ := cmdTranscript.StdoutPipe()
		stderr, _ := cmdTranscript.StderrPipe()
		err = cmdTranscript.Start()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"cmdTranscript": cmdTranscript.String(),
			}).Error(err)
		}

		utils.StreamCmdTranscript(stdout, stderr)

		err = cmdTranscript.Wait()
		fmt.Printf("\n")
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"cmdTranscript": cmdTranscript.String(),
			}).Error(err)
		}

		state.Status = "audio_transcripted"
		saveState(state)
	}

	//MARK: 8. Translate transcript.vtt -> transcript_translated.vtt
	// As of now this will be done manually using google genAI and manually write the file
	logrus.Info("8. TRANSLATE TRANSCRIPT TO TARGET LANGUAGE")
	transcriptVttPath := fmt.Sprintf("%s/%s", taskDir, "transcript.vtt")
	transcriptTranslatedPath := fmt.Sprintf("%s/%s", taskDir, "transcript_translated.vtt")
	if state.Status == "audio_transcripted" {
		// logrus.Info("TODO: manual transcript on GenAI")
		// TODO: Implement logic
		vttContentByte, err := os.ReadFile(transcriptVttPath)
		if err != nil {
			logrus.Error(err)
			return
		}
		vttContent := string(vttContentByte)

		subObj, _ := astisub.OpenFile(transcriptVttPath)

		bar := progressbar.Default(int64(len(subObj.Items)), "Translating")
		for _, subItem := range subObj.Items {
			translated, err := gtranslate.TranslateWithParams(
				subItem.String(),
				gtranslate.TranslationParams{
					From: "en", To: "id",
				},
			)
			if err != nil {
				logrus.Error(err)
				return
			}

			vttContent = strings.ReplaceAll(vttContent, subItem.String(), translated)
			bar.Add(1)
		}

		err = os.WriteFile(transcriptTranslatedPath, []byte(vttContent), 0644)
		if err != nil {
			logrus.Error(err)
			return
		}

		state.Status = "transcript_translated"
		saveState(state)
	}

	//MARK: 9. Generate audio for the transcript
	logrus.Info("9. GENERATE AUDIO FOR TRANSLATED TRANSCRIPT")
	generatedSpeechDir := fmt.Sprintf("%s/generated_speech", taskDir)
	if state.Status == "transcript_translated" {
		cmd = exec.Command("mkdir", "-p", generatedSpeechDir)
		_, err = cmd.Output()
		if err != nil {
			logrus.Error(err)
			return
		}

		subObj, _ := astisub.OpenFile(transcriptTranslatedPath)
		bar := progressbar.Default(int64(len(subObj.Items)), "Generating Audio")
		for idx, subItem := range subObj.Items {
			// fmt.Printf("%v - [%v --> %v] %v\n", idx, subItem.StartAt, subItem.EndAt, subItem.String())

			genSpeechPath := fmt.Sprintf("%s/%v.wav", generatedSpeechDir, idx)
			cmd = exec.Command(
				"edge-tts",
				"--text", fmt.Sprintf("\"%s\"", subItem.String()),
				"--write-media", genSpeechPath,
				"-v", "id-ID-ArdiNeural",
				"--rate=-10%",
				"--pitch=-10Hz",
			)
			cmd.Stderr = &stderr
			_, err = cmd.Output()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"genSpeechPath": genSpeechPath,
					"cmd":           cmd.String(),
				}).Errorf("%v - %v", err.Error(), stderr.String())
				return
			}
			bar.Add(1)
		}

		state.Status = "audio_generated"
		saveState(state)
	}

	//MARK: 9.2 TODO: Do a voice synthesize

	//MARK: 10. Merge video with instrument only audio
	logrus.Info("10. MERGE VIDEO WITH INSTRUMENT ONLY")
	instrumentPath := fmt.Sprintf("%s_Instruments.wav", strings.TrimSuffix(rawVideoAudioPath, ".wav"))
	instrumentVideoPath := fmt.Sprintf("%s/%s", taskDir, "instrument_video.mp4")
	if state.Status == "audio_generated" {
		cmd = exec.Command("ffmpeg", "-i", rawVideoPath, "-i", instrumentPath, "-c:v", "copy", "-c:a", "aac", "-map", "0:v:0", "-map", "1:a:0", instrumentVideoPath)
		cmd.Stderr = &stderr
		_, err = cmd.Output()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"vocalPath": vocalPath,
			}).Errorf("%v - %v", err.Error(), stderr.String())
			return
		}

		state.Status = "video_with_instrument_generated"
		saveState(state)
	}

	//MARK: 11. Adjust generated audio duration
	logrus.Info("11. ADJUST AUDIO DURATION")
	adjustedSpeechDir := fmt.Sprintf("%s/adjusted_speech", taskDir)
	if state.Status == "video_with_instrument_generated" {
		cmd = exec.Command("mkdir", "-p", adjustedSpeechDir)
		_, err = cmd.Output()
		if err != nil {
			logrus.Error(err)
			return
		}

		subObj, _ := astisub.OpenFile(transcriptTranslatedPath)
		bar := progressbar.Default(int64(len(subObj.Items)), "Adjusting Audio")
		for idx, subItem := range subObj.Items {
			genSpeechPath := fmt.Sprintf("%s/%v.wav", generatedSpeechDir, idx)
			adjustedSpeechPath := fmt.Sprintf("%s/%v.wav", adjustedSpeechDir, idx)

			originalDuration := subItem.EndAt - subItem.StartAt
			translatedDuration, _ := getWavDuration(genSpeechPath)
			aTempo := toFixed(translatedDuration.Seconds()/originalDuration.Seconds(), 6)
			if aTempo < 1 {
				aTempo = 1.1
			} else if aTempo > 100 {
				aTempo = 100
			}

			cmd = exec.Command(
				"ffmpeg", "-i", genSpeechPath, "-codec:a", "libmp3lame", "-filter:a", fmt.Sprintf("atempo=%v", aTempo), "-b:a", "320k", adjustedSpeechPath,
			)
			cmd.Stderr = &stderr
			_, err = cmd.Output()
			if err != nil {
				logrus.Errorf("%v - %v", err.Error(), stderr.String())
				return
			}
			bar.Add(1)
		}

		state.Status = "audio_adjusted"
		saveState(state)
	}

	//MARK: 12. Merge instrument video with generated audio
	logrus.Info("12. MERGE VIDEO WITH GENERATED AUDIO")
	dubbedVideoPath := fmt.Sprintf("%s/%s", taskDir, "dubbed_video.mp4")
	if state.Status == "audio_adjusted" {
		subObj, _ := astisub.OpenFile(transcriptTranslatedPath)

		ffmpegArgs := []string{
			"-i", instrumentVideoPath,
		}

		for idx := range subObj.Items {
			adjustedSpeechPath := fmt.Sprintf("%s/%v.wav", adjustedSpeechDir, idx)
			ffmpegArgs = append(ffmpegArgs, "-i", adjustedSpeechPath)
		}

		filterComplexes := []string{
			"[0]volume=40dB[video0]",
		}
		filterComplexCloser := ""
		for idx, subItem := range subObj.Items {
			audioIdx := fmt.Sprintf("[audio%v]", idx)

			// logrus.Infof("DEBUG DELAY %v: %v - %v", idx, subItem.StartAt.Milliseconds(), subItem.EndAt.Milliseconds())

			filter := fmt.Sprintf(
				"[%v]volume=40dB,adelay=%v%s",
				idx+1,
				subItem.StartAt.Milliseconds(),
				audioIdx,
			)
			filterComplexes = append(filterComplexes, filter)

			filterComplexCloser += audioIdx
		}
		filterComplexCloserFormatted := fmt.Sprintf("[video0]%samix=%v", filterComplexCloser, len(subObj.Items)+1)
		filterComplexes = append(filterComplexes, filterComplexCloserFormatted)

		ffmpegArgs = append(ffmpegArgs, "-filter_complex", fmt.Sprintf("%s", strings.Join(filterComplexes, ";")))

		ffmpegArgs = append(
			ffmpegArgs,
			"-c:v", "copy",
			dubbedVideoPath,
		)

		cmd = exec.Command("ffmpeg", ffmpegArgs...)
		logrus.Infof("EXECUTING MERGING DUB CMD: %+v", cmd.String())
		cmd.Stderr = &stderr
		_, err = cmd.Output()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"ffmpegArgs": ffmpegArgs,
			}).Errorf("%v - %v", err.Error(), stderr.String())
			return
		}

		state.Status = "dubbed_video_generated"
		saveState(state)
	}

	fmt.Printf("TASK [%s] DONE\n", taskName)
}

func saveState(state model.TaskState) {
	stateJson, _ := json.Marshal(state)

	err := os.WriteFile(statePath, stateJson, 0644)
	if err != nil {
		logrus.Errorf("%s - %s", "SAVE STATE FAIL", err.Error())
	}
}

func getWavDuration(filename string) (time.Duration, error) {
	cmd := exec.Command("ffprobe", "-i", filename, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		logrus.Errorf("%v - %v", err.Error(), stderr.String())
	}

	seconds, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		logrus.Error(err)
	}

	// Convert seconds to time.Duration (in nanoseconds)
	duration := time.Duration(seconds * float64(time.Second))

	return duration, nil
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
