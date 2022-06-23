package webrequest

type CategoryUpdateRequest struct {
	Id   int64  `validate:"required" json:"id"`
	Name string `validate:"required,max=200,min=1" json:"name"`
}
