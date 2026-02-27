package cos

import (
	"context"
	"encoding/base64"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"mediahub/pkg/storage"
	"mime"
	"net/http"
	url2 "net/url"
	"path"
)

type cosStorageFactory struct {
	bucketUrl string
	secretId  string
	secretKey string
	cdnDomain string
}

func (f *cosStorageFactory) CreateStorage() storage.Storage {
	return &cosStorage{
		bucketUrl: f.bucketUrl,
		secretKey: f.secretKey,
		secretId:  f.secretId,
		cdnDomain: f.cdnDomain,
	}
}

func NewCosStorageFactory(bucketUrl, secretId, secretKey, cdnDomain string) storage.StorageFactory {
	return &cosStorageFactory{
		bucketUrl: bucketUrl,
		secretKey: secretKey,
		secretId:  secretId,
		cdnDomain: cdnDomain,
	}
}

type cosStorage struct {
	bucketUrl string
	secretId  string
	secretKey string
	cdnDomain string
}

func (s *cosStorage) Upload(r io.Reader, md5Digest []byte, dstPath string) (url string, err error) {
	u, _ := url2.Parse(s.bucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  s.secretId,
			SecretKey: s.secretKey,
		},
	})

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: mime.TypeByExtension(path.Ext(dstPath)),
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{},
	}
	if len(md5Digest) != 0 {
		opt.ObjectPutHeaderOptions.ContentMD5 = base64.StdEncoding.EncodeToString(md5Digest)
	}

	_, err = client.Object.Put(context.Background(), dstPath, r, opt)
	if err != nil {
		return "", err
	}
	url = s.bucketUrl + dstPath
	if s.cdnDomain != "" {
		url = s.cdnDomain + dstPath
	}
	return
}
