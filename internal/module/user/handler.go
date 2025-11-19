package user

import (
	model "core-ledger/model/core-ledger"
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

type UserHandler struct {
	logger         logger.CustomLogger
	service        *UserService
	userRepo       repo.UserRepo
	roleRepo       repo.RoleRepo
	permissionRepo repo.PermissionRepo
	dispatcher     queue.Dispatcher
}

func NewUserHandler(service *UserService, userRepo repo.UserRepo, roleRepo repo.RoleRepo, permissionRepo repo.PermissionRepo, dispatcher queue.Dispatcher) *UserHandler {
	return &UserHandler{
		logger:         logger.NewSystemLog("UserHandler"),
		service:        service,
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		dispatcher:     dispatcher,
	}
}

// User Management - CRUD
func (h *UserHandler) ListUsers(c *gin.Context) {
	filter := &repo.UserFilter{}
	if err := c.ShouldBindQuery(&filter); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.userRepo.Paginate(filter)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: res,
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetOneByFields(c, map[string]interface{}{"id": id}, "Roles", "Permissions")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginhp.RespondError(c, http.StatusNotFound, "User not found")
		} else {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: user,
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Email         string   `json:"email" binding:"required"`
		Password      string   `json:"password" binding:"required"`
		FullName      string   `json:"full_name"`
		GuardName     string   `json:"guard_name"`
		RoleIDs       []uint64 `json:"role_ids"`
		PermissionIDs []uint64 `json:"permission_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.GuardName == "" {
		req.GuardName = "web"
	}

	user := &model.User{
		Email:     req.Email,
		Password:  req.Password, // Note: Should hash password in production
		FullName:  req.FullName,
		GuardName: req.GuardName,
	}

	err := h.userRepo.Create(user)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Assign roles if provided
	if len(req.RoleIDs) > 0 {
		err = h.service.SyncUserRolesByIDs(c, user.ID, req.RoleIDs)
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Assign permissions directly if provided
	if len(req.PermissionIDs) > 0 {
		// Convert IDs to names for SyncUserPermissions
		permissions, err := h.permissionRepo.GetManyByFields(c, map[string]interface{}{"id": req.PermissionIDs})
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
			return
		}
		permissionNames := make([]string, 0, len(permissions))
		for _, perm := range permissions {
			permissionNames = append(permissionNames, perm.Name)
		}
		err = h.service.SyncUserPermissions(c, user.ID, permissionNames, req.GuardName)
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Reload user with relations
	user, err = h.userRepo.GetOneByFields(c, map[string]interface{}{"id": user.ID}, "Roles", "Permissions")
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: user,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetOneByFields(c, map[string]interface{}{"id": id}, "Roles", "Permissions")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ginhp.RespondError(c, http.StatusNotFound, "User not found")
		} else {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	var req struct {
		Email          *string  `json:"email"`
		Password       *string  `json:"password"`
		FullName       *string  `json:"full_name"`
		GuardName      string   `json:"guard_name"`
		RoleIDs        []uint64 `json:"role_ids"`
		PermissionIDs  []uint64 `json:"permission_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	updateFields := make(map[string]interface{})
	if req.Email != nil {
		updateFields["email"] = *req.Email
	}
	if req.Password != nil {
		updateFields["password"] = *req.Password // Note: Should hash password in production
	}
	if req.FullName != nil {
		updateFields["full_name"] = *req.FullName
	}
	if req.GuardName != "" {
		updateFields["guard_name"] = req.GuardName
	} else {
		req.GuardName = user.GuardName
	}

	if len(updateFields) > 0 {
		err = h.userRepo.UpdateSelectField(user, updateFields)
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Sync roles if provided
	if req.RoleIDs != nil {
		err = h.service.SyncUserRolesByIDs(c, user.ID, req.RoleIDs)
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Sync permissions directly if provided
	if req.PermissionIDs != nil {
		// Convert IDs to names for SyncUserPermissions
		permissions, err := h.permissionRepo.GetManyByFields(c, map[string]interface{}{"id": req.PermissionIDs})
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
			return
		}
		permissionNames := make([]string, 0, len(permissions))
		for _, perm := range permissions {
			permissionNames = append(permissionNames, perm.Name)
		}
		err = h.service.SyncUserPermissions(c, user.ID, permissionNames, req.GuardName)
		if err != nil {
			ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Reload user with relations
	user, err = h.userRepo.GetOneByFields(c, map[string]interface{}{"id": user.ID}, "Roles", "Permissions")
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.PreResponse{
		Data: user,
	})
}

// User Role & Permission Management
func (h *UserHandler) SyncUserRoles(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		RoleIDs []uint64 `json:"role_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.SyncUserRolesByIDs(c, uint64(id), req.RoleIDs)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginhp.RespondOK(c, "User roles synced successfully")
}

func (h *UserHandler) SyncUserPermissions(c *gin.Context) {
	id, err := utils.ParseIntIdParam(c.Param("id"))
	if err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		PermissionIDs []uint64 `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Convert IDs to names for SyncUserPermissions
	permissions, err := h.permissionRepo.GetManyByFields(c, map[string]interface{}{"id": req.PermissionIDs})
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	permissionNames := make([]string, 0, len(permissions))
	for _, perm := range permissions {
		permissionNames = append(permissionNames, perm.Name)
	}

	// Get user to get guard_name
	user, err := h.userRepo.GetByID(c, id)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.service.SyncUserPermissions(c, uint64(id), permissionNames, user.GuardName)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginhp.RespondOK(c, "User permissions synced successfully")
}

// Legacy endpoints for backward compatibility
func (h *UserHandler) GivePermissionToUser(c *gin.Context) {
	var req struct {
		UserID         uint64 `json:"user_id" binding:"required"`
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

	err := h.service.GivePermissionToUser(c, req.UserID, req.PermissionName, req.GuardName)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginhp.RespondOK(c, "Permission assigned to user successfully")
}

func (h *UserHandler) RevokePermissionFromUser(c *gin.Context) {
	var req struct {
		UserID         uint64 `json:"user_id" binding:"required"`
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

	err := h.service.RevokePermissionFromUser(c, req.UserID, req.PermissionName, req.GuardName)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginhp.RespondOK(c, "Permission revoked from user successfully")
}

func (h *UserHandler) AssignRoleToUser(c *gin.Context) {
	var req struct {
		UserID    uint64 `json:"user_id" binding:"required"`
		RoleName  string `json:"role_name" binding:"required"`
		GuardName string `json:"guard_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.GuardName == "" {
		req.GuardName = "web"
	}

	err := h.service.AssignRoleToUser(c, req.UserID, req.RoleName, req.GuardName)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginhp.RespondOK(c, "Role assigned to user successfully")
}

func (h *UserHandler) RemoveRoleFromUser(c *gin.Context) {
	var req struct {
		UserID    uint64 `json:"user_id" binding:"required"`
		RoleName  string `json:"role_name" binding:"required"`
		GuardName string `json:"guard_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ginhp.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.GuardName == "" {
		req.GuardName = "web"
	}

	err := h.service.RemoveRoleFromUser(c, req.UserID, req.RoleName, req.GuardName)
	if err != nil {
		ginhp.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginhp.RespondOK(c, "Role removed from user successfully")
}

