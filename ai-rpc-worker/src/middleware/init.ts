import type { HonoEnv } from "@/pkg/hono/env";
import type { MiddlewareHandler } from "hono";

export function init(): MiddlewareHandler<HonoEnv> {
  return async (_c, next) => {
    await next();
  };
}
