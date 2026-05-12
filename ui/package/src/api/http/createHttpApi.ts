import { HttpConnectionError, HttpInvalidJsonError, HttpRequestError } from "./errors";
import type { CreateHttpApiOptions, HttpApi, HttpApiAuthMode, HttpApiRequest } from "./types";

const joinUrl = (baseUrl: string, path: string) => {
  const normalizedBase = baseUrl.endsWith("/") ? baseUrl.slice(0, -1) : baseUrl;
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;
  return `${normalizedBase}${normalizedPath}`;
};

const appendQuery = (url: URL, query?: Record<string, unknown>) => {
  if (!query) return;

  for (const [key, value] of Object.entries(query)) {
    if (value == null) continue;

    if (Array.isArray(value)) {
      for (const item of value) {
        if (item != null) url.searchParams.append(key, String(item));
      }
      continue;
    }

    url.searchParams.set(key, String(value));
  }
};

const mergeHeaders = (baseHeaders?: HeadersInit, nextHeaders?: HeadersInit) => {
  const headers = new Headers(baseHeaders ?? {});
  const overrides = new Headers(nextHeaders ?? {});

  overrides.forEach((value, key) => {
    headers.set(key, value);
  });

  return headers;
};

const normalizeAuthorizationToken = (token: string | null) => {
  if (!token) return null;

  const normalizedToken = token.trim();
  if (normalizedToken === "") return null;
  if (/^Bearer\s+/i.test(normalizedToken)) return normalizedToken;

  return `Bearer ${normalizedToken}`;
};

const resolveToken = async (
  auth: HttpApiAuthMode,
  options: CreateHttpApiOptions,
  tokenOverride?: string,
) => {
  if (tokenOverride) return normalizeAuthorizationToken(tokenOverride);
  if (auth === "none") return null;
  if (auth === "system") return normalizeAuthorizationToken((await options.getSystemToken?.()) ?? null);
  return normalizeAuthorizationToken((await options.getTenantToken?.()) ?? null);
};

const parseJson = async <T>(response: Response) => {
  try {
    return (await response.json()) as T;
  } catch (error) {
    if (error instanceof SyntaxError) {
      throw new HttpInvalidJsonError();
    }

    throw error;
  }
};

export const createHttpApi = (options: CreateHttpApiOptions): HttpApi => {
  const request = async <T>({
    path,
    method = "GET",
    auth = "tenant",
    query,
    data,
    formData,
    headers,
    signal,
    responseShape = "data",
    tokenOverride,
  }: HttpApiRequest): Promise<T> => {
    const baseUrl = await options.getBaseUrl();
    const url = new URL(joinUrl(baseUrl, path));
    appendQuery(url, query);

    const mergedHeaders = mergeHeaders(options.defaultHeaders, headers);
    const token = await resolveToken(auth, options, tokenOverride);

    if (token) {
      mergedHeaders.set("Authorization", token);
    }

    let body: BodyInit | undefined;

    if (formData) {
      body = formData;
      mergedHeaders.delete("Content-Type");
    } else if (data !== undefined) {
      body = JSON.stringify(data);

      if (!mergedHeaders.has("Content-Type")) {
        mergedHeaders.set("Content-Type", "application/json");
      }
    }

    let response: Response;

    try {
      response = await fetch(url, {
        method,
        headers: mergedHeaders,
        body,
        signal,
      });
    } catch (error) {
      if (error instanceof TypeError) {
        throw new HttpConnectionError();
      }

      throw error;
    }

    if (responseShape === "raw") {
      if (!response.ok) {
        if ((response.status === 401 || response.status === 403) && auth !== "none") {
          await options.onUnauthorized?.(auth);
        }

        throw new HttpRequestError(response.statusText || "Request error", response.status);
      }

      return response as T;
    }

    const decoded = await parseJson<{ data?: T; message?: string } & T>(response);

    if (!response.ok) {
      if ((response.status === 401 || response.status === 403) && auth !== "none") {
        await options.onUnauthorized?.(auth);
      }

      const message = typeof decoded === "object" && decoded && "message" in decoded && typeof decoded.message === "string"
        ? decoded.message
        : response.statusText || "Request error";

      throw new HttpRequestError(message, response.status, decoded);
    }

    if (responseShape === "json") {
      return decoded as T;
    }

    return (decoded as { data: T }).data;
  };

  return {
    request,
    get: (input) => request({ ...input, method: "GET" }),
    post: (input) => request({ ...input, method: "POST" }),
    put: (input) => request({ ...input, method: "PUT" }),
    patch: (input) => request({ ...input, method: "PATCH" }),
    delete: (input) => request({ ...input, method: "DELETE" }),
  };
};
