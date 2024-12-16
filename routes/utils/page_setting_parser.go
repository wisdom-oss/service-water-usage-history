package routeUtils

import (
	"microservice/structs"

	"github.com/gin-gonic/gin"

	apiErrors "microservice/internal/errors"
)

const (
	KeyPageOffset = "query.offset"
	KeyPageSize   = "query.page-size"
)

func ReadPageSettings(c *gin.Context) {
	var pageSettings structs.PageSettings

	err := c.ShouldBind(&pageSettings)
	if err != nil {
		c.Abort()
		apiErrors.ErrInvalidPageSettings.Emit(c)
		return
	}

	offset := pageSettings.Size * (pageSettings.Page - 1)

	c.Set(KeyPageOffset, offset)
	c.Set(KeyPageSize, pageSettings.Size)
}
