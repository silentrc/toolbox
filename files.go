package toolbox

import (
	"crypto/sha1"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type filesUtils struct {
}

// 文件类
func (u *utils) NewFilesUtils() *filesUtils {
	return &filesUtils{}
}

// 根据Render计算filehash
func (f *filesUtils) FileHashRenderSum(render io.Reader) (hex string, err error) {
	b, err := io.ReadAll(render)
	if err != nil {
		return
	}
	return f.FileHashBytesSum(b)
}

// 根据Bytes计算filehash
func (f *filesUtils) FileHashBytesSum(fileByte []byte) (hex string, err error) {
	h := sha1.New()
	_, err = h.Write(fileByte)
	if err != nil {
		return
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// 检测目录是否存在,不存在则创建
func (f *filesUtils) CheckOSFilePath(filePath string) error {
	_, err := os.Stat(filePath)
	if err == nil {
		return nil
	}
	return os.MkdirAll(filePath, 0755)
}

// 计算文件大小
func (f *filesUtils) FormatFileSize(fileSize int64) (size string) {
	switch {
	case fileSize < 1024:
		return fmt.Sprintf("%.0fB", float64(fileSize)/float64(1))
	case fileSize < (1024 * 1024):
		return fmt.Sprintf("%.0fK", float64(fileSize)/float64(1024))
	case fileSize < (1024 * 1024 * 1024):
		return fmt.Sprintf("%.0fM", float64(fileSize)/float64(1024*1024))
	case fileSize < (1024 * 1024 * 1024 * 1024):
		return fmt.Sprintf("%.0fG", float64(fileSize)/float64(1024*1024*1024))
	case fileSize < (1024 * 1024 * 1024 * 1024 * 1024):
		return fmt.Sprintf("%.0fT", float64(fileSize)/float64(1024*1024*1024*1024))
	default:
		return fmt.Sprintf("%.0fE", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

// 下载文件获取进度条
func (f *filesUtils) DownloadFileProgress(fileHash, fileUrl, localPath string, callback IOCallback) error {
	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()
	//兼容https为ip情况
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cli := http.Client{Transport: tr}
	resp, err := cli.Get(fileUrl)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	//获取文件总大小
	fs, err := strconv.ParseUint(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		return err
	}
	//opts := &ioWriter{
	//	callback: callback,
	//	progress: progress{
	//		NowSize:  int64(fs),
	//		FileHash: fileHash,
	//	},
	//}
	//counter := NewUtils().NewIOUtils().NewIOWriteCallback(fileHash, int64(fs), opts, callback)
	counter := NewUtils().NewIOUtils().NewIOWriteCallback(fileHash, int64(fs), callback)

	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}
	return nil
}

func (f *filesUtils) UploadFileProgress() {

}
