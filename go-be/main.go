package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
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
		// TODO: Implement logic

		state.Status = "transcript_translated"
		saveState(state)
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
		// ffmpeg -i testvideo.mp4 -i test_Instruments.wav -c:v copy -c:a aac -map 0:v:0 -map 1:a:0 testvideo2.mp4
		// fmt.Printf("INSTRUMENT PATH: %s\n", instrumentPath)
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
	if state.Status == "video_with_instrument_generated" {
		ffmpegArgs := []string{
			"-i", instrumentVideoPath,
		}

		subObj, _ := astisub.OpenFile(transcriptTranslatedPath)
		for idx := range subObj.Items {
			genSpeechPath := fmt.Sprintf("%s/%v.wav", generatedSpeechDir, idx)
			ffmpegArgs = append(ffmpegArgs, "-i", genSpeechPath)
		}

		filterComplexes := []string{}
		filterComplexCloser := ""
		for idx, subItem := range subObj.Items {
			audioIdx := fmt.Sprintf("[audio%v]", idx)

			// [1:a]atrim=start=3.040,end=6.920,asetpts=PTS-STARTPTS[audio1]
			// filter := fmt.Sprintf(
			// 	"[%v:a]atrim=start=%v,end=%v,asetpts=PTS-STARTPTS%s",
			// 	idx,
			// 	formatDuration(subItem.StartAt),
			// 	formatDuration(subItem.EndAt),
			// 	audioIdx,
			// )
			// filterComplexes = append(filterComplexes, filter)

			// 	ffmpeg -i in.avi -i audio.wav -filter_complex
			// "[0:a]adelay=62000|62000[aud];[0][aud]amix" -c:v copy out.avi
			filter := fmt.Sprintf(
				"[%v:a]adelay=%v|%v%s",
				idx,
				subItem.StartAt.Milliseconds(),
				subItem.StartAt.Milliseconds(),
				audioIdx,
			)
			filterComplexes = append(filterComplexes, filter)

			filterComplexCloser += audioIdx
		}
		// filterComplexCloser += fmt.Sprintf("concat=n=%v:v=0:a=1[out]", len(subObj.Items))
		filterComplexCloserFormatted := fmt.Sprintf("[0]%samix=%v", filterComplexCloser, len(subObj.Items)+1)
		// filterComplexes = append(filterComplexes, filterComplexCloser)
		filterComplexes = append(filterComplexes, filterComplexCloserFormatted)

		ffmpegArgs = append(ffmpegArgs, "-filter_complex", fmt.Sprintf("\"%s\"", strings.Join(filterComplexes, ";")))

		ffmpegArgs = append(
			ffmpegArgs,
			// "-map", "0:v",
			// "-map", "\"[out]\"",
			"-c:v", "copy",
			// "-c:a", "aac",
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

// [00:00:00.000 --> 00:00:03.040] Hai kamu, senang sekali kamu bisa bergabung dengan kami.
// [00:00:03.040 --> 00:00:06.920] Kami ingin menceritakan sesuatu yang mengubah seni Kotska selamanya.
// [00:00:06.920 --> 00:00:11.600] Seni Kotska dimulai sebagai proyek hobi berskala kecil, tetapi membuat video sains animasi
// [00:00:11.600 --> 00:00:14.920] yang gratis untuk semua orang tidaklah menghasilkan uang.
// [00:00:14.920 --> 00:00:16.600] Sialan, kenyataan.
// transcripts := []Transcript{
// 	{
// 		Start:   "",
// 		Stop:    "",
// 		Text:    "",
// 		AudFile: "",
// 	},
// }

// ffmpegCmd := `ffmpeg`
// ffmpegCmd += ` -i testvideo2.mp4`
// ffmpegCmd += ` ` + strings.Join(
// 	[]string{
// 		"edge-tts-voice-1.wav",
// 		"edge-tts-voice-2.wav",
// 		"edge-tts-voice-3.wav",
// 		"edge-tts-voice-4.wav",
// 		"edge-tts-voice-5.wav",
// 		"edge-tts-voice-6.wav",
// 	}, " ",
// )
// ffmpegCmd += ` -filter_complex`
// ffmpegCmd += ` "`
// ffmpegCmd += strings.Join([]string{
// 	"[0:a]atrim=end=3.040,asetpts=PTS-STARTPTS[audio0]",
// 	"[1:a]atrim=start=3.040,end=6.920,asetpts=PTS-STARTPTS[audio1]",
// 	"[2:a]atrim=start=6.920,end=11.600,asetpts=PTS-STARTPTS[audio2]",
// 	"[audio0][audio1][audio2]concat=n=3:v=0:a=1[out]",
// }, ";")
// ffmpegCmd += `"`

// ffmpegCmd += ` -map 0:v -map "[out]" -c:v copy -c:a aac output.mp4`

// fmt.Println(ffmpegCmd)

// ffmpeg -i video.mp4 -i audio1.wav -i audio2.wav -i audio3.wav \
// -filter_complex "[0:a]atrim=end=3.040,asetpts=PTS-STARTPTS[audio0]; \
// [1:a]atrim=start=3.040,end=6.920,asetpts=PTS-STARTPTS[audio1]; \
// [2:a]atrim=start=6.920,end=11.600,asetpts=PTS-STARTPTS[audio2]; \
// [audio0][audio1][audio2]concat=n=3:v=0:a=1[out]" \
// -map 0:v -map "[out]" -c:v copy -c:a aac output.mp4
