# Skill: `errs`

Paquete para construir errores de aplicación con:

- Código HTTP asociado.
- Mensaje público (apto para API).
- Error interno opcional (para logging/diagnóstico).

También incluye mapeo de errores comunes de Postgres (pgx/gorm) a mensajes y códigos HTTP.

## Tipo base

`*errs.Err` implementa `error` y expone:

- `Message() string` mensaje público.
- `Code() int` código HTTP.
- `Wrapped() error` error interno (puede ser `nil`).

`answer.UnwrapErr` reconoce `*errs.Err` y lo convierte a respuesta HTTP.

## Constructores comunes

Versiones que envuelven un error interno:

- `BadRequestError(err, format, args...)` → 400.
- `NotFoundError(err, format, args...)` → 404.
- `UnauthorizedError(err, format, args...)` → 401.
- `ForbiddenError(err, format, args...)` → 403.
- `UnsupportedMediaTypeError(err, format, args...)` → 415.
- `InternalError(err, format, args...)` → 500.

Versiones sin error interno (solo mensaje):

- `BadRequestf(format, args...)`, `NotFoundf(...)`, `InternalErrorf(...)`.
- `BadRequestDirect(message)`, `NotFoundDirect(message)`, `InternalErrorDirect(message)`, `UnauthorizedDirect(message)`, `ForbiddenDirect(message)`, `UnsupportedMediaTypeDirect(message)`.

Otras utilidades:

- `NewWithMessage(err, message)` cambia el mensaje manteniendo el código (si `err` ya es `*Err`); si no, asume 400.
- `WrapError(err, message, httpCode)` retorna `nil` si `err == nil`.
- `ContainsMessage(err, substr)` busca en `Message()` cuando `err` es `*Err`.
- `IsBadRequest(err)`, `IsInternalError(err)`, `IsErr(err)`.
- `ToSummary(err)` produce un resumen estable (`"Error <code>: <message>"`).

## Postgres / Gorm (`postgres_errors.go`)

`Pgf(err)` convierte errores de DB a `*errs.Err` con mensajes amigables:

- `gorm.ErrRecordNotFound` → 400 con `ErrRecordNotFound`.
- Errores `*pgconn.PgError` (pgx) se mapean por código SQLSTATE (`PGCode`) a mensaje/código HTTP.
- Caso especial: `23503` con mensaje de `insert or update` usa un mensaje más específico (`message23503`).

Seguridad y logging:

- Cada código tiene `loggable bool`.
- Si `loggable == true` o `devmode == true`, el `*Err` conserva el error interno envuelto.
- Si no, se omite el error interno (`wrapped=nil`) para evitar filtrar detalles.

Personalización:

- `AddPgErrs(code, message, httpCode, loggable)` agrega/override del mapeo.
- `Devmode()` fuerza `devmode=true`.
- `IsPgErrCode(err, code)` detecta el SQLSTATE.

## Ejemplo (servicio)

```go
func (s *Service) Get(id string) (*User, error) {
  u, err := s.repo.FindByID(id)
  if err != nil {
    return nil, errs.Pgf(err) // o errs.NotFoundError/BadRequestError según aplique
  }
  return u, nil
}
```

## Notas

- Si el error va a salir por HTTP, construye un `*errs.Err` con mensaje pensado para el cliente.
- Para detalles internos (stack, SQL, etc.), envuelve el error pero controla `loggable/devmode`.
