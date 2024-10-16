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
