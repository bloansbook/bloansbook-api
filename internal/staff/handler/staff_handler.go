package handler

import (
	"strings"

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
	filter := staff.StaffFilter{
		Search:     fiber.Query(c, "search", ""),
		Status:     fiber.Query(c, "status", ""),
		Department: fiber.Query(c, "department", ""),
		SortBy:     fiber.Query(c, "sortBy", "createdAt"),
		SortOrder:  fiber.Query(c, "sortOrder", "desc"),
		Limit:      fiber.Query(c, "limit", 20),
		Offset:     fiber.Query(c, "offset", 0),
	}

	totalCount, err := h.usecase.GetStaffCount(c.Context())
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	staffList, err := h.usecase.GetAllStaff(c.Context(), filter)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	data := models.DataWithPagination{
		Data:       staffList,
		Count:      len(staffList),
		TotalCount: totalCount,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
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

func (h *StaffHandler) FireStaff(c fiber.Ctx) error {
	staffID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	caller, ok := middleware.CallerStaffID(c)
	if !ok {
		return response.Error(c, sysmsg.Unauthorized, fiber.StatusUnauthorized)
	}
	callerID := caller.(uuid.UUID)

	var payload staff.FireStaffPayload
	if err := c.Bind().Body(&payload); err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	result, err := h.usecase.FireStaff(c.Context(), staffID, &payload, callerID)
	if err != nil {
		if strings.Contains(err.Error(), "already terminated") {
			return response.Error(c, sysmsg.StaffAlreadyFired, fiber.StatusConflict)
		}
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.StaffTerminated, result, fiber.StatusOK)
}

func (h *StaffHandler) OverrideTermination(c fiber.Ctx) error {
	staffID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	caller, ok := middleware.CallerStaffID(c)
	if !ok {
		return response.Error(c, sysmsg.Unauthorized, fiber.StatusUnauthorized)
	}
	callerID := caller.(uuid.UUID)

	var payload staff.OverrideTerminationPayload
	if err := c.Bind().Body(&payload); err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	result, err := h.usecase.OverrideTermination(c.Context(), staffID, &payload, callerID)
	if err != nil {
		if strings.Contains(err.Error(), "no active termination record") {
			return response.Error(c, sysmsg.StaffNotFound, fiber.StatusNotFound)
		}
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.StaffTerminationOverride, result, fiber.StatusOK)
}
