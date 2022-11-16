package bizmodel

import validation "github.com/go-ozzo/ozzo-validation"

type IngressActionStaticPageRequest struct {
	PageRequest
	NamespaceId string `json:"namespace_id"`
	Title       string `json:"title"`
	State       string `bson:"state"`
}

type IngressActionStaticAddRequest struct {
	NamespaceId string `json:"namespace_id"`
	Title       string `json:"title"`
	ContentType string `json:"content_type"`
	DataType    string `json:"data_type"`
	Data        string `json:"data"`
}

type IngressActionStaticSaveRequest struct {
	AppDataIdRequest
	NamespaceId string `json:"namespace_id"`
	Title       string `json:"title"`
	ContentType string `json:"content_type"`
	DataType    string `json:"data_type"`
	Data        string `json:"data"`
	State       string `bson:"state"`
}

func (this *IngressActionStaticAddRequest) Validate() error {

	xError := validation.ValidateStruct(this,
		validation.Field(&this.NamespaceId, validation.Required.Error("请选择网关空间!")),
		validation.Field(&this.Title, validation.Required.Error("请输入资源名称!"), validation.Length(4, 50).Error("资源名称有效长度4-50!")),
		validation.Field(&this.ContentType, validation.Required.Error("请选择资源类型!")),
		validation.Field(&this.DataType, validation.Required.Error("请选择数据类型!"), validation.In("PlainText", "Base64Data", "HttpResContent", "HttpResZip").Error("请选择正确的数据类型!")),
		validation.Field(&this.Data, validation.Required.Error("请输入数据内容!"), validation.Length(0, 1024).Error("数据内容最大长度为1024!")),
	)

	return GetValidateError(xError)

}

func (this *IngressActionStaticSaveRequest) Validate() error {

	xError := validation.ValidateStruct(this,
		validation.Field(&this.NamespaceId, validation.Required.Error("请选择网关空间!")),
		validation.Field(&this.Title, validation.Required.Error("请输入资源名称!"), validation.Length(4, 50).Error("资源名称有效长度4-50!")),
		validation.Field(&this.ContentType, validation.Required.Error("请选择资源类型!")),
		validation.Field(&this.DataType, validation.Required.Error("请选择数据类型!"), validation.In("PlainText", "Base64Data", "HttpResContent", "HttpResZip").Error("请选择正确的数据类型!")),
		validation.Field(&this.Data, validation.Required.Error("请输入数据内容!"), validation.Length(0, 1024).Error("数据内容最大长度为1024!")),
		validation.Field(&this.State, validation.Required.Error("请选择状态!"), validation.In("enable", "disable").Error("请选择正确的状态!")),
	)

	return GetValidateError(xError)

}
