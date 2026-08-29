package createtask

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"main/internal/domain/contracts"
	"main/internal/domain/models"
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
	var request dto.CreateTaskRequest

	err := ginCtx.ShouldBindJSON(&request)
	if err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	err = validateConfig(models.TaskType(request.Type), request.Config)
	if err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	interval, err := time.ParseDuration(request.Interval)
	if err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	t := models.TaskCreation{
		Name:     request.Name,
		Type:     models.TaskType(request.Type),
		Active:   request.Active,
		Interval: interval,
		Config:   request.Config,
	}

	task, err := h.tasksRepo.Insert(ginCtx.Request.Context(), t)
	if err != nil {
		h.apiErrorHandler.Handle(ginCtx, err)

		return
	}

	ginCtx.JSON(http.StatusCreated, dto.ToTaskResponse(task))
}

func validateConfig(taskType models.TaskType, config json.RawMessage) error {
	switch taskType {
	case models.TaskTypeCeneoCatcher:
		var cfg models.CeneoCatcherConfig

		if err := json.Unmarshal(config, &cfg); err != nil {
			return fmt.Errorf("invalid ceneo catcher config: %w", err)
		}

	case models.TaskTypeDoctorReminder:
		var cfg models.DoctorReminderConfig

		if err := json.Unmarshal(config, &cfg); err != nil {
			return fmt.Errorf("invalid doctor reminder config: %w", err)
		}

	default:
		return fmt.Errorf("unknown task type: %q", taskType)
	}

	return nil
}
