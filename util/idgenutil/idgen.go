package idgenutil

import (
	"github.com/mojocn/base64Captcha"
)

func ID() string {
	return base64Captcha.RandText(17, "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func IDNum() string {
	return base64Captcha.RandText(5, "0123456789")
}
