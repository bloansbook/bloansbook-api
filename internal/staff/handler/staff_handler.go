package handler

import (
	"github.com/bloansbook/bloansbook-api/internal/auth/middleware"
	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/bloansbook/bloansbook-api/internal/models/staff"
	"github.com/bloansbook/bloansbook-api/internal/staff/usecase"
	"github.com/bloansbook/bloansbook-api/pkg/response"
	"github.com/bloansbook/bloansbook-api/pkg/sysmsg"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type StaffHandler struct {
	usecase *usecase.StaffUsecase
}

func NewStaffHandler(u *usecase.StaffUsecase) *StaffHandler {
	return &StaffHandler{
		usecase: u,
	}
}

func (h *StaffHandler) CreateStaff(c fiber.Ctx) error {
	var payload staff.CreateStaffPayload
	createdBy, ok := middleware.CallerStaffID(c)
	if !ok {
		return response.Error(c, sysmsg.Unauthorized, fiber.StatusUnauthorized)
	}

	id, ok := createdBy.(uuid.UUID)
	if !ok {
		return response.Error(c, sysmsg.Unauthorized, fiber.StatusUnauthorized)
	}

	if err := c.Bind().Body(&payload); err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	staff, err := h.usecase.CreateStaff(c.Context(), id, &payload)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.StaffCreated, &staff, fiber.StatusCreated)
}

func (h *StaffHandler) GetStaffById(c fiber.Ctx) error {
	staffId := c.Params("id")
	if staffId == "" {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	id, err := uuid.Parse(staffId)
	if err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	staff, err := h.usecase.GetStaffByID(c.Context(), id)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.StaffFetched, staff, fiber.StatusOK)
}

func (h *StaffHandler) GetAllStaff(c fiber.Ctx) error {
	totalCount, err := h.usecase.GetStaffCount(c.Context())
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	limit := fiber.Query(c, "limit", 10)
	offset := fiber.Query(c, "offset", 0)

	staffList, err := h.usecase.GetAllStaff(c.Context(), limit, offset)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	data := models.DataWithPagination{
		Data:       staffList,
		Count:      len(staffList),
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
	}

	return response.Success(c, sysmsg.StaffListFetched, data, fiber.StatusOK)
}

func (h *StaffHandler) UpdateStaff(c fiber.Ctx) error {
	staffID := c.Params("id")

	var payload staff.UpdateStaffPayload
	if err := c.Bind().Body(&payload); err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	if staffID == "" {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	id, err := uuid.Parse(staffID)
	if err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	staff, err := h.usecase.UpdateStaff(c.Context(), id, &payload)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.StaffUpdated, staff, fiber.StatusOK)
}

// --- Staff Role Handlers ---

func (h *StaffHandler) AssignRole(c fiber.Ctx) error {
	staffID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	performedBy, ok := middleware.CallerStaffID(c)
	if !ok {
		return response.Error(c, sysmsg.Unauthorized, fiber.StatusUnauthorized)
	}
	callerID := performedBy.(uuid.UUID)

	var payload staff.AssignRolePayload
	if err := c.Bind().Body(&payload); err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	result, err := h.usecase.AssignRole(c.Context(), staffID, payload.RoleID, callerID, payload.Reason)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.RoleAssigned, result, fiber.StatusOK)
}

func (h *StaffHandler) RevokeRole(c fiber.Ctx) error {
	staffID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	performedBy, ok := middleware.CallerStaffID(c)
	if !ok {
		return response.Error(c, sysmsg.Unauthorized, fiber.StatusUnauthorized)
	}
	callerID := performedBy.(uuid.UUID)

	var payload staff.RevokeRolePayload
	if err := c.Bind().Body(&payload); err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	result, err := h.usecase.RevokeRole(c.Context(), staffID, payload.RoleID, callerID, payload.Reason)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.RoleRevoked, result, fiber.StatusOK)
}

func (h *StaffHandler) UpdateRole(c fiber.Ctx) error {
	staffID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	performedBy, ok := middleware.CallerStaffID(c)
	if !ok {
		return response.Error(c, sysmsg.Unauthorized, fiber.StatusUnauthorized)
	}
	callerID := performedBy.(uuid.UUID)

	var payload staff.UpdateRolePayload
	if err := c.Bind().Body(&payload); err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	result, err := h.usecase.UpdateRole(c.Context(), staffID, callerID, &payload)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.RoleAssigned, result, fiber.StatusOK)
}

func (h *StaffHandler) GetRoleHistory(c fiber.Ctx) error {
	staffID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	history, err := h.usecase.GetRoleHistory(c.Context(), staffID)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.RoleFetched, history, fiber.StatusOK)
}
