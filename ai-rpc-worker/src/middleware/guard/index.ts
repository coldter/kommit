import type { MiddlewareHandler } from "hono";

export const isPublicAccess: MiddlewareHandler = async (_, next) => {
  await next();
};
