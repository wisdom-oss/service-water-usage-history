package routes

import (
	"microservice/internal/db"
	apiErrors "microservice/internal/errors"
	"microservice/structs"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
)

func PagedUsages(c *gin.Context) {
	var parameter structs.PageSettings
	if err := c.ShouldBind(&parameter); err != nil {
		c.Abort()
		apiErrors.ErrInvalidPageSettings.Emit(c)
		return
	}

	query, err := db.Queries.Raw("get-paginated")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	offset := parameter.Size * (parameter.Page - 1)

	var records []structs.UsageRecord
	err = pgxscan.Select(c, db.Pool, &records, query, parameter.Size, offset)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.JSON(200, records)
}
