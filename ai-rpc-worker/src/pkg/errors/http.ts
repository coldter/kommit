import type { HonoEnv } from "@/pkg/hono/env";
import { parseZodErrorMessage } from "@/pkg/utils/zod-error";
import { z } from "@hono/zod-openapi";
import type { Context } from "hono";
import { HTTPException } from "hono/http-exception";
import type { StatusCode } from "hono/utils/http-status";
import { ZodError } from "zod";

const ErrorCode = z.enum([
  "BAD_REQUEST",
  "FORBIDDEN",
  "INTERNAL_SERVER_ERROR",
  "USAGE_EXCEEDED",
  "DISABLED",
  "NOT_FOUND",
  "NOT_UNIQUE",
  "RATE_LIMITED",
  "UNAUTHORIZED",
  "PRECONDITION_FAILED",
  "INSUFFICIENT_PERMISSIONS",
  "METHOD_NOT_ALLOWED",
  "EXPIRED",
  "ROUTE_NOT_FOUND",
]);

export type TErrorCode = z.infer<typeof ErrorCode>;
function codeToStatus(code: TErrorCode) {
  switch (code) {
    case "BAD_REQUEST":
      return 400;
    case "FORBIDDEN":
    case "DISABLED":
    case "UNAUTHORIZED":
    case "INSUFFICIENT_PERMISSIONS":
    case "USAGE_EXCEEDED":
    case "EXPIRED":
      return 403;
    case "NOT_FOUND":
    case "ROUTE_NOT_FOUND":
      return 404;
    case "METHOD_NOT_ALLOWED":
      return 405;
    case "NOT_UNIQUE":
      return 409;
    case "PRECONDITION_FAILED":
      return 412;
    case "RATE_LIMITED":
      return 429;
    case "INTERNAL_SERVER_ERROR":
      return 500;
  }
}

function statusToCode(status: StatusCode): z.infer<typeof ErrorCode> {
  switch (status) {
    case 400:
      return "BAD_REQUEST";
    case 401:
      return "UNAUTHORIZED";
    case 403:
      return "FORBIDDEN";
    case 404:
      return "NOT_FOUND";
    case 405:
      return "METHOD_NOT_ALLOWED";
    case 500:
      return "INTERNAL_SERVER_ERROR";
    default:
      return "INTERNAL_SERVER_ERROR";
  }
}

export class ApiError extends HTTPException {
  public readonly code: TErrorCode;

  constructor({ code, message }: { code: z.infer<typeof ErrorCode>; message: string }) {
    super(codeToStatus(code), { message });
    this.code = code;
  }
}

export function handleZodError(
  result:
    | {
        success: true;
        data: any;
      }
    | {
        success: false;
        error: ZodError;
      },
  c: Context,
) {
  if (!result.success) {
    return c.json<z.infer<typeof ErrorSchema>>(
      {
        success: false,
        error: {
          code: "BAD_REQUEST",
          message: parseZodErrorMessage(result.error),
        },
        requestId: c.get("requestId"),
      },
      { status: 400 },
    );
  }
  return;
}

export async function handleError(err: Error, c: Context<HonoEnv>): Promise<Response> {
  if (err instanceof ZodError) {
    return handleZodError({ success: false, error: err }, c) as Response;
  }

  /**
   * We can handle this very well, as it is something we threw ourselves
   */
  if (err instanceof ApiError) {
    if (err.status >= 500) {
      console.error(err.message, {
        name: err.name,
        code: err.code,
        status: err.status,
      });
    }

    return c.json<z.infer<typeof ErrorSchema>>(
      {
        success: false,
        error: {
          code: err.code,
          message: err.message,
        },
      },
      { status: err.status },
    );
  }

  /**
   * HTTPExceptions from hono at least give us some idea of what to do as they provide a status and
   * message
   */
  if (err instanceof HTTPException) {
    if (err.status >= 500) {
      console.error("HTTPException", {
        message: err.message,
        status: err.status,
      });
    }
    const code = statusToCode(err.status);
    return c.json<z.infer<typeof ErrorSchema>>(
      {
        success: false,
        error: {
          code,
          message: err.message,
        },
      },
      { status: err.status },
    );
  }

  /**
   * We're lost here, all we can do is return a 500 and log it to investigate
   */
  console.error("unhandled exception", {
    name: err?.name,
    message: err?.message,
    cause: err?.cause,
    stack: err?.stack,
    constructor: err?.constructor.name,
  });

  return c.json(
    {
      success: false,
      error: {
        code: "INTERNAL_SERVER_ERROR",
        message: err.message ?? "something unexpected happened",
      },
    },
    { status: 500 },
  );
}

export const ErrorSchema = z.object({
  success: z.boolean(),
  // success: z.literal(false),
  error: z.object({
    code: ErrorCode.openapi({
      description: "A machine readable error code.",
      example: "INTERNAL_SERVER_ERROR",
    }),
    message: z.string().openapi({ description: "A human readable explanation of what went wrong" }),
  }),
  requestId: z.string().optional().openapi({
    description: "Please always include the requestId in your error report",
    example: "req_1234",
  }),
});

export type TErrorSchema = z.infer<typeof ErrorSchema>;
