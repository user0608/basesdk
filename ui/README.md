# UI workspace

Workspace aislado para la libreria UI y su playground local.

## Estructura

- `package/`: libreria de componentes UI.
- `example/`: app de ejemplo para probar la libreria localmente.

## Uso

```bash
npm install
npm run dev
```

El `example` usa `VITE_API_URL` y por defecto queda configurado en `ui/example/.env` con:

```bash
VITE_API_URL=http://localhost:7622
```

Comandos utiles:

```bash
npm run dev
npm run build
npm run build:package
npm run build:example
```
