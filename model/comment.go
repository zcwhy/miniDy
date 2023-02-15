package model

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	Id         int64     `json:"id"`
	UserInfoId int64     `json:"-"` //用于一对多关系的id
	VideoId    int64     `json:"-"` //一对多，视频对评论
	User       UserInfo  `json:"user" gorm:"-"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"-"`
	CreateDate string    `json:"create_date" gorm:"-"`
}

type CommentDAO struct{}

var commentDAO CommentDAO

func NewCommentDAO() *CommentDAO {
	return &commentDAO
}

func (c *CommentDAO) AddComment(comment *Comment) error {
	if comment == nil {
		return errors.New("receiving nullptr when adding comment")
	}
	//Execution Services
	return DB.Transaction(func(tx *gorm.DB) error {
		//add comment
		if err := tx.Create(comment).Error; err != nil {
			// 返回任何错误都会回滚事务
			//return err
		}
		//update video comment-count
		err := tx.Model(&Video{}).Where("id = ?", comment.VideoId).
			Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
		if err != nil {
			//return err
		}
		return nil
	})
}

func (c *CommentDAO) DeleteComment(comment *Comment) error {
	if comment == nil {
		return errors.New("receiving nullptr when deleting comment")
	}
	//Execution Services
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Comment{}).Delete(comment).Error; err != nil {
			//return err
		}
		err := tx.Model(&Comment{}).Where("id = ?", comment.VideoId).
			Update("comment_count", gorm.Expr("comment_count + ?", -1)).Error
		if err != nil {
			//return err
		}
		return nil
	})
}

//func (c *CommentDAO) QueryCommentById(id int64, comment *Comment) error {
//	if comment == nil {
//		return errors.New("receiving nullptr when querying comment by id")
//	}
//	//Execution Services
//	_ = DB.First(comment, id).Error
//	return nil
//}

func (c *CommentDAO) QueryCommentListByVideoId(id int64, comment *[]*Comment) error {
	if comment == nil {
		return errors.New("receiving nullptr when querying comment by video")
	}
	//Execution Services
	_ = DB.First(comment).Where("video_id=?", id).Error
	return nil
}
