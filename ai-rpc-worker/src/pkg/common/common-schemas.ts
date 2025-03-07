import { ErrorSchema } from "@/pkg/errors/http";
import { z } from "@hono/zod-openapi";
import type { BaseMime } from "hono/utils/mime";

export const failWithErrorSchema = ErrorSchema;

/**
 * @description
 */
export const fileRequestSchema = <const L extends number = 1>(
  limit: L = 1 as L,
  allowedMimeTypes?: BaseMime[],
) =>
  z
    .custom<L extends 1 ? File : File[]>((v) => {
      if (Array.isArray(v)) {
        return v.every((file) => file instanceof File);
      }

      return v instanceof File;
    })
    .refine(
      (v) => {
        if (Array.isArray(v)) {
          if (!allowedMimeTypes) {
            return true;
          }
          return v.every((file) => allowedMimeTypes.includes(file.type as any));
        }
        //
        if (!allowedMimeTypes) {
          return true;
        }

        return allowedMimeTypes.includes(v.type as any);
      },
      {
        message: `Invalid file type supplied. Allowed types: ${allowedMimeTypes?.join(", ")}`,
      },
    )
    .refine(
      (v) => {
        if (Array.isArray(v)) {
          if (limit === 0) {
            return true;
          }
          if (limit === 1) {
            return false;
          }
          return v.length <= limit;
        }

        return limit === 1;
      },
      {
        message: `Only ${limit} file(s) allowed, or only one file provided`,
      },
    )
    .openapi({
      type: "string",
      format: "binary",
    });
