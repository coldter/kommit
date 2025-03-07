import { init } from "@/middleware/init";
import { docs } from "@/pkg/docs";
import { ApiError } from "@/pkg/errors/http";
import { type Context, newApp } from "@/pkg/hono/app";
import type { Env, HonoEnv } from "@/pkg/hono/env";
import { setupApiRoutes } from "@/routes";
import { WorkersKVStore } from "@hono-rate-limiter/cloudflare";
import type { Next } from "hono";
import { rateLimiter } from "hono-rate-limiter";
import { cors } from "hono/cors";
import { secureHeaders } from "hono/secure-headers";

const app = newApp();

app.use("*", init());
app.use("*", cors({ credentials: true, origin: "*" }));
app.use(
  "*",
  secureHeaders({
    crossOriginResourcePolicy: false,
  }),
);

// * Rate limiter
const authTokenKVKey = "auth:bypass:rate-limit:tokens";
app.use((c: Context, next: Next) =>
  rateLimiter<HonoEnv>({
    windowMs: 60 * 60 * 1000, // 1 hour
    // windowMs: 60 * 1000 * 2, // 2 minute
    limit: 5,
    keyGenerator: (c) => c.req.header("cf-connecting-ip") ?? "", // Method to generate custom identifiers for clients.
    store: new WorkersKVStore({ namespace: c.env.RATE_LIMIT_CACHE }), // Here CACHE is your WorkersKV Binding.,
    async skip(c) {
      const bearerToken = c.req.header("authorization")?.split(" ")[1];
      if (!bearerToken) {
        return false;
      }
      const tokens = await c.env.RATE_LIMIT_CACHE.get(authTokenKVKey);

      if (!tokens) {
        return false;
      }
      const tokensArray = tokens.split(",");
      return tokensArray.includes(bearerToken);
    },
    handler: () => {
      throw new ApiError({
        code: "RATE_LIMITED",
        message: "Too many requests, please try again later.",
      });
    },
  })(c, next),
);

// * Register API routes
setupApiRoutes(app);

// Init OpenAPI docs
docs(app, true);

export default {
  async fetch(request: Request, env: Env, context?: ExecutionContext) {
    return app.fetch(request, env, context);
  },
} satisfies ExportedHandler<Env>;
