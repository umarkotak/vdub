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
  sox

WORKDIR /root

SHELL ["/bin/bash", "-c"]

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

RUN git clone https://github.com/IAHispano/Applio.git

WORKDIR /root/Applio

RUN bash -i -c "source ~/.bashrc && python -m venv .venv"

RUN bash -i -c "source ~/.bashrc && . .venv/bin/activate"

RUN bash -i -c "source ~/.bashrc && python -m ensurepip"

RUN sed -i 's/praat-parselmouth/praat-parselmouth @ git\+https:\/\/github.com\/YannickJadoul\/Parselmouth/g' requirements.txt

RUN bash -i -c "source ~/.bashrc && python -m pip install -r requirements.txt"

RUN bash -i -c "source ~/.bashrc && python -m pip install edge-tts"

# RUN

RUN mkdir -p /root/shared

ENV PYTORCH_ENABLE_MPS_FALLBACK=1

ENV PYTORCH_MPS_HIGH_WATERMARK_RATIO=0.0

ENV GRADIO_SERVER_NAME="0.0.0.0"

EXPOSE 6969

ENTRYPOINT [ "python3" ]

CMD ["app.py"]
