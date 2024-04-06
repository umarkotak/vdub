# WHISPER

## APPS
- python3.12
- python3 pip
- yt-dlp // youtube downloader
- whisper

## Ref
- https://github.com/openai/whisper#available-models-and-languages

## Build Command
- docker build --platform="linux/amd64" -t vdub-whisper -f vdub-whisper/Dockerfile .
- docker run vdub-whisper

- docker build -t vdub-whisper -f vdub-whisper/Dockerfile .
- docker run vdub-whisper

## Debug Run
- // this allow docker process to keep running
- docker run -dit --name vdub-whisper vdub-whisper
- docker exec -it vdub-whisper bash
- docker rm -f vdub-whisper
- yt-dlp -x --audio-format mp3 https://www.youtube.com/watch?v=rlf2OGUTvJg -o test
