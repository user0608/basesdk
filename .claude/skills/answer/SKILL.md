# Skill: `answer`

Paquete para estandarizar respuestas HTTP JSON en handlers.

## Qué resuelve

- Respuestas de éxito consistentes (`data` o `message`).
- Respuestas de error consistentes (`message`) a partir de errores de dominio.
- Mapeo de errores a códigos HTTP.
- Log de errores internos inesperados.

## Contrato

El paquete opera sobre un `Target` (adaptable a tu framework HTTP):

```go
type Target interface {
  JSON(code int, i any) error
  NoContent(code int) error
}
```

Para integrarlo, normalmente envuelves el contexto de tu framework para que implemente `JSON`/`NoContent`.

## API principal

- `Ok(c, payload)` → `200` con `{ data: payload }`.
- `Created(c, payload)` → `201` con `{ data: payload }`.
- `Accepted(c, payload)` → `202` con `{ data: payload }`.
- `Message(c, msg)` → `200` con `{ message: msg }`.
- `CreatedMessage(c, msg)` → `201` con `{ message: msg }`.
- `AcceptedMessage(c, msg)` → `202` con `{ message: msg }`.
- `Success(c)` → `200` con mensaje estándar.
- `CreatedSuccess(c)` → `201` con mensaje estándar.
- `AcceptedSuccess(c)` → `202` con mensaje estándar.
- `NoContent(c)` → `204` sin cuerpo.

Errores:

- `Err(c, err)` → responde según `UnwrapErr(err)`.
- `Auto(c, err)` → si `err != nil` usa `Err`, si no usa `Success`.

## Mapeo de errores (`UnwrapErr`)

`UnwrapErr(err)` retorna `(code, message)` con estas reglas:

- Si `err` es `*errs.Err`: usa `err.Code()` y `err.Message()`. Si además tiene `Wrapped()`, ese error interno se loguea con `slog.Error`.
- Si `err.Error()` empieza con `":"` (prefijo): se trata como error de cliente: `400` y el mensaje público es el contenido después de `":"`.
- Cualquier otro error no esperado: devuelve `500` con mensaje genérico y loguea el error con `slog.Error`.

## Ejemplo

```go
func (h *Handler) CreateUser(c answer.Target, req CreateUserReq) error {
  user, err := h.svc.CreateUser(req)
  if err != nil {
    return answer.Err(c, err)
  }
  return answer.Created(c, user)
}
```

## Notas

- Evita devolver errores arbitrarios con mensajes sensibles. Si necesitas un mensaje público con código HTTP, usa `errs`.
- El prefijo `":"` es un escape rápido para 400, pero para APIs grandes suele ser preferible `errs.BadRequest...` por consistencia.
