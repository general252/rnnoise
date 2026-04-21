package rnnoise

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	FrameSize   = 480                    // RNNoise 固定每帧采样点数
	SampleSize  = 4                      // float32 占用 4 字节
	FrameLength = FrameSize * SampleSize // 1920 字节
)

// BytesToFrames 安全地将 ffmpeg 的 f32le 字节流转换为 float32 切片
func BytesToFrames(buf []byte) ([]float32, error) {
	if len(buf) != FrameLength {
		return nil, fmt.Errorf("数据长度错误: 期望 %d, 实际 %d", FrameLength, len(buf))
	}

	frames := make([]float32, FrameSize)
	for i := 0; i < FrameSize; i++ {
		// 读取 4 字节并转换为小端序 uint32
		bits := binary.LittleEndian.Uint32(buf[i*SampleSize : (i+1)*SampleSize])
		// 将 uint32 位的表示转换为 float32
		frames[i] = math.Float32frombits(bits)
	}
	return frames, nil
}

// FramesToBytes 安全地将处理后的 float32 切片转换回字节流
func FramesToBytes(frames []float32) ([]byte, error) {
	if len(frames) != FrameSize {
		return nil, fmt.Errorf("帧大小错误: 期望 %d, 实际 %d", FrameSize, len(frames))
	}

	buf := make([]byte, FrameLength)
	for i := 0; i < FrameSize; i++ {
		// 将 float32 转换为 uint32 位表示
		bits := math.Float32bits(frames[i])
		// 以小端序写入字节数组
		binary.LittleEndian.PutUint32(buf[i*SampleSize:(i+1)*SampleSize], bits)
	}
	return buf, nil
}
