FROM ubuntu:24.04

RUN apt-get update && apt-get install -y \
  wget \
  curl \
  git \
  build-essential \
  ffmpeg \
  xdg-utils \
  w3m \
  yt-dlp

WORKDIR /root

SHELL ["/bin/bash", "-c"]

RUN curl -OL https://go.dev/dl/go1.22.2.linux-arm64.tar.gz

RUN mkdir -p go/bin

RUN tar -C /usr/local -xvf go1.22.2.linux-arm64.tar.gz

# RUN echo 'export PATH=$PATH:/root/go/bin' >> /root/.profile

# RUN echo 'export PATH=$PATH:/root/go/bin' >> /root/.bashrc

# RUN source /root/.profile

ENV GOROOT=/usr/local/go

ENV GOPATH=/root/go

ENV GOBIN=$GOPATH/bin

ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

RUN git clone https://github.com/kkdai/youtube.git

WORKDIR /root/youtube

RUN go mod tidy

RUN go build -v ./cmd/youtubedr

RUN chmod +x youtubedr

RUN mv youtubedr /root/go/bin/youtubedr

WORKDIR /root

RUN rm -rf go1.22.2.linux-arm64.tar.gz

RUN git clone https://github.com/ggerganov/whisper.cpp.git

WORKDIR /root/whisper.cpp

RUN bash ./models/download-ggml-model.sh base.en

RUN make

WORKDIR /root
