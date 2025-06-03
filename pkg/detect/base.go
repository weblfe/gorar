package detect

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// 压缩格式魔数签名库
var signatures = []struct {
	name    string
	offset  int
	magic   []byte
	checkFn func([]byte) bool // 特殊格式的验证函数
}{
	// 标准文件头签名
	{"7z", 0, []byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C}, nil},
	{"rar", 0, []byte{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x00}, nil},
	{"rar", 0, []byte{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x01}, nil}, // RAR5 格式
	{"zip", 0, []byte{0x50, 0x4B, 0x03, 0x04}, nil},
	{"gz", 0, []byte{0x1F, 0x8B}, nil},
	{"bz2", 0, []byte{0x42, 0x5A, 0x68}, nil},
	{"xz", 0, []byte{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00}, nil},
	{"zst", 0, []byte{0x28, 0xB5, 0x2F, 0xFD}, nil},
	{"lz4", 0, []byte{0x04, 0x22, 0x4D, 0x18}, nil}, // LZ4 帧格式
	{"lz", 0, []byte{0x4C, 0x5A, 0x49, 0x50}, nil},  // LZIP 格式
	{"s2", 0, []byte{0xFF, 0x06, 0x00, 0x00}, nil},  // S2 压缩格式
	// 需要特殊偏移检查的格式
	{"tar", 257, []byte{0x75, 0x73, 0x74, 0x61, 0x72}, nil}, // ustar 标识
	// 需要自定义验证逻辑的格式
	{"br", 0, nil, func(data []byte) bool { // Brotli 格式验证
		return len(data) >= 3 &&
			data[0] == 0x1B &&
			data[1] == 0x1F &&
			(data[2] == 0x00 || (data[2] >= 0x01 && data[2] <= 0x03))
	}},
	{"zz", 0, nil, func(data []byte) bool { // Zlib 格式验证
		if len(data) < 2 || data[0] != 0x78 {
			return false
		}
		// 验证第二个字节范围值
		return data[1] >= 0x01 && data[1] <= 0xDA
	}},
	{"sz", 0, nil, func(data []byte) bool { // Snappy 帧验证
		return len(data) >= 4 &&
			data[0] == 0xFF &&
			data[1] == 0x06 &&
			data[2] == 0x00 &&
			data[3] == 0x00
	}},
}

// Detect 检测压缩格式
//
//	参数：filePath string 文件路径
//	可选参数: fss  fs.FS 文件系统
func Detect(filePath string, fss ...fs.FS) (string, error) {
	var (
		ext    string
		source = filePath
	)
	for strings.Contains(source, ".") {
		ext = filepath.Ext(source)[1:]
		source = ext
	}
	if ext != "" && ext == source {
		return ext, nil
	}

	var (
		err  error
		file fs.File
	)
	if len(fss) > 0 && fss[0] != nil {
		if file, err = fss[0].Open(filePath); err != nil {
			return "", err
		}
	} else {
		if file, err = os.Open(filePath); err != nil {
			return "", err
		}
	}
	defer func(file fs.File) {
		_ = file.Close()
	}(file)

	// 读取前512字节（覆盖所有格式签名范围）
	var (
		n      int
		buffer = make([]byte, 512)
	)
	if n, err = file.Read(buffer); err != nil {
		return "", err
	}
	if n == 0 {
		return "", errors.New("empty file")
	}

	// 遍历所有签名规则进行检测
	for _, sig := range signatures {
		// 检查偏移量是否超出文件范围
		start := sig.offset
		if start >= n {
			continue
		}

		// 优先使用自定义验证函数
		if sig.checkFn != nil {
			if sig.checkFn(buffer[start:]) {
				return sig.name, nil
			}
			continue
		}

		// 标准魔数匹配
		end := start + len(sig.magic)
		if end > n {
			continue
		}

		if bytes.Equal(buffer[start:end], sig.magic) {
			return sig.name, nil
		}
	}

	return "", errors.New("unknown compress format")
}
