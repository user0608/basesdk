package usecases

import (
	"basesdk/errs"
	"basesdk/security/dtos"
	"basesdk/security/models"
	"basesdk/security/repositories"
	"context"
	"strings"
	"time"
)

type TenantsUsecase struct {
	repository *repositories.TenantRepository
}

func NewTenantsUsecase(repository *repositories.TenantRepository) *TenantsUsecase {
	return &TenantsUsecase{repository: repository}
}

func (u *TenantsUsecase) List(ctx context.Context) ([]dtos.TenantResponse, error) {
	tenants, err := u.repository.FindTenants(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dtos.TenantResponse, 0, len(tenants))
	for _, tenant := range tenants {
		response = append(response, toTenantResponse(tenant))
	}

	return response, nil
}

func (u *TenantsUsecase) Find(ctx context.Context, codigo string) (*dtos.TenantResponse, error) {
	tenant, err := u.repository.FindTenant(ctx, strings.TrimSpace(codigo))
	if err != nil {
		return nil, err
	}

	response := toTenantResponse(*tenant)
	return &response, nil
}

func (u *TenantsUsecase) Create(ctx context.Context, input dtos.CreateTenantInput, createdBy string) error {
	tenant := &models.Tenant{
		Codigo:         strings.TrimSpace(input.Codigo),
		Name:           strings.TrimSpace(input.Name),
		Timezone:       strings.TrimSpace(input.Timezone),
		MaxActiveUsers: input.MaxActiveUsers,
		ExpiresAt:      input.ExpiresAt,
		Disabled:       false,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
	}
	if tenant.Codigo == "" || tenant.Name == "" || tenant.Timezone == "" {
		return errs.BadRequestDirect("codigo, nombre y zona horaria son requeridos")
	}
	if tenant.MaxActiveUsers < 1 {
		return errs.BadRequestDirect("maxActiveUsers debe ser mayor a cero")
	}

	return u.repository.CreateTenant(ctx, tenant)
}

func (u *TenantsUsecase) Update(ctx context.Context, codigo string, input dtos.UpdateTenantInput, updatedBy string) error {
	now := time.Now()
	tenant := &models.Tenant{
		Codigo:         strings.TrimSpace(codigo),
		Name:           strings.TrimSpace(input.Name),
		Timezone:       strings.TrimSpace(input.Timezone),
		MaxActiveUsers: input.MaxActiveUsers,
		Disabled:       input.Disabled,
		ExpiresAt:      input.ExpiresAt,
		UpdatedBy:      &updatedBy,
		UpdatedAt:      &now,
	}
	if tenant.Codigo == "" || tenant.Name == "" || tenant.Timezone == "" {
		return errs.BadRequestDirect("codigo, nombre y zona horaria son requeridos")
	}
	if tenant.MaxActiveUsers < 1 {
		return errs.BadRequestDirect("maxActiveUsers debe ser mayor a cero")
	}

	return u.repository.UpdateTenant(ctx, tenant)
}

func (u *TenantsUsecase) Enable(ctx context.Context, codigos []string, updatedBy string) error {
	if len(codigos) == 0 {
		return errs.BadRequestDirect("tenants requeridos")
	}
	return u.repository.SetTenantsDisabled(ctx, codigos, false, updatedBy)
}

func (u *TenantsUsecase) Disable(ctx context.Context, codigos []string, updatedBy string) error {
	if len(codigos) == 0 {
		return errs.BadRequestDirect("tenants requeridos")
	}
	return u.repository.SetTenantsDisabled(ctx, codigos, true, updatedBy)
}
