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

- PROGRESS_NO_TRUNC=1 docker build --progress plain -t vdub-applio -f vdub-applio/Dockerfile .
- docker run -dit -p 6969:6969 -v /Users/umarramadhana/umar/personal_projects/vdub/shared:/root/shared --name vdub-applio vdub-applio
- docker run -d --network host --name vdub-applio vdub-applio
- docker exec -it vdub-applio bash
- docker rm -f vdub-applio
- ./run-applio.sh
- edge-tts --text "Hello, world!" --write-media hello.mp3 --write-subtitles hello.vtt

- docker run -dit --name vdub-whisper-cpp vdub-whisper-cpp
- docker exec -it vdub-whisper-cpp bash
- docker rm -f vdub-whisper-cpp
- youtubedr download -o test.mp4 https://www.youtube.com/watch?v=rlf2OGUTvJg
- ffmpeg -i test.mp3 -acodec pcm_u8 -ar 16000 -ac 1 -acodec pcm_s16le test.wav
- ./main test.wav

- PROGRESS_NO_TRUNC=1 docker build --progress plain -t applio -f Dockerfile .
- docker run -dit -p 6969:6969 -p 7070:7070 --expose 6969 --expose 7070 -v /Users/umarramadhana/umar/personal_projects/vdub:/root/shared --name applio applio
- docker run -dit --network host -v /Users/umarramadhana/umar/personal_projects/vdub:/root/shared --name applio applio
- docker exec -it applio bash
- docker rm -f applio
- ./run-applio.sh
