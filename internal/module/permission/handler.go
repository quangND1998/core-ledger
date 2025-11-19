package permission

import (
	"core-ledger/model/dto"
	"core-ledger/pkg/ginhp"
	"core-ledger/pkg/logger"
	"core-ledger/pkg/queue"
	"core-ledger/pkg/repo"
	"core-ledger/pkg/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PermissionHandler struct {
	logger         logger.CustomLogger
	permissionRepo repo.PermissionRepo
	dispatcher     queue.Dispatcher
}

func NewPermissionHandler(permissionRepo repo.PermissionRepo, dispatcher queue.Dispatcher) *PermissionHandler {
	return &PermissionHandler{
		logger:         logger.NewSystemLog("PermissionHandler"),
		permissionRepo: permissionRepo,
		dispatcher:     dispatcher,
	}
}

// Permission Management - CRUD
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	q := &dto.BasePaginationQuery{}
	if err := c.ShouldBindQuery(&q); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.permissionRepo.Paginate(q)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: res,
	})
}

func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	permission, err := h.permissionRepo.GetByID(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginhp.RespondError(c, http.StatusNotFound, "Permission not found")
		} else {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: permission,
	})
}

func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		GuardName string `json:"guard_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.GuardName == "" {
		req.GuardName = "web"
	}

	permission, err := h.permissionRepo.FindOrCreate(c, req.Name, req.GuardName)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: permission,
	})
}

func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	permission, err := h.permissionRepo.GetByID(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginhp.RespondError(c, http.StatusNotFound, "Permission not found")
		} else {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	var req struct {
		Name      string `json:"name" binding:"required"`
		GuardName string `json:"guard_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	updateFields := make(map[string]interface{})
	updateFields["name"] = req.Name
	if req.GuardName != "" {
		updateFields["guard_name"] = req.GuardName
	}

	err = h.permissionRepo.UpdateSelectField(permission, updateFields)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Reload permission
	permission, err = h.permissionRepo.GetByID(c, id)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: permission,
	})
}

func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	permission, err := h.permissionRepo.GetByID(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginhp.RespondError(c, http.StatusNotFound, "Permission not found")
		} else {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	err = h.permissionRepo.Delete(permission)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginhp.RespondOK(c, "Permission deleted successfully")
}
