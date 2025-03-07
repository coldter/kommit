import type { createRoute } from "@hono/zod-openapi";
import { z } from "@hono/zod-openapi";
import { failWithErrorSchema } from "./common-schemas";

type ResponseConfig = Parameters<typeof createRoute>[0]["responses"];

export const successWithoutDataSchema = z.object({
  success: z.boolean(),
});

export const withPaginationSchema = <T extends z.ZodTypeAny>(schema: T) => {
  return z.object({
    items: schema.array(),
    total: z.number(),
  });
};

export const withOutPaginationSchema = <T extends z.ZodTypeAny>(schema: T) => {
  return z.object({
    items: schema.array(),
  });
};

export const successWithDataSchema = <T extends z.ZodTypeAny>(schema: T) =>
  z.object({ success: z.boolean(), data: schema });

export const errorResponses = {
  400: {
    description: "Bad request: problem processing request.",
    content: {
      "application/json": {
        schema: failWithErrorSchema,
      },
    },
  },
  401: {
    description: "Unauthorized: authentication required.",
    content: {
      "application/json": {
        schema: failWithErrorSchema,
      },
    },
  },
  403: {
    description: "Forbidden: insufficient permissions.",
    content: {
      "application/json": {
        schema: failWithErrorSchema,
      },
    },
  },
  404: {
    description: "Not found: resource does not exist.",
    content: {
      "application/json": {
        schema: failWithErrorSchema,
      },
    },
  },
  500: {
    description: "Server error: something went wrong.",
    content: {
      "application/json": {
        schema: failWithErrorSchema,
      },
    },
  },
} satisfies ResponseConfig;
