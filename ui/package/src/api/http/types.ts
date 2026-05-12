export type HttpApiAuthMode = "tenant" | "system" | "none";

export type HttpApiMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export type HttpApiResponseShape = "data" | "json" | "raw";

export type HttpApiRequest = {
  path: string;
  method?: HttpApiMethod;
  auth?: HttpApiAuthMode;
  query?: Record<string, unknown>;
  data?: unknown;
  formData?: FormData;
  headers?: HeadersInit;
  signal?: AbortSignal;
  responseShape?: HttpApiResponseShape;
  tokenOverride?: string;
};

export type CreateHttpApiOptions = {
  getBaseUrl: () => string | Promise<string>;
  getTenantToken?: () => string | null | Promise<string | null>;
  getSystemToken?: () => string | null | Promise<string | null>;
  onUnauthorized?: (auth: Exclude<HttpApiAuthMode, "none">) => void | Promise<void>;
  defaultHeaders?: HeadersInit;
};

export type HttpApi = {
  request: <T>(request: HttpApiRequest) => Promise<T>;
  get: <T>(request: Omit<HttpApiRequest, "method">) => Promise<T>;
  post: <T>(request: Omit<HttpApiRequest, "method">) => Promise<T>;
  put: <T>(request: Omit<HttpApiRequest, "method">) => Promise<T>;
  patch: <T>(request: Omit<HttpApiRequest, "method">) => Promise<T>;
  delete: <T>(request: Omit<HttpApiRequest, "method">) => Promise<T>;
};
