package bizmodel

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"strings"
)

type UserInfoUpdateRequest struct {
	DataId   string `json:"data_id"`
	Nickname string `json:"nickname"`
	RealName string `json:"real_name"`
}

type UserSecurityUpdateRequest struct {
	DataId             string `json:"data_id"`
	PasswordOld        string `json:"password_old"`
	PasswordNew        string `json:"password_new"`
	PasswordNewConfirm string `json:"password_new_confirm"`
}

func (this *UserInfoUpdateRequest) Validate() error {

	xError := validation.ValidateStruct(this,
		validation.Field(&this.DataId, validation.Required.Error("请选择正确的用户!"), validation.Length(1, 50).Error("请选择正确的用户!")),
		validation.Field(&this.Nickname, validation.Required.Error("请输入用户昵称!"), validation.Length(4, 50).Error("用户昵称有效长度4-50!")),
		validation.Field(&this.RealName, validation.Required.Error("请输入用户姓名!"), validation.Length(4, 50).Error("用户姓名有效长度4-50!")),
	)

	if xError != nil {
		return xError
	}

	return xError

}

func (this *UserSecurityUpdateRequest) Validate() error {

	xError := validation.ValidateStruct(this,
		validation.Field(&this.DataId, validation.Required.Error("请选择正确的用户!"), validation.Length(1, 50).Error("请选择正确的用户!")),
		validation.Field(&this.PasswordOld, validation.Required.Error("请输入当前登录密码!"), validation.Length(6, 50).Error("当前登录密码有效长度6-50!")),
		validation.Field(&this.PasswordNew, validation.Required.Error("请输入修改后密码!"), validation.Length(6, 50).Error("修改后密码有效长度6-50!")),
	)

	if xError != nil {
		return xError
	}

	if !strings.EqualFold(this.PasswordNew, this.PasswordNewConfirm) {
		xError = errors.New("修改后密码两次输入不一致!")
		return xError
	}

	return xError

}
