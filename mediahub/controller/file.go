package controller

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mediahub/pkg/config"
	"mediahub/pkg/log"
	"mediahub/pkg/storage"
	"mediahub/pkg/zerror"
	"mediahub/services"
	"mediahub/services/shorturl"
	"mediahub/services/shorturl/proto"
	"net/http"
	"path"
)

type Controller struct {
	sf     storage.StorageFactory
	log    log.ILogger
	config *config.Config
}

func NewController(sf storage.StorageFactory, log log.ILogger, cnf *config.Config) *Controller {
	return &Controller{
		sf:     sf,
		log:    log,
		config: cnf,
	}
}
func (c *Controller) Upload(ctx *gin.Context) {
	userId := ctx.GetInt64("User.ID")
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		c.log.Error(zerror.NewByErr(err))
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.log.Error(zerror.NewByErr(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		c.log.Error(zerror.NewByErr(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	//校验图片格式
	if !isImage(io.NopCloser(bytes.NewReader(content))) {
		err = zerror.NewByMsg("仅支持jpg、png、gif格式")
		c.log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	md5Digest := calMD5Digest(content)
	filename := fmt.Sprintf("%x%s", md5Digest, path.Ext(fileHeader.Filename))
	filePath := "/public/" + filename
	if userId != 0 {
		filePath = fmt.Sprintf("/%d/%s", userId, filename)
	}
	s := c.sf.CreateStorage()
	url, err := s.Upload(io.NopCloser(bytes.NewReader(content)), md5Digest, filePath)
	if err != nil {
		c.log.Error(zerror.NewByErr(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	shortPool := shorturl.NewShortUrlClientPool()
	clientConn := shortPool.Get()
	defer shortPool.Put(clientConn)

	client := proto.NewShortUrlClient(clientConn)
	in := &proto.Url{
		Url:      url,
		UserID:   userId,
		IsPublic: userId == 0,
	}
	outGoingCtx := context.Background()
	outGoingCtx = services.AppendBearerTokenToContext(outGoingCtx, c.config.DependOn.ShortUrl.AccessToken)
	outUrl, err := client.GetShortUrl(outGoingCtx, in)
	if err != nil {
		c.log.Error(zerror.NewByErr(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url": outUrl.Url,
	})
	return
}

func isImage(r io.Reader) bool {
	_, _, err := image.Decode(r)
	if err != nil {
		return false
	}
	return true
}
func calMD5Digest(msg []byte) []byte {
	m := md5.New()
	m.Write(msg)
	bs := m.Sum(nil)
	return bs
}
