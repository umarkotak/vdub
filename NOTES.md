# Command snippet

## Build Command
- docker build --platform="linux/amd64" -t vdub-whisper -f vdub-whisper/Dockerfile .

## Debug Run
- docker build -t vdub-whisper -f vdub-whisper/Dockerfile .
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

- docker build -t vdub-whisper-cpp -f vdub-whisper-cpp/Dockerfile .
- docker run -dit -v /Users/umarramadhana/umar/personal_projects/vdub/shared:/root/shared --name vdub-whisper-cpp vdub-whisper-cpp
- docker exec -it vdub-whisper-cpp bash
- docker rm -f vdub-whisper-cpp
- youtubedr download -o test.mp4 https://www.youtube.com/watch?v=rlf2OGUTvJg
- youtubedr download -o test.mp4 https://www.youtube.com/watch?v=yDMZJ7LgrGY
- ffmpeg -i test.mp3 -acodec pcm_u8 -ar 16000 -ac 1 -acodec pcm_s16le test.wav
- ffmpeg -i test.mp4 -vn -acodec pcm_s16le -ar 44100 -ac 2 test.wav
- ffmpeg -i test_Vocals.wav -acodec pcm_s16le -ac 1 -ar 16000 test_Vocals_16khz.wav
- ./main test_Vocals_16khz.wav

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
-filter_complex "
[0:a]atrim=end=3,asetpts=PTS-STARTPTS[audio0]; \
[1:a]atrim=start=10,asetpts=PTS-STARTPTS[audio1]; \
[audio0][audio1]concat=n=2:v=0:a=1[out]" \
-map 0:v -map "[out]" -c:v copy -c:a aac testvideo3.mp4

---

- docker build -t vdub-core -f vdub-core/Dockerfile .
- docker rm -f vdub-core
- youtubedr download -o test.mp4 https://www.youtube.com/watch?v=rlf2OGUTvJg
- youtubedr download -o test.mp4 https://www.youtube.com/watch?v=yDMZJ7LgrGY
- ffmpeg -i test.mp3 -acodec pcm_u8 -ar 16000 -ac 1 -acodec pcm_s16le test.wav
- ffmpeg -i test.mp4 -vn -acodec pcm_s16le -ar 44100 -ac 2 test.wav
- ffmpeg -i test_Vocals.wav -acodec pcm_s16le -ac 1 -ar 16000 test_Vocals_16khz.wav
- ./main test_Vocals_16khz.wav

---

edge-tts --text "Hai kamu, senang sekali kamu bisa bergabung dengan kami." --write-media edge-tts-voice-1.wav -v id-ID-ArdiNeural --rate=-10% --pitch=-10Hz

edge-tts --text "Kami ingin menceritakan sesuatu yang mengubah seni Kotska selamanya." --write-media edge-tts-voice-2.wav -v id-ID-ArdiNeural --rate=-10% --pitch=-10Hz

edge-tts --text "Seni Kotska dimulai sebagai proyek hobi berskala kecil, tetapi membuat video sains animasi" --write-media edge-tts-voice-3.wav -v id-ID-ArdiNeural --rate=-10% --pitch=-10Hz

edge-tts --text "yang gratis untuk semua orang tidaklah menghasilkan uang." --write-media edge-tts-voice-4.wav -v id-ID-ArdiNeural --rate=-10% --pitch=-10Hz

edge-tts --text "Sialan, kenyataan." --write-media edge-tts-voice-5.wav -v id-ID-ArdiNeural --rate=-10% --pitch=-10Hz

edge-tts --text "Untuk memenuhi kebutuhan, kami harus bekerja di siang hari dan membuat video seni Kotska di malam hari." --write-media edge-tts-voice-6.wav -v id-ID-ArdiNeural --rate=-10% --pitch=-10Hz

ffmpeg -i video.mp4 -i edge-tts-voice-1.wav -i edge-tts-voice-2.wav -i edge-tts-voice-3.wav -i edge-tts-voice-4.wav -i edge-tts-voice-5.wav -i edge-tts-voice-6.wav \
-filter_complex "[0:a]atrim=end=3,asetpts=PTS-STARTPTS[audio0]; \
[1:a]atrim=start=10,asetpts=PTS-STARTPTS[audio1]; \
[audio0][audio1]concat=n=2:v=0:a=1[out]" \
-map 0:v -map "[out]" -c:v copy -c:a aac testvideo3.mp4

