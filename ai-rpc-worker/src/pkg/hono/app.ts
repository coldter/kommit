import { handleError, handleZodError } from "@/pkg/errors/http";
import type { HonoEnv } from "@/pkg/hono/env";
import { OpenAPIHono } from "@hono/zod-openapi";
import type { Context as GenericContext } from "hono";
import { contextStorage } from "hono/context-storage";
import { prettyJSON } from "hono/pretty-json";
import { trimTrailingSlash } from "hono/trailing-slash";

function setupCommonIgnoreRoutes<T extends HonoEnv>(app: OpenAPIHono<T>) {
  app.get("/favicon.ico", (c) =>
    c.body(null, {
      status: 204,
      headers: {
        "Cache-Control": "public, max-age=604800", // 1 week
      },
    }),
  );
}

export function newApp() {
  // defaultHook
  const app = new OpenAPIHono<HonoEnv>({
    defaultHook: handleZodError,
  });

  app.use(contextStorage());
  setupCommonIgnoreRoutes(app);

  app.use(trimTrailingSlash());
  app.use(prettyJSON());
  app.onError(handleError);

  app.notFound((c) => {
    return c.json(
      {
        success: false,
        error: {
          code: "NOT_FOUND",
          message: "Not Found",
        },
      },
      404,
    );
  });

  return app;
}

export type App = ReturnType<typeof newApp>;
export type Context = GenericContext<HonoEnv>;
