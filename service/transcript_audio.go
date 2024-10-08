package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/asticode/go-astisub"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/config"
	"github.com/umarkotak/vdub-go/utils"
)

const (
	BEAM_SIZE      = "6"
	ENTHROPY_THOLD = "2.8"
	MAX_CONTEXT    = "128"
)

func TranscriptAudio(ctx context.Context, audioPath, transcriptPath string) error {
	cmdTranscript := exec.Command(
		config.Get().WhisperBinary,
		"-m", config.Get().WhisperModelPath,
		"-ovtt",
		"--beam-size", BEAM_SIZE,
		"--entropy-thold", ENTHROPY_THOLD,
		"--max-context", MAX_CONTEXT,
		"-of", transcriptPath,
		"--translate",
		audioPath,
	)

	stdout, _ := cmdTranscript.StdoutPipe()
	stderr, _ := cmdTranscript.StderrPipe()
	err := cmdTranscript.Start()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"cmd": cmdTranscript.String(),
		}).Error(err)
	}

	utils.StreamCmdTranscript(stdout, stderr)

	err = cmdTranscript.Wait()
	fmt.Printf("\n")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cmd": cmdTranscript.String(),
		}).Error(err)
	}

	return nil
}

func TranscriptAudioWithDiarization(ctx context.Context, taskDir, audioPath, transcriptPath, segmentedSpeechDir string) error {
	cmd := exec.Command("mkdir", "-p", segmentedSpeechDir)
	_, err := cmd.Output()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	diarizationVtt, err := astisub.OpenFile(fmt.Sprintf("%s/diarization.vtt", taskDir))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	bar := progressbar.Default(int64(len(diarizationVtt.Items)), "Translating")

	for idx, subItem := range diarizationVtt.Items {
		cmd = exec.Command(
			"ffmpeg", "-y",
			"-i", audioPath,
			"-ss", fmt.Sprintf("%v", subItem.StartAt.Seconds()),
			"-to", fmt.Sprintf("%v", subItem.EndAt.Seconds()),
			"-c", "copy",
			fmt.Sprintf("%s/%v.wav", segmentedSpeechDir, idx),
		)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		_, err := cmd.Output()
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"audio_path": audioPath,
				"std_err":    stderr.String(),
			}).Error(err)
			return err
		}

		segmentedVoicePath := fmt.Sprintf("%s/segmented_speech/%v.wav", taskDir, idx)
		cmd = exec.Command(
			config.Get().WhisperBinary,
			"--no-prints", "--output-txt",
			"--model", config.Get().WhisperModelPath,
			"--translate", segmentedVoicePath,
		)
		_, err = cmd.Output()
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return err
		}

		textContent, _ := readFileAsOneLine(fmt.Sprintf("%s.txt", segmentedVoicePath))
		if textContent == "" {
			textContent = "missing text"
		}
		subItem.Lines[0].Items[0].Text = textContent

		bar.Add(1)
	}

	diarizationVtt.Write(utils.GenTranscriptVttPath(taskDir))

	return nil
}

func readFileAsOneLine(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	return strings.TrimSpace(strings.ReplaceAll(string(content), "\n", " ")), err
}
