package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// generateID 生成一个随机的 32 字节 ID
// 使用随机数生成器填充缓冲区，然后将其编码为 16 进制字符串
func generateID() string {
	buf := make([]byte, 32)
	io.ReadFull(rand.Reader, buf)
	return hex.EncodeToString(buf)
}

// hashKey 使用 MD5 哈希算法对键进行哈希
// 将键转换为字节数组，计算 MD5 哈希，然后将其编码为 16 进制字符串
func hashKey(key string) string {
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

// newEncryptionKey 生成一个新的 32 字节加密密钥
// 使用随机数生成器填充缓冲区，返回字节数组作为密钥
func newEncryptionKey() []byte {
	keyBuf := make([]byte, 32)
	io.ReadFull(rand.Reader, keyBuf)
	return keyBuf
}

// copyStream 使用给定的流加密器复制数据
// 从源读取数据，使用流加密器加密，然后写入目标
func copyStream(stream cipher.Stream, blockSize int, src io.Reader, dst io.Writer) (int, error) {
	var (
		buf = make([]byte, 32*1024)
		nw  = blockSize
	)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			nn, err := dst.Write(buf[:n])
			if err != nil {
				return 0, err
			}
			nw += nn
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}
	return nw, nil
}

// copyDecrypt 使用给定的密钥解密数据流
// 从源读取数据，解密后写入目标
func copyDecrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key) // 创建 AES 加密块
	if err != nil {
		return 0, err
	}

	// 从源读取初始化向量（IV）
	iv := make([]byte, block.BlockSize()) // 创建与块大小相同的 IV 缓冲区
	if _, err := src.Read(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv) // 创建 CTR 模式的流加密器
	return copyStream(stream, block.BlockSize(), src, dst) // 使用流加密器复制数据
}

// copyEncrypt 使用给定的密钥加密数据流
// 从源读取数据，加密后写入目标
func copyEncrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, block.BlockSize()) // 16 bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	// prepend the IV to the file.
	if _, err := dst.Write(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dst)
}
