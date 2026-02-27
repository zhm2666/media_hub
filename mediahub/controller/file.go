package controller

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mediahub/pkg/config"
	"mediahub/pkg/log"
	"mediahub/pkg/storage"
	"mediahub/pkg/utils"
	"mediahub/pkg/zerror"
	"mediahub/services"
	"mediahub/services/shorturl"
	"mediahub/services/shorturl/proto"
	"net/http"
	"path"
)

type Controller struct {
	log    log.ILogger
	config *config.Config
	sf     storage.StorageFactory
}

func NewController(sf storage.StorageFactory, log log.ILogger, cnf *config.Config) *Controller {
	return &Controller{
		log:    log,
		config: cnf,
		sf:     sf,
	}
}

func (c *Controller) Upload(ctx *gin.Context) {
	userId := ctx.GetInt64("User.ID")
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		c.log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		c.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	//校验格式
	if !utils.IsImage(io.NopCloser(bytes.NewReader(content))) {
		err = zerror.NewByMsg("仅支持jpg、png、git格式")
		c.log.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	md5Digest := utils.MD5(content)
	filename := fmt.Sprintf("%x%s", md5Digest, path.Ext(fileHeader.Filename))
	filePath := "/public/" + filename
	if userId != 0 {
		filePath = fmt.Sprintf("/%d/%s", userId, filename)
	}

	s := c.sf.CreateStorage()
	url, err := s.Upload(io.NopCloser(bytes.NewReader(content)), md5Digest, filePath)
	if err != nil {
		c.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	//短链转化
	shortPool := shorturl.NewShortUrlClientPool()
	conn := shortPool.Get()
	defer shortPool.Put(conn)
	client := proto.NewShortUrlClient(conn)
	in := &proto.Url{
		Url:      url,
		UserID:   userId,
		IsPublic: userId == 0,
	}
	outGoingCtx := context.Background()
	outGoingCtx = services.AppendBearerTokenToContext(outGoingCtx, c.config.DependOn.ShortUrl.AccessToken)
	outUrl, err := client.GetShortUrl(outGoingCtx, in)
	if err != nil {
		c.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url": outUrl.Url,
	})
	return
}
