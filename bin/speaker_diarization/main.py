from pyannote.audio import Pipeline
import torch
import argparse

def diarize_audio(file_path):
    # Load the pipeline
    pipeline = Pipeline.from_pretrained("pyannote/speaker-diarization-3.1", use_auth_token="HF_DIARIZATION_TOKEN")

    # Send to GPU if available
    # if torch.cuda.is_available():
    #     pipeline.to(torch.device("cuda"))

    pipeline.to(torch.device("cuda"))

    # Apply the pipeline to the audio file
    diarization = pipeline(file_path)

    return diarization

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Perform speaker diarization on an audio file.")
    parser.add_argument("--file_path", required=True, help="Path to the audio file.")
    parser.add_argument("--auth_token", help="Hugging Face access token (optional).")
    args = parser.parse_args()

    diarization = diarize_audio(args.file_path)

    # Print the results
    for turn, _, speaker in diarization.itertracks(yield_label=True):
        print(f"start={turn.start:.1f}s stop={turn.end:.1f}s speaker_{speaker}")
