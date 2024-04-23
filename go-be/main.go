package main

import (
	"bufio"
	"bytes"
	"context"
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
	"github.com/bregydoc/gtranslate"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

type (
	Config struct {
		GOOGLE_AI_STUDIO_KEY string
	}

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
	genaiGiminiProVision = &genai.GenerativeModel{}
	config               = Config{}
	baseDir              = "/root/shared"

	whisperModelPath = "/root/whisper.cpp/models/ggml-medium.en-q5_0.bin"

	// youtubeVideoURL = "https://www.youtube.com/watch?v=yDMZJ7LgrGY"
	// youtubeVideoURL = "https://www.youtube.com/watch?v=pQWd9YqvloU"
	// youtubeVideoURL = "https://www.youtube.com/watch?v=pjoQdz0nxf4"
	youtubeVideoURL = "https://www.youtube.com/watch?v=_Zc-NE8pmtg"
	taskName        = "test4-what-if-blackhole"

	taskDir   = fmt.Sprintf("%s/%s", baseDir, taskName)
	stateName = "state.json"
	statePath = fmt.Sprintf("%s/%s", taskDir, stateName)
)

func initialize() {
	logrus.SetReportCaller(true)

	godotenv.Load()

	config = Config{
		GOOGLE_AI_STUDIO_KEY: os.Getenv("GOOGLE_AI_STUDIO_KEY"),
	}

	genaiClient, err := genai.NewClient(
		context.TODO(),
		option.WithAPIKey(config.GOOGLE_AI_STUDIO_KEY),
	)
	if err != nil {
		logrus.Error(err)
		return
	}

	// genaiGiminiProVision = genaiClient.GenerativeModel("gemini-pro-vision")
	genaiGiminiProVision = genaiClient.GenerativeModel("gemini-pro")
	// genaiGiminiProVision.SafetySettings = []*genai.SafetySetting{
	// 	{
	// 		Category:  genai.HarmCategoryUnspecified,
	// 		Threshold: genai.HarmBlockNone,
	// 	},
	// }
}

func main() {
	initialize()

	state := TaskState{}

	//MARK: 0. Load existing state
	stateJson, err := os.ReadFile(statePath)
	if err == nil {
		json.Unmarshal(stateJson, &state)
	}

	//MARK: 1. Prepare directory for task
	logrus.Info("1. PREPARING DIR")
	cmd := exec.Command("mkdir", "-p", taskDir)
	_, err = cmd.Output()
	if err != nil {
		logrus.Error(err)
		return
	}

	//MARK: 2. Initializing state.json for state management
	logrus.Info("2. INITIALIZING TASK")
	if state.Status == "" {
		state.Status = "initialized"
		saveState(state)
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
	vocalRemoverExe := "/root/vocal-remover/inference.py"
	modelPath := "/root/vocal-remover/baseline.pth"
	logrus.Info("5. SEPARATE VIDEO VOCAL AND SOUND")
	if state.Status == "video_audio_generated" {
		cmd = exec.Command("python", vocalRemoverExe, "--input", rawVideoAudioPath, "-P", modelPath, "-o", taskDir)

		// stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()
		err = cmd.Start()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"cmd": cmd.String(),
			}).Error(err)
			return
		}

		streamStd(stderr)

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
	whisperExe := "/root/shared/main"
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

		streamCmdTranscript(stdout, stderr)

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

		// genAIPrompt := fmt.Sprintf(`
		// 	Help translate this transcript to Bahasa Indonesia.
		// 	I want the output in vtt. Give me all the result without break.

		// 	%s
		// `, string(vttContentByte))

		// resp, err := genaiGiminiProVision.GenerateContent(context.TODO(), genai.Text(genAIPrompt))
		// if err != nil {
		// 	logrus.Error(err)
		// 	return
		// }

		// if resp.Candidates == nil {
		// 	logrus.Error("nil genai candidates")
		// 	return
		// }

		// vttContent := ""
		// for _, candidate := range resp.Candidates {
		// 	vttContent += fmt.Sprintf("%+v", candidate.Content.Parts)
		// 	if vttContent != "" {
		// 		break
		// 	}
		// }
		// if vttContent == "" {
		// 	logrus.Error("empty string genai candidates")
		// 	return
		// }

		// vttContent = strings.TrimPrefix(vttContent, "[")
		// vttContent = strings.TrimSuffix(vttContent, "]")

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
			"[0]volume=30dB[video0]",
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

func saveState(state TaskState) {
	stateJson, _ := json.Marshal(state)

	err := os.WriteFile(statePath, stateJson, 0644)
	if err != nil {
		logrus.Errorf("%s - %s", "SAVE STATE FAIL", err.Error())
	}
}

// func formatDuration(d time.Duration) string {
// 	hours := int(d.Hours())
// 	minutes := int(d.Minutes()) % 60
// 	seconds := int(d.Seconds()) % 60
// 	milliseconds := int(d.Milliseconds()) % 1000
// 	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
// }

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

func streamCmd(stdout, stderr io.ReadCloser) {
	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
}

func streamStd(std io.ReadCloser) {
	scanner := bufio.NewScanner(std)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
}

func streamCmdTranscript(stdout, stderr io.ReadCloser) {
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
}
