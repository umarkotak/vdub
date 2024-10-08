from pyannote.audio import Pipeline
# from pyannote.audio.pipelines.utils.hook import Hook, ProgressHook, TimingHook
import torch
import argparse

def diarize_audio_to_vtt(file_path, output_path="diarization.vtt", auth_token="your_auth_token_here"):
    # Load the pipeline
    pipeline = Pipeline.from_pretrained("pyannote/speaker-diarization-3.1", use_auth_token=auth_token)

    # Send to GPU if available
    if torch.cuda.is_available():
        pipeline.to(torch.device("cuda"))

    # pipeline.to(torch.device("cuda"))

    # Apply the pipeline to the audio file
    diarization = pipeline(file_path)

    # with Hook(ProgressHook(), TimingHook()) as hook:
    #   diarization = pipeline(file_path, hook=hook)

    # Generate VTT content
    vtt_content = "WEBVTT\n\n"
    for i, (turn, _, speaker) in enumerate(diarization.itertracks(yield_label=True)):
        start_time = format_timestamp(turn.start)
        end_time = format_timestamp(turn.end)
        vtt_content += f"{i+1}\n"
        vtt_content += f"{start_time} --> {end_time}\n"
        vtt_content += f"Speaker {speaker}\n\n"

    # Write to VTT file
    with open(output_path, "w") as f:
        f.write(vtt_content)

    return diarization

def format_timestamp(seconds):
    milliseconds = int(round((seconds - int(seconds)) * 1000))
    return f"{int(seconds // 3600):02}:{int((seconds % 3600) // 60):02}:{int(seconds % 60):02}.{milliseconds:03}"

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Perform speaker diarization and save output in VTT format.")
    parser.add_argument("--file_path", required=True, help="Path to the audio file.")
    parser.add_argument("--output_path", default="diarization.vtt", help="Path to the output VTT file (default: diarization.vtt).")
    parser.add_argument("--auth_token", help="Hugging Face access token (optional).")
    args = parser.parse_args()

    diarize_audio_to_vtt(args.file_path, args.output_path, args.auth_token)
    print(f"Diarization saved to: {args.output_path}")
