package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/bloansbook/bloansbook-api/internal/models/suppliers"
	"github.com/bloansbook/bloansbook-api/internal/suppliers/repository"
	"github.com/bloansbook/bloansbook-api/pkg/config"
	"github.com/bloansbook/bloansbook-api/pkg/idgen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplierUsecase struct {
	repository *repository.SupplierRepository
	db         *pgxpool.Pool
	config     *config.Config
}

func NewSupplierUsecase(db *pgxpool.Pool, repo *repository.SupplierRepository, config *config.Config) *SupplierUsecase {
	return &SupplierUsecase{
		repository: repo,
		db:         db,
		config:     config,
	}
}

func (u *SupplierUsecase) GetSupplierCount(ctx context.Context) (int, error) {
	return u.repository.CountSuppliers(ctx)
}

// CreateSupplier generates the sequential SPL-000x id and inserts the supplier.
// created_by is the authenticated caller. category/currency values are left to the
// DB CHECK constraints; we only guard that a caller and the required fields are present.
func (u *SupplierUsecase) CreateSupplier(ctx context.Context, createdBy uuid.UUID, payload *suppliers.CreateSupplierPayload) (*suppliers.CreateSupplierResponse, error) {
	if createdBy == uuid.Nil {
		return nil, fmt.Errorf("authentication required: createdBy is missing")
	}

	payload.Name = strings.TrimSpace(payload.Name)
	payload.Phone = strings.TrimSpace(payload.Phone)
	payload.Category = suppliers.SupplierCategory(strings.TrimSpace(string(payload.Category)))

	if payload.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if payload.Phone == "" {
		return nil, fmt.Errorf("phone is required")
	}
	if payload.Category == "" {
		return nil, fmt.Errorf("category is required")
	}

	supplierID, err := idgen.GenerateSequentialID(ctx, u.db, "suppliers", "supplier_id", u.config.IDGen.SupplierPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to generate supplier id: %w", err)
	}

	m, err := u.repository.CreateSupplier(ctx, supplierID, createdBy, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}

	return &suppliers.CreateSupplierResponse{
		ID:        m.ID,
		Supplier:  m.ToSummary(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

func (u *SupplierUsecase) GetSupplierByID(ctx context.Context, id uuid.UUID) (*suppliers.SupplierDTO, error) {
	m, err := u.repository.GetSupplierByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("supplier not found: %w", err)
	}

	dto := m.ToDTO()
	return &dto, nil
}

func (u *SupplierUsecase) GetAllSuppliers(ctx context.Context, filter suppliers.SupplierFilter) ([]suppliers.SupplierDTO, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	list, err := u.repository.GetAllSuppliers(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get suppliers: %w", err)
	}

	dtos := make([]suppliers.SupplierDTO, len(list))
	for i, m := range list {
		dtos[i] = m.ToDTO()
	}
	return dtos, nil
}

func (u *SupplierUsecase) UpdateSupplier(ctx context.Context, id uuid.UUID, payload *suppliers.UpdateSupplierPayload) (*suppliers.UpdateSupplierResponse, error) {
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
	if payload.Category != nil {
		trimmed := suppliers.SupplierCategory(strings.TrimSpace(string(*payload.Category)))
		if trimmed == "" {
			return nil, fmt.Errorf("category cannot be empty")
		}
		payload.Category = &trimmed
	}

	if _, err := u.repository.GetSupplierByID(ctx, id); err != nil {
		return nil, fmt.Errorf("cannot update non-existent supplier: %w", err)
	}

	m, err := u.repository.UpdateSupplier(ctx, id, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update supplier: %w", err)
	}

	return &suppliers.UpdateSupplierResponse{
		ID:        m.ID,
		Supplier:  m.ToSummary(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}
