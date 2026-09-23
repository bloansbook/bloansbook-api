package handler

import (
	"strings"

	"github.com/bloansbook/bloansbook-api/internal/auth/middleware"
	"github.com/bloansbook/bloansbook-api/internal/customers/usecase"
	"github.com/bloansbook/bloansbook-api/internal/models"
	"github.com/bloansbook/bloansbook-api/internal/models/customers"
	"github.com/bloansbook/bloansbook-api/pkg/response"
	"github.com/bloansbook/bloansbook-api/pkg/sysmsg"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CustomerHandler struct {
	usecase *usecase.CustomerUsecase
}

func NewCustomerHandler(u *usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		usecase: u,
	}
}

func (h *CustomerHandler) CreateCustomer(c fiber.Ctx) error {
	var payload customers.CreateCustomerPayload

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

	customer, err := h.usecase.CreateCustomer(c.Context(), id, &payload)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.CustomerCreated, customer, fiber.StatusCreated)
}

func (h *CustomerHandler) GetCustomerById(c fiber.Ctx) error {
	customerID := c.Params("id")
	if customerID == "" {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	id, err := uuid.Parse(customerID)
	if err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	customer, err := h.usecase.GetCustomerByID(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.Error(c, sysmsg.CustomerNotFound, fiber.StatusNotFound)
		}
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.CustomerFetched, customer, fiber.StatusOK)
}

func (h *CustomerHandler) GetAllCustomers(c fiber.Ctx) error {
	filter := customers.CustomerFilter{
		Search:    fiber.Query(c, "search", ""),
		Type:      fiber.Query(c, "type", ""),
		Status:    fiber.Query(c, "status", ""),
		SortBy:    fiber.Query(c, "sortBy", "createdAt"),
		SortOrder: fiber.Query(c, "sortOrder", "desc"),
		Limit:     fiber.Query(c, "limit", 20),
		Offset:    fiber.Query(c, "offset", 0),
	}

	totalCount, err := h.usecase.GetCustomerCount(c.Context())
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	customerList, err := h.usecase.GetAllCustomers(c.Context(), filter)
	if err != nil {
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	data := models.DataWithPagination{
		Data:       customerList,
		Count:      len(customerList),
		TotalCount: totalCount,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
	}

	return response.Success(c, sysmsg.CustomerListFetched, data, fiber.StatusOK)
}

func (h *CustomerHandler) UpdateCustomer(c fiber.Ctx) error {
	customerID := c.Params("id")

	var payload customers.UpdateCustomerPayload
	if err := c.Bind().Body(&payload); err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	if customerID == "" {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	id, err := uuid.Parse(customerID)
	if err != nil {
		return response.Error(c, sysmsg.BadRequest, fiber.StatusBadRequest)
	}

	customer, err := h.usecase.UpdateCustomer(c.Context(), id, &payload)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.Error(c, sysmsg.CustomerNotFound, fiber.StatusNotFound)
		}
		return response.Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	return response.Success(c, sysmsg.CustomerUpdated, customer, fiber.StatusOK)
}
