

```bash
ffmpeg -i stereo.wav -ar 48000 -ac 1 -f f32le stereo_48k.pcm


ffplay -f f32le -ar 48000 stereo_48k.pcm
ffplay -f f32le -ar 48000 stereo_48k_denoised.pcm
```