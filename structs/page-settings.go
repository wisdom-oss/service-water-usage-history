package structs

type PageSettings struct {
	Size int `form:"pageSize,default=10000" binding:"min=1,max=100000"`
	Page int `form:"page,default=1" binding:"min=1"`
}
