package repository

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/8treenet/cdp-service/server/conf"
	"github.com/8treenet/freedom"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *ClondRepository {
			return &ClondRepository{}
		})
	})
}

// ClondRepository .
type ClondRepository struct {
	freedom.Repository
}

// NewUptoken .
func (repo *ClondRepository) NewUptoken(key ...string) (string, error) {
	bucket := conf.Get().System.QiniuBucket
	if len(key) > 0 {
		bucket = fmt.Sprintf("%s:%s", bucket, key[0])
	}
	putPolicy := storage.PutPolicy{
		Scope:   bucket,
		Expires: conf.Get().System.QiniuExpire,
	}

	mac := qbox.NewMac(conf.Get().System.QiniuAccessKey, conf.Get().System.QiniuSecretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken, nil
}

// UpLoad .
func (repo *ClondRepository) UpLoad(key string, data []byte) error {
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone, _ = storage.GetZone(conf.Get().System.QiniuAccessKey, conf.Get().System.QiniuBucket)
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	dataLen := int64(len(data))

	putExtra := storage.PutExtra{}

	token, err := repo.NewUptoken()
	if err != nil {
		return err
	}
	return formUploader.Put(context.Background(), &ret, token, key, bytes.NewReader(data), dataLen, &putExtra)
}

// PublicDownload .
func (repo *ClondRepository) PublicDownload(key string) ([]byte, error) {
	publicAccessURL := storage.MakePublicURL(conf.Get().System.QiniuDomain, key)
	data, rep := repo.NewHTTPRequest(publicAccessURL).Get().ToBytes()
	if rep.Error != nil {
		return nil, rep.Error
	}
	return data, nil
}

// PrivateDownload .
func (repo *ClondRepository) PrivateDownload(key string) ([]byte, error) {
	mac := qbox.NewMac(conf.Get().System.QiniuAccessKey, conf.Get().System.QiniuSecretKey)
	deadline := time.Now().Add(time.Second * 7200).Unix()
	privateAccessURL := storage.MakePrivateURL(mac, conf.Get().System.QiniuDomain, key, deadline)

	data, rep := repo.NewHTTPRequest(privateAccessURL).Get().ToBytes()
	if rep.Error != nil {
		return nil, rep.Error
	}
	return data, nil
}
