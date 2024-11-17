# VDUB Project
This app receive youtube video url as an input then it will generate a new video using Bahasa Indonesia as an output

## How To Install

```
For amd (windows, linux, wsl, etc)
make build_amd

For arm (mac m1, etc)
make build_arm

To run the backend
make docker_run_win_gpu
```

## How To Run
1. After the backend server is running, you can open https://vdubb.vercel.app for the web gui
2. This web will by default pointing to localhost:29000 so make sure that port is available
3. Enjoy the app

## Notes
1. As of now the project only receive youtube url to be translated
2. You need to copy and rename `cookies.txt.sample` to `cookies.txt`

## Tech used
All tech used is open sourced, you can adjust specific part you need - for the complete list, see Dockerfile
1. [Whisper] - Speech to text for transcripting
2. [Golang] - GO programming language
3. [pyenv] - Python env management
4. [Python dependencies]
5. [Yt DLP] - Youtube video downloader
6. [Edge TTS] - Text to speech
7. [Audio separator] - Separate audio file into Instruments and Vocals
8. [Pyannote] - Speaker diarization, identify who speak when
9. [Rvc Cli] - For voice cloning
