# Skill: `connection`

Package for database connection access and transaction propagation using Gorm.

Use this skill when creating repositories, use cases that need transactions, migrations, or any code that needs database access.

## Main API

The main abstraction is `connection.StorageManager`.

```go
type StorageManager interface {
	Conn(ctx context.Context) *gorm.DB
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
```

Use `Conn(ctx)` to retrieve the current `*gorm.DB` connection.

Use `WithTx(ctx, fn)` to run a unit of work inside a transaction.

## How `Conn` Works

`Conn(ctx)` is context-aware.

- If `ctx` contains an active transaction, it returns that transaction.
- If `ctx` does not contain a transaction, it returns the base database connection.
- In both cases, the returned `*gorm.DB` is bound to the provided context.

This means repository methods do not need to know whether they are running inside a transaction.

```go
tx := r.manager.Conn(ctx)
```

Always use the `ctx` received by the current method. Do not replace it with `context.Background()` or a new unrelated context.

## Repository Pattern

Repositories should receive `connection.StorageManager` through their constructor and store it as `manager`.

```go
type ProductRepository struct {
	manager connection.StorageManager
}

func NewProductRepository(manager connection.StorageManager) *ProductRepository {
	return &ProductRepository{
		manager: manager,
	}
}
```

Every repository method should call `manager.Conn(ctx)` at the start of the method.

```go
func (r *ProductRepository) Find(ctx context.Context, id string) (*models.Product, error) {
	tx := r.manager.Conn(ctx)

	var product models.Product
	rs := tx.Where("id = ?", id).First(&product)
	if rs.Error != nil {
		return nil, errs.Pgf(rs.Error)
	}

	return &product, nil
}
```

Convert Gorm and Postgres errors with `errs.Pgf` before returning them.

## Create Pattern

```go
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	tx := r.manager.Conn(ctx)

	rs := tx.Create(product)
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	return nil
}
```

## Update Pattern

For updates that must fail when the row does not exist, check `RowsAffected`.

```go
func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	tx := r.manager.Conn(ctx)

	rs := tx.Model(&models.Product{}).
		Where("id = ?", product.ID).
		Select("*").
		Updates(product)

	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	if rs.RowsAffected == 0 {
		return errs.NotFoundDirect("product not found")
	}

	return nil
}
```

## Delete Pattern

```go
func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	tx := r.manager.Conn(ctx)

	rs := tx.Delete(&models.Product{}, "id = ?", id)
	if rs.Error != nil {
		return errs.Pgf(rs.Error)
	}

	return nil
}
```

Check `RowsAffected` when deleting a missing row should be reported as `404`.

## Transactions

Use `manager.WithTx(ctx, fn)` in use cases when multiple repository operations must commit or rollback together.

```go
func (u *ProductUsecase) Create(ctx context.Context, req CreateProductRequest) error {
	return u.manager.WithTx(ctx, func(ctx context.Context) error {
		product := &models.Product{
			Name: req.Name,
		}

		if err := u.products.Create(ctx, product); err != nil {
			return err
		}

		return u.audit.Create(ctx, "product_created", product.ID)
	})
}
```

The callback receives a new context containing the transaction. Pass that callback context to every repository method that must participate in the transaction.

## Nested Transactions

`WithTx` is safe to call when a transaction already exists in the context.

If `ctx` already contains a transaction, `WithTx` does not create a new one. It simply executes `fn(ctx)`.

This allows use cases to call other use cases without accidentally creating nested Gorm transactions.

```go
return u.manager.WithTx(ctx, func(ctx context.Context) error {
	if err := u.firstOperation(ctx); err != nil {
		return err
	}

	return u.secondOperation(ctx)
})
```

## Transaction Propagation Rule

The transaction is propagated only through `context.Context`.

Correct:

```go
return u.manager.WithTx(ctx, func(ctx context.Context) error {
	return u.repo.Create(ctx, entity)
})
```

Incorrect:

```go
return u.manager.WithTx(ctx, func(ctx context.Context) error {
	return u.repo.Create(context.Background(), entity)
})
```

Using a new context breaks transaction propagation because `manager.Conn(ctx)` will not find the transaction.

## Getting `*sql.DB`

Only infrastructure code should normally retrieve the underlying `*sql.DB`, for example migrations.

```go
tx := manager.Conn(ctx)
db, err := tx.DB()
if err != nil {
	return err
}
```

Application repositories should use Gorm through `manager.Conn(ctx)` instead.

## Configuration

`NewConnection(config configs.DatabaseConfig)` creates the storage manager.

The Gorm connection uses:

- Postgres driver.
- `SkipDefaultTransaction: true`.
- Log level from database config.
- Singular table names through `schema.NamingStrategy{SingularTable: true}`.

`NewConnection` is provided through Fx in setup code.

```go
fx.Provide(connection.NewConnection)
```

## SkipStorage

`connection.SkipStorage()` returns a `StorageManager` that disables persistence.

- `Conn(ctx)` returns `nil`.
- `WithTx(ctx, fn)` just executes `fn(ctx)` without a transaction.

Use it only for flows where storage is intentionally disabled. Do not use it in normal repository code unless the caller is prepared for `Conn(ctx)` to be `nil`.

## Recommendations

- Inject `connection.StorageManager`, not `*gorm.DB`.
- Store the dependency as `manager connection.StorageManager`.
- Call `manager.Conn(ctx)` inside each repository method.
- Use `manager.WithTx(ctx, fn)` at the use-case layer for multi-step writes.
- Pass the callback context from `WithTx` to every repository call inside the transaction.
- Convert database errors with `errs.Pgf` before returning them.
- Check `RowsAffected` when update or delete operations must detect missing records.
- Keep transaction boundaries out of repositories unless there is a very specific infrastructure reason.

## Common Mistakes

- Do not pass `*gorm.DB` manually across layers.
- Do not call `db.Transaction` directly in repositories.
- Do not use `context.Background()` inside a transaction callback.
- Do not ignore the callback context passed by `WithTx`.
- Do not return raw Gorm or Postgres errors from repositories.
- Do not assume `WithTx` creates nested transactions; it reuses the existing transaction when one is already present.
- Do not call `Conn(ctx)` from a different context if the operation must be part of the active transaction.
