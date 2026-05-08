# Skill: `kcheck`

Validador ligero para structs en Go basado en tags.

## Qué resuelve

- Validación declarativa de DTOs con tags (por defecto `chk`).
- Validación de structs anidados.
- Opciones para validar solo algunos campos (select) o excluir campos (skip).
- Errores agregados por campo con formato estable.

## Uso rápido

```go
type CreateUser struct {
  Name  string `chk:"required min=2 max=50"`
  Email string `chk:"required email"`
  Age   int    `chk:"gte=18 lte=120"`
}

if err := kcheck.Valid(CreateUser{...}); err != nil {
  // err.Error(): "Name: ...; Email: ..."
}
```

## API

- `kcheck.Valid(input, skips...)` valida el struct usando el validador por defecto, omitiendo campos.
- `kcheck.ValidSelect(input, selected...)` valida solo los campos indicados.
- `kcheck.Struct(input)` valida todo (equivalente al validador por defecto).

También puedes crear un validador aislado:

```go
v := kcheck.New()
v.Register("startsx", func(f kcheck.Field) error { ... })
err := v.Struct(dto)
```

## Tags y reglas

El tag (por defecto `chk`) admite varias reglas separadas por espacios:

- Sin parámetro: `required`, `email`, `uuid`, `url`, `ip`, `lower`, `upper`, etc.
- Con parámetro: `min=2`, `max=50`, `len=6`, `oneof=a,b,c`, `prefix=USR-`, `gte=18`.

Si una regla no está registrada, el error será: `"validador [X] no registrado"`.

## Structs anidados

- Se hace "dive" a campos `struct` (incluyendo punteros a struct no nulos).
- Excepción: `time.Time` no se considera struct anidado para recorrer.

## Skip vs Select

- Skip: `kcheck.Valid(user, "Email")` omite `Email`.
- Select: `kcheck.ValidSelect(user, "Address.City")` valida solo ese path.
- Puedes pasar tanto nombres (`"Email"`) como paths (`"Address.City"`).

## Errores

El error agregado es `kcheck.Errors` (con `[]FieldError`) y su string:

- Formato: `"<Path>: <Message>; <Path>: <Message>"`.
- `Errors.Err()` retorna `nil` si no hubo errores.

## Validadores incluidos

Registrados por defecto en `RegisterDefaults()`:

- Requerido: `required`, `nonil`.
- Longitud/tamaño: `len`, `min`, `max`.
- Comparadores numéricos: `gt`, `gte`, `lt`, `lte`.
- Strings: `alpha`, `alphanum`, `num`, `decimal`, `lower`, `upper`.
- Formato: `email`, `uuid` (v4), `url`, `ip`, `ipv4`, `ipv6`.
- Strings avanzados: `oneof`, `prefix`, `suffix`, `contains`.
- Fechas: `date` (DateOnly), `time` (TimeOnly), `datetime` (DateTime), `utc` (RFC3339 en UTC o `time.Time` en UTC).

## Notas

- Input inválido retorna `kcheck.ErrInvalidInput` (nil, puntero nil, o no-struct).
- `required` soporta `string`, `[]string` y punteros nil (a través de `IsNil`).
