package utils

import (
	"fmt"

	"github.com/umarkotak/vdub-go/config"
)

func GenTaskName(username, taskName string) string {
	if username == "" {
		username = "public"
	}
	return fmt.Sprintf("task-%s-%s", username, taskName)
}

func GenTaskDir(taskName string) string {
	return fmt.Sprintf("%s/%s", config.Get().BaseDir, taskName)
}

func GenVideoScreenshotPath(taskDir string) string {
	return fmt.Sprintf("%s/%s", taskDir, "video_snapshot.jpg")
}

func GenRawVideoPath(taskDir, rawVideoName string) string {
	return fmt.Sprintf("%s/%s", taskDir, rawVideoName)
}

func GenRawVideoAudioPath(taskDir, rawVideoAudioName string) string {
	return fmt.Sprintf("%s/%s", taskDir, rawVideoAudioName)
}

func GenAudioInstrumentPath(taskDir string) string {
	return fmt.Sprintf("%s/raw_video_audio_(Instrumental)_UVR-MDX-NET-Voc_FT.wav", taskDir)
}

func GenAudioVocalPath(taskDir string) string {
	return fmt.Sprintf("%s/raw_video_audio_(Vocals)_UVR-MDX-NET-Voc_FT.wav", taskDir)
}

func GenVocal16KHzPath(taskDir string) string {
	return fmt.Sprintf("%s/raw_video_audio_(Vocals)_UVR-MDX-NET-Voc_FT_16KHz.wav", taskDir)
}

func GenInstrumentVideoPath(taskDir string) string {
	return fmt.Sprintf("%s/%s", taskDir, "instrument_video.mp4")
}

func GenTranscriptPath(taskDir string) string {
	return fmt.Sprintf("%s/%s", taskDir, "transcript")
}

func GenTranscriptVttPath(taskDir string) string {
	return fmt.Sprintf("%s/%s", taskDir, "transcript.vtt")
}

func GenTranscriptTranslatedPath(taskDir string) string {
	return fmt.Sprintf("%s/%s", taskDir, "transcript_translated.vtt")
}

func GenTranscriptTranslatedTestPath(taskDir string) string {
	return fmt.Sprintf("%s/%s", taskDir, "transcript_translated_test.vtt")
}

func GenGeneratedSpeechDir(taskDir string) string {
	return fmt.Sprintf("%s/generated_speech", taskDir)
}

func GenSegmentedSpeechDir(taskDir string) string {
	return fmt.Sprintf("%s/segmented_speech", taskDir)
}

func GenSpeechAdjustedDir(taskDir string) string {
	return fmt.Sprintf("%s/adjusted_speech", taskDir)
}

func GenDubbedVideoPath(taskDir string) string {
	return fmt.Sprintf("%s/%s", taskDir, "dubbed_video.mp4")
}
