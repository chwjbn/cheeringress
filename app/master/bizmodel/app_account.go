package bizmodel

import validation "github.com/go-ozzo/ozzo-validation"

type CheckCodeImageData struct {
	CodeId    string `json:"code_id"`
	ImageData string `json:"image_data"`
}

type AccountLoginRequest struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	ImageCodeId   string `json:"image_code_id"`
	ImageCodeData string `json:"image_code_data"`
}

func (this *AccountLoginRequest) Validate() error {

	xError := validation.ValidateStruct(this,
		validation.Field(&this.Username, validation.Required.Error("请输入登录用户名!"), validation.Length(4, 50).Error("登录用户名有效长度4-50!")),
		validation.Field(&this.Password, validation.Required.Error("请输入登录密码!"), validation.Length(6, 50).Error("登录密码有效长度6-50!")),
		validation.Field(&this.ImageCodeId, validation.Required.Error("请先刷新图形验证码!")),
		validation.Field(&this.ImageCodeData, validation.Required.Error("请输入图形验证码!")),
	)

	return GetValidateError(xError)

}
