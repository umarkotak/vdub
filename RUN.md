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

---

edge-tts --text "Hei kamu, senang sekali kamu bergabung dengan kami! Kami ingin memberi tahu Anda tentang sesuatu yang mengubah kurzgesagt selamanya. Kurzgesagt dimulai sebagai proyek gairah skala kecil." --write-media edge-tts-voice-1.wav -v id-ID-ArdiNeural --rate=-10% --pitch=-10Hz

python inference.py --input path/to/an/audio/file

yt-dlp -S ext mp4 https://www.youtube.com/watch?v=yDMZJ7LgrGY -o testvideo

ffmpeg -i testvideo.mp4 -i test_Instruments.wav -c:v copy -c:a aac -map 0:v:0 -map 1:a:0 testvideo2.mp4

ffmpeg -y -i testvideo2.mp4 -i edge-tts-voice-1.wav -filter_complex "[0:a]aformat=channel_layouts=mono[0a];[1:a]adelay=7000|7000,aformat=channel_layouts=mono[tmp];[0a][tmp]amerge=inputs=2,aformat=channel_layouts=stereo[audio_out]" -map 0:v -map [audio_out] testvideo3.mp4

ffmpeg -i video.mp4 -i edge-tts-voice-1.wav -i edge-tts-voice-1.wav \
-filter_complex "[0:a]atrim=end=3,asetpts=PTS-STARTPTS[audio0]; \
[1:a]atrim=start=10,asetpts=PTS-STARTPTS[audio1]; \
[audio0][audio1]concat=n=2:v=0:a=1[out]" \
-map 0:v -map "[out]" -c:v copy -c:a aac testvideo3.mp4
