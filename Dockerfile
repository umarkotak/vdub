# FROM ubuntu:24.04
FROM nvidia/cuda:12.6.2-cudnn-devel-ubuntu24.04

RUN apt-get update && apt-get install -y \
  wget \
  curl \
  git \
  ffmpeg \
  build-essential \
  libssl-dev \
  zlib1g-dev \
  libbz2-dev \
  libreadline-dev \
  libsqlite3-dev \
  llvm \
  libncursesw5-dev \
  xz-utils \
  tk-dev \
  libxml2-dev \
  libxmlsec1-dev \
  libffi-dev \
  liblzma-dev \
  libasound2-dev \
  libgtk2.0-dev \
  libx11-dev \
  vim \
  sox \
  xdg-utils \
  w3m \
  gcc \
  g++ \
  make \
  python3 \
  python3-dev \
  python3-pip \
  python3-venv \
  python3-wheel \
  espeak-ng \
  libsndfile1-dev \
  yt-dlp \
  kmod

WORKDIR /root

SHELL ["/bin/bash", "-c"]

## Install [pyenv] - Python env management

RUN curl https://pyenv.run | bash

ENV PYENV_ROOT="/root/.pyenv"

ENV PATH="$PYENV_ROOT/bin:$PATH"

RUN echo 'eval "$(pyenv init -)"' >> /root/.bashrc

RUN echo 'eval "$(pyenv virtualenv-init -)"' >> /root/.bashrc

RUN pyenv install 3.12.7

RUN pyenv global 3.12.7

RUN ln -s /root/.pyenv/shims/python /usr/bin/python

RUN mv /usr/bin/python3 /usr/bin/python3.original

RUN ln -s /root/.pyenv/shims/python3 /usr/bin/python3

## Install [Whisper] - Speech to text for transcripting

WORKDIR /root

RUN git clone https://github.com/ggerganov/whisper.cpp.git

WORKDIR /root/whisper.cpp

RUN make

RUN make large-v3-turbo

RUN make quantize

RUN ./quantize models/ggml-large-v3-turbo.bin models/ggml-large-v3-turbo-q5_0.bin q5_0

## Install [Golang] - GO programming language

WORKDIR /root

RUN mkdir -p go/bin

RUN curl -OL https://go.dev/dl/go1.22.2.linux-arm64.tar.gz

RUN curl -OL https://go.dev/dl/go1.22.2.linux-amd64.tar.gz

ARG ARCH

RUN if [ "$ARCH" == "arm64" ]; then \
    tar -C /usr/local -xvf go1.22.2.linux-arm64.tar.gz; \
  else \
    tar -C /usr/local -xvf go1.22.2.linux-amd64.tar.gz; \
  fi

RUN rm -rf go1.22.2.linux-arm64.tar.gz

RUN rm -rf go1.22.2.linux-amd64.tar.gz

ENV GOROOT=/usr/local/go

ENV GOPATH=/root/go

ENV GOBIN=$GOPATH/bin

ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

## Install [Python dependencies]



RUN bash -i -c "source ~/.bashrc && python -m pip install wget"
RUN bash -i -c "source ~/.bashrc && python -m pip install llvmlite torch torchaudio --extra-index-url https://download.pytorch.org/whl/cu121"
RUN bash -i -c "source ~/.bashrc && python -m pip install tensorboard"
RUN bash -i -c "source ~/.bashrc && python -m pip install bs4"
RUN bash -i -c "source ~/.bashrc && python -m pip install pydub"
RUN bash -i -c "source ~/.bashrc && python -m pip install transformers"
RUN bash -i -c "source ~/.bashrc && python -m pip install noisereduce"
RUN bash -i -c "source ~/.bashrc && python -m pip install torchcrepe"
RUN bash -i -c "source ~/.bashrc && python -m pip install pedalboard"
RUN bash -i -c "source ~/.bashrc && python -m pip install onnxruntime"

RUN bash -i -c "source ~/.bashrc && python -m ensurepip --upgrade"
RUN bash -i -c "source ~/.bashrc && python -m pip install --upgrade pip"
RUN bash -i -c "source ~/.bashrc && python -m pip install --upgrade setuptools"
RUN bash -i -c "source ~/.bashrc && python -m pip install audio_upscaler"
# RUN bash -i -c "source ~/.bashrc && python -m pip install faiss"
RUN bash -i -c "source ~/.bashrc && python -m pip install faiss-cpu faiss-gpu-cu12"

## Install [Yt DLP] - Youtube video downloader

WORKDIR /root

RUN bash -i -c "source ~/.bashrc && python -m pip install yt-dlp"

## Install [Edge TTS] - Text to speech

WORKDIR /root

RUN bash -i -c "source ~/.bashrc && python -m pip install edge-tts"

## Install [Audio separator] - Separate audio file into Instruments and Vocals

WORKDIR /root

RUN bash -i -c "source ~/.bashrc && python -m pip install audio-separator"

## Install [Pyannote] - Speaker diarization, identify who speak when

RUN bash -i -c "source ~/.bashrc && python -m pip install pyannote.audio"

## Install [Rvc Cli] - For voice cloning

WORKDIR /root

RUN git clone https://github.com/blaisewf/rvc-cli.git

WORKDIR /root/rvc-cli

RUN chmod +x install.sh

RUN ./install.sh

RUN python rvc_cli.py -h

RUN wget https://huggingface.co/ORVC/Ov2Super/resolve/main/f0Ov2Super40kD.pth?download=true -O /root/rvc-cli/rvc/models/pretraineds/pretrained_v2/f0Ov2Super40kD.pth
RUN wget https://huggingface.co/lj1995/VoiceConversionWebUI/resolve/main/rmvpe.pt?download=true -O /root/rvc-cli/rvc/models/predictors/rmvpe.pt
RUN wget https://huggingface.co/lj1995/VoiceConversionWebUI/resolve/main/pretrained_v2/f0G48k.pth -O /root/rvc-cli/rvc/models/pretraineds/pretrained_v2/f0G48k.pth
RUN wget https://huggingface.co/lj1995/VoiceConversionWebUI/resolve/main/pretrained_v2/f0D48k.pth -O /root/rvc-cli/rvc/models/pretraineds/pretrained_v2/f0D48k.pth

## DONE

WORKDIR /root

RUN mkdir -p /root/shared

RUN mkdir -p /root/vdub

WORKDIR /root/vdub

EXPOSE 29000
