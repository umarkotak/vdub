# Command snippet

### RVC-CLI

python rvc_cli.py infer \
  --input_path /root/vdub/shared/task-public-nak-1/raw_video_audio_Vocals_16KHz.wav \
  --output_path /root/vdub/shared/task-public-nak-1/generated_speech/0.wav \
  --pth_path "/root/vdub/shared/alya.pth" \
  --index_path "/root/vdub/shared/added_IVF777_Flat_nprobe_1_alya_v2.index"

python rvc_cli.py infer \
  --input_path /root/vdub/shared/task-public-nak-1/generated_speech/0.wav \
  --output_path /root/vdub/shared/task-public-nak-1/generated_speech/0-alya.wav \
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

/usr/bin/python /root/vocal-remover/inference.py \
  --input /root/shared/task-public-kurz-1/raw_video_audio.wav \
  -P /root/vocal-remover/baseline.pth \
  -o /root/shared/task-public-kurz-1

git clone https://github.com/SWivid/F5-TTS.git

cd F5-TTS

python inference-cli.py \
--model "F5-TTS" \
--output_dir /root/vdub/shared/task-public-nak-2/e2_generated_speech \
--ref_audio "/root/vdub/shared/task-public-nak-2/segmented_speech/0.wav" \
--ref_text "The human beings seem like a bad choice, but there's a better alternative." \
--gen_text "Kami menampilkan diri kami sebagai alternatif. Anda mengerti? Itulah yang dikatakan malaikat itu. Dan ketika mereka mengatakan ini, siapa yang mendengarkan?"

python inference-cli.py \
--model "E2-TTS" \
--output_dir /root/vdub/shared/task-public-nak-2/e2_generated_speech \
--ref_audio "/root/vdub/shared/task-public-nak-2/segmented_speech/1.wav" \
--ref_text "One. We're presenting ourselves as an alternative. You get it? That's what the angel said. And when they said this, who was listening?" \
--gen_text "Kami menampilkan diri kami sebagai alternatif. Anda mengerti? Itulah yang dikatakan malaikat itu. Dan ketika mereka mengatakan ini, siapa yang mendengarkan?"

python inference-cli.py \
--model "E2-TTS" \
--ref_audio "tests/ref_audio/test_zh_1_ref_short.wav" \
--ref_text "对，这就是我，万人敬仰的太乙真人。" \
--gen_text "突然，身边一阵笑声。我看着他们，意气风发地挺直了胸膛，甩了甩那稍显肉感的双臂，轻笑道，我身上的肉，是为了掩饰我爆棚的魅力，否则，岂不吓坏了你们呢？"
