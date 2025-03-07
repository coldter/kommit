/**
 * @description
 */
export function getDomainFromRequest(request: Request) {
  // Check various headers in order of preference
  const forwardedHost = request.headers.get("X-Forwarded-Host");
  const originalHost = request.headers.get("X-Original-Host");
  const host = request.headers.get("Host");

  // Return the first non-null/non-empty header
  return host || forwardedHost || originalHost || new URL(request.url).host;
}
