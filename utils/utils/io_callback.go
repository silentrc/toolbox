package utils

import "io"

type ioUtils struct {
}

// io options
func (io utils) NewIOUtils() *ioUtils {
	return &ioUtils{}
}

// 回调参数
// total总大小 单位为B
// nowSize当前大小 单位为B
// fileHash文件唯一标识
type IOCallback func(total int64, nowSize int64, fileHash string)

type progress struct {
	Total    int64
	NowSize  int64
	FileHash string
}

type ioWriter struct {
	callback IOCallback
	progress
}

type ioRender struct {
	callback IOCallback
	progress
}

// io render
// fileHash文件唯一标识
// total总大小 单位为字节
// callback IOCallback
func (io *ioUtils) NewIORenderCallback(fileHash string, size int64, callback IOCallback) io.Reader {
	return &ioRender{
		callback: callback,
		progress: progress{
			Total:    size,
			FileHash: fileHash,
		},
	}
}

// io write
// fileHash文件唯一标识
// total总大小 单位为字节
// callback IOCallback
func (io *ioUtils) NewIOWriteCallback(fileHash string, size int64, callback IOCallback) io.Writer {
	return &ioWriter{
		callback: callback,
		progress: progress{
			Total:    size,
			FileHash: fileHash,
		},
	}
}

// 读取字节长度并返回
func (io *ioRender) Read(p []byte) (n int, err error) {
	n = len(p)
	io.NowSize += int64(n)
	if err == nil && io.callback != nil {
		io.callback(io.Total, io.NowSize, io.FileHash)
	}
	return n, err
}

// 读取写入字节长度并返回
func (io *ioWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	io.NowSize += int64(n)
	if err == nil && io.callback != nil {
		io.callback(io.Total, io.NowSize, io.FileHash)
	}
	return n, err
}
