import type { App } from "@/pkg/hono/app";
import type { HonoEnv } from "@/pkg/hono/env";
import { getDomainFromRequest } from "@/pkg/utils/uri";
import type { OpenAPIObjectConfigure } from "@hono/zod-openapi";
import { apiReference } from "@scalar/hono-api-reference";

type ExtractConfig<T> = T extends (...args: any[]) => any ? never : T;
type OpenAPIObjectConfig = ExtractConfig<OpenAPIObjectConfigure<any, any>>;
type TagObject = NonNullable<OpenAPIObjectConfig["tags"]>[number];

export const docs = (app: App, enable: boolean) => {
  if (!enable) {
    return;
  }

  const registry = app.openAPIRegistry;

  registry.registerComponent("securitySchemes", "bearerAuth", {
    scheme: "bearer",
    description: "session token issued by server",
    type: "http",
    bearerFormat: "JWT",
  });

  const commonTags: TagObject[] = [];

  app.doc31("/openapi.json", {
    servers: [{ url: "http://localhost:3100" }],
    info: {
      title: "Api Reference",
      version: "v1",
    },
    openapi: "3.1.0",
    tags: commonTags,
    security: [{ bearerAuth: [] }],
  });

  app.get("/docs", (c) => {
    const currentRequestDomain = getDomainFromRequest(c.req.raw);
    return apiReference<HonoEnv>({
      spec: {
        url: "openapi.json",
      },
      theme: "deepSpace",
      servers: [
        {
          url: `http://${currentRequestDomain}`,
          description: "Current",
        },
        {
          url: `https://${currentRequestDomain}`,
          description: "Current HTTPS",
        },
        {
          url: "http://localhost:3100",
          description: "Localhost",
        },
        {
          url: "{CUSTOM_URL}",
          description: "Custom",
          variables: {
            CUSTOM_URL: {
              default: "http://localhost:3100",
            },
          },
        },
      ],
    })(c, async () => {});
  });
};