trim mode
ffmpeg -i video.mp4 -i audio1.wav -i audio2.wav -i audio3.wav \
-filter_complex "[0:a]atrim=end=3.040,asetpts=PTS-STARTPTS[audio0]; \
[1:a]atrim=start=3.040,end=6.920,asetpts=PTS-STARTPTS[audio1]; \
[2:a]atrim=start=6.920,end=11.600,asetpts=PTS-STARTPTS[audio2]; \
[audio0][audio1][audio2]concat=n=3:v=0:a=1[out]" \
-map 0:v -map "[out]" -c:v copy -c:a aac output.mp4

[0:a]atrim=start=00:00:00.000,end=00:00:03.080,asetpts=PTS-STARTPTS[audio0];

speed mode
ffmpeg -i video.mp4 -i audio1.wav -i audio2.wav -i audio3.wav \
-filter_complex "[0:a]atempo=0.33333[audio0]; \
[1:a]atempo=2.0,asetpts=PTS-STARTPTS[audio1]; \
[2:a]atempo=1.5,asetpts=PTS-STARTPTS[audio2]; \
[audio0][audio1][audio2]concat=n=3:v=0:a=1[out]" \
-map 0:v -map "[out]" -c:v copy -c:a aac output.mp4

docker update --cpus="2048" vdub-core

---

  git clone https://github.com/neonbjb/tortoise-tts.git

  cd tortoise-tts/

  python setup.py install

  python tortoise/do_tts.py --text "I'm going to speak this" --voice random --preset fast

### RVC-CLI

python rvc_cli.py infer \
  --input_path /root/vdub/shared/task-public-nak-1/raw_video_audio_Vocals_16KHz.wav \
  --output_path /root/vdub/shared/task-public-nak-1/generated_speech/0.wav \
  --pth_path "/root/vdub/shared/alya.pth" \
  --index_path "/root/vdub/shared/added_IVF777_Flat_nprobe_1_alya_v2.index"

python rvc_cli.py preprocess \
  --model_name "nak_1" \
  --dataset_path /root/vdub/shared/task-public-nak-1/rvc_cli_train_dataset \
  --sample_rate 48000

python rvc_cli.py extract \
  --model_name "nak_1" \
  --sample_rate 48000 \
  --rvc_version v2 \
  --gpu 0 --cpu_cores 1

python rvc_cli.py train \
  --model_name "nak_1" \
  --rvc_version v2 \
  --gpu 1 \
  --sample_rate 48000 --total_epoch 500 --save_every_epoch 100

python rvc_cli.py index \
  --model_name "nak_1" \
  --rvc_version v2

python main.py --file_path /root/vdub/shared/task-public-nak-1/raw_video_audio_Vocals_16KHz.wav

python main.py \
  --file_path /root/vdub/shared/task-public-coba1/raw_video_audio_Vocals_16KHz.wav
  --output_path /root/vdub/shared/task-public-coba1/diarization.vtt
  --auth_token blablabla

apt-get install kmod

ARG nvidia_binary_version="550.120"
ARG nvidia_binary="NVIDIA-Linux-x86_64-${nvidia_binary_version}.run"
RUN wget https://us.download.nvidia.com/XFree86/Linux-x86_64/${nvidia_binary_version}/${nvidia_binary} &&
chmod +x ${nvidia_binary} &&
./${nvidia_binary} --accept-license --ui=none --no-kernel-module --no-questions &&
rm -rf ${nvidia_binary}

https://us.download.nvidia.com/XFree86/Linux-x86_64/550.120/NVIDIA-Linux-x86_64-550.120.run

export CUDA_HOME=/usr/local/cuda-12.6/bin
export PATH=$PATH:/usr/local/cuda-12.6/bin
export LD_LIBRARY_PATH=/usr/local/cuda-12.6/lib64
export CUDA_VISIBLE_DEVICES=0

sudo modprobe -r nvidia_uvm && sudo modprobe nvidia_uvm

apt install linux-modules-nvidia-565.90-5.15.153.1-microsoft-standard-WSL2
