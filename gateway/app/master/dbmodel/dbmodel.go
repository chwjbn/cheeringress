package dbmodel

type PageData struct {
	DataList   []IDbModelMongo `json:"data_list"`
	PageCount  int64           `json:"page_count"`
	PageSize   int64           `json:"page_size"`
	PageNo     int64           `json:"page_no"`
	TotalCount int64           `json:"total_count"`
}

func (this *PageData) Calc() {

	if this.PageNo < 1 {
		this.PageNo = 1
	}

	if this.PageSize < 1 {
		this.PageSize = 100
	}

	this.PageCount = int64(this.TotalCount / this.PageSize)
	if this.TotalCount%this.PageSize > 0 {
		this.PageCount = this.PageCount + 1
	}

	if this.PageNo > this.PageCount {
		this.PageNo = this.PageCount
	}

	this.DataList = []IDbModelMongo{}
}
