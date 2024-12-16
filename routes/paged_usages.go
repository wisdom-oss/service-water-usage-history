package routes

import (
	"microservice/internal/db"
	"microservice/structs"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
)

func PagedUsages(c *gin.Context) {
	query, err := db.Queries.Raw("get-paginated")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	var records []structs.UsageRecord
	err = pgxscan.Select(c, db.Pool, &records, query, c.GetInt(KeyPageSize), c.GetInt(KeyPageOffset))
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.JSON(200, records)
}
