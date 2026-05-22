package repositories

import (
	"basesdk/connection"
	"basesdk/errs"
	"basesdk/security/models"
	"context"
	"time"
)

type TenantRepository struct {
	manager connection.StorageManager
}

func NewTenantRepository(manager connection.StorageManager) *TenantRepository {
	return &TenantRepository{manager: manager}
}

func (r *TenantRepository) FindTenants(ctx context.Context) ([]models.Tenant, error) {
	var tenants []models.Tenant
	rs := r.manager.Conn(ctx).Raw(`
		select
			t.*,
			(select count(*) from app_user u where u.tenant_codigo = t.codigo) as users_count,
			(select count(*) from role r where r.tenant_codigo = t.codigo) as roles_count,
			(select count(*) from app_group g where g.tenant_codigo = t.codigo) as groups_count
		from tenant t
		order by t.codigo
	`).Scan(&tenants)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return tenants, nil
}

func (r *TenantRepository) FindTenant(ctx context.Context, codigo string) (*models.Tenant, error) {
	var tenant models.Tenant
	rs := r.manager.Conn(ctx).Raw(`
		select
			t.*,
			(select count(*) from app_user u where u.tenant_codigo = t.codigo) as users_count,
			(select count(*) from role r where r.tenant_codigo = t.codigo) as roles_count,
			(select count(*) from app_group g where g.tenant_codigo = t.codigo) as groups_count
		from tenant t
		where t.codigo = ?
	`, codigo).Scan(&tenant)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return nil, errs.NotFoundDirect("tenant no encontrado")
	}

	return &tenant, nil
}

func (r *TenantRepository) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	return errs.Pgf(r.manager.Conn(ctx).Create(tenant).Error)
}

func (r *TenantRepository) UpdateTenant(ctx context.Context, tenant *models.Tenant) error {
	rs := r.manager.Conn(ctx).Model(&models.Tenant{}).
		Where("codigo = ?", tenant.Codigo).
		Updates(map[string]any{
			"name":             tenant.Name,
			"timezone":         tenant.Timezone,
			"max_active_users": tenant.MaxActiveUsers,
			"disabled":         tenant.Disabled,
			"expires_at":       tenant.ExpiresAt,
			"updated_by":       tenant.UpdatedBy,
			"updated_at":       tenant.UpdatedAt,
		})
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}
	if rs.RowsAffected == 0 {
		return errs.NotFoundDirect("tenant no encontrado")
	}

	return nil
}

func (r *TenantRepository) SetTenantsDisabled(ctx context.Context, codigos []string, disabled bool, updatedBy string) error {
	now := time.Now()
	rs := r.manager.Conn(ctx).Model(&models.Tenant{}).
		Where("codigo in ?", codigos).
		Updates(map[string]any{
			"disabled":   disabled,
			"updated_by": updatedBy,
			"updated_at": now,
		})
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	return nil
}
