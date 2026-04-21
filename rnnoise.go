package rnnoise

//#include <stdio.h>
//#include "rnnoise.h"
//#cgo pkg-config: rnnoise
import "C"
import (
	"fmt"
	"unsafe"
)

type RNNoise struct {
	state *C.DenoiseState
}

func NewRNNoise() *RNNoise {
	return &RNNoise{
		state: C.rnnoise_create(nil),
	}
}

func (rnn *RNNoise) Close() {
	C.rnnoise_destroy(rnn.state)
}

// ProcessFrame 处理单个音频帧（10毫秒）
//
// 参数:
//   - frame: 480个float32样本（10ms @ 48kHz），范围-1.0到1.0
//
// 返回:
//   - float32: 语音概率 (0.0-1.0)
//   - []float32: 降噪后的音频帧（480个样本）
//   - error: 处理失败时的错误信息
//
// 注意: 输入帧必须恰好包含480个样本，对应48kHz采样率下的10毫秒音频
func (rnn *RNNoise) ProcessFrame(frame []float32) (float32, []float32, error) {
	if len(frame) != 480 {
		return 0, nil, fmt.Errorf("帧大小必须为480个样本（10ms @ 48kHz），当前为%d", len(frame))
	}

	// 将float32样本转换为RNNoise期望的格式（16位整数范围的float）
	rnnoiseInput := make([]float32, 480)
	rnnoiseOutput := make([]float32, 480)

	for i, sample := range frame {
		// 将-1.0到1.0范围转换为-32768到32767范围
		rnnoiseInput[i] = sample * 32768.0
	}

	// 调用RNNoise处理函数
	voiceProb := C.rnnoise_process_frame(
		rnn.state,
		(*C.float)(unsafe.Pointer(&rnnoiseOutput[0])),
		(*C.float)(unsafe.Pointer(&rnnoiseInput[0])),
	)

	// 将输出转换回-1.0到1.0范围
	output := make([]float32, 480)
	for i, sample := range rnnoiseOutput {
		output[i] = sample / 32768.0
	}

	return float32(voiceProb), output, nil
}
