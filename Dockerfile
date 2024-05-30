FROM ubuntu:24.04

ARG ARCH

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
  yt-dlp

WORKDIR /root

SHELL ["/bin/bash", "-c"]

## Install [pyenv] - Python env management

RUN curl https://pyenv.run | bash

ENV PYENV_ROOT="/root/.pyenv"

ENV PATH="$PYENV_ROOT/bin:$PATH"

RUN echo 'eval "$(pyenv init -)"' >> /root/.bashrc

RUN echo 'eval "$(pyenv virtualenv-init -)"' >> /root/.bashrc

RUN pyenv install 3.10.14

RUN pyenv global 3.10.14

RUN ln -s /root/.pyenv/shims/python /usr/bin/python

RUN mv /usr/bin/python3 /usr/bin/python3.original

RUN ln -s /root/.pyenv/shims/python3 /usr/bin/python3

## Install [Yt DLP] - Youtube video downloader

RUN bash -i -c "source ~/.bashrc && python -m pip install yt-dlp llvmlite torch torchaudio --extra-index-url https://download.pytorch.org/whl/cu118"

## Install [Edge TTS] - Text to speech

RUN bash -i -c "source ~/.bashrc && python -m pip install edge-tts"

## Install [Vocal remover] - Separate audio file into Instruments and Vocals

WORKDIR /root

RUN git clone https://github.com/tsurumeso/vocal-remover.git

WORKDIR /root/vocal-remover

RUN bash -i -c "source ~/.bashrc && python -m pip install -r requirements.txt"

RUN curl -OL https://huggingface.co/fabiogra/baseline_vocal_remover/resolve/main/baseline.pth?download=true

## Install [Golang] - GO programming language

WORKDIR /root

RUN mkdir -p go/bin

RUN if [ "$ARCH" == "arm64" ]; then \
    curl -OL https://go.dev/dl/go1.22.2.linux-arm64.tar.gz; \
    tar -C /usr/local -xvf go1.22.2.linux-arm64.tar.gz; \
    rm -rf go1.22.2.linux-arm64.tar.gz; \
  else \
    curl -OL https://go.dev/dl/go1.22.2.linux-amd64.tar.gz; \
    tar -C /usr/local -xvf go1.22.2.linux-amd64.tar.gz; \
    rm -rf go1.22.2.linux-amd64.tar.gz; \
  fi

ENV GOROOT=/usr/local/go

ENV GOPATH=/root/go

ENV GOBIN=$GOPATH/bin

ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

## Install [Whisper] - Speech to text for transcripting

RUN git clone https://github.com/ggerganov/whisper.cpp.git

WORKDIR /root/whisper.cpp

RUN make

RUN make medium.en

RUN make quantize

RUN ./quantize models/ggml-medium.en.bin models/ggml-medium.en-q5_0.bin q5_0

WORKDIR /root

# RUN

RUN mkdir -p /root/shared

ENV PYTORCH_ENABLE_MPS_FALLBACK=1

ENV PYTORCH_MPS_HIGH_WATERMARK_RATIO=0.0

EXPOSE 29000
