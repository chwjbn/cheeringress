package bizmodel

import validation "github.com/go-ozzo/ozzo-validation"

type IngressNamespacePageRequest struct {
	PageRequest
	Title string `json:"title"`
	State string `bson:"state"`
}

type IngressNamespaceAddRequest struct {
	Title string `json:"title"`
}

type IngressNamespaceSaveRequest struct {
	AppDataIdRequest
	Title string `json:"title"`
	State string `bson:"state"`
}

func (this *IngressNamespaceAddRequest) Validate() error {

	xError := validation.ValidateStruct(this,
		validation.Field(&this.Title, validation.Required.Error("请输入空间名称!"), validation.Length(4, 50).Error("空间名称有效长度4-50!")),
	)

	return GetValidateError(xError)

}

func (this *IngressNamespaceSaveRequest) Validate() error {

	xError := validation.ValidateStruct(this,
		validation.Field(&this.Title, validation.Required.Error("请输入空间名称!"), validation.Length(4, 50).Error("空间名称有效长度4-50!")),
		validation.Field(&this.State, validation.Required.Error("请选择状态!"), validation.In("enable", "disable").Error("请选择正确的状态!")),
	)

	return GetValidateError(xError)

}
