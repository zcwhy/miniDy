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

	//执行dao层，查询用户发布的视频数量用来构造存储时的文件名， 并查询当前用户是否已经发布过此标题的视频
	videoDao := model.NewVideoDao()
	isTitleExist, err := videoDao.IsVideoTitleExistById(userId, p.request.Title)

	if err != nil {
		return err
	}

	if isTitleExist == true {
		return errors.New("当前发布视频标题重复")
	}

	var videoCount int64
	err = videoDao.CountUserVideoById(userId, &videoCount)

	if err != nil {
		return err
	}

	//根据videoCount来生成存储的文件名
	suffix := filepath.Ext(p.request.Video.Filename)
	saveVideoName := fmt.Sprintf("%d-%d%s", userId, videoCount, suffix)
	if err != nil {
		return err
	}

	//本地磁盘存储视频
	saveFilePath := filepath.Join("static", saveVideoName)
	if err := p.request.Context.SaveUploadedFile(file, saveFilePath); err != nil {
		return err
	}

	//截取视频的一帧作为封面，本地磁盘存储封面
	saveCoverName := saveVideoName + "_cover"
	if err := util.GetCoverFromVideo(saveVideoName, saveCoverName); err != nil {
		return err
	}

	//数据库存储视频的元信息（meta data）
	videoUrl := util.GetVideoURL(saveVideoName)
	coverUrl := util.GetCoverURL(saveCoverName)
	if err := videoDao.CreateVideo(userId, videoUrl, coverUrl, p.request.Title); err != nil {
		return err
	}

	//更新user的work_count视频数量的信息
	userInfoDao := model.NewUserInfoDao()
	if err := userInfoDao.AddUserWorkCount(userId); err != nil {
		return err
	}

	return nil
}
