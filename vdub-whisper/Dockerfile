FROM ubuntu:24.04

RUN apt-get update && apt-get install -y \
  wget \
  curl \
  git \
  ffmpeg \
  python3.12 \
  python3-pip \
  python3.12-venv

WORKDIR /root

SHELL ["/bin/bash", "-c"]

RUN mkdir -p installations/whisper

# RUN python3 -m venv installations/whisper

# RUN source installations/whisper/bin/activate

# RUN installations/whisper/bin/pip install -U openai-whisper

RUN pip install setuptools-rust --break-system-packages

RUN pip install -U openai-whisper --break-system-packages

RUN pip install git+https://github.com/openai/whisper.git  --break-system-packages

RUN pip install --upgrade youtube-dl --break-system-packages

RUN pip install yt-dlp --break-system-packages

