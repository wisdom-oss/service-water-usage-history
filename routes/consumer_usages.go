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

func ConsumerUsages(c *gin.Context) {
	consumerID := strings.ReplaceAll(strings.TrimSpace(c.Param("consumerID")), "/", "")

	if consumerID == "" {
		c.Abort()
		apiErrors.ErrEmptyConsumerID.Emit(c)
		return
	}

	if err := uuid.Validate(consumerID); err != nil {
		c.Abort()
		apiErrors.ErrInvalidConsumerID.Emit(c)
		return
	}

	var exists bool
	q, err := db.Queries.Raw("consumer-exists")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}
	err = pgxscan.Get(c, db.Pool, &exists, q, consumerID)
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	if !exists {
		c.Abort()
		apiErrors.ErrUnknownConsumer.Emit(c)
		return
	}

	q, err = db.Queries.Raw("consumer-usages")
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}
	var records []structs.UsageRecord
	err = pgxscan.Select(c, db.Pool, &records, q, consumerID, c.GetInt(KeyPageSize), c.GetInt(KeyPageOffset))
	if err != nil {
		c.Abort()
		_ = c.Error(err)
		return
	}

	c.JSON(200, records)

}
