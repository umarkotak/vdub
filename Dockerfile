FROM ubuntu:24.04

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

## Install [Whisper] - Speech to text for transcripting

WORKDIR /root

RUN git clone https://github.com/ggerganov/whisper.cpp.git

WORKDIR /root/whisper.cpp

RUN make

RUN make medium.en

RUN make quantize

RUN ./quantize models/ggml-medium.en.bin models/ggml-medium.en-q5_0.bin q5_0

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

# [TEMP DEV] RVC CLI

# WORKDIR /root

# RUN git clone https://github.com/blaisewf/rvc-cli.git

# WORKDIR /root/rvc-cli

# RUN chmod +x install.sh

# RUN ./install.sh

# RUN pip install tensorboard wget bs4 pydub transformers faiss-cpu noisereduce torchcrepe

# RUN python rvc_cli.py -h

# RUN wget https://huggingface.co/ORVC/Ov2Super/resolve/main/f0Ov2Super40kD.pth?download=true -O /root/rvc-cli/rvc/models/pretraineds/pretrained_v2/f0Ov2Super40kD.pth

# RUN wget https://huggingface.co/lj1995/VoiceConversionWebUI/resolve/main/rmvpe.pt?download=true -O rvc/models/predictors/rmvpe.pt
# RUN wget https://huggingface.co/lj1995/VoiceConversionWebUI/resolve/main/pretrained_v2/f0G48k.pth -O rvc/models/pretraineds/pretrained_v2/f0G48k.pth
# RUN wget https://huggingface.co/lj1995/VoiceConversionWebUI/resolve/main/pretrained_v2/f0D48k.pth -O rvc/models/pretraineds/pretrained_v2/f0D48k.pth

# RUN

RUN mkdir -p /root/shared

RUN mkdir -p /root/vdub

ENV PYTORCH_ENABLE_MPS_FALLBACK=1

ENV PYTORCH_MPS_HIGH_WATERMARK_RATIO=0.0

WORKDIR /root/vdub

EXPOSE 29000
