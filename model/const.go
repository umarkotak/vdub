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

	STATE_IDX_MAP = map[string]STATE_DETAIL{
		STATE_INITIALIZED:                     {Idx: 0, StatusHuman: "Downloading Video"},
		STATE_VIDEO_DOWNLOADED:                {Idx: 1, StatusHuman: "Generating Audio"},
		STATE_VIDEO_AUDIO_GENERATED:           {Idx: 2, StatusHuman: "Separating Audio Instrument"},
		STATE_VIDEO_AUDIO_SEPARATED:           {Idx: 3, StatusHuman: "Converting Audio"},
		STATE_AUDIO_16KHZ_GENERATED:           {Idx: 4, StatusHuman: "Generating Video Instrument"},
		STATE_VIDEO_WITH_INSTRUMENT_GENERATED: {Idx: 5, StatusHuman: "Transcripting"},
		STATE_AUDIO_TRANSCRIPTED:              {Idx: 6, StatusHuman: "Translating"},
		STATE_TRANSCRIPT_TRANSLATED:           {Idx: 7, StatusHuman: "Generating Voice"},
		STATE_AUDIO_GENERATED:                 {Idx: 8, StatusHuman: "Adjusting Voice"},
		STATE_AUDIO_ADJUSTED:                  {Idx: 9, StatusHuman: "Merging Video"},
		STATE_DUBBED_VIDEO_GENERATED:          {Idx: 10, StatusHuman: "Completed"},
	}
)

type (
	STATE_DETAIL struct {
		Idx         int
		StatusHuman string
	}
)
