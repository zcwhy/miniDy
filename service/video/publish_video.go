package video

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"miniDy/model"
	"miniDy/util"
	"path/filepath"
)

// 空struct构造set
var videoIndexMap = map[string]struct{}{
	".mp4":  {},
	".avi":  {},
	".wmv":  {},
	".flv":  {},
	".mpeg": {},
	".mov":  {},
}

type PublishVideoRequest struct {
	Context *gin.Context

	UserId int64
	Video  *multipart.FileHeader
	Title  string
}

type PublishVideoService struct {
	request *PublishVideoRequest
}

func PublishVideo(r *PublishVideoRequest) error {
	return NewPublishVideoService(r).Do()
}

func NewPublishVideoService(r *PublishVideoRequest) *PublishVideoService {
	return &PublishVideoService{request: r}
}

func (p *PublishVideoService) Do() error {
	if err := p.checkParam(); err != nil {
		return err
	}

	if err := p.publishVideo(); err != nil {
		return err
	}

	return nil
}

func (p *PublishVideoService) checkParam() error {
	fileName := p.request.Video.Filename
	suffix := filepath.Ext(fileName)

	if _, ok := videoIndexMap[suffix]; !ok {
		return errors.New("不支持的视频格式")
	}

	return nil
}

func (p *PublishVideoService) publishVideo() error {
	file := p.request.Video
	userId := p.request.UserId

	//执行dao层，查询用户发布的视频数量用来构造存储时的文件名
	videoDao := model.NewVideoDao()
	videoCount, err := videoDao.CountUserVideoById(userId)

	if err != nil {
		return err
	}

	suffix := filepath.Ext(p.request.Video.Filename)
	saveName := fmt.Sprintf("%d-%d%s", userId, videoCount, suffix)
	if err != nil {
		return err
	}

	//本地磁盘存储视频
	savePath := filepath.Join("static", saveName)
	if err := p.request.Context.SaveUploadedFile(file, savePath); err != nil {
		return err
	}

	//数据库存储视频的元信息（meta data）
	videoUrl := util.GetFileURL(saveName)
	if err := videoDao.CreateVideo(userId, videoUrl, p.request.Title); err != nil {
		return err
	}

	return nil
}
