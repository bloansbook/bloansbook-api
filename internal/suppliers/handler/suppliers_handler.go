package handler

import (
	"strings"

	"github.com/bloansbook/bloansbook-api/internal/auth/middleware"
	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/bloansbook/bloansbook-api/internal/models/suppliers"
	"github.com/bloansbook/bloansbook-api/internal/suppliers/usecase"
	"github.com/bloansbook/bloansbook-api/pkg/response"
	"github.com/bloansbook/bloansbook-api/pkg/sysmsg"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type SupplierHandler struct {
	usecase *usecase.SupplierUsecase
}

func NewSupplierHandler(u *usecase.SupplierUsecase) *SupplierHandler {
	return &SupplierHandler{
		usecase: u,
	}
}

func (h *SupplierHandler) CreateSupplier(c fiber.Ctx) error {
	var payload suppliers.CreateSupplierPayload

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

	supplier, err := h.usecase.CreateSupplier(c.Context(), id, &payload)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.SupplierCreated, supplier, fiber.StatusCreated)
}

func (h *SupplierHandler) GetSupplierById(c fiber.Ctx) error {
	supplierID := c.Params("id")
	if supplierID == "" {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	id, err := uuid.Parse(supplierID)
	if err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	supplier, err := h.usecase.GetSupplierByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.Error(c, sysmsg.SupplierNotFound, fiber.StatusNotFound)
		}
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.SupplierFetched, supplier, fiber.StatusOK)
}

func (h *SupplierHandler) GetAllSuppliers(c fiber.Ctx) error {
	filter := suppliers.SupplierFilter{
		Search:    fiber.Query(c, "search", ""),
		Category:  fiber.Query(c, "category", ""),
		Status:    fiber.Query(c, "status", ""),
		SortBy:    fiber.Query(c, "sortBy", "createdAt"),
		SortOrder: fiber.Query(c, "sortOrder", "desc"),
		Limit:     fiber.Query(c, "limit", 20),
		Offset:    fiber.Query(c, "offset", 0),
	}

	totalCount, err := h.usecase.GetSupplierCount(c.Context())
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	supplierList, err := h.usecase.GetAllSuppliers(c.Context(), filter)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	data := models.DataWithPagination{
		Data:       supplierList,
		Count:      len(supplierList),
		TotalCount: totalCount,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
	}

	return response.Success(c, sysmsg.SupplierListFetched, data, fiber.StatusOK)
}

func (h *SupplierHandler) UpdateSupplier(c fiber.Ctx) error {
	supplierID := c.Params("id")

	var payload suppliers.UpdateSupplierPayload
	if err := c.Bind().Body(&payload); err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	if supplierID == "" {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	id, err := uuid.Parse(supplierID)
	if err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	supplier, err := h.usecase.UpdateSupplier(c.Context(), id, &payload)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.Error(c, sysmsg.SupplierNotFound, fiber.StatusNotFound)
		}
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.SupplierUpdated, supplier, fiber.StatusOK)
}
