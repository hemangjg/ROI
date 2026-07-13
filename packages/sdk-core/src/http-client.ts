import { createBearerHeader, resolveAuthToken } from "./auth.js";
import { FinOpsHttpError } from "./errors.js";

export type RequestAuth = "api-key" | "access-token" | "none";

export type HttpClientConfig = {
  baseUrl: string;
  apiKey?: string;
  accessToken?: string;
  timeoutMs?: number;
  fetch?: typeof fetch;
  userAgent?: string;
};

export type RequestOptions = {
  auth?: RequestAuth;
};

export class HttpClient {
  private readonly baseUrl: string;
  private readonly apiKey?: string;
  private readonly accessToken?: string;
  private readonly timeoutMs: number;
  private readonly fetchFn: typeof fetch;
  private readonly userAgent: string;

  constructor(config: HttpClientConfig) {
    if (!config.baseUrl) {
      throw new Error("baseUrl is required");
    }

    this.baseUrl = config.baseUrl.replace(/\/$/, "");
    this.apiKey = config.apiKey;
    this.accessToken = config.accessToken;
    this.timeoutMs = config.timeoutMs ?? 30_000;
    this.fetchFn = config.fetch ?? globalThis.fetch;

    if (!this.fetchFn) {
      throw new Error("fetch is not available in this runtime");
    }

    this.userAgent = config.userAgent ?? "ai-finops-sdk-ts/0.0.0";
  }

  getBaseUrl(): string {
    return this.baseUrl;
  }

  async get<T>(path: string, options?: RequestOptions): Promise<T> {
    return this.request<T>("GET", path, undefined, options);
  }

  async post<T>(path: string, body?: unknown, options?: RequestOptions): Promise<T> {
    return this.request<T>("POST", path, body, options);
  }

  async request<T>(
    method: string,
    path: string,
    body?: unknown,
    options?: RequestOptions,
  ): Promise<T> {
    const url = `${this.baseUrl}${path.startsWith("/") ? path : `/${path}`}`;
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), this.timeoutMs);

    try {
      const headers: Record<string, string> = {
        Accept: "application/json",
        "User-Agent": this.userAgent,
      };

      if (body !== undefined) {
        headers["Content-Type"] = "application/json";
      }

      const authMode = options?.auth ?? "api-key";
      if (authMode !== "none") {
        const token =
          authMode === "access-token"
            ? this.accessToken
            : resolveAuthToken(this.apiKey, this.accessToken);

        if (token) {
          headers.Authorization = createBearerHeader(token);
        }
      }

      const response = await this.fetchFn(url, {
        method,
        headers,
        body: body !== undefined ? JSON.stringify(body) : undefined,
        signal: controller.signal,
      });

      if (!response.ok) {
        throw await FinOpsHttpError.fromResponse(response);
      }

      if (response.status === 204 || response.headers.get("content-length") === "0") {
        return undefined as T;
      }

      return (await response.json()) as T;
    } catch (error) {
      if (error instanceof Error && error.name === "AbortError") {
        throw new FinOpsHttpError(`Request timed out after ${this.timeoutMs}ms`, 408);
      }

      throw error;
    } finally {
      clearTimeout(timer);
    }
  }
}
