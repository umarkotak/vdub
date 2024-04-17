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
  liblzma-dev

WORKDIR /root

SHELL ["/bin/bash", "-c"]

RUN curl https://pyenv.run | bash

ENV PYENV_ROOT="/root/.pyenv"

ENV PATH="$PYENV_ROOT/bin:$PATH"

RUN eval "$(pyenv init -)"

RUN echo 'PYENV_ROOT="/root/.pyenv"' >> /root/.bashrc
RUN echo '[[ -d $PYENV_ROOT/bin ]] && export PATH="$PYENV_ROOT/bin:$PATH"'  >> /root/.bashrc
RUN echo 'eval "$(pyenv init -)"' >> /root/.bashrc

RUN echo 'PYENV_ROOT="/root/.pyenv"' >> /root/.profile
RUN echo '[[ -d $PYENV_ROOT/bin ]] && export PATH="$PYENV_ROOT/bin:$PATH"'  >> /root/.profile
RUN echo 'eval "$(pyenv init -)"' >> /root/.profile

RUN source /root/.bashrc

RUN pyenv install 3.10.14

RUN pyenv global 3.10.14

RUN ln -s /root/.pyenv/shims/python /usr/bin/python

RUN mv /usr/bin/python3 /usr/bin/python3.original

RUN ln -s /root/.pyenv/shims/python3 /usr/bin/python3

RUN pyenv virtualenv 3.10.14 vdub-applio

RUN pyenv activate vdub-applio

RUN python -m ensurepip --upgrade

RUN apt-get install -y pip

RUN pip install torch torchvision torchaudio praat-parselmouth

RUN git clone https://github.com/IAHispano/Applio.git

WORKDIR /root/Applio

RUN chmod +x run-install.sh

RUN ./run-install.sh

RUN chmod +x run-applio.sh

RUN ./run-applio.sh
