# Command snippet

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

- docker run -dit --name vdub-applio vdub-applio
- docker exec -it vdub-applio bash
- docker rm -f vdub-applio
- ./run-applio.sh

- docker run -dit --name vdub-whisper-cpp vdub-whisper-cpp
- docker exec -it vdub-whisper-cpp bash
- docker rm -f vdub-whisper-cpp
- youtubedr download -o test.mp4 https://www.youtube.com/watch?v=rlf2OGUTvJg
- ./main test.wav
- ffmpeg -i test.mp3 -acodec pcm_u8 -ar 16000 -ac 1 -acodec pcm_s16le test.wav