package routes

import (
	"microservice/internal/db"
	apiErrors "microservice/internal/errors"
	"microservice/structs"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TypedUsages(c *gin.Context) {
	usageTypeID := strings.ReplaceAll(strings.TrimSpace(c.Param("usageTypeID")), "/", "")

	if usageTypeID == "" {
		c.Abort()
		apiErrors.ErrEmptyUsageTypeID.Emit(c)
		return
	}

	if err := uuid.Validate(usageTypeID); err != nil {
		c.Abort()
		apiErrors.ErrInvalidUsageTypeID.Emit(c)
		return
	}

	q, err := db.Queries.Raw("typed-usages")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}
	var records []structs.UsageRecord
	err = pgxscan.Select(c, db.Pool, &records, q, usageTypeID, c.GetInt(KeyPageSize), c.GetInt(KeyPageOffset))
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.JSON(200, records)
}
