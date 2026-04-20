package handler

import (
	"github.com/bloansbook/bloansbook-api/internal/models/roles"
	"github.com/bloansbook/bloansbook-api/internal/roles/usecase"
	"github.com/bloansbook/bloansbook-api/pkg/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type RolesHandler struct {
	usecase *usecase.RolesUsecase
}

func NewRolesHandler(u *usecase.RolesUsecase) *RolesHandler {
	return &RolesHandler{
		usecase: u,
	}
}

func (r *RolesHandler) CreateRole(c fiber.Ctx) error {
	// Parse request body
	var payload roles.CreateRolePayload
	if err := c.Bind().Body(&payload); err != nil {
		return response.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	// Create role with usecase
	role, err := r.usecase.CreateRole(c.Context(), &payload)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	// Return response
	return response.Success(c, "Role created successfully", role, fiber.StatusCreated)
}

func (r *RolesHandler) CreatePermission(c fiber.Ctx) error {
	// Parse request body
	var payload roles.CreatePermissionPayload
	if err := c.Bind().Body(&payload); err != nil {
		return response.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	// Create Permission with usecase
	permission, err := r.usecase.CreatePermission(c.Context(), &payload)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	// Return response
	return response.Success(c, "Permission created successfully", permission, fiber.StatusCreated)
}

func (r *RolesHandler) AssignPermissionToRole(c fiber.Ctx) error {
	// Parse request body
	var payload roles.CreateRolePermissionPayload
	if err := c.Bind().Body(&payload); err != nil {
		return response.Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	// Assign permission to role with usecase
	rolePermission, err := r.usecase.AssignPermissionToRole(c.Context(), &payload)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	// Return response
	return response.Success(c, "Permission assigned to role successfully", rolePermission, fiber.StatusCreated)
}

func (r *RolesHandler) GetAllRoles(c fiber.Ctx) error {
	limit := fiber.Query[int](c, "limit", 10)
	offset := fiber.Query[int](c, "offset", 0)

	roles, err := r.usecase.GetAllRoles(c.Context(), limit, offset)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, "Roles fetched successfully", roles, fiber.StatusOK)
}

func (r *RolesHandler) GetRoleWithPermissions(c fiber.Ctx) error {
	roleID := c.Params("id")
	if roleID == "" {
		return response.Error(c, "RoleID is required", fiber.StatusBadRequest)
	}

	id, err := uuid.Parse(roleID)
	if err != nil {
		return response.Error(c, "Invalid RoleID format", fiber.StatusBadRequest)
	}

	role, err := r.usecase.GetRoleWithPermissions(c.Context(), id)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, "Role fetched successfully", role, fiber.StatusOK)
}
