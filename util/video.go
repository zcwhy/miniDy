package util

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
)

var ip = "192.168.138.31:8080"

func GetVideoURL(videoName string) string {
	videoURL := fmt.Sprintf("http://%s/static/%s", ip, videoName)
	return videoURL
}

func GetCoverURL(coverName string) string {
	coverURL := fmt.Sprintf("http://%s/static/%s", ip, coverName)
	return coverURL
}

func GetCoverFromVideo(videoName, coverName string) error {
	videoPath := "static/" + videoName
	coverPath := "static/covers/" + coverName
	buf := bytes.NewBuffer(nil)

	err := ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}). //1 表示截取视频的第一帧作为封面
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败1：", err)
		return err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败2：", err)
		return err
	}

	err = imaging.Save(img, coverPath+".png")
	if err != nil {
		log.Fatal("生成缩略图失败3：", err)
		return err
	}

	return nil
}
