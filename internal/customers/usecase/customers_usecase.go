package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/bloansbook/bloansbook-api/internal/customers/repository"
	"github.com/bloansbook/bloansbook-api/internal/models/customers"
	"github.com/bloansbook/bloansbook-api/pkg/config"
	"github.com/bloansbook/bloansbook-api/pkg/idgen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerUsecase struct {
	repository *repository.CustomerRepository
	db         *pgxpool.Pool
	config     *config.Config
}

func NewCustomerUsecase(db *pgxpool.Pool, repo *repository.CustomerRepository, config *config.Config) *CustomerUsecase {
	return &CustomerUsecase{
		repository: repo,
		db:         db,
		config:     config,
	}
}

func (u *CustomerUsecase) GetCustomerCount(ctx context.Context) (int, error) {
	return u.repository.CountCustomers(ctx)
}

// CreateCustomer generates the sequential CUST-000x id and inserts the customer.
// created_by is the authenticated caller. type/currency values are left to the
// DB CHECK constraints; we only guard that a caller is present.
func (u *CustomerUsecase) CreateCustomer(ctx context.Context, createdBy uuid.UUID, payload *customers.CreateCustomerPayload) (*customers.CreateCustomerResponse, error) {
	if createdBy == uuid.Nil {
		return nil, fmt.Errorf("authentication required: createdBy is missing")
	}

	payload.Name = strings.TrimSpace(payload.Name)
	payload.Phone = strings.TrimSpace(payload.Phone)
	payload.Type = customers.CustomerType(strings.TrimSpace(string(payload.Type)))

	if payload.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if payload.Phone == "" {
		return nil, fmt.Errorf("phone is required")
	}
	if payload.Type == "" {
		return nil, fmt.Errorf("type is required")
	}

	customerID, err := idgen.GenerateSequentialID(ctx, u.db, "customers", "customer_id", u.config.IDGen.CustomerPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to generate customer id: %w", err)
	}

	m, err := u.repository.CreateCustomer(ctx, customerID, createdBy, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return &customers.CreateCustomerResponse{
		ID:        m.ID,
		Customer:  m.ToSummary(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

func (u *CustomerUsecase) GetCustomerByID(ctx context.Context, id uuid.UUID) (*customers.CustomerDTO, error) {
	m, err := u.repository.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	dto := m.ToDTO()
	return &dto, nil
}

func (u *CustomerUsecase) GetAllCustomers(ctx context.Context, filter customers.CustomerFilter) ([]customers.CustomerDTO, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	list, err := u.repository.GetAllCustomers(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}

	dtos := make([]customers.CustomerDTO, len(list))
	for i, m := range list {
		dtos[i] = m.ToDTO()
	}
	return dtos, nil
}

func (u *CustomerUsecase) UpdateCustomer(ctx context.Context, id uuid.UUID, payload *customers.UpdateCustomerPayload) (*customers.UpdateCustomerResponse, error) {
	if payload.Name != nil {
		trimmed := strings.TrimSpace(*payload.Name)
		if trimmed == "" {
			return nil, fmt.Errorf("name cannot be empty")
		}
		payload.Name = &trimmed
	}
	if payload.Phone != nil {
		trimmed := strings.TrimSpace(*payload.Phone)
		if trimmed == "" {
			return nil, fmt.Errorf("phone cannot be empty")
		}
		payload.Phone = &trimmed
	}
	if payload.Type != nil {
		trimmed := customers.CustomerType(strings.TrimSpace(string(*payload.Type)))
		if trimmed == "" {
			return nil, fmt.Errorf("type cannot be empty")
		}
		payload.Type = &trimmed
	}

	if _, err := u.repository.GetCustomerByID(ctx, id); err != nil {
		return nil, fmt.Errorf("cannot update non-existent customer: %w", err)
	}

	m, err := u.repository.UpdateCustomer(ctx, id, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return &customers.UpdateCustomerResponse{
		ID:        m.ID,
		Customer:  m.ToSummary(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}
