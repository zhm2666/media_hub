package server

import (
	"context"
	"fmt"
	"shorturl/pkg/config"
	"shorturl/pkg/constants"
	"shorturl/pkg/log"
	"shorturl/pkg/utils"
	"shorturl/pkg/zerror"
	"shorturl/proto"
	"shorturl/shorturl-server/cache"
	"shorturl/shorturl-server/data"
	"strconv"
	"time"
)

type shortUrlService struct {
	proto.UnimplementedShortUrlServer
	config            *config.Config
	log               log.ILogger
	urlMapDataFactory data.IUrlMapDataFactory
	kvCacheFactory    cache.CacheFactory
}

func NewService(config *config.Config, log log.ILogger, urlMapDataFactory data.IUrlMapDataFactory, kvCacheFactory cache.CacheFactory) proto.ShortUrlServer {
	return &shortUrlService{
		config:            config,
		log:               log,
		urlMapDataFactory: urlMapDataFactory,
		kvCacheFactory:    kvCacheFactory,
	}
}

func (s *shortUrlService) GetShortUrl(ctx context.Context, in *proto.Url) (*proto.Url, error) {
	isPublic := in.IsPublic
	if in.UserID != 0 {
		isPublic = false
	}
	if in.Url == "" {
		err := zerror.NewByMsg("参数检查失败")
		s.log.Error(err)
		return nil, err
	}
	if !utils.IsUrl(in.Url) {
		err := zerror.NewByMsg("参数检查失败")
		s.log.Error(err)
		return nil, err
	}
	// 根据长链接查询数据库，记录是否已存在
	data := s.urlMapDataFactory.NewUrlMapData(isPublic)
	entity, err := data.GetByOriginal(in.Url)
	if err != nil {
		s.log.Error(zerror.NewByErr(err))
		return nil, err
	}
	if entity.ShortKey == "" {
		//新增记录
		id, err := data.GenerateID(in.GetUserID(), time.Now().Unix())
		if err != nil {
			s.log.Error(zerror.NewByErr(err))
			return nil, err
		}
		entity.ShortKey = utils.ToBase62(id)
		entity.OriginalUrl = in.Url
		entity.ID = id
		entity.UpdateAt = time.Now().Unix()
		err = data.Update(entity)
		if err != nil {
			s.log.Error(zerror.NewByErr(err))
			return nil, err
		}
	}
	keyPrefix := ""
	domain := s.config.ShortDomain
	if !isPublic {
		keyPrefix = "user_"
		domain = s.config.UserShortDomain
	}
	kvCache := s.kvCacheFactory.NewKVCache()
	defer kvCache.Destroy()
	key := keyPrefix + entity.ShortKey
	err = kvCache.Set(key, entity.OriginalUrl, cache.DefaultTTL)
	if err != nil {
		s.log.Error(zerror.NewByErr(err))
		return nil, err
	}
	return &proto.Url{
		Url:    domain + entity.ShortKey,
		UserID: in.UserID,
	}, nil
}
func (s *shortUrlService) GetOriginalUrl(ctx context.Context, in *proto.ShortKey) (*proto.Url, error) {
	isPublic := in.IsPublic
	if in.UserID != 0 {
		isPublic = false
	}
	if in.Key == "" {
		err := zerror.NewByMsg("参数检查失败")
		s.log.Error(err)
		return nil, err
	}
	id := utils.ToBase10(in.Key)
	if id == 0 {
		err := zerror.NewByMsg("参数检查失败")
		s.log.Error(err)
		return nil, err
	}

	keyPrefix := ""
	if !isPublic {
		keyPrefix = "user_"
	}
	kvCache := s.kvCacheFactory.NewKVCache()
	defer kvCache.Destroy()
	key := keyPrefix + in.Key

	data := s.urlMapDataFactory.NewUrlMapData(isPublic)
	originalUrl, err := kvCache.Get(key)
	if err != nil {
		s.log.Error(err)
		return nil, zerror.NewByErr(err)
	}
	if originalUrl == "" {
		// 缓存穿透过滤
		err = s.idFilter(id, kvCache, isPublic)
		if err != nil {
			s.log.Error(err)
			return nil, err
		}

		entity, err := data.GetByID(id)
		if err != nil {
			s.log.Error(err)
			return nil, zerror.NewByErr(err)
		}
		originalUrl = entity.OriginalUrl
	}
	err = kvCache.Set(key, originalUrl, cache.DefaultTTL)
	if err != nil {
		s.log.Error(err)
		return nil, zerror.NewByErr(err)
	}
	err = data.IncrementTimes(id, 1, time.Now().Unix())
	if err != nil {
		s.log.Warning(err)
		err = nil
	}
	return &proto.Url{
		Url:    originalUrl,
		UserID: in.UserID,
	}, nil
}

func (s *shortUrlService) idFilter(id int64, kvCache cache.KVCache, isPublic bool) error {
	key := fmt.Sprintf("%s_%s", constants.TABLENAME_URL_MAP, "max_id")
	if !isPublic {
		key = fmt.Sprintf("%s_%s", constants.TABLENAME_URL_MAP_USER, "max_id")
	}
	idStr, err := kvCache.Get(key)
	if err != nil {
		s.log.Error(err)
		return err
	}
	var res int64
	if idStr != "" {
		res, err = strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			s.log.Error(err)
			return err
		}
	}
	if res < id {
		err = zerror.NewByMsg("短链非法")
		s.log.Error(err)
		return err
	}
	return nil
}
