package types

type HttpItem struct {
	//
	Id int64 `json:"id" form:"id" gorm:"primaryKey" `
	// string format of http
	Item string `json:"item" form:"item" `
}
