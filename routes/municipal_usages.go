package routes

import (
	"microservice/internal/db"
	apiErrors "microservice/internal/errors"
	"microservice/structs"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
)

func MunicipalUsages(c *gin.Context) {
	ars := strings.ReplaceAll(strings.TrimSpace(c.Param("ars")), "/", "")

	if ars == "" {
		c.Abort()
		apiErrors.ErrEmptyARS.Emit(c)
		return
	}

	if len(ars) != 12 {
		c.Abort()
		apiErrors.ErrInvalidARS.Emit(c)
		return
	}

	q, err := db.Queries.Raw("municipal-usages")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var records []structs.UsageRecord
	err = pgxscan.Select(c, db.Pool, &records, q, ars, c.GetInt(KeyPageSize), c.GetInt(KeyPageOffset))
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.JSON(200, records)

}
