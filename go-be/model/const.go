package model

const (
	STATE_INITIALIZED                     = "initialized"
	STATE_VIDEO_DOWNLOADED                = "video_downloaded"
	STATE_VIDEO_AUDIO_GENERATED           = "video_audio_generated"
	STATE_VIDEO_AUDIO_SEPARATED           = "video_audio_separated"
	STATE_AUDIO_16KHZ_GENERATED           = "audio_16khz_generated"
	STATE_VIDEO_WITH_INSTRUMENT_GENERATED = "video_with_instrument_generated"
	STATE_AUDIO_TRANSCRIPTED              = "audio_transcripted"
	STATE_TRANSCRIPT_TRANSLATED           = "transcript_translated"
	STATE_AUDIO_GENERATED                 = "audio_generated"
	STATE_AUDIO_ADJUSTED                  = "audio_adjusted"
	STATE_DUBBED_VIDEO_GENERATED          = "dubbed_video_generated"
)

var (
	STATE_IDX_ARR = []string{
		STATE_INITIALIZED,
		STATE_VIDEO_DOWNLOADED,
		STATE_VIDEO_AUDIO_GENERATED,
		STATE_VIDEO_AUDIO_SEPARATED,
		STATE_AUDIO_16KHZ_GENERATED,
		STATE_VIDEO_WITH_INSTRUMENT_GENERATED,
		STATE_AUDIO_TRANSCRIPTED,
		STATE_TRANSCRIPT_TRANSLATED,
		STATE_AUDIO_GENERATED,
		STATE_AUDIO_ADJUSTED,
		STATE_DUBBED_VIDEO_GENERATED,
	}

	STATE_IDX_MAP = map[string]int{
		STATE_INITIALIZED:                     0,
		STATE_VIDEO_DOWNLOADED:                1,
		STATE_VIDEO_AUDIO_GENERATED:           2,
		STATE_VIDEO_AUDIO_SEPARATED:           3,
		STATE_AUDIO_16KHZ_GENERATED:           4,
		STATE_VIDEO_WITH_INSTRUMENT_GENERATED: 5,
		STATE_AUDIO_TRANSCRIPTED:              6,
		STATE_TRANSCRIPT_TRANSLATED:           7,
		STATE_AUDIO_GENERATED:                 8,
		STATE_AUDIO_ADJUSTED:                  9,
		STATE_DUBBED_VIDEO_GENERATED:          10,
	}
)
