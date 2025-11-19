package role

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

type RoleHandler struct {
	logger         logger.CustomLogger
	service        *RoleService
	roleRepo       repo.RoleRepo
	permissionRepo repo.PermissionRepo
	dispatcher     queue.Dispatcher
}

func NewRoleHandler(service *RoleService, roleRepo repo.RoleRepo, permissionRepo repo.PermissionRepo, dispatcher queue.Dispatcher) *RoleHandler {
	return &RoleHandler{
		logger:         logger.NewSystemLog("RoleHandler"),
		service:        service,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		dispatcher:     dispatcher,
	}
}

// Role Management - CRUD
func (h *RoleHandler) ListRoles(c *gin.Context) {
	q := &dto.BasePaginationQuery{}
	if err := c.ShouldBindQuery(&q); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.roleRepo.Paginate(q)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: res,
	})
}

func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.roleRepo.GetByID(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginhp.RespondError(c, http.StatusNotFound, "Role not found")
		} else {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: role,
	})
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name          string   `json:"name" binding:"required"`
		GuardName     string   `json:"guard_name"`
		PermissionIDs []uint64 `json:"permission_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.GuardName == "" {
		req.GuardName = "web"
	}

	role, err := h.roleRepo.FindOrCreate(c, req.Name, req.GuardName)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Sync permissions if provided
	if len(req.PermissionIDs) > 0 {
		err = h.service.SyncRolePermissionsByIDs(c, role.ID, req.PermissionIDs)
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Reload role with permissions
	role, err = h.roleRepo.GetByID(c, int64(role.ID))
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: role,
	})
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.roleRepo.GetByID(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginhp.RespondError(c, http.StatusNotFound, "Role not found")
		} else {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	var req struct {
		Name          string   `json:"name" binding:"required"`
		GuardName     string   `json:"guard_name"`
		PermissionIDs []uint64 `json:"permission_ids"`
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

	err = h.roleRepo.UpdateSelectField(role, updateFields)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Sync permissions if provided
	if req.PermissionIDs != nil {
		err = h.service.SyncRolePermissionsByIDs(c, role.ID, req.PermissionIDs)
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Reload role with permissions
	role, err = h.roleRepo.GetByID(c, int64(role.ID))
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: role,
	})
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.roleRepo.GetByID(c, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginhp.RespondError(c, http.StatusNotFound, "Role not found")
		} else {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	err = h.roleRepo.Delete(role)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginhp.RespondOK(c, "Role deleted successfully")
}

// Role Permission Management
func (h *RoleHandler) GivePermissionToRole(c *gin.Context) {
	var req struct {
		RoleID         uint64 `json:"role_id" binding:"required"`
		PermissionName string `json:"permission_name" binding:"required"`
		GuardName      string `json:"guard_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.GuardName == "" {
		req.GuardName = "web"
	}

	err := h.service.GivePermissionToRole(c, req.RoleID, req.PermissionName, req.GuardName)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginhp.RespondOK(c, "Permission assigned to role successfully")
}

func (h *RoleHandler) RevokePermissionFromRole(c *gin.Context) {
	var req struct {
		RoleID         uint64 `json:"role_id" binding:"required"`
		PermissionName string `json:"permission_name" binding:"required"`
		GuardName      string `json:"guard_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.GuardName == "" {
		req.GuardName = "web"
	}

	err := h.service.RevokePermissionFromRole(c, req.RoleID, req.PermissionName, req.GuardName)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginhp.RespondOK(c, "Permission revoked from role successfully")
}

func (h *RoleHandler) SyncRolePermissions(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var req struct {
		PermissionIDs []uint64 `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.SyncRolePermissionsByIDs(c, uint64(id), req.PermissionIDs)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginhp.RespondOK(c, "Role permissions synced successfully")
}

func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid role ID")
		return
	}

	permissions, err := h.roleRepo.GetPermissions(c, uint64(id))
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: permissions,
	})
}

