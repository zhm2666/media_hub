package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"mime"
	"net/http"
	url2 "net/url"
	"os"
	"path"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// COS 配置
var (
	bucketUrl = flag.String("bucket", "", "COS 存储桶地址，如: https://example-123456789.cos.ap-guangzhou.myqcloud.com")
	secretId  = flag.String("id", "", "SecretID")
	secretKey = flag.String("key", "", "SecretKey")
	cdnDomain = flag.String("cdn", "", "CDN 加速域名，如: https://cdn.example.com")
	filePath  = flag.String("file", "", "要上传的文件路径")
	dstPath   = flag.String("dst", "", "目标路径，如: /test/image.jpg")
)

// 上传文件到 COS
func uploadToCOS(fileContent []byte, md5Digest []byte, dst string) (url string, err error) {
	u, _ := url2.Parse(*bucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  *secretId,
			SecretKey: *secretKey,
		},
	})

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: mime.TypeByExtension(path.Ext(dst)),
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{},
	}
	if len(md5Digest) != 0 {
		opt.ObjectPutHeaderOptions.ContentMD5 = base64.StdEncoding.EncodeToString(md5Digest)
	}

	_, err = client.Object.Put(context.Background(), dst, bytes.NewReader(fileContent), opt)
	if err != nil {
		return "", err
	}

	url = *bucketUrl + dst
	if *cdnDomain != "" {
		url = *cdnDomain + dst
	}
	return
}

// 计算文件 MD5
func calcMD5(data []byte) []byte {
	// 简单模拟，实际使用 crypto/md5
	return nil
}

func main() {
	flag.Parse()

	// 检查必需参数
	if *bucketUrl == "" || *secretId == "" || *secretKey == "" {
		fmt.Println("错误: 请提供 COS 配置参数")
		fmt.Println("Usage: cos_test -bucket <COS桶地址> -id <SecretID> -key <SecretKey> -cdn <CDN域名> -file <文件路径> -dst <目标路径>")
		fmt.Println("")
		fmt.Println("Example:")
		fmt.Println("  cos_test -bucket https://example-1256487221.cos.ap-guangzhou.myqcloud.com -id AKxxx -key xxx -cdn https://cdn.example.com -file ./test.jpg -dst /test/image.jpg")
		os.Exit(1)
	}

	// 如果没有指定文件，则创建一个测试图片
	var fileContent []byte
	if *filePath != "" && fileExists(*filePath) {
		var err error
		fileContent, err = os.ReadFile(*filePath)
		if err != nil {
			fmt.Printf("读取文件失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("读取文件: %s, 大小: %d bytes\n", *filePath, len(fileContent))
	} else {
		// 创建一个简单的测试图片 (1x1 白色 PNG)
		fileContent = createTestImage()
		*dstPath = "/test/test_image.png"
		fmt.Println("使用默认测试图片，大小:", len(fileContent), "bytes")
	}

	// 计算 MD5
	md5Digest := calcMD5(fileContent)

	// 上传到 COS
	fmt.Println("\n开始上传到 COS...")
	fmt.Printf("目标路径: %s\n", *dstPath)

	resultURL, err := uploadToCOS(fileContent, md5Digest, *dstPath)
	if err != nil {
		fmt.Printf("上传失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n========== 上传成功! ==========")
	fmt.Printf("文件访问 URL: %s\n", resultURL)
	fmt.Println("================================\n")

	// 打印 CDN 配置信息
	fmt.Println("配置信息:")
	fmt.Printf("  Bucket: %s\n", *bucketUrl)
	fmt.Printf("  CDN:    %s\n", *cdnDomain)
}

// 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// 创建测试图片 (1x1 白色 PNG)
func createTestImage() []byte {
	// 最小 PNG 图片 (1x1 白色)
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG 签名
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR 块
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, 0xDE,
		0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, 0x54, // IDAT 块
		0x08, 0xD7, 0x63, 0xF8, 0xFF, 0xFF, 0x3F, 0x00, 0x05, 0xFE, 0x02, 0xFE, 0xDC, 0xCC, 0x59, 0xE7,
		0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, // IEND 块
		0xAE, 0x42, 0x60, 0x82,
	}
	return pngData
}

// ========== 高级用法示例 ==========

// 示例1: 使用自定义 Content-Type
func uploadWithCustomContentType() {
	content := []byte("Hello COS!")
	dstPath := "/test/custom.txt"

	u, _ := url2.Parse(*bucketUrl)
	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  *secretId,
			SecretKey: *secretKey,
		},
	})

	// 强制指定 Content-Type
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType:        "text/plain; charset=utf-8",
			ContentDisposition: "attachment; filename=custom.txt",
		},
	}

	_, err := client.Object.Put(context.Background(), dstPath, bytes.NewReader(content), opt)
	if err != nil {
		fmt.Printf("上传失败: %v\n", err)
	}
}

// 示例2: 设置文件访问权限 (ACL)
func uploadWithACL() {
	content := []byte("Private file")
	dstPath := "/test/private.txt"

	u, _ := url2.Parse(*bucketUrl)
	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  *secretId,
			SecretKey: *secretKey,
		},
	})

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			// 私有读: X-Cos-ACL: private
			// 公有读: X-Cos-ACL: public-read
			XCosACL: "private",
		},
	}

	_, err := client.Object.Put(context.Background(), dstPath, bytes.NewReader(content), opt)
	if err != nil {
		fmt.Printf("上传失败: %v\n", err)
	}
}

//// 示例3: 上传大文件 (分片上传)
//func uploadLargeFile(filePath string) error {
//	u, _ := url2.Parse(*bucketUrl)
//	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
//		Transport: &cos.AuthorizationTransport{
//			SecretID:  *secretId,
//			SecretKey: *secretKey,
//		},
//	})
//
//	key := "/test/large_" + path.Base(filePath)
//	f, err := os.Open(filePath)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//
//	// 分片上传
//	_, _, err = client.Object.Upload(
//		context.Background(),
//		key,
//		f,
//		nil,
//	)
//
//	return err
//}

// 示例4: 下载文件
func downloadFile(dstPath string) ([]byte, error) {
	u, _ := url2.Parse(*bucketUrl)
	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  *secretId,
			SecretKey: *secretKey,
		},
	})

	resp, err := client.Object.Get(context.Background(), dstPath, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// 示例5: 删除文件
func deleteFile(dstPath string) error {
	u, _ := url2.Parse(*bucketUrl)
	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  *secretId,
			SecretKey: *secretKey,
		},
	})

	_, err := client.Object.Delete(context.Background(), dstPath, nil)
	return err
}

// 示例6: 查询文件是否存在
func doesFileExist(dstPath string) (bool, error) {
	u, _ := url2.Parse(*bucketUrl)
	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  *secretId,
			SecretKey: *secretKey,
		},
	})

	_, err := client.Object.Head(context.Background(), dstPath, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
