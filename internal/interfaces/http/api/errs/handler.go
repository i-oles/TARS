package errs

import (
	"errors"
	"net/http"

	"main/internal/domain/errs/api"

	"github.com/gin-gonic/gin"
)

type IErrorHandler interface {
	Handle(ctx *gin.Context, err error)
}

type errorHandler struct{}

func NewErrorHandler() errorHandler {
	return errorHandler{}
}

func (e errorHandler) Handle(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, api.ErrTaskNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
