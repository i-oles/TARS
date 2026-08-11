package updatetask

import (
	"errors"
	"net/http"

	"main/internal/domain/contracts"
	"main/internal/interfaces/http/api/dto"
	apiErrs "main/internal/interfaces/http/api/errs"

	"github.com/gin-gonic/gin"
)

type handler struct {
	tasksRepo       contracts.ITasks
	apiErrorHandler apiErrs.IErrorHandler
}

func NewHandler(
	tasksRepo contracts.ITasks,
	apiErrorHandler apiErrs.IErrorHandler,
) *handler {
	return &handler{
		tasksRepo:       tasksRepo,
		apiErrorHandler: apiErrorHandler,
	}
}

func (h *handler) Handle(ginCtx *gin.Context) {
	var request dto.UpdateTaskRequest

	err := ginCtx.ShouldBindJSON(&request)
	if err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	var uri dto.UpdateTaskURI

	if err := ginCtx.ShouldBindUri(&uri); err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	update, err := getDataForTaskUpdate(request)
	if err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	task, err := h.tasksRepo.Update(ginCtx.Request.Context(), uri.ID, update)
	if err != nil {
		h.apiErrorHandler.Handle(ginCtx, err)

		return
	}

	ginCtx.JSON(http.StatusCreated, dto.ToTaskResponse(task))
}

func getDataForTaskUpdate(req dto.UpdateTaskRequest) (map[string]any, error) {
	updateData := map[string]any{}

	if req.Active != nil {
		updateData["active"] = *req.Active
	}

	if req.LastRunAt != nil {
		updateData["last_run_at"] = *req.LastRunAt
	}

	if len(updateData) == 0 {
		return nil, errors.New("no fields to update class")
	}

	return updateData, nil
}
