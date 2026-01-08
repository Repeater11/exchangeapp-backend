package service

import "errors"

var ErrForbidden = errors.New("无权限")
var ErrThreadNotFound = errors.New("帖子不存在")
var ErrReplyNotFound = errors.New("回复不存在")
