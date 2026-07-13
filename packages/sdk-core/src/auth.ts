const API_KEY_PREFIX = "aif_";

export function isApiKey(value: string): boolean {
  return value.startsWith(API_KEY_PREFIX);
}

export function createBearerHeader(token: string): string {
  return `Bearer ${token}`;
}

export function resolveAuthToken(apiKey?: string, accessToken?: string): string | undefined {
  return apiKey ?? accessToken;
}
