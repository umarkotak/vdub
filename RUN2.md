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