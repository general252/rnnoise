package main

import (
	"log"
	"os"

	"github.com/general252/rnnoise"
)

func main() {
	rnn := rnnoise.NewRNNoise()
	defer rnn.Close()

	bitsPCM, err := os.ReadFile("stereo_48k.pcm")
	if err != nil {
		return
	}

	f, _ := os.Create("stereo_48k_denoised.pcm")
	defer f.Close()

	for i := 0; i+rnnoise.FrameLength < len(bitsPCM); i += rnnoise.FrameLength {
		float3Frame, _ := rnnoise.BytesToFrames(bitsPCM[i : i+rnnoise.FrameLength])

		voiceProb, denoisedFrame, err := rnn.ProcessFrame(float3Frame)
		if err != nil {
			return
		}
		log.Printf("语音概率: %f\n", voiceProb)
		log.Printf("降噪后的音频帧: %v\n", denoisedFrame)
		log.Println("-----------------")

		denoisedBytes, _ := rnnoise.FramesToBytes(denoisedFrame)
		_, _ = f.Write(denoisedBytes)
	}

}
