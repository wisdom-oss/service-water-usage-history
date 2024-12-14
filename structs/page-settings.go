package structs

type PageSettings struct {
	Size int `form:"pageSize,default=10000" default:"10000" validate:"min=-1"`
	Page int `form:"page,default=1" validate:"gte=1" default:"1"`
}
