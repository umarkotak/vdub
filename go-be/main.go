package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	astisub "github.com/asticode/go-astisub"
	"github.com/sirupsen/logrus"
)

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

var (
	baseDir = "/root/shared"

	taskName = "test1"
	taskDir  = fmt.Sprintf("%s/%s", baseDir, taskName)

	stateName = "state.json"
	statePath = fmt.Sprintf("%s/%s", taskDir, stateName)

	youtubeVideoURL = "https://www.youtube.com/watch?v=yDMZJ7LgrGY"
)

func main() {
	logrus.SetReportCaller(true)

	state := TaskState{}

	// 0. Load existing state
	stateJson, err := os.ReadFile(statePath)
	if err == nil {
		json.Unmarshal(stateJson, &state)
	}

	// 1. Prepare directory for task
	logrus.Info("1. PREPARING DIR")
	cmd := exec.Command("mkdir", "-p", taskDir)
	_, err = cmd.Output()
	if err != nil {
		logrus.Error(err)
		return
	}

	// 2. Initializing state.json for state management
	logrus.Info("2. INITIALIZING TASK")
	if state.Status == "" {
		state.Status = "initialized"
		saveState(state)
	}

	// 3. Download youtube video
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

	// 4. Generate video audio in wav format
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

	// 5. Separate video vocal and sound
	vocalRemoverExe := "/root/vocal-remover/inference.py"
	modelPath := "/root/shared/baseline.pth"
	logrus.Info("5. SEPARATE VIDEO VOCAL AND SOUND")
	if state.Status == "video_audio_generated" {
		cmd = exec.Command("python", vocalRemoverExe, "--input", rawVideoAudioPath, "-P", modelPath, "-o", taskDir)
		cmd.Stderr = &stderr
		_, err = cmd.Output()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"rawVideoAudioPath": rawVideoAudioPath,
			}).Errorf("%v - %v", err.Error(), stderr.String())
			return
		}

		state.Status = "video_audio_separated"
		saveState(state)
	}

	// 6. Convert vocal sound to 16KHz
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

	// 7. Transcript vocal
	logrus.Info("7. TRANSCRIPTING VOCAL TO TEXT")
	whisperModelPath := "/root/shared/models/ggml-base.en.bin"
	transcriptPath := fmt.Sprintf("%s/%s", taskDir, "transcript")
	whisperExe := "/root/shared/main"
	if state.Status == "audio_16khz_generated" {
		// ./main
		cmdTranscript := exec.Command(whisperExe, "-m", whisperModelPath, "-oj", "-otxt", "-ovtt", "-of", transcriptPath, vocal16KHzPath)
		// cmd = exec.Command(whisperExe, "-m", whisperModelPath, vocal16KHzPath)

		stdout, _ := cmdTranscript.StdoutPipe()
		stderr, _ := cmdTranscript.StderrPipe()
		_ = cmdTranscript.Start()

		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			m := scanner.Text()

			prefix := ""
			if strings.Contains(m, "[") {
				prefix = "\n"
			} else {
				prefix = " "
			}

			fmt.Printf("%s%s", prefix, m)
		}
		_ = cmdTranscript.Wait()
		fmt.Printf("\n")

		state.Status = "audio_transcripted"
		saveState(state)
	}

	// 8. TODO: Translate transcript.vtt -> transcript_translated.vtt
	// As of now this will be done manually using google genAI and manually write the file
	logrus.Info("8. TRANSLATE TRANSCRIPT TO TARGET LANGUAGE")
	// transcriptPath := fmt.Sprintf("%s/%s", taskDir, "transcript.vtt")
	transcriptTranslatedPath := fmt.Sprintf("%s/%s", taskDir, "transcript_translated.vtt")
	if state.Status == "audio_transcripted" {
		logrus.Info("TODO: manual transcript on GenAI")
		// TODO: Implement logic

		// state.Status = "transcript_translated"
		// saveState(state)
	}

	// 9. Generate audio for the transcript
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
		for idx, subItem := range subObj.Items {
			fmt.Printf("%v - [%v --> %v] %v\n", idx, subItem.StartAt, subItem.EndAt, subItem.String())

			genSpeechPath := fmt.Sprintf("%s/%v.wav", generatedSpeechDir, idx)
			cmd = exec.Command("edge-tts", "--text", subItem.String(), "--write-media", genSpeechPath, "-v", "id-ID-ArdiNeural", "--rate=-10%", "--pitch=-10Hz")
			cmd.Stderr = &stderr
			_, err = cmd.Output()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"genSpeechPath": genSpeechPath,
				}).Errorf("%v - %v", err.Error(), stderr.String())
				return
			}
		}

		state.Status = "audio_generated"
		saveState(state)
	}

	// 9.2 TODO: Do a voice synthesize

	// 10. Merge video with instrument only audio
	logrus.Info("9. MERGE VIDEO WITH INSTRUMENT ONLY")
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

	// 11. Merge instrument video with generated audio
	dubbedVideoPath := fmt.Sprintf("%s/%s", taskDir, "dubbed_video.mp4")
	adjustedSpeechDir := fmt.Sprintf("%s/adjusted_speech", taskDir)
	if state.Status == "video_with_instrument_generated" {
		cmd = exec.Command("mkdir", "-p", adjustedSpeechDir)
		_, err = cmd.Output()
		if err != nil {
			logrus.Error(err)
			return
		}

		ffmpegArgs := []string{
			"-i", instrumentVideoPath,
		}

		subObj, _ := astisub.OpenFile(transcriptTranslatedPath)
		for idx, subItem := range subObj.Items {
			genSpeechPath := fmt.Sprintf("%s/%v.wav", generatedSpeechDir, idx)
			adjustedSpeechPath := fmt.Sprintf("%s/%v.wav", adjustedSpeechDir, idx)

			originalDuration := subItem.EndAt - subItem.StartAt
			translatedDuration, _ := getWavDuration(genSpeechPath)
			aTempo := toFixed(translatedDuration.Seconds()/originalDuration.Seconds(), 6)
			// logrus.Infof(
			// 	"Diff Duration %v: %v - %v = %v",
			// 	idx, originalDuration.String(), translatedDuration.String(), aTempo,
			// )
			if aTempo < 0.5 {
				aTempo = 0.5
			} else if aTempo > 100 {
				aTempo = 100
			}

			// ffmpeg -i input.wav -codec:a libmp3lame -filter:a "atempo=0.9992323" -b:a 320K output.mp3
			cmd = exec.Command(
				"ffmpeg", "-i", genSpeechPath, "-codec:a", "libmp3lame", "-filter:a", fmt.Sprintf("atempo=%v", aTempo), "-b:a", "320k", adjustedSpeechPath,
			)
			cmd.Stderr = &stderr
			_, err = cmd.Output()
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"ffmpegArgs": ffmpegArgs,
				}).Errorf("%v - %v", err.Error(), stderr.String())
				return
			}
		}

		for idx := range subObj.Items {
			adjustedSpeechPath := fmt.Sprintf("%s/%v.wav", adjustedSpeechDir, idx)
			ffmpegArgs = append(ffmpegArgs, "-i", adjustedSpeechPath)
		}

		filterComplexes := []string{}
		filterComplexCloser := ""
		for idx, subItem := range subObj.Items {
			audioIdx := fmt.Sprintf("[audio%v]", idx)

			logrus.Infof("DEBUG DELAY %v: %v - %v", idx, subItem.StartAt.Milliseconds(), subItem.EndAt.Milliseconds())

			filter := fmt.Sprintf(
				"[%v]adelay=%v%s",
				idx+1,
				subItem.StartAt.Milliseconds(),
				audioIdx,
			)
			filterComplexes = append(filterComplexes, filter)

			filterComplexCloser += audioIdx
		}
		filterComplexCloserFormatted := fmt.Sprintf("[0]%samix=%v", filterComplexCloser, len(subObj.Items)+1)
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

		// state.Status = "dubbed_video_generated"
		// saveState(state)
	}

	fmt.Printf("TASK [%s] DONE\n", taskName)
}

func saveState(state TaskState) {
	stateJson, _ := json.Marshal(state)

	err := os.WriteFile(statePath, stateJson, 0644)
	if err != nil {
		logrus.Errorf("%s - %s", "SAVE STATE FAIL", err.Error())
	}
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}

func getWavDuration(filename string) (time.Duration, error) {
	// ffprobe -i /root/shared/test1/generated_speech/32.wav -show_entries format=duration -v quiet -of csv="p=0"

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