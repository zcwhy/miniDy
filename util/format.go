package util

import (
	"miniDy/model"
	"strconv"
)

func DateFormat(c *model.Comment) {
	c.CreateDate = c.CreatedAt.Format("1-2")
}

func StringToInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}
