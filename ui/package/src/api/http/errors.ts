export class HttpConnectionError extends Error {
  constructor(message = "No se pudo establecer conexion con el servicio") {
    super(message);
    this.name = "HttpConnectionError";
  }
}

export class HttpInvalidJsonError extends Error {
  constructor(message = "La respuesta del servicio no contiene un JSON valido") {
    super(message);
    this.name = "HttpInvalidJsonError";
  }
}

export class HttpRequestError extends Error {
  status: number;
  payload: unknown;

  constructor(message: string, status: number, payload?: unknown) {
    super(message);
    this.name = "HttpRequestError";
    this.status = status;
    this.payload = payload;
  }
}
